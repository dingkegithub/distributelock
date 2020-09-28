package zk

import (
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/dingkegithub/distrubutelock/utils/log"
	"github.com/samuel/go-zookeeper/zk"
)

type ZkClient struct {
	cli    *zk.Conn
	option *zkOptions
	logger log.Logging
}

func NewZkClient(opt *zkOptions) *ZkClient {
	return &ZkClient{
		cli:    nil,
		option: opt,
		logger: opt.logger,
	}
}

func (zc *ZkClient) Open(retry int) error {
	cli, ev, err := zk.Connect(zc.option.addrs, time.Second)
	if err != nil {
		if retry > 0 {
			return zc.Open(retry - 1)
		} else {
			return err
		}
	}

	if s, ok := <-ev; ok {
		zc.logger.Log("file", "zkclient.go", "func", "Open", "msg", "open connect event", "type", s.Type.String())
	} else {
		zc.logger.Log("file", "zkclient.go", "func", "Open", "msg", "event closed")
	}

	zc.cli = cli
	return nil
}

func (zc *ZkClient) Close() {
	if zc.cli != nil {
		zc.cli.Close()
		zc.cli = nil
	}
}

func (zc *ZkClient) CreateFirstEphemeral(dir, value string, block bool, tmoutable bool, timeout time.Duration) (bool, string, error) {
	node, err := zc.createEphemeralNode(dir, value)
	if err != nil {
		return false, "", err
	}

	interval := time.After(timeout)

	for {
		nodes, _, event, err := zc.cli.ChildrenW(dir)
		if err != nil {
			zc.logger.Log("file", "zkclient.go", "func", "TimeoutWaitFirstEphemeral", "msg", "children watch", "error", err)
			return false, node, err
		}

		ok, err := zc.isValidMinNode(node, nodes, value)
		if err != nil {
			return false, node, err
		}

		if ok {
			return true, node, nil
		}

		if !block && !tmoutable {
			return false, node, nil
		}

		select {
		case ev := <-event:
			zc.logger.Log("file", "zkclient.go", "func", "BlockWaitFirstEphemeral", "msg", "recv watch event", "type", ev.Type.String(), "ok", ok)

			nodes, _, err := zc.cli.Children(dir)
			if err != nil {
				zc.logger.Log("file", "zkclient.go", "func", "TimeoutWaitFirstEphemeral", "msg", "get children", "error", err)
				return false, node, err
			}

			ok, err := zc.isValidMinNode(node, nodes, value)
			if err != nil {
				return false, node, err
			}

			if ok {
				return true, node, nil
			}
		case <-interval:
			if tmoutable {
				return false, node, nil
			}
			interval = time.After(timeout)
		}
	}
}

func (zc *ZkClient) createEphemeralNode(dir string, value string) (string, error) {
	ok, _, err := zc.cli.Exists(dir)
	if err != nil {
		zc.logger.Log("file", "zkclient.go", "func", "CreateEphemeralNode", "msg", "zk exist failed", "error", err)
		return "", err
	}

	if !ok {
		res, err := zc.cli.Create(dir, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil && err != zk.ErrNodeExists {
			zc.logger.Log("file", "zkclient.go", "func", "CreateEphemeralNode", "msg", "zk create failed", "error", err)
			return "", err
		}

		if dir != res {
			zc.logger.Log("file", "zkclient.go", "func", "CreateEphemeralNode", "msg", "create node failed", "want", dir, "actual", res)
			return "", zk.ErrInvalidPath
		}
	}

	emphemeralNode := path.Join(dir, value)

	str, err := zc.cli.CreateProtectedEphemeralSequential(emphemeralNode, nil, zk.WorldACL(zk.PermAll))
	if err != nil {
		zc.logger.Log("file", "zkclient.go", "func", "CreateEphemeralNode", "msg", "create ephemeral node failed", "error", err)
		return "", err
	}

	if len(str) <= 0 {
		zc.logger.Log("file", "zkclient.go", "func", "CreateEphemeralNode", "msg", "ephemeral node name blank")
		return "", zk.ErrNoNode
	}

	return path.Base(str), nil
}

func (zc *ZkClient) DeleteEphemeralNode(dir string, value string) {

	zc.logger.Log("file", "zkclient.go", "func", "DeleteEphemeralNode", "msg", "clear node", "value", value)
	if len(value) <= 0 {
		return
	}
	ephemeralNode := path.Join(dir, value)
	zc.logger.Log("file", "zkclient.go", "func", "DeleteEphemeralNode", "msg", ephemeralNode)
	zc.cli.Delete(ephemeralNode, 0)
}

func (zc *ZkClient) isValidMinNode(node string, nodes []string, commonField string) (bool, error) {
	if len(node) <= 0 {
		zc.logger.Log("file", "zkclient.go", "func", "isMinNode", "msg", "node blank exception")
		return false, zk.ErrNodeExists
	}

	if !strings.Contains(node, commonField) {
		zc.logger.Log("file", "zkclient.go", "func", "isMinNode", "msg", "invalid node has not common field")
		return false, zk.ErrNodeExists
	}

	if !zc.nodeInNodes(node, nodes) {
		zc.logger.Log("file", "zkclient.go", "func", "isMinNode", "msg", "node node exists in children node")
		return false, zk.ErrNodeExists
	}

	minNum, err := zc.nodeSeqNum(node, commonField)
	if err != nil {
		zc.logger.Log("file", "zkclient.go", "func", "isMinNode", "msg", "get node seq num failed", "error", err)
		return false, err
	}

	for _, n := range nodes {
		if !strings.Contains(n, commonField) {
			zc.logger.Log("file", "zkclient.go", "func", "isMinNode", "msg", "unexpect emphemeral node exists in children", "node", n)
			continue
		}

		num, err := zc.nodeSeqNum(n, commonField)
		if err != nil {
			zc.logger.Log("file", "zkclient.go", "func", "isMinNode", "msg", "children seq num failed", "node", n, "error", err)
			continue
		}

		if num < minNum {
			return false, nil
		}
	}

	return true, nil
}

func (zc *ZkClient) nodeInNodes(node string, nodes []string) bool {
	for _, n := range nodes {
		if node == n {
			return true
		}
	}
	return false
}

func (zc *ZkClient) nodeSeqNum(node string, commonField string) (uint64, error) {
	nodeFields := strings.Split(node, commonField)
	if len(nodeFields) != 2 {
		return 0, zk.ErrNoNode
	}

	nodeNum, err := strconv.ParseUint(nodeFields[1], 10, 64)
	if err != nil {
		return 0, zk.ErrNoNode
	}

	return nodeNum, nil
}

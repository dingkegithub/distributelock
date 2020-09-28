package etcdv2

import "fmt"

type Node struct {
	Key           string `json:"key"`
	Value         string `json:"value"`
	ModifiedIndex int    `json:"modifiedIndex"`
	CreatedIndex  int    `json:"createdIndex"`
}

func (n Node) String() string {
	return fmt.Sprintf("key: %s, Value: %s, ModifiedIndex: %d, CreatedIndex: %d",
		n.Key, n.Value, n.ModifiedIndex, n.CreatedIndex)
}

type AtomicCreateResp struct {
	Action string `json:"action"`
	Node   *Node  `json:"node"`
}

func (a AtomicCreateResp) String() string {

	if a.Node == nil {
		return fmt.Sprintf("Action: %s, Node: nil", a.Action)
	}
	return fmt.Sprintf("Action: %s, Node: %s", a.Action, a.Node.String())
}

type AtomicDeleteResp struct {
	Action   string `json:"action"`
	Node     *Node  `json:"node"`
	PrevNode *Node  `json:"prevNode"`
}

// {"action":"update","node":{"key":"/foo","value":"","expiration":"2020-09-23T08:32:05.0389929Z","ttl":30,"modifiedIndex":15,"createdIndex":7},"prevNode":{"key":"/foo","value":"","expiration":"2020-09-23T08:32:03.8015632Z","ttl":29,"modifiedIndex":14,"createdIndex":7}}

type ErrResp struct {
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message"`
	Cause     string `json:"cause"`
	Index     int    `json:"index"`
}

package etcdv2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	cluserutils "github.com/dingkegithub/distributelock/pkg/clusterutils"
)

type EtcdClient struct {
	mutex *sync.Mutex
	nm    *cluserutils.ClusterNodeManager
	opts  *EtcdOptions
}

func NewEtcdClient(opts *EtcdOptions) (*EtcdClient, error) {
	rand.Seed(time.Now().Unix())

	nm, err := cluserutils.NewClusterNodeManager(opts.heartInterval, opts.logger, opts.Addrs...)
	if err != nil {
		return nil, err
	}

	client := &EtcdClient{
		nm:    nm,
		opts:  opts,
		mutex: &sync.Mutex{},
	}
	return client, nil
}

// AtomicLease: lease key by flush ttl
// curl http://x.x.x.x:x/v2/keys/key -XPUT -d ttl=5 -d value=value refresh=true -d prevExist=true
func (ec *EtcdClient) AtomicLease(key string, value string, ttl uint64) (*AtomicCreateResp, error) {
	form := url.Values(make(map[string][]string))
	form.Add("prevExist", "true")
	form.Add("refresh", "true")
	form.Add("value", value)
	form.Add("ttl", fmt.Sprintf("%d", ttl))
	return ec.put(key, form.Encode(), ec.opts.retry)
}

// AtomicCreate: compare and set
// prevExist - checks existence of the key:
// if prevExist is true, it is an update request;
// if prevExist is false, it is a create request
// curl http://x.x.x.x:x/v2/keys/key?prevExist=false -XPUT -d value=value ttl=5
func (ec *EtcdClient) AtomicCreate(key string, value string, ttl uint64) (*AtomicCreateResp, error) {
	form := url.Values{}
	form.Add("value", value)
	form.Add("prevExist", "false")
	form.Add("ttl", fmt.Sprintf("%d", ttl))
	return ec.put(key, form.Encode(), ec.opts.retry)
}

// Atomic: atomic operation, compare and delete
// curl http://x.x.x.x:x/v2/keys/key?prevValue=preValue -XDELETE
func (ec *EtcdClient) AtomicDelete(key string, preValue string) (*AtomicCreateResp, error) {
	return ec.del(key, preValue, 0)
}

func (ec *EtcdClient) del(key, preValue string, retry int) (*AtomicCreateResp, error) {
	node, err := ec.nm.Random()
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("http://%s/v2/keys/%s?prevValue=%s", node, key, preValue)

	req, err := http.NewRequest(http.MethodDelete, addr, nil)
	if err == nil {
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()

			buffer, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				if strings.Contains(string(buffer), "errorCode") {
					return nil, ErrorExist
				} else if strings.Contains(string(buffer), "action") {
					var respMsg *AtomicCreateResp
					err = json.Unmarshal(buffer, &respMsg)
					if err != nil {
						return nil, err
					}
					return respMsg, nil
				} else {
					ec.opts.logger.Log("file", "etcdhttpclient.go", "func", "put", "msg", ErrorUnknownEtcdMsg.Error(), "msg", string(buffer))
					return nil, ErrorUnknownEtcdMsg
				}
			} else {
				ec.opts.logger.Log("file", "etcdhttpclient.go", "func", "put", "msg", "read http response body failed", "error", err)
			}
		} else {
			ec.opts.logger.Log("file", "etcdhttpclient.go", "func", "put", "msg", "send http request failed", "error", err)
		}
	}

	if retry > 0 {
		time.Sleep(3 * time.Millisecond)
		return ec.del(key, preValue, retry-1)
	}

	return nil, ErrorCommit
}

func (ec *EtcdClient) put(key, form string, retry int) (*AtomicCreateResp, error) {
	node, err := ec.nm.Random()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s/v2/keys/%s?%s", node, key, form)

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err == nil {
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			defer resp.Body.Close()

			buffer, err := ioutil.ReadAll(resp.Body)
			//fmt.Println("file", "etcdhttpclient.go", "func", "put", "msg", "buffer readed", "buffer: ", string(buffer), "url: ", url, "err", err)
			if err == nil {
				if strings.Contains(string(buffer), "errorCode") {
					return nil, ErrorExist
				} else if strings.Contains(string(buffer), "action") {
					var respMsg *AtomicCreateResp
					err = json.Unmarshal(buffer, &respMsg)
					if err != nil {
						return nil, err
					}
					return respMsg, nil
				} else {
					ec.opts.logger.Log("file", "etcdhttpclient.go", "func", "put", "msg", ErrorUnknownEtcdMsg.Error(), "msg", string(buffer))
					return nil, ErrorUnknownEtcdMsg
				}
			}
		}
	}

	if retry > 0 {
		time.Sleep(3 * time.Millisecond)
		return ec.put(key, form, retry-1)
	}

	return nil, ErrorCommit
}

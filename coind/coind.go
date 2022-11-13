package coind

import (
	"encoding/json"
	_ "errors"
	_ "strconv"
)

const (
	// RpcclientTimeout VERSION represents bicoind package version
	// Version = 0.1
	// RpcclientTimeout DEFAULT_RPCCLIENT_TIMEOUT represent http timeout for rcp client
	RpcclientTimeout = 15
)

// Coind A Bitcoind represents a Bitcoind client
type Coind struct {
	client *rpcClient
}

// New return a new bitcoind
func New(host string, port int, user, passwd string, useSSL bool, timeoutParam ...int) (*Coind, error) {
	var timeout = RpcclientTimeout
	// If the timeout is specified in timeoutParam, allow it.
	if len(timeoutParam) != 0 {
		timeout = timeoutParam[0]
	}

	rpcClient, err := newClient(host, port, user, passwd, useSSL, timeout)
	if err != nil {
		return nil, err
	} else {
		return &Coind{rpcClient}, nil
	}
}

func (b *Coind) Call(command string, par ...any) ([]byte, error) {
	var r rpcResponse
	var err error

	if len(par) > 0 {
		switch x := par[0].(type) {
		case []interface{}:
			if len(x) == 0 {
				r, err = b.client.call(command, interface{}(nil))
			} else {
				r, err = b.client.call(command, par[0])
			}
		default:
			r, err = b.client.call(command, par)
		}
	} else {
		r, err = b.client.call(command, par)
	}
	if err != nil {
		return nil, err
	}
	if command != "sendtoaddress" {
		marshal, err := json.Marshal(r.Result)
		if err != nil {
			return nil, err
		}
		return marshal, nil
	}
	return r.Result, nil
}

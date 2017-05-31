package rpc

import (
	etcd "github.com/coreos/etcd/client"
	wlib "github.com/wothing/wonaming/lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

// EtcdWatcher is the implementaion of grpc.naming.Watcher
type EtcdWatcher struct {
	key     string
	keyapi  etcd.KeysAPI
	watcher etcd.Watcher
	addrs   []string
	ctx     context.Context
	cancel  context.CancelFunc
}

func (w *EtcdWatcher) Close() {
	w.cancel()
	<-w.ctx.Done()
}

func newEtcdWatcher(key string, cli etcd.Client) naming.Watcher {

	api := etcd.NewKeysAPI(cli)
	watcher := api.Watcher(key, &etcd.WatcherOptions{Recursive: true})
	ctx, cancel := context.WithCancel(context.Background())

	w := &EtcdWatcher{
		key:     key,
		keyapi:  api,
		watcher: watcher,
		ctx:     ctx,
		cancel:  cancel,
	}
	return w
}

func (w *EtcdWatcher) Next() ([]*naming.Update, error) {

	if w.addrs == nil {
		resp, _ := w.keyapi.Get(context.Background(), w.key, &etcd.GetOptions{Recursive: true})
		addrs := extractAddrs(resp)

		if len(addrs) != 0 {
			w.addrs = addrs
			return wlib.GenUpdates([]string{}, addrs), nil
		}
	}

	for {
		_, err := w.watcher.Next(w.ctx)
		if err == nil {
			resp, err := w.keyapi.Get(w.ctx, w.key, &etcd.GetOptions{Recursive: true})
			if err != nil {
				continue
			}

			addrs := extractAddrs(resp)

			updates := wlib.GenUpdates(w.addrs, addrs)
			w.addrs = addrs
			if len(updates) != 0 {
				return updates, nil
			}
		}
	}

	return []*naming.Update{}, nil
}

func extractAddrs(resp *etcd.Response) (addrs []string) {
	addrs = []string{}

	if resp == nil || resp.Node == nil || resp.Node.Nodes == nil || len(resp.Node.Nodes) == 0 {
		return addrs
	}

	for _, node := range resp.Node.Nodes {
		addr := node.Value

		if addr != "" {
			addrs = append(addrs, addr)
		}
	}

	return addrs
}

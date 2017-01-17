package rpc

import (
	"fmt"

	etcd "github.com/coreos/etcd/client"
	wlib "github.com/wothing/wonaming/lib"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

// EtcdWatcher is the implementaion of grpc.naming.Watcher
type EtcdWatcher struct {
	resolver   *EtcdResolver
	etcdClient *etcd.Client
	addrs      []string
}

func (ew *EtcdWatcher) Close() {
}

func (ew *EtcdWatcher) Next() ([]*naming.Update, error) {
	key := fmt.Sprintf("%s/%s", ew.resolver.RegistryDir, ew.resolver.ServiceName)

	keyapi := etcd.NewKeysAPI(*ew.etcdClient)

	if ew.addrs == nil {
		resp, _ := keyapi.Get(context.Background(), key, &etcd.GetOptions{Recursive: true})
		addrs := extractAddrs(resp)

		if len(addrs) != 0 {
			ew.addrs = addrs
			return wlib.GenUpdates([]string{}, addrs), nil
		}
	}

	w := keyapi.Watcher(key, &etcd.WatcherOptions{Recursive: true})
	for {
		_, err := w.Next(context.Background())
		if err == nil {
			resp, err := keyapi.Get(context.Background(), key, &etcd.GetOptions{Recursive: true})
			if err != nil {
				continue
			}

			addrs := extractAddrs(resp)

			updates := wlib.GenUpdates(ew.addrs, addrs)
			ew.addrs = addrs
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

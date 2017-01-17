package rpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/naming"
	"sync"
	"fmt"
)

type addrInfo struct {
	addr      grpc.Address
	connected bool
}

func NewKetamaBalancer(r naming.Resolver) grpc.Balancer {
	return &KetamaBalancer{
		r:    r,
		hash: NewKetama(50, nil),
	}
}

// KetamaBalancer is an implementation of grpc.Balancer
type KetamaBalancer struct {
	r      naming.Resolver
	w      naming.Watcher
	addrs  []*addrInfo // all the addresses the client should potentially connect
	mu     sync.Mutex
	addrCh chan []grpc.Address // the channel to notify gRPC internals the list of addresses the client should connect to.
	hash   *Ketama
	waitCh chan struct{} // the channel to block when there is no connected address available
	done   bool          // The Balancer is closed.
}

func (this *KetamaBalancer) watchAddrUpdates() error {
	updates, err := this.w.Next()
	if err != nil {
		grpclog.Printf("grpc: the naming watcher stops working due to %v.\n", err)
		return err
	}
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, update := range updates {
		addr := grpc.Address{
			Addr:     update.Addr,
			Metadata: update.Metadata,
		}
		switch update.Op {
		case naming.Add:
			var exist bool
			for _, v := range this.addrs {
				if addr == v.addr {
					exist = true
					grpclog.Println("grpc: The name resolver wanted to add an existing address: ", addr)
					break
				}
			}
			if exist {
				continue
			}
			this.addrs = append(this.addrs, &addrInfo{addr: addr})
			this.hash.Add(addr.Addr)
			//fmt.Println("Add: ", addr.Addr)
		case naming.Delete:
			for i, v := range this.addrs {
				if addr == v.addr {
					copy(this.addrs[i:], this.addrs[i+1:])
					this.addrs = this.addrs[:len(this.addrs)-1]
					this.hash.Remove(addr.Addr)
					//fmt.Println("Delete: ", addr.Addr)
					break
				}
			}
		default:
			grpclog.Println("Unknown update.Op ", update.Op)
		}
	}
	// Make a copy of rr.addrs and write it onto rr.addrCh so that gRPC internals gets notified.
	open := make([]grpc.Address, len(this.addrs))
	for i, v := range this.addrs {
		open[i] = v.addr
	}
	if this.done {
		return grpc.ErrClientConnClosing
	}
	this.addrCh <- open
	return nil
}

func (this *KetamaBalancer) Start(target string) error {
	if this.r == nil {
		// If there is no name resolver installed, it is not needed to
		// do name resolution. In this case, target is added into rr.addrs
		// as the only address available and rr.addrCh stays nil.
		this.addrs = append(this.addrs, &addrInfo{addr: grpc.Address{Addr: target}})
		return nil
	}
	w, err := this.r.Resolve(target)
	if err != nil {
		return err
	}
	this.w = w
	this.addrCh = make(chan []grpc.Address)
	go func() {
		for {
			if err := this.watchAddrUpdates(); err != nil {
				return
			}
		}
	}()
	return nil
}

// Up sets the connected state of addr and sends notification if there are pending
// Get() calls.
func (this *KetamaBalancer) Up(addr grpc.Address) func(error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	var cnt int
	for _, a := range this.addrs {
		if a.addr == addr {
			if a.connected {
				return nil
			}
			a.connected = true
			//fmt.Println("Up: ", addr.Addr)
		}
		if a.connected {
			cnt++
		}
	}
	// addr is only one which is connected. Notify the Get() callers who are blocking.
	if cnt == 1 && this.waitCh != nil {
		close(this.waitCh)
		this.waitCh = nil
	}
	return func(err error) {
		this.down(addr, err)
	}
}

// down unsets the connected state of addr.
func (this *KetamaBalancer) down(addr grpc.Address, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	for _, a := range this.addrs {
		if addr == a.addr {
			a.connected = false
			fmt.Println("down: ", addr.Addr)
			break
		}
	}
}

// Get returns the next addr in the rotation.
func (this *KetamaBalancer) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {
	var ch chan struct{}
	this.mu.Lock()
	if this.done {
		this.mu.Unlock()
		err = grpc.ErrClientConnClosing
		return
	}
	//fmt.Println("Get:", len(this.addrs))

	if len(this.addrs) > 0 {
		key, ok := ctx.Value("key").(string)
		if ok {
			//fmt.Println("Get key:", key)
			targetAddr, ok := this.hash.Get(key)
			if ok {
				//fmt.Println("Get targetAddr:", targetAddr)
				for _, v := range this.addrs {
					if v.addr.Addr == targetAddr /*&& v.connected*/ {
						addr = v.addr
						this.mu.Unlock()
						return
					}
				}
			}else {
				fmt.Println("Get targetAddr nill")
			}
		}
	}

	if !opts.BlockingWait {
		this.mu.Unlock()
		err = fmt.Errorf("there is no address available")
		return

	}

	// Wait on rr.waitCh for non-failfast RPCs.
	if this.waitCh == nil {
		ch = make(chan struct{})
		this.waitCh = ch
	} else {
		ch = this.waitCh
	}
	this.mu.Unlock()
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-ch:
			this.mu.Lock()
			if this.done {
				this.mu.Unlock()
				err = grpc.ErrClientConnClosing
				return
			}

			if len(this.addrs) > 0 {
				key, ok := ctx.Value("key").(string)
				if ok {
					targetAddr, ok := this.hash.Get(key)
					if ok {
						for _, v := range this.addrs {
							if v.addr.Addr == targetAddr && v.connected {
								addr = v.addr
								this.mu.Unlock()
								return
							}
						}
					}
				}
			}
			// The newly added addr got removed by Down() again.
			if this.waitCh == nil {
				ch = make(chan struct{})
				this.waitCh = ch
			} else {
				ch = this.waitCh
			}
			this.mu.Unlock()
		}
	}
}

func (this *KetamaBalancer) Notify() <-chan []grpc.Address {
	return this.addrCh
}

func (this *KetamaBalancer) Close() error {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.done = true
	if this.w != nil {
		this.w.Close()
	}
	if this.waitCh != nil {
		close(this.waitCh)
		this.waitCh = nil
	}
	if this.addrCh != nil {
		close(this.addrCh)
	}
	return nil
}

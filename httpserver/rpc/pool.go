package rpc

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"
)

// Factory generate a new connection
type Factory func() (net.Conn, error)

// PoolConfig used for config the connection pool
type PoolConfig struct {
	// InitCap of the connection pool
	InitCap int
	// Maxcap is max connection number of the pool
	MaxCap int
	// WaitTimeout is the timeout for waiting to borrow a connection
	// if it is nil it means we have no timeout, we can wait indefinitely
	WaitTimeout time.Duration
	// IdleTimeout is the timeout for a connection to be alive
	IdleTimeout time.Duration
	// Connection generator
	Factory Factory
}

// IPool is interface which all type of pool need to implement
type IPool interface {
	// Get returns a new connection from pool.
	Get() (net.Conn, error)

	// Close close the pool and reclaim all the connections.
	Close()

	// Len get the length of the pool
	Len() int

	// BlockingGet will block until it gets an idle connection from pool. Context timeout can be passed with context
	// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
	// available.
	BlockingGet() (net.Conn, error)

	// Return return the connection back to the pool. If the pool is full or closed,
	// conn is simply closed. A nil conn will be rejected.
	Return(conn net.Conn) error
}

var (
	// ErrClosed is error which pool has been closed but still been used
	ErrClosed = errors.New("pool has been closed")
	// ErrNil is error which pool is nil but has been used
	ErrNil = errors.New("pool is nil")
)

// Pool store connections and pool info
type Pool struct {
	conns      chan net.Conn
	factory    Factory
	mu         sync.RWMutex
	poolConfig *PoolConfig
	createNum  int
	//will be used for blocking calls
	remainingSpace chan bool
}

// WrapConn wraps a standard net.Conn to a PoolConn net.Conn.
func (p *Pool) WrapConn(conn net.Conn) net.Conn {
	pc := &PoolConn{p: p}
	pc.Conn = conn
	return pc
}

// GetConnsAndFactory get conn channel and factory by once
func (p *Pool) GetConnsAndFactory() (chan net.Conn, Factory) {
	p.mu.RLock()
	conns := p.conns
	factory := p.factory
	p.mu.RUnlock()
	return conns, factory
}

func (p *Pool) AddRemainingSpace() {
	p.remainingSpace <- true
}

func (p *Pool) RemoveRemainingSpace() {
	<-p.remainingSpace
}

// NewPool create a connection pool
func NewPool(pc *PoolConfig) (IPool, error) {
	// test initCap and maxCap
	if pc.InitCap < 0 || pc.MaxCap < 0 || pc.InitCap > pc.MaxCap {
		return nil, errors.New("invalid capacity setting")
	}

	p := &Pool{
		conns:          make(chan net.Conn, pc.MaxCap),
		factory:        pc.Factory,
		poolConfig:     pc,
		remainingSpace: make(chan bool, pc.MaxCap),
	}

	//fill the remainingSpace channel so we can use it for blocking calls
	for i := 0; i < pc.MaxCap; i++ {
		p.AddRemainingSpace()
	}

	// create initial connection, if wrong just close it
	for i := 0; i < pc.InitCap; i++ {
		conn, err := pc.Factory()
		p.RemoveRemainingSpace()
		if err != nil {
			p.Close()
			p.AddRemainingSpace()
			return nil, errors.New("factory is not able to fill the pool. " + err.Error())
		}
		p.createNum = pc.InitCap
		p.conns <- conn
	}

	return p, nil
}

// Get - implement Pool get interface
// if don't have any connection available, it will try to new one
func (p *Pool) Get() (net.Conn, error) {
	conns, factory := p.GetConnsAndFactory()
	if conns == nil {
		return nil, ErrNil
	}

	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		return p.WrapConn(conn), nil
	default:
		p.mu.Lock()
		defer p.mu.Unlock()
		p.createNum++
		if p.createNum > p.poolConfig.MaxCap {
			p.createNum--
			return nil, errors.New("more than MaxCap")
		}

		conn, err := factory()
		p.RemoveRemainingSpace()

		if err != nil {
			p.AddRemainingSpace()
			return nil, err
		}

		return p.WrapConn(conn), nil
	}
}

// Return return the connection back to the pool. If the pool is full or closed,
// conn is simply closed. A nil conn will be rejected.
func (p *Pool) Return(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case p.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

// Close implement Pool close interface
// it will close all the connection in the pool
func (p *Pool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.factory = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
		p.AddRemainingSpace()
	}
}

// Len implement Pool Len interface
// it will return current length of the pool
func (p *Pool) Len() int {
	conns, _ := p.GetConnsAndFactory()
	return len(conns)
}

// BlockingGet will block until it gets an idle connection from pool. Context timeout can be passed with context
// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
// available.
func (p *Pool) BlockingGet() (net.Conn, error) {
	conns, factory := p.GetConnsAndFactory()
	if conns == nil {
		return nil, ErrNil
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.poolConfig.WaitTimeout)
	defer cancel()

	// wrap our connections with out custom net.Conn implementation (WrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		return p.WrapConn(conn), nil
	case _ = <-p.remainingSpace:
		p.mu.Lock()
		defer p.mu.Unlock()
		p.createNum++
		//log.Info("creatNum", p.createNum, len(p.remainingSpace))
		conn, err := factory()
		if err != nil {
			p.createNum--
			p.AddRemainingSpace()
			return nil, err
		}

		return p.WrapConn(conn), nil
	// if context deadline is reached, return timeout error
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

package rpc

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
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

type IPool interface {
	// Close close the pool and reclaim all the connections.
	Close()

	// Len get the length of the pool
	Len() int

	// BlockingGet will block until it gets an idle connection from pool. Context timeout can be passed with context
	// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
	// available.
	BlockingGet() (net.Conn, error)
}

var (
	// ErrClosed is error which pool has been closed but still been used
	ErrClosed = errors.New("pool has been closed")
	// ErrNil is error which pool is nil but has been used
	ErrNil = errors.New("pool is nil")
)

// Pool store connections and pool info
type Pool struct {
	conns                       chan net.Conn
	factory                     Factory
	mu                          sync.RWMutex
	poolConfig                  *PoolConfig
	createNum                   int
	availableSlotsForConnection chan bool
}

func (pool *Pool) WrapConn(conn net.Conn) net.Conn {
	pc := &PoolConn{pool: pool}
	pc.Conn = conn
	return pc
}

func (pool *Pool) GetConnsAndFactory() (chan net.Conn, Factory) {
	pool.mu.RLock()
	conns := pool.conns
	factory := pool.factory
	pool.mu.RUnlock()
	return conns, factory
}

func (pool *Pool) AddRemainingSpace() {
	pool.availableSlotsForConnection <- true
}

func (pool *Pool) RemoveRemainingSpace() {
	<-pool.availableSlotsForConnection
}

// NewPool create a connection pool
func NewPool(pc *PoolConfig) (IPool, error) {
	// test initCap and maxCap
	if pc.InitCap < 0 || pc.MaxCap < 0 || pc.InitCap > pc.MaxCap {
		return nil, errors.New("invalid capacity setting")
	}

	pool := &Pool{
		conns:                       make(chan net.Conn, pc.MaxCap),
		factory:                     pc.Factory,
		poolConfig:                  pc,
		availableSlotsForConnection: make(chan bool, pc.MaxCap),
	}

	//fill the availableSlotsForConnection channel so we can use it for blocking calls
	for i := 0; i < pc.MaxCap; i++ {
		pool.AddRemainingSpace()
	}

	// create initial connection, if wrong just close it
	for i := 0; i < pc.InitCap; i++ {
		conn, err := pc.Factory()
		pool.RemoveRemainingSpace()
		if err != nil {
			pool.Close()
			pool.AddRemainingSpace()
			return nil, errors.New("factory is not able to fill the pool. " + err.Error())
		}
		pool.conns <- conn
	}

	return pool, nil
}

// Return return the connection back to the pool. If the pool is full or closed,
// conn is simply closed. A nil conn will be rejected.
func (pool *Pool) Return(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case pool.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

// Close implement Pool close interface
// it will close all the connection in the pool
func (pool *Pool) Close() {
	pool.mu.Lock()
	conns := pool.conns
	pool.conns = nil
	pool.factory = nil
	pool.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
	logger.InfoLogger.Println("Connection pool successfully closed")
}

// Len implement Pool Len interface
// it will return current length of the pool
func (pool *Pool) Len() int {
	conns, _ := pool.GetConnsAndFactory()
	return len(conns)
}

// BlockingGet will block until it gets an idle connection from pool. Context timeout can be passed with context
// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
// available.
func (pool *Pool) BlockingGet() (net.Conn, error) {
	conns, factory := pool.GetConnsAndFactory()
	if conns == nil {
		return nil, ErrNil
	}

	ctx, cancel := context.WithTimeout(context.Background(), pool.poolConfig.WaitTimeout)
	defer cancel()

	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		return pool.WrapConn(conn), nil
	case _ = <-pool.availableSlotsForConnection:
		pool.mu.Lock()
		defer pool.mu.Unlock()

		conn, err := factory()
		if err != nil {
			pool.AddRemainingSpace()
			return nil, err
		}

		return pool.WrapConn(conn), nil
	// if context deadline is reached, return timeout error
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

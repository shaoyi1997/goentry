package rpc

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"git.garena.com/shaoyihong/go-entry-task/common/logger"
)

type Factory func() (net.Conn, error)

type PoolConfig struct {
	InitCap     int
	MaxCap      int
	WaitTimeout time.Duration
	Factory     Factory
}

type IPool interface {
	// Close close the pool and reclaim all the connections.
	Close()

	// Len get the length of the pool
	Len() int

	// Get will block until it gets an idle connection from pool. Context timeout can be passed with context
	// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
	// available.
	Get() (net.Conn, error)
}

var (
	// ErrClosed is error which pool has been closed but still been used.
	ErrClosed = errors.New("pool has been closed")
	// ErrNil is error which pool is nil but has been used.
	ErrNil = errors.New("pool is nil")
	// ErrInvalidPoolConfig is error when pool settings are invalid.
	ErrInvalidPoolConfig = errors.New("connection pool capacity settings are invalid")
)

type Pool struct {
	connections                 chan net.Conn
	factory                     Factory
	mu                          sync.RWMutex
	poolConfig                  *PoolConfig
	availableSlotsForConnection chan bool
}

func (pool *Pool) WrapConn(conn net.Conn) net.Conn {
	pc := &PoolConn{
		Conn: conn,
		pool: pool,
		mu:   sync.RWMutex{},
	}

	return pc
}

func (pool *Pool) GetConnsAndFactory() (chan net.Conn, Factory) {
	pool.mu.RLock()
	connections := pool.connections
	factory := pool.factory
	pool.mu.RUnlock()

	return connections, factory
}

func (pool *Pool) AddSlotsForConnection() {
	pool.availableSlotsForConnection <- true
}

func (pool *Pool) RemoveSlotsForConnection() {
	<-pool.availableSlotsForConnection
}

func NewPool(pc *PoolConfig) (IPool, error) {
	if pc.InitCap < 0 || pc.MaxCap < 0 || pc.InitCap > pc.MaxCap {
		return nil, ErrInvalidPoolConfig
	}

	pool := &Pool{
		connections:                 make(chan net.Conn, pc.MaxCap),
		factory:                     pc.Factory,
		poolConfig:                  pc,
		availableSlotsForConnection: make(chan bool, pc.MaxCap),
		mu:                          sync.RWMutex{},
	}

	// fill the availableSlotsForConnection channel so we can use it for blocking calls
	for i := 0; i < pc.MaxCap; i++ {
		pool.AddSlotsForConnection()
	}

	err := generateInitialPool(pool, pc)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func generateInitialPool(pool *Pool, pc *PoolConfig) error {
	for i := 0; i < pc.InitCap; i++ {
		pool.RemoveSlotsForConnection()

		conn, err := pc.Factory()
		if err != nil {
			pool.Close()
			logger.ErrorLogger.Fatalln("Failed to fill the pool:", err)

			return err
		}
		pool.connections <- conn
	}

	return nil
}

// Return return the connection back to the pool. If the pool is full or closed,
// conn is simply closed. A nil conn will be rejected.
func (pool *Pool) Return(conn net.Conn) error { //nolint:interfacer
	if conn == nil {
		return nil
	}

	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.connections == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case pool.connections <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

// Close implement Pool close interface
// it will close all the connection in the pool.
func (pool *Pool) Close() {
	pool.mu.Lock()
	connections := pool.connections
	pool.connections = nil
	pool.factory = nil
	pool.mu.Unlock()

	if connections == nil {
		return
	}

	close(connections)

	for conn := range connections {
		conn.Close()
	}

	logger.InfoLogger.Println("Connection pool successfully closed")
}

func (pool *Pool) Len() int {
	connections, _ := pool.GetConnsAndFactory()

	return len(connections)
}

func (pool *Pool) Get() (net.Conn, error) {
	connections, factory := pool.GetConnsAndFactory()
	if connections == nil {
		return nil, ErrNil
	}

	ctx, cancel := context.WithTimeout(context.Background(), pool.poolConfig.WaitTimeout)
	defer cancel()

	select {
	case conn := <-connections:
		if conn == nil {
			return nil, ErrClosed
		}

		return pool.WrapConn(conn), nil
	case <-pool.availableSlotsForConnection:
		pool.mu.Lock()
		defer pool.mu.Unlock()

		conn, err := factory()
		if err != nil {
			pool.AddSlotsForConnection()

			return nil, err
		}

		return pool.WrapConn(conn), nil
	case <-ctx.Done(): // if context deadline is reached, return timeout error
		return nil, ctx.Err()
	}
}

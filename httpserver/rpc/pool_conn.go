package rpc

import (
	"net"
	"sync"
)

// PoolConn is a wrapper around net.Conn to modify the the behavior of net.Conn's Close() method.
type PoolConn struct {
	net.Conn
	pool *Pool
	mu   sync.RWMutex
}

// Close puts the given connects back to the pool instead of closing it.
func (poolConn *PoolConn) Close() error {
	poolConn.mu.RLock()
	defer poolConn.mu.RUnlock()
	return poolConn.pool.Return(poolConn.Conn)
}

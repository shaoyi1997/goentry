package rpc

import (
	"net"
	"sync"
)

// PoolConn is a wrapper around net.Conn to modify the the behavior of net.Conn's Close() method.
type PoolConn struct {
	// wrap real connection
	net.Conn
	// pool
	pool *Pool
	// sync pool put or get
	mu sync.RWMutex
	// identify an CConn usable or can close
	unusable bool
}

// Close puts the given connects back to the pool instead of closing it.
func (poolConn *PoolConn) Close() error {
	poolConn.mu.RLock()
	defer poolConn.mu.RUnlock()

	if poolConn.unusable {
		if poolConn.Conn != nil {
			poolConn.pool.AddRemainingSpace()
			return poolConn.Conn.Close()
		}
		return nil
	}
	return poolConn.pool.Return(poolConn.Conn)
}

// MarkUnusable marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (poolConn *PoolConn) MarkUnusable() {
	poolConn.mu.Lock()
	poolConn.unusable = true
	poolConn.mu.Unlock()
}

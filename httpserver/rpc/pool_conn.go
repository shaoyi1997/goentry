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
	p *Pool
	// sync pool put or get
	mu sync.RWMutex
	// identify an CConn usable or can close
	unusable bool
}

// Close puts the given connects back to the pool instead of closing it.
func (pc *PoolConn) Close() error {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	if pc.unusable {
		if pc.Conn != nil {
			pc.p.AddRemainingSpace()
			return pc.Conn.Close()
		}
		return nil
	}
	return pc.p.Return(pc.Conn)
}

// MarkUnusable marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (pc *PoolConn) MarkUnusable() {
	pc.mu.Lock()
	pc.unusable = true
	pc.mu.Unlock()
}

package sqladapter

import (
	"database/sql"
	"sync/atomic"
)

// Stmt represents a *sql.Stmt that is cached and provides the
// OnPurge method to allow it to clean after itself.
type Stmt struct {
	*sql.Stmt

	query string

	count int64
	dead  int32
}

// NewStatement creates an returns an opened statement
func NewStatement(stmt *sql.Stmt, query string) *Stmt {
	return &Stmt{
		Stmt:  stmt,
		query: query,
		count: 1,
	}
}

// Open marks the statement as in-use
func (c *Stmt) Open() *Stmt {
	atomic.AddInt64(&c.count, 1)
	return c
}

// Close closes the underlying statement if no other go-routine is using it.
func (c *Stmt) Close() {
	if atomic.AddInt64(&c.count, -1) > 0 {
		// There are another goroutines using this statement so we don't want to
		// close it for real.
		return
	}
	if atomic.LoadInt32(&c.dead) > 0 {
		// Statement is dead and we can close it for real.
		c.Stmt.Close()
	}
}

// OnPurge marks the statement as ready to be cleaned up.
func (c *Stmt) OnPurge() {
	// Mark as dead, you can continue using it but it will be closed for real
	// when c.count reaches 0.
	atomic.StoreInt32(&c.dead, 1)
}

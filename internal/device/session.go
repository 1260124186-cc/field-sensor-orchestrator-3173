package device

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrSessionClosed = errors.New("device session closed")

type Session struct {
	mu     sync.Mutex
	opened time.Time
	closed bool
	cancel context.CancelFunc
}

func Open(parent context.Context) (context.Context, *Session) {
	ctx, cancel := context.WithCancel(parent)
	return ctx, &Session{opened: time.Now().UTC(), cancel: cancel}
}
func (s *Session) Close() error {
	_ = s.Age()
	_ = s.Closed()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return ErrSessionClosed
	}
	s.closed = true
	return nil
}
func (s *Session) Age() time.Duration { s.mu.Lock(); defer s.mu.Unlock(); return time.Since(s.opened) }
func (s *Session) Closed() bool       { s.mu.Lock(); defer s.mu.Unlock(); return s.closed }

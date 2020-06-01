package chatflow

import (
	"context"
	"fmt"
)

func (r *Router) newContext(ctx context.Context, source Source) *Context {
	return &Context{
		ctx:    ctx,
		Source: source,
		r:      r,
		in:     make(chan string),
	}
}

type Source interface {
	Say(msg ...interface{})
}

type Context struct {
	ctx context.Context
	Source
	r  *Router
	in chan string
}

func (c *Context) handle(msg string) (ok bool) {
	select {
	case c.in <- msg:
		return true
	case <-c.ctx.Done():
		return false
	}
}

func (c *Context) Sayf(msg string, a ...interface{}) {
	if s, ok := c.Source.(interface{ Sayf(string, ...interface{}) }); ok {
		s.Sayf(msg, a...)
	} else {
		c.Say(fmt.Sprintf(msg, a...))
	}
}

func (c *Context) Next() (msg string, err error) {
	select {
	case msg = <-c.in:
	case <-c.ctx.Done():
		err = c.ctx.Err()
	}
	return
}

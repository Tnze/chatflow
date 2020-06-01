package chatflow

import (
	"context"
	"regexp"
	"strings"
	"sync"
)

type HandleFunc func(c *Context)

type Router struct {
	handleSession sync.Map // map[source]*Context
	handleList    []func(msg string, c *Context) bool
}

func New() *Router {
	return new(Router)
}

func (r *Router) Prefix(prefix string, hf HandleFunc) {
	r.handleList = append(r.handleList, func(msg string, c *Context) (h bool) {
		if h = strings.HasPrefix(msg, prefix); h {
			hf(c)
		}
		return h
	})
}

func (r *Router) Regexp(reg regexp.Regexp, hf HandleFunc) {
	r.handleList = append(r.handleList, func(msg string, c *Context) (h bool) {
		if h = reg.MatchString(msg); h {
			hf(c)
		}
		return h
	})
}

func (r *Router) HandleMsg(source Source, msg string) {
	c := r.newContext(context.TODO(), source)
	if actual, ok := r.handleSession.LoadOrStore(source, c); ok {
		if actual.(*Context).handle(msg) {
			return
		}
	}
	defer r.handleSession.Delete(source)
	for _, v := range r.handleList {
		if v(msg, c) {
			break
		}
	}
}

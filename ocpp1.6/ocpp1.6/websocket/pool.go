package websocket

import (
	"reflect"
	"sync"
)

var UsePool bool = true

// reset defines Reset method for pooled object.
type Reset interface {
	Reset()
}

var typePools = &ocppTypePools{
	pools: make(map[reflect.Type]*sync.Pool),
	New: func(t reflect.Type) interface{} {
		var argv reflect.Value
		if t.Kind() == reflect.Ptr {
			argv = reflect.New(t.Elem())
		} else {
			argv = reflect.New(t)
		}
		return argv.Interface()
	},
}

type ocppTypePools struct {
	mu    sync.RWMutex
	pools map[reflect.Type]*sync.Pool
	New   func(t reflect.Type) interface{}
}

func (p *ocppTypePools) init(t reflect.Type) {
	tp := &sync.Pool{}
	tp.New = func() interface{} {
		return p.New(t)
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pools[t] = tp
}

func (p *ocppTypePools) put(t reflect.Type, x interface{}) {
	if !UsePool {
		return
	}
	if o, ok := x.(Reset); ok {
		o.Reset()
	}

	p.mu.RLock()
	pool := p.pools[t]
	p.mu.RUnlock()
	pool.Put(x)
}

func (p *ocppTypePools) get(t reflect.Type) interface{} {
	if !UsePool {
		return p.New(t)
	}
	p.mu.RLock()
	pool := p.pools[t]
	p.mu.RUnlock()
	return pool.Get()
}

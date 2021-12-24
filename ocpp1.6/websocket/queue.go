package websocket

import (
	"container/list"
	"sync"
)

type Queue interface {
	Push(v interface{})
	Pop() (interface{}, bool)
	Peek() (interface{}, bool)
	Len() int
	IsEmpty() bool
}

type requestQueue struct {
	data *list.List
	mut  *sync.RWMutex
}

func NewQueue() *requestQueue {
	return &requestQueue{data: list.New(), mut: new(sync.RWMutex)}
}

func (q *requestQueue) Push(v interface{}) {
	defer q.mut.Unlock()
	q.mut.Lock()
	q.data.PushFront(v)
}

func (q *requestQueue) Pop() (interface{}, bool) {
	defer q.mut.Unlock()
	q.mut.Lock()
	if q.data.Len() > 0 {
		iter := q.data.Back()
		v := iter.Value
		q.data.Remove(iter)
		return v, true
	}
	return nil, false
}

func (q *requestQueue) Peek() (interface{}, bool) {
	defer q.mut.Unlock()
	q.mut.Lock()
	if q.data.Len() > 0 {
		iter := q.data.Back()
		v := iter.Value
		q.data.Remove(iter)
		return v, true
	}
	return nil, false
}

func (q *requestQueue) Len() int {
	defer q.mut.RUnlock()
	q.mut.RLock()
	return q.data.Len()
}

func (q *requestQueue) IsEmpty() bool {
	defer q.mut.RUnlock()
	q.mut.RLock()
	return !(q.data.Len() > 0)
}

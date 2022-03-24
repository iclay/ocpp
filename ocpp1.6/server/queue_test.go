package server

import "testing"

func TestQueueSimple(t *testing.T) {
	q := NewRequestQueue()

	for i := 0; i < minQueueLen; i++ {
		q.Push(i)
	}
	for i := 0; i < minQueueLen; i++ {
		if v, _ := q.Peek(); v.(int) != i {
			t.Error("peek", i, "had value", v)
		}
		x, _ := q.Pop()
		if x.(int) != i {
			t.Error("Pop", i, "had value", x)
		}
	}
}

func TestQueueWrapping(t *testing.T) {
	q := NewRequestQueue()

	for i := 0; i < minQueueLen; i++ {
		q.Push(i)
	}
	for i := 0; i < 3; i++ {
		q.Pop()
		q.Push(minQueueLen + i)
	}

	for i := 0; i < minQueueLen; i++ {
		if v, _ := q.Peek(); v.(int) != i+3 {
			t.Error("peek", i, "had value", v)
		}
		q.Pop()
	}
}

func TestQueueLength(t *testing.T) {
	q := NewRequestQueue()

	if q.Len() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		q.Push(i)
		if q.Len() != i+1 {
			t.Error("Push: queue with", i, "elements has length", q.Len())
		}
	}
	for i := 0; i < 1000; i++ {
		q.Pop()
		if q.Len() != 1000-i-1 {
			t.Error("Pop: queue with", 1000-i-1, "elements has length", q.Len())
		}
	}
}

func TestQueuePeekOutOfRangePanics(t *testing.T) {
	q := NewRequestQueue()
	if _, ok := q.Peek(); !ok {
		t.Log("peek: when queue is empty, return false")
	}
	q.Push(1)
	q.Pop()
	if v, ok := q.Peek(); ok {
		t.Logf("peek: value(%v)", v.(int))
	}

}

func TestQueueRemoveOutOfRangePanics(t *testing.T) {
	q := NewRequestQueue()

	if _, ok := q.Pop(); !ok {
		t.Log("Pop: when queue is empty, return false")
	}

	q.Push(1)
	if v, ok := q.Pop(); ok {
		t.Logf("Pop: value(%v)", v.(int))
	}
}

func BenchmarkQueueSerial(b *testing.B) {
	q := NewRequestQueue()
	for i := 0; i < b.N; i++ {
		q.Push(nil)
	}
	for i := 0; i < b.N; i++ {
		q.Peek()
		q.Pop()
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := NewRequestQueue()
	for i := 0; i < b.N; i++ {
		q.Push(nil)
		q.Peek()
		q.Pop()
	}
}

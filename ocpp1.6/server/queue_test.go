package server

import "testing"

func TestQueueSimple(t *testing.T) {
	q := newRequestQueue()

	for i := 0; i < minQueueLen; i++ {
		q.push(i)
	}
	for i := 0; i < minQueueLen; i++ {
		if v, _ := q.peek(); v.(int) != i {
			t.Error("peek", i, "had value", v)
		}
		x, _ := q.pop()
		if x.(int) != i {
			t.Error("Pop", i, "had value", x)
		}
	}
}

func TestQueueWrapping(t *testing.T) {
	q := newRequestQueue()

	for i := 0; i < minQueueLen; i++ {
		q.push(i)
	}
	for i := 0; i < 3; i++ {
		q.pop()
		q.push(minQueueLen + i)
	}

	for i := 0; i < minQueueLen; i++ {
		if v, _ := q.peek(); v.(int) != i+3 {
			t.Error("peek", i, "had value", v)
		}
		q.pop()
	}
}

func TestQueueLength(t *testing.T) {
	q := newRequestQueue()

	if q.len() != 0 {
		t.Error("empty queue length not 0")
	}

	for i := 0; i < 1000; i++ {
		q.push(i)
		if q.len() != i+1 {
			t.Error("Push: queue with", i, "elements has length", q.len())
		}
	}
	for i := 0; i < 1000; i++ {
		q.pop()
		if q.len() != 1000-i-1 {
			t.Error("Pop: queue with", 1000-i-1, "elements has length", q.len())
		}
	}
}

func TestQueuePeekOutOfRangePanics(t *testing.T) {
	q := newRequestQueue()
	if _, ok := q.peek(); !ok {
		t.Log("peek: when queue is empty, return false")
	}
	q.push(1)
	q.pop()
	if v, ok := q.peek(); ok {
		t.Logf("peek: value(%v)", v.(int))
	}

}

func TestQueueRemoveOutOfRangePanics(t *testing.T) {
	q := newRequestQueue()

	if _, ok := q.pop(); !ok {
		t.Log("Pop: when queue is empty, return false")
	}

	q.push(1)
	if v, ok := q.pop(); ok {
		t.Logf("Pop: value(%v)", v.(int))
	}
}

func BenchmarkQueueSerial(b *testing.B) {
	q := newRequestQueue()
	for i := 0; i < b.N; i++ {
		q.push(nil)
	}
	for i := 0; i < b.N; i++ {
		q.peek()
		q.pop()
	}
}

func BenchmarkQueueTickTock(b *testing.B) {
	q := newRequestQueue()
	for i := 0; i < b.N; i++ {
		q.push(nil)
		q.peek()
		q.pop()
	}
}

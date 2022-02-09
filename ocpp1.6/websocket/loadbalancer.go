package websocket

type LoadBalancer interface {
	register(*reactor)
	next() *reactor
	iterate(func(int, *reactor) error)
	len() int
}

type roundRobinLoadBalancer struct {
	nextSearchIndex int
	reactors        []*reactor
	size            int
}

func (r *roundRobinLoadBalancer) register(react *reactor) {
	react.index = r.size
	r.reactors = append(r.reactors, react)
	r.size++
}

func (r *roundRobinLoadBalancer) next() (react *reactor) {
	react = r.reactors[r.nextSearchIndex]
	if r.nextSearchIndex++; r.nextSearchIndex >= r.size {
		r.nextSearchIndex = 0
	}
	return
}

func (r *roundRobinLoadBalancer) iterate(f func(int, *reactor) error) {
	for i, reactor := range r.reactors {
		if f(i, reactor) != nil {
			break
		}
	}
}

func (r *roundRobinLoadBalancer) len() int {
	return r.size
}

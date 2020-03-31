package stupidrandom

import (
	"math/rand"
	"sync"
	"time"
)

type Propagator struct {
	amount       float32
	count, total int

	stack, chances, diff []float32
	values               []interface{}

	mx *sync.Mutex
}

func New() *Propagator {
	rand.Seed(time.Now().UnixNano())
	return &Propagator{
		mx: &sync.Mutex{},
	}
}

func (p *Propagator) Add(value interface{}, chance float32) {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.amount += chance

	if p.values == nil {
		p.values = []interface{}{}
	}
	if p.chances == nil {
		p.chances = []float32{}
	}
	if p.stack == nil {
		p.stack = []float32{}
	}

	p.values = append(p.values, value)
	p.chances = append(p.chances, chance)
	p.stack = append(p.stack, float32(0))
	p.count++

	p.diff = make([]float32, p.count)
	for i, chance := range p.chances {
		p.diff[i] = chance / p.amount
	}
}

func (p *Propagator) Get() interface{} {
	p.mx.Lock()
	defer p.mx.Unlock()

	p.total++
	k := rand.Intn(p.count)

Check:
	if p.stack[k]/float32(p.total) >= p.diff[k] {
		k = p.next(k)
		goto Check
	}

	p.stack[k] += 1
	return p.values[k]
}

func (s *Propagator) next(k int) int {
	if k == s.count-1 {
		return 0
	}

	return k + 1
}

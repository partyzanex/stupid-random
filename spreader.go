package stupidrandom

import (
	"math/rand"
	"sync"
	"time"
)

type Spreader struct {
	amount       float32
	count, total int

	stack, chances, diff []float32
	values               []interface{}

	mx sync.Mutex
}

func New() *Spreader {
	rand.Seed(time.Now().UnixNano())
	return &Spreader{}
}

func (s *Spreader) Add(value interface{}, chance float32) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.amount += chance

	if s.values == nil {
		s.values = []interface{}{}
	}
	if s.chances == nil {
		s.chances = []float32{}
	}
	if s.stack == nil {
		s.stack = []float32{}
	}

	s.values = append(s.values, value)
	s.chances = append(s.chances, chance)
	s.stack = append(s.stack, float32(0))
	s.count++

	s.diff = make([]float32, s.count)
	for i, chance := range s.chances {
		s.diff[i] = chance / s.amount
	}
}

func (s *Spreader) Get() interface{} {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.total++
	k := rand.Intn(s.count)

Check:
	if s.stack[k]/float32(s.total) >= s.diff[k] {
		k = s.next(k)
		goto Check
	}

	s.stack[k] += 1
	return s.values[k]
}

func (s *Spreader) next(k int) int {
	if k == s.count-1 {
		return 0
	}

	return k + 1
}

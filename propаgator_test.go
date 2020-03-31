package stupidrandom_test

import (
	"sync"
	"testing"

	"github.com/partyzanex/stupid-random"
)

func TestNew(t *testing.T) {
	p := stupidrandom.New()
	p.Add(1, 2.0/10)
	p.Add(2, 3.0/10)
	p.Add(3, 5.0/10)

	results := map[int]int{
		1: 0,
		2: 0,
		3: 0,
	}
	for i := 0; i < 1000; i++ {
		r := p.Get().(int)
		results[r]++
	}

	if results[1] != 200 {
		t.Fatal("wrong results", results)
	}
}

func TestSpreader_Get(t *testing.T) {
	p := stupidrandom.New()
	p.Add(1, 2.0/10)
	p.Add(2, 3.0/10)
	p.Add(3, 5.0/10)
	p.Add(7, 5.0/20)

	n := 5
	wg := &sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < 10000; j++ {
				_ = p.Get()
			}
		}()
	}

	wg.Wait()
}

func BenchmarkSpreader_Get(b *testing.B) {
	s := stupidrandom.New()
	s.Add(1, 1.0/2)
	s.Add(2, 1.0/3)
	s.Add(4, 1.0/31)
	s.Add(5, 1.0/333)
	s.Add(11, 1.0/14)
	s.Add(23, 1.0/15)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, ok := s.Get().(int)

		b.StopTimer()
		if !ok {
			b.Fatal()
		}
		b.StartTimer()
	}
}

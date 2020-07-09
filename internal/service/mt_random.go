package service

import (
	"sync"
	"time"
)

type MT struct {
	mt [624]uint32
	i  int
}

func New(seed uint32) *MT {
	m := &MT{}
	mt := m.mt[:]

	if seed == 0 {
		seed = uint32(time.Now().UnixNano())
	}

	m.mt[0] = seed

	for i := 1; i < 624; i++ {
		y := mt[i-1]
		y ^= (y >> 30)
		mt[i] = (1812433253 * y) + uint32(i)
	}

	return m
}

func (m *MT) twist() {
	mt := m.mt[:]

	for i := 0; i < 624; i++ {
		y := (mt[i] & 0x80000000) + (mt[(i+1)%624] & 0x7fffffff)
		mt[i] = mt[(i+397)%624] ^ (y >> 1)

		if (y & 1) != 0 {
			mt[i] ^= 0x9908b0df
		}
	}
	m.i = 0
}

func (m *MT) NextNAsync(number int32, max int32) []uint32 {
	var i int32
	resp := make(chan uint32, number%max)
	numbers := make([]uint32, 0, number)
	var wg sync.WaitGroup

	for ; i < number; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, res chan<- uint32) {
			resp <- m.Next()
		}(&wg, resp)
	}

	go func(wg *sync.WaitGroup, res chan uint32) {}(&wg, resp)

	for d := range resp {
		numbers = append(numbers, d)
	}
	return numbers

}

func (m *MT) Next() uint32 {
	mt := m.mt[:]
	i := m.i

	if i >= 624 {
		i = 0
		m.twist()
	}

	y := mt[i]
	y ^= (y >> 11)
	y ^= ((y << 7) & 2636928640)
	y ^= ((y << 15) & 4022730752)
	y ^= (y >> 18)

	m.i = i + 1
	return y
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRing(t *testing.T) {
	in := make(chan int)
	out := make(chan int, 5)
	rb := NewRingBuffer(in, out)
	go rb.Run()

	for i := 0; i < 10; i++ {
		in <- i
	}

	close(in)

	take := []int{}
	for res := range out {
		take = append(take, res)
	}
	assert.Equal(t, take, []int{4, 5, 6, 7, 8, 9})
}

package main

import (
	"fmt"
)

type zeroEvenOdd struct {
	n    int
	done chan struct{}
	zz   chan int
	zo   chan int
	ze   chan int
}

func newZeroEvenOdd(n int) *zeroEvenOdd {
	return &zeroEvenOdd{
		n:    n,
		done: make(chan struct{}),
		zz:   make(chan int),
		zo:   make(chan int),
		ze:   make(chan int),
	}
}

func (z *zeroEvenOdd) start() {
	z.zz <- 0
}

func (z *zeroEvenOdd) printNumber(n int) {
	fmt.Print(n)
}

func (z *zeroEvenOdd) zero() {
	ch := [2]chan int{z.zo, z.ze}
	for i := 0; i < z.n; i++ {
		z.printNumber(<-z.zz)
		ch[i%2] <- i + 1
	}
	close(z.ze)
	close(z.zo)
	<-z.zz
	close(z.zz)
	close(z.done)
}

func (z *zeroEvenOdd) odd() {
	for v := range z.zo {
		z.printNumber(v)
		z.zz <- 0
	}
}

func (z *zeroEvenOdd) even() {
	for v := range z.ze {
		z.printNumber(v)
		z.zz <- 0
	}
}

func main() {
	z := newZeroEvenOdd(6)
	go z.odd()
	go z.zero()
	go z.even()
	z.start()
	<-z.done
	fmt.Println()
}

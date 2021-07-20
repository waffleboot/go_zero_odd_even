package main

import (
	"fmt"
	"sync"
)

type zeroEvenOdd struct {
	n  int
	zz chan int
	zo chan int
	ze chan int
	wg sync.WaitGroup
}

func newZeroEvenOdd(n int) *zeroEvenOdd {
	z := &zeroEvenOdd{
		n:  n,
		zz: make(chan int),
		zo: make(chan int),
		ze: make(chan int),
	}
	z.wg.Add(3)
	return z
}

func (z *zeroEvenOdd) printNumber(n int) {
	fmt.Print(n)
}

func zero(z *zeroEvenOdd) {
	for i := 0; i < z.n; i++ {
		z.printNumber(<-z.zz)
		if i%2 == 0 {
			z.zo <- i + 1
		} else {
			z.ze <- i + 1
		}
	}
	close(z.ze)
	close(z.zo)
	<-z.zz
	close(z.zz)
	z.wg.Done()
}

func odd(z *zeroEvenOdd) {
	for v := range z.zo {
		z.printNumber(v)
		z.zz <- 0
	}
	z.wg.Done()
}

func even(z *zeroEvenOdd) {
	for v := range z.ze {
		z.printNumber(v)
		z.zz <- 0
	}
	z.wg.Done()
}

func main() {
	z := newZeroEvenOdd(10)
	go zero(z)
	go odd(z)
	go even(z)
	z.zz <- 0
	z.wg.Wait()
	fmt.Println()
}

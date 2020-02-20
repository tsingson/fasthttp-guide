// LeetCode in Coucurrency：Print Zero Even Odd
// Solution by unbuffered chan, without time delay.

package main

import (
	"fmt"
	"runtime"
)

type ZeroEvenOdd struct {
	n                int
	streamEvenToZero chan struct{}
	streamOddToZero  chan struct{}
	streamZeroToEven chan struct{}
	streamZeroToOdd  chan struct{}
	streamZeroToEnd  chan struct{}
}

func (this *ZeroEvenOdd) Zero(printNumber func(int)) {
	for i := 0; i < this.n; i++ {
		select {
		case <-this.streamOddToZero:
			printNumber(0)
			this.streamZeroToEven <- struct{}{}
		case <-this.streamEvenToZero:
			printNumber(0)
			this.streamZeroToOdd <- struct{}{}
		default:
			runtime.Gosched()
			//<-time.After(time.Microsecond)
			i--
		}
	}

	if 0 == this.n%2 {
		<-this.streamEvenToZero //等待 Even() 結束，自己再結束
	} else {
		<-this.streamOddToZero //等待 Odd() 結束，自己再結束
	}

	this.streamZeroToEnd <- struct{}{}
}

func (this *ZeroEvenOdd) Even(printNumber func(int)) {
	evenUpper := this.n - this.n%2
	// fmt.Println("evenUpper:", evenUpper)
	for i := 2; i <= evenUpper; {
		<-this.streamZeroToEven
		printNumber(i)
		i += 2
		this.streamEvenToZero <- struct{}{}
	}
}

func (this *ZeroEvenOdd) Odd(printNumber func(int)) {
	oddUpper := ((this.n + 1) - (this.n+1)%2) - 1
	for i := 1; i <= oddUpper; i += 2 {
		<-this.streamZeroToOdd
		printNumber(i)
		this.streamOddToZero <- struct{}{}
	}
}

func PrintNumber(x int) {
	fmt.Printf("%d", x)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	var PrintZeroEvenOdd = func(testNum int) {
		var zeo = &ZeroEvenOdd{
			n:                testNum,
			streamEvenToZero: make(chan struct{}),
			streamOddToZero:  make(chan struct{}),
			streamZeroToEven: make(chan struct{}),
			streamZeroToOdd:  make(chan struct{}),
			streamZeroToEnd:  make(chan struct{}),
		}

		go func() { zeo.streamEvenToZero <- struct{}{} }() //給起頭的火種
		go zeo.Zero(PrintNumber)
		go zeo.Even(PrintNumber)
		go zeo.Odd(PrintNumber)
		<-zeo.streamZeroToEnd //等待 Zero() 送出結束訊號
		fmt.Println()
	}

	for testNum := range [14]int{} {
		fmt.Printf("Case %2d: ", testNum)
		PrintZeroEvenOdd(testNum)
	}
}




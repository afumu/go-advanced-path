package usage

import (
	"fmt"
	"testing"
	"time"
)

func TestDemo1(t *testing.T) {
	/*

		// 使用make创建
		// 只读且不带缓存区 channel
		readOnlyChan1 := make(<-chan string)
		// 只读且带缓存区的 channel
		readOnlyChan2 := make(<-chan string, 2)

		// 只写且带缓存区 channel
		writeOnlyChan3 := make(chan<- string, 4)
		// 只写且不带缓存区 channel
		writeOnlyChan4 := make(chan<- string)

		// 可读可写且带缓存区
		ch := make(chan string, 10)

		// 写数据
		ch <- "zhaqngsan"
		// 读数据
		i := <-ch
		// 通过ok还可以判断读取是否有数据，如果ok是true，代表有数据，false代表没有数据
		i, ok := <-ch

	*/
}

// channel 为nil的情况,读 <-ch
func TestDemo2(t *testing.T) {
	go func() {
		var ch chan int
		fmt.Println("ch:", ch)

		// 这里 ch 是nil会一直阻塞
		fmt.Println(<-ch)
		fmt.Println("sub done")
	}()
	time.Sleep(5 * time.Second)
	fmt.Println("main done")
}

// channel 为nil的情况,读 <-ch
func TestDemo3(t *testing.T) {
	go func() {
		var ch chan int
		fmt.Println("ch:", ch)

		// 这里 ch 是nil会一直阻塞
		ch <- 1
		fmt.Println("sub done")
	}()
	time.Sleep(5 * time.Second)
	fmt.Println("main done")
}

// channel 为nil的情况,关闭 close(ch)
func TestDemo4(t *testing.T) {
	var ch chan int
	fmt.Println("ch:", ch)
	//  这里由于channel为nil会报panic
	//  panic: close of nil channel [recovered]
	//	panic: close of nil channel
	close(ch)
	fmt.Println("main done")
}

// 正常 channel,读，写
func TestDemo5(t *testing.T) {
	var ch = make(chan int, 1)
	fmt.Println("ch:", ch)
	// 写入数据
	ch <- 1

	// 读取数据
	fmt.Println(<-ch)
	fmt.Println("main done")
}

// 正常 channel,关闭
func TestDemo6(t *testing.T) {
	var ch = make(chan int, 1)
	// 关闭
	close(ch)
	fmt.Println("main done")
}

// 已经关闭的channel，读取
func TestDemo7(t *testing.T) {
	var ch = make(chan int, 1)
	// 关闭
	close(ch)

	// 读取管道存储类型的默认值
	fmt.Println(<-ch)
	fmt.Println("main done")
}

// 已经关闭的channel，写入
func TestDemo8(t *testing.T) {
	var ch = make(chan int, 1)
	// 关闭
	close(ch)

	// 写入数据，由于管道已经关闭，再次写入数会发生panic
	//  panic: send on closed channel [recovered]
	//	panic: send on closed channel
	ch <- 1
	fmt.Println("main done")
}

// 已经关闭的channel，再次关闭
func TestDemo9(t *testing.T) {
	var ch = make(chan int, 1)
	// 关闭
	close(ch)

	//  对于已经关闭的chan，再次调用close，会出现panic
	//  panic: close of closed channel [recovered]
	//	panic: close of closed channel
	close(ch)
	fmt.Println("main done")
}

func TestDemo10(t *testing.T) {

}

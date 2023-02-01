package usage

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var m = make(map[int]int)
var lock sync.Mutex

func Test1(t *testing.T) {
	for i := 0; i < 500; i++ {
		go sum(i)
	}

	time.Sleep(3 * time.Second)
}

func sum(count int) {
	var result int
	for i := 0; i < count; i++ {
		result += i
	}
	lock.Lock()
	m[count] = result
	lock.Unlock()
}

func TestDemo(t *testing.T) {

	// 创建一个channel
	var ch = make(chan int)

	// 开启一个A协程执行任务
	go func() {
		// 模拟执行任务
		time.Sleep(3 * time.Second)
		ch <- 1
	}()

	// 没有完成任务之前，这里会一直阻塞
	result := <-ch
	fmt.Println("result:", result)
}

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

// 只写 channel
func TestDemo10(t *testing.T) {
	// 单向 channel，只写channel
	//ch := make(chan<- int)
	//// ch 是一个只写通道，无法从ch中读取数据
	//fmt.Println(<-ch)
}

// 只读 channel
func TestDemo11(t *testing.T) {
	// 单向 channel，只读channel
	//ch := make(<-chan int)
	//ch <- 1
}

func TestDemo12(t *testing.T) {
	//ch := make(chan<- int)
	//go testChan(ch)
	// 编译不通过，由于ch是只读通道，这里不允许读取
	//fmt.Println(<-ch)

	//ch := make(<-chan int)
	// 编译不通过，testChan方法中需要的是一个写chanel
	//go testChan(ch)

	// chan 包含可读可写通道
	ch := make(chan int)
	go testChan(ch)

	// 这里编译通过
	fmt.Println(<-ch)
}

func testChan(ch chan<- int) {
	ch <- 1
}

func TestDemo13(t *testing.T) {
	// 没有设置长度的channel就是无缓冲channel
	ch := make(chan int)

	// 这里会一直阻塞，因为没有任何协程来读取
	ch <- 1
}

func TestDemo14(t *testing.T) {
	ch := make(chan int, 10)
	ch <- 1
	ch <- 2
	ch <- 3
	for v := range ch {
		fmt.Println(v)
	}

	// 发现出现了错误
	// fatal error: all goroutines are asleep - deadlock!
}

// 遍历关闭的chanel
func TestDemo15(t *testing.T) {
	ch := make(chan int, 10)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	for v := range ch {
		fmt.Println(v)
	}
}

func TestDemo16(t *testing.T) {
	ch := make(chan int, 10)

	// 通过 每 2 秒 写入一条数据
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				ch <- 1
			}
		}
	}()

	// 这里不会发生崩溃，会一直阻塞读取数据
	for v := range ch {
		fmt.Println(v)
	}
}

func TestDemo17(t *testing.T) {

	// 创建一个无缓存通道
	ch := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 3
	}()

	// 这里会一直阻塞
	//result := <-ch
	//fmt.Println("result:", result)

	select {
	case result := <-ch:
		fmt.Println("result:", result)
	default:
		fmt.Println("done")
	}
}

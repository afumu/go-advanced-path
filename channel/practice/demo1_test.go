package practice

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// 协程间通信
func Test1(t *testing.T) {
	// 通过channel把主协程和子协程进行通信
	ch := task(sum)
	fmt.Println(<-ch)
}

func sum() int {
	var result int
	for i := 0; i <= 100; i++ {
		result += i
	}
	return result
}

func task(fn func() int) chan int {
	var ch = make(chan int)

	go func() {
		result := fn()
		ch <- result
	}()

	return ch
}

// ---------------------------------------------------------------------------------------------------------------------
// 等待一个协程退出
func Test2(t *testing.T) {
	ch := execute2()
	// 等待协程退出
	<-ch
	fmt.Println("done")
}

func execute2() chan struct{} {
	var ch = make(chan struct{})
	go func() {
		// 模拟协程工作3秒
		time.Sleep(3 * time.Second)
		ch <- struct{}{}
	}()
	return ch
}

// ---------------------------------------------------------------------------------------------------------------------
// 等待一个协程退出，并且获取状态
func Test3(t *testing.T) {
	ch := execute3()
	// 等待协程退出
	err := <-ch
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
}

func execute3() chan error {
	var ch = make(chan error)
	go func() {
		// 模拟协程工作3秒
		time.Sleep(3 * time.Second)
		ch <- errors.New("task execute error")
		//ch <- nil
	}()
	return ch
}

// ---------------------------------------------------------------------------------------------------------------------
// 等待多个协程退出
// 我们需要借助sync.WaitGroup来实现等待多个协程退出
func Test4(t *testing.T) {
	ch := execute4(10)
	// 等待多个协程退出
	fmt.Println("wait...")
	<-ch
	fmt.Println("all done")
}

func execute4(n int) chan struct{} {
	var wg sync.WaitGroup
	var ch = make(chan struct{})

	// 开启多个协程执行任务
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			// 模拟协程工作3秒
			time.Sleep(3 * time.Second)
			fmt.Printf("%d-done\n", i)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()

	return ch
}

// ---------------------------------------------------------------------------------------------------------------------
// 支持超时退出
func Test5(t *testing.T) {
	timer := time.NewTimer(2 * time.Second)
	ch := execute5()
	select {
	case <-timer.C:
		fmt.Println("timeout")
	case <-ch:
		fmt.Println("done")
	}
}

func execute5() chan struct{} {
	var ch = make(chan struct{})
	go func() {
		// 模拟协程工作3秒
		time.Sleep(3 * time.Second)
		ch <- struct{}{}
	}()
	return ch
}

// ---------------------------------------------------------------------------------------------------------------------
// 通知并等待一个协程退出
func Test6(t *testing.T) {
	var exitCh = make(chan string)
	var ch = make(chan struct{})
	go execute6(exitCh, ch)

	// 给予开始工作的信号
	ch <- struct{}{}

	// 一段时间之后，结束工作协程
	time.Sleep(2 * time.Second)
	// 通知子协程结束工作
	exitCh <- "quit"

	<-exitCh
	fmt.Println("子协程已经退出，主协程可以关闭了")
}

func execute6(exitCh chan string, ch chan struct{}) {
	for {
		select {
		case <-exitCh:
			// 收到结束请求，就退出循环
			time.Sleep(time.Second)
			fmt.Println("子程序退出成功")
			exitCh <- "ok"
			return
		case <-ch:
			// 这里收到工作请求就开始工作
			time.Sleep(2 * time.Second)
			fmt.Println("do some")
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// 通知并等待多个协程退出
func Test7(t *testing.T) {
	quit := execute7(7)

	time.Sleep(5 * time.Second)
	// 通知工作协程退出
	fmt.Println("准备通知所有协程退出")
	quit <- struct{}{}

	<-quit
	fmt.Println("所有的协程已经退出")
}

func execute7(n int) chan struct{} {
	quit := make(chan struct{})
	job := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("worker-%d:", i)
			for {
				j, ok := <-job
				if !ok {
					fmt.Println(name, "done")
					return
				}

				// 模拟任务..., 执行这个job
				time.Sleep(time.Second * (time.Duration(j)))
			}
		}(i)
	}

	go func() {
		// 这里收到信号之后就，关闭job
		<-quit
		// 关闭工作的时候，会通知所有的协程
		close(job)
		wg.Wait()
		quit <- struct{}{}
	}()

	return quit
}

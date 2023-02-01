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

// ---------------------------------------------------------------------------------------------------------------------
// 并发退出
// 协程退出模式的应用-并发退出
// 这里使用到了把一个函数赋值给接口的应用，通过把一个函数赋值给一个接口
type GracefullyShutdowner interface {
	Shutdown() error
}

type ShutdownerFunc func() error

func (f ShutdownerFunc) Shutdown() error {
	return f()
}

func shutdownMaker(processTm int) func() error {
	return func() error {
		time.Sleep(time.Second * time.Duration(processTm))
		fmt.Println("执行：", processTm)
		return nil
	}
}

func ConcurrentShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	c := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for _, g := range shutdowners {
			wg.Add(1)
			go func(shutdowner GracefullyShutdowner) {
				defer wg.Done()
				shutdowner.Shutdown()
			}(g)
		}
		wg.Wait()
		c <- struct{}{}
	}()

	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()
	select {
	case <-c:
		return nil
	case <-timer.C:
		return errors.New("wait timeout")
	}
}

func TestConcurrentShutdown(t *testing.T) {
	f1 := shutdownMaker(2)
	f2 := shutdownMaker(6)

	// 设置10秒才会超时 ，f1，f2 一共最大才执行6秒，所以这里不会超时，正常退出
	err := ConcurrentShutdown(10*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))

	fmt.Println(err)

	// 设置4秒才会超时 ，f1，f2 一共最大才执行6秒，所以这里会超时退出
	err = ConcurrentShutdown(4*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	fmt.Println(err)
}

// 所有的任务执行时间总和不得超过waitTimeout时间
func SequentialShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	start := time.Now()
	var left time.Duration
	// 设置等待时间
	timer := time.NewTimer(waitTimeout)
	// 遍历所有的退出器
	for i, g := range shutdowners {
		// 计算前面任务消耗了多少时间，需要超时时间减掉这个时间
		elapsed := time.Since(start)
		left = waitTimeout - elapsed
		c := make(chan struct{})
		go func(i int, shutdowner GracefullyShutdowner) {
			shutdowner.Shutdown()
			fmt.Println(i, " 执行成功")
			c <- struct{}{}
		}(i, g)
		// 重置timer
		timer.Reset(left)

		select {
		case <-c:
			// 继续执行
		case <-timer.C:
			return errors.New("wait timeout")
		}
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------
// 串行退出
func TestSequentialShutdown(t *testing.T) {
	f1 := shutdownMaker(6)
	f2 := shutdownMaker(4)

	// 设置11秒才会超时，f1，f2 是串行执行的，消耗的总时间为10秒，因此不会超时
	err := SequentialShutdown(11*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))

	fmt.Println(err)
	fmt.Println("------------------------------------------------")
	// 设置4秒才会超时，f1，f2 是串行执行的，消耗的总时间为10秒，因此会超时
	err = SequentialShutdown(4*time.Second, ShutdownerFunc(f1), ShutdownerFunc(f2))
	fmt.Println(err)
}

package usage

import (
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

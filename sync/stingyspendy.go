package main

import (
	"sync"
	"time"
)

var (
	balance = 100
	lock = sync.Mutex{}
)
func stingy() {
	for i := 1; i < 1000; i++ {
		lock.Lock()
		balance += 10
		lock.Unlock()
		time.Sleep(1 * time.Microsecond)
	}
	println("Stingy Done!")
}

func spendy() {
	for i := 1; i < 1000; i++ {
		lock.Lock()
		balance -= 10
		lock.Unlock()
		time.Sleep(1 * time.Microsecond)
	}
	println("Spendy Done!")
}

func main() {
	go stingy()
	go spendy()
	time.Sleep(3000 * time.Millisecond)
	print(balance)

}
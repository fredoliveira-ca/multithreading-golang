package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	balance = 100
	lock    = sync.Mutex{}
	moneyDeposited = sync.NewCond(&lock)
)

func stingy() {
	for i := 1; i < 1000; i++ {
		lock.Lock()
		balance += 10
		fmt.Println("Stingy sees balance of ", balance)
		moneyDeposited.Signal()
		lock.Unlock()
		time.Sleep(1 * time.Microsecond)
	}
	println("Stingy Done!")
}

func spendy() {
	for i := 1; i < 1000; i++ {
		lock.Lock()
		for balance-20 < 0 {
			moneyDeposited.Wait()
		}
		balance -= 20
		fmt.Println("Spendy sees balance of ", balance)
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

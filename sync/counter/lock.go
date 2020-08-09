package main

import (
	"fmt"
	"sync"
	"time"
)

var lock = sync.Mutex{}
var rwlock = sync.RWMutex{}

func oneTwoThreeA() {
	lock.Lock()
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Millisecond)
	}
	lock.Unlock()
}

func StartThreadsA() {
	for i := 1; i <= 2; i++ {
		go oneTwoThreeA()
	}
	time.Sleep(1 * time.Second)
}

func oneTwoThreeB() {
	rwlock.RLock()
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
		time.Sleep(1 * time.Millisecond)
	}
	rwlock.RLock()
}

func StartThreadsB() {
	for i := 1; i <= 2; i++ {
		go oneTwoThreeB()
	}
	time.Sleep(1 * time.Second)
}

func callLockTwice() {
	lock.Lock()
	lock.Lock()
	fmt.Print("Hello there")
}

func RunAndWait() {
	go callLockTwice()
	time.Sleep(10 * time.Second)
}

func main() {
	fmt.Println("Starting ThreadsA")
	StartThreadsA()
	//Since the readers lock allows multiple threads to execute the same read locked block of code
	//it's impossible to know the exact order that the threads will be interleaved.
	fmt.Println("Starting ThreadsB")
	StartThreadsB()

	//Interesting fact: The mutexes in GO are not re-entrant. Using re-entrant lock,
	//a thread calling lock() while already holding the same lock would have no affect and would return immediately.
	fmt.Println("Starting callLockTwice")
	// The second time called is called, the call will block forever. Tje program will exit after 10 secondes.
	// Hello there is never printed
	RunAndWait()
}

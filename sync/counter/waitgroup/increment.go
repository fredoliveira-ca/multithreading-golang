package main

import (
	"fmt"
	"sync"
)

// The main thread creates 5 threads and waits for them to complete. 
// Each thread adds 100 to the counter. 
// However threads might overwrite each other, since access to x is not exclusive. 
// The final answer of X can be various numbers due to this.
func count() {
	wg := sync.WaitGroup{}
	x := 0
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go increment(&x, &wg)
	}
	wg.Wait()
	fmt.Printf("%d\n", x)
}

func increment(x *int, wg *sync.WaitGroup) {
	for i := 0; i < 100; i++ {
		*x++
	}
	wg.Done()
}

func main() {
	count()
}

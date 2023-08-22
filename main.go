package main

import (
	"fmt"
	"postinit/logger"
	"postinit/prog"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	// initializing logger instance
	go logger.Init(&wg)

	// simulating calls to function before init
	for i := 0; i < 5; i++ {
		// service/caller function which uses logs, assuming that logger is already initialized
		err1 := prog.Caller1(fmt.Sprintf("func: caller1 value-%v", i+1))
		if err1 != nil {
			fmt.Printf("Error on caller1 %v: %v\n", i+1, err1)
		}
		// service/caller function which uses logs, safe handling initialization errors
		err2 := prog.Caller2(fmt.Sprintf("func: caller2 value-%v", i+1))
		if err2 != nil {
			fmt.Printf("Error on caller2 %v: %v\n", i+1, err2)
		}
	}

	wg.Wait()

}

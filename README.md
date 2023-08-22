# postinit
A sample application to demonstrate post execution using [channels](https://pkg.go.dev/github.com/eapache/channels) (by creating an internal queue) in GO.

## Problem

This is a common problem which occurs when you try to access functions where your instance is un-initialized.

### Example 

In a basic webservice application if you are initializing logger on web-server start, you may not be able to log environment configuration read through flags/config files before starting the server.

You may get errors/panics due to unavailability of config files and it would unsafe & inconvenient to track down such issues.

Let's create a small application to regenerate the problem!!!

Consider a simple Go Application having following packages:

* `logger` - This would act as a logging package.
* `prog` - This would act as a caller/consumer package.

To depict the logging instance, I have created a struct inside `logger` package.

```
type Logger struct {
	instance string
}
```

Just to pretend that.

```
if len(logger.instance) > 0 {
    // initialized state
} else {
    // un-initialized state
}
```
Here is the init function of logger.

```
func Init(wg *sync.WaitGroup) {
	time.Sleep(5 * time.Second) // artificial delay
	fmt.Println("Init")
	logger = Logger{
		instance: "Logger",
	}
    wg.Done()
}
```

here is what Log function looks like.

```
func Log(msg string) error {
	if len(logger.instance) > 0 {
		fmt.Printf("got message to log `%v`\n", msg)
	} else {
		return errors.New("logger instance not initialized")
	}
	return nil
}
```

**NOTE:** This is not an actual logger but just a mere depiction.


and here is what our caller/consumer function looks like (`prog` package).

```
func Caller1(msg string) error {
	return logger.Log(msg)
}
```


Then in our `main.go`
```
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
    }
}
```

The above code would give us uninitialized error because initialization takes up to 5 seconds and happening in a separate [Go Routine](https://go.dev/doc/effective_go#goroutines) and we started consuming the logger instance (directly or indirectly).


## Solution

Let's extend the logger by creating a private struct for holding incoming log data.

```
type logDetail struct {
	message string
}
```

For the sake of simplicity, I am just using `message` field but in real world scenario it is much more than that.

Create a collection of channels, of type `logDetail` on package level.

```
var chann = make([]chan logDetail, 0)
```

Add a `safelog` function with same definition as the `log` function.


Transform the data into `logDetail` structure, send it into the channel and append it to the package level collection.

```
func SafeLog(msg string) error {
	if len(logger.instance) > 0 {
		return Log(msg)
	} else {
		ld := logDetail{
			message: msg,
		}
		ch := make(chan logDetail, 1)
		ch <- ld
		chann = append(chann, ch)
	}
	return nil
}
```

You can add instance check and delegate it back to original `Log` function if you have the instance already initialized.

This would form a queue inside the logger package to be consumed later.

**NOTE:** I have used a buffered channel here because we are not speciying a receiver at this point, otherwise this would become a ***blocking call***

Now finally, edit the init function to consume the pending logs *(if any)*

```
// fake delaying logger initialization
func Init(wg *sync.WaitGroup) {
	time.Sleep(5 * time.Second)
	fmt.Println("Init")
	logger = Logger{
		instance: "Logger",
	}
	for _, v := range chann {
		val := <-v
		Log(val.message)
		close(v)
	}
	wg.Done()
}
```


To test this, add another function in consumer package and call the `safeLog` function

```
func Caller2(msg string) error {
	return logger.SafeLog(msg)
}
```
Then in `main.go` add a call to this new consumer function

```
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
```

Here is the output from debug console

```
Error on caller1 1: logger instance not initialized
Error on caller1 2: logger instance not initialized
Error on caller1 3: logger instance not initialized
Error on caller1 4: logger instance not initialized
Error on caller1 5: logger instance not initialized
Init
got message to log `func: caller2 value-1`
got message to log `func: caller2 value-2`
got message to log `func: caller2 value-3`
got message to log `func: caller2 value-4`
got message to log `func: caller2 value-5`
```

As you can see, requests to caller1 function failed but requests to caller2 function gets pushed into the queue to be consumed after logger initialization.


This is how you create an internal queue using [channels](https://pkg.go.dev/github.com/eapache/channels)


## Feel free to edit/expand/explore this repository

For feedback and queries, reach me on LinkedIn at [here](https://www.linkedin.com/in/usama28232/?original_referer=)
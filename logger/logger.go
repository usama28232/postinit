package logger

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var logger Logger

type logDetail struct {
	message string
}

type Logger struct {
	instance string
}

var chann = make([]chan logDetail, 0)

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

func Log(msg string) error {
	if len(logger.instance) > 0 {
		fmt.Printf("got message to log `%v`\n", msg)
	} else {
		return errors.New("logger instance not initialized")
	}
	return nil
}

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

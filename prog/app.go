package prog

import (
	"postinit/logger"
)

// simply calls logger.Log assuming that logger is already initialized
func Caller1(msg string) error {
	return logger.Log(msg)
}

// calls logger.SafeLog checking if logger is initialized
// if not, it queues the message into a channel to be consumed later
func Caller2(msg string) error {
	return logger.SafeLog(msg)
}

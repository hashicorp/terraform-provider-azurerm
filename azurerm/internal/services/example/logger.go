package example

import (
	"fmt"
	"log"
)

type ExampleLogger struct {
}

func (ExampleLogger) Info(message string) {
	log.Printf(message)
}

func (ExampleLogger) InfoF(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf(format, args))
}

func (ExampleLogger) Warn(message string) {
	log.Printf(message)
}

func (ExampleLogger) WarnF(format string, args ...interface{}) {
	log.Printf(fmt.Sprintf(format, args))
}

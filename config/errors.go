package config

import (
	"fmt"
	"os"
	"time"
)

func NewErrorHandler() ErrorHandler {
	handler := ErrorHandler{
		make(chan chan error),
		make(chan struct{}),
	}
	return handler
}

type ErrorHandler struct {
	newErrors chan chan error
	quit      chan struct{}
}

func (handler *ErrorHandler) Run() {
	go func() {
	loop:
		for {
			select {
			case errors := <-handler.newErrors:
				go handler.handle(errors)
			case <-handler.quit:
				break loop
			default:
				time.Sleep(time.Millisecond * 10)
			}
		}
	}()
}

func (handler *ErrorHandler) Handle(errors chan error) {
	handler.newErrors <- errors
}

func (handler *ErrorHandler) Quit() {
	handler.quit <- struct{}{}
}

func (handler *ErrorHandler) handle(errors chan error) {
	// todo: Add handler interface for handling each error
	f, _ := os.Create("error.log")
	defer f.Close()
	for err := range errors {
		a := fmt.Sprintf("%v", err)
		f.WriteString(a)
	}
}

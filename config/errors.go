package config

import (
	"fmt"
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
	for err := range errors {
		fmt.Println(err)
	}
}

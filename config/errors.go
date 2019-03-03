package config

import (
	"fmt"
	"time"
)

func NewErrorHandler() errorHandler {
	handler := errorHandler{
		make(chan chan error),
		make(chan struct{}),
	}
	return handler
}

type errorHandler struct {
	newErrors chan chan error
	quit      chan struct{}
}

func (handler *errorHandler) Run() {
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

func (handler *errorHandler) Handle(errors chan error) {
	handler.newErrors <- errors
}

func (handler *errorHandler) Quit() {
	handler.quit <- struct{}{}
}

func (handler *errorHandler) handle(errors chan error) {
	// todo: Add handler interface for handling each error
	for err := range errors {
		fmt.Println(err)
	}
}

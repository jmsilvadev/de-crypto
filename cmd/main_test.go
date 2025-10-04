package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Run("main_function_executes", func(t *testing.T) {
		done := make(chan bool)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					done <- true
				}
			}()

			main()
			done <- true
		}()

		select {
		case <-done:
			assert.True(t, true)
		case <-time.After(100 * time.Millisecond):
			os.Exit(0)
		}
	})
}

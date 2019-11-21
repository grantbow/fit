package issues

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("caught %v\n", r)
		}
	}()
	e := errors.New("test")
	check(e)
	time.Sleep(1 * time.Second) // ensure time to recover
}

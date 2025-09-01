package safe

import (
	"fmt"

	"github.com/pkg/errors"
)

// defaultRecover default recover function to deal with panic.
func defaultRecover(r any) {
	fmt.Printf("recover a panic: %v", r)
}

// defaultRecover default recover function to deal with panic
// and transmit error to external of goroutine.
func defaultRecoverErr(r any, err *error) {
	var ok bool
	if *err, ok = r.(error); !ok {
		*err = errors.Errorf("unknown error: %v", r)
	}
}

// GoSafe call goroutine with safe(recover) mode. we can use
// external channel or variable to receive error message.
func GoSafe(fn func(), rfn func(any)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if rfn != nil {
					rfn(r)
				} else {
					defaultRecover(r)
				}
			}
		}()
		fn()
	}()
}

// GoSafeErr call goroutine with safe(recover) mode, and fulfill error to input args.
// notice: we have to init the argument of err before call the function.
func GoSafeErr(fn func(), rfn func(any, *error), err *error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if rfn != nil {
					rfn(r, err)
				} else {
					defaultRecoverErr(r, err)
				}
			}
		}()
		fn()
	}()
}

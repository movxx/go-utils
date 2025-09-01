package safe

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSafeGo(t *testing.T) {
	testText := "test panic!"
	testErr := errors.New(testText)
	var err error
	ch := make(chan error, 1)
	errfn := func(r any) {
		if err, ok := r.(error); ok {
			ch <- err
		} else {
			ch <- errors.Errorf("unknown panic: %v", r)
		}
	}

	errfn2 := func(r any) {
		var ok bool
		if err, ok = r.(error); !ok {
			err = errors.Errorf("unknown panic: %v", r)
		}
	}

	for _, ts := range []struct {
		name   string
		withCh bool
		fn     func()
		errfn  func(any)
		expect error
	}{
		{
			name:   "happy path - panic string and recover error",
			withCh: true,
			fn:     func() { panic(testText) },
			errfn:  errfn,
			expect: errors.Errorf("unknown panic: %v", testText),
		},
		{
			name:   "happy path - panic error and recover error",
			withCh: true,
			fn:     func() { panic(testErr) },
			errfn:  errfn,
			expect: testErr,
		},
		{
			name:   "happy path - panic string and recover error",
			withCh: false,
			fn:     func() { panic(testText) },
			errfn:  errfn2,
			expect: errors.Errorf("unknown panic: %v", testText),
		},
		{
			name:   "happy path - panic error and recover error",
			withCh: false,
			fn:     func() { panic(testErr) },
			errfn:  errfn2,
			expect: testErr,
		},
	} {
		t.Run(ts.name, func(t *testing.T) {
			GoSafe(ts.fn, ts.errfn)
			if ts.withCh {
				// use the external channel of goroutine to receive error
				panicErr := <-ch
				assert.Equal(t, ts.expect.Error(), panicErr.Error())
			} else {
				time.Sleep(time.Second)
				// use the external variable of goroutine to receive error
				assert.Equal(t, ts.expect.Error(), err.Error())
			}

		})
	}
}

func TestSafeGoErr(t *testing.T) {
	testText := "test panic!"
	testErr := errors.New(testText)
	errfn := func(r any, err *error) {
		var ok bool
		if *err, ok = r.(error); !ok {
			*err = errors.Errorf("unknown error: %v", r)
		}
	}

	for _, ts := range []struct {
		name   string
		fn     func()
		errfn  func(any, *error)
		expect error
	}{
		{
			name:   "happy path - panic string and recover error",
			fn:     func() { panic(testText) },
			errfn:  errfn,
			expect: errors.Errorf("unknown error: %v", testText),
		},
		{
			name:   "happy path - panic error and recover error",
			fn:     func() { panic(testErr) },
			errfn:  errfn,
			expect: testErr,
		},
	} {
		t.Run(ts.name, func(t *testing.T) {
			var err error
			GoSafeErr(ts.fn, ts.errfn, &err)
			time.Sleep(time.Second)
			// t.Logf("err: %v", *err)
			assert.Equal(t, ts.expect.Error(), err.Error())
		})
	}
}

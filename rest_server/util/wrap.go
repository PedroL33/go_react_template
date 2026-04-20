package util

import (
	"fmt"
	"runtime"
	"strings"
)

// Wrap annotates err with the caller's function name.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return err
	}
	fn := runtime.FuncForPC(pc).Name()
	if i := strings.LastIndex(fn, "/"); i >= 0 {
		fn = fn[i+1:]
	}
	return fmt.Errorf("%s: %w", fn, err)
}

// Wrapf adds a message too.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	pc, _, _, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()
	if i := strings.LastIndex(fn, "/"); i >= 0 {
		fn = fn[i+1:]
	}
	return fmt.Errorf("%s: "+format+": %w", append([]any{fn}, append(args, err)...)...)
}

// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"fmt"
	"runtime"
)

// NewFrameError wraps the specified error with an error that provides stack frame information.
// Call this at the inner error's origin to provide file name and line number info with the error.
// You MUST supply an inner error.
// DO NOT ARBITRARILY CALL THIS TO WRAP ERRORS!  There MUST be only ONE error of this type in the chain.
func NewFrameError(inner error, stackTrace bool, skipFrames, totalFrames int) error {
	fe := FrameError{inner: inner, info: "stack trace unavailable"}
	if stackTrace {
		// the skipFrames+3 is to skip runtime.Callers(), StackTrace and ourselves
		fe.info = StackTrace(skipFrames+3, totalFrames)
	} else if pc, file, line, ok := runtime.Caller(skipFrames + 1); ok {
		// the skipFrames + 1 is to skip ourselves
		frame := runtime.FuncForPC(pc)
		fe.info = fmt.Sprintf("%s()\n\t%s:%d\n", frame.Name(), file, line)
	}
	return &fe
}

// FrameError associates an error with stack frame information.
// Exported for testing purposes, use NewFrameError().
type FrameError struct {
	inner error
	info  string
}

// Error implements the error interface for type FrameError.
func (f *FrameError) Error() string {
	return fmt.Sprintf("%s:\n%s\n", f.inner.Error(), f.info)
}

// Unwrap returns the inner error.
func (f *FrameError) Unwrap() error {
	return f.inner
}

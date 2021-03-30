// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package runtime

import (
	"fmt"
	"runtime"
	"strings"
)

// StackTrace returns a formatted stack trace string.
// skipFrames - the number of stack frames to skip before composing the trace string.
// totalFrames - the maximum number of stack frames to include in the trace string.
func StackTrace(skipFrames, totalFrames int) string {
	sb := strings.Builder{}
	pcCallers := make([]uintptr, totalFrames)
	runtime.Callers(skipFrames, pcCallers)
	frames := runtime.CallersFrames(pcCallers)
	for {
		frame, more := frames.Next()
		sb.WriteString(frame.Function)
		sb.WriteString("()\n\t")
		sb.WriteString(frame.File)
		sb.WriteRune(':')
		sb.WriteString(fmt.Sprintf("%d\n", frame.Line))
		if !more {
			break
		}
	}
	return sb.String()
}

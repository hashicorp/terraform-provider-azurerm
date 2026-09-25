//go:build !darwin || !arm64
// +build !darwin !arm64

package gomonkey

func buildPrivateMethodDirective(_, double, _ uintptr) ([]byte, uintptr) {
	return buildJmpDirective(double), 0
}

func releasePrivateMethodTrampoline(uintptr) {}

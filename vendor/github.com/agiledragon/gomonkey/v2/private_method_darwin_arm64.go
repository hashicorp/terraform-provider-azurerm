//go:build darwin && arm64
// +build darwin,arm64

package gomonkey

import (
	"encoding/binary"
	"fmt"
	"sync"
)

const (
	privateMethodTrampolineSize  = 24
	privateMethodTrampolineCount = 128
	arm64BranchRange             = 128 * 1024 * 1024
)

var privateMethodTrampolines struct {
	sync.Mutex
	used [privateMethodTrampolineCount]bool
}

// privateMethodTrampolineBase returns the address of the executable slot area
// declared in private_method_darwin_arm64.s.
func privateMethodTrampolineBase() uintptr

func buildPrivateMethodDirective(target, double, trampoline uintptr) ([]byte, uintptr) {
	acquired := false
	if trampoline == 0 {
		trampoline = acquirePrivateMethodTrampoline()
		acquired = true
	}
	completed := false
	defer func() {
		if acquired && !completed {
			releasePrivateMethodTrampoline(trampoline)
		}
	}()

	code, err := buildArm64BranchDirective(target, trampoline)
	if err != nil {
		panic(err)
	}

	// Keep the complete 24-byte function-value jump in an isolated executable
	// slot. The target itself only receives one 4-byte B instruction, so a short
	// method cannot overwrite the function laid out immediately after it.
	modifyBinary(trampoline, buildJmpDirective(double))
	completed = true
	return code, trampoline
}

func buildArm64BranchDirective(from, to uintptr) ([]byte, error) {
	delta := int64(to) - int64(from)
	if delta%4 != 0 || delta < -arm64BranchRange || delta >= arm64BranchRange {
		return nil, fmt.Errorf("private method trampoline %#x is out of ARM64 branch range from target %#x", to, from)
	}

	instruction := uint32(0x14000000) | (uint32(delta>>2) & 0x03ffffff)
	code := make([]byte, 4)
	binary.LittleEndian.PutUint32(code, instruction)
	return code, nil
}

func acquirePrivateMethodTrampoline() uintptr {
	privateMethodTrampolines.Lock()
	defer privateMethodTrampolines.Unlock()

	base := privateMethodTrampolineBase()
	for i := range privateMethodTrampolines.used {
		if !privateMethodTrampolines.used[i] {
			privateMethodTrampolines.used[i] = true
			return base + uintptr(i*privateMethodTrampolineSize)
		}
	}
	panic(fmt.Sprintf("private method trampoline capacity exceeded: %d active patches", privateMethodTrampolineCount))
}

func releasePrivateMethodTrampoline(trampoline uintptr) {
	base := privateMethodTrampolineBase()
	if trampoline < base {
		return
	}
	offset := trampoline - base
	if offset%privateMethodTrampolineSize != 0 {
		return
	}
	index := offset / privateMethodTrampolineSize
	if index >= privateMethodTrampolineCount {
		return
	}

	privateMethodTrampolines.Lock()
	privateMethodTrampolines.used[index] = false
	privateMethodTrampolines.Unlock()
}

package gomonkey

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

func PtrOf(val []byte) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&val)).Data
}

func modifyBinary(target uintptr, bytes []byte) {
	pageSize := syscall.Getpagesize()
	if err := validateWritePageSize(pageSize); err != nil {
		panic(err)
	}

	targetPage := pageStart(target)
	res := write(target, PtrOf(bytes), len(bytes), targetPage,
		protectSize(target, len(bytes)), syscall.PROT_READ|syscall.PROT_EXEC)
	if res != 0 {
		panic(fmt.Errorf("failed to write memory, code %v", res))
	}
}

// validateWritePageSize ensures that write's static assembly padding is wide
// enough to keep its executable instructions outside the page being modified.
// Without this check, a larger page size could make write remove execute
// permission from its own current page and terminate the process with SIGBUS.
func validateWritePageSize(pageSize int) error {
	if pageSize > writeIsolationSize {
		return fmt.Errorf("unsupported system page size %d: write isolation supports at most %d bytes",
			pageSize, writeIsolationSize)
	}
	return nil
}

// protectSize returns the number of bytes that must be made writable so that a
// patch of n bytes starting at target is fully covered, even when the patch
// spans a page boundary.
func protectSize(target uintptr, n int) int {
	targetPage := pageStart(target)
	endPage := pageStart(target + uintptr(n) - 1)
	size := syscall.Getpagesize()
	if endPage > targetPage {
		size = int(endPage-targetPage) + syscall.Getpagesize()
	}
	return size
}

//go:cgo_import_dynamic mach_task_self mach_task_self "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic mach_vm_protect mach_vm_protect "/usr/lib/libSystem.B.dylib"
//go:cgo_import_dynamic sys_icache_invalidate sys_icache_invalidate "/usr/lib/libSystem.B.dylib"
func write(target, data uintptr, len int, page uintptr, pageSize, oriProt int) int

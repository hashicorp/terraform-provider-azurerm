package testclient

// Maps goroutines to their active tests so VCR recorder always know which test context it's in, mostly for running
// in parallel.

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var goroutineTests sync.Map

// RegisterTestT associates t with the current goroutine.
// Called at the entry point of each acceptance test run.
// Note: Ensure a defer testclient.UnregisterTestT() is called immediately after this!
func RegisterTestT(t *testing.T) {
	goroutineTests.Store(currentGoroutineID(), t)
}

// UnregisterTestT removes the association for the current goroutine.
// Note: Call in a defer immediately after RegisterTest()
func UnregisterTestT() {
	goroutineTests.Delete(currentGoroutineID())
}

// CurrentTestName returns the name of the test running in this goroutine,
// or "" if none is registered.
func CurrentTestName() string {
	if v, ok := goroutineTests.Load(currentGoroutineID()); ok {
		return v.(*testing.T).Name()
	}
	return ""
}

func currentGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	fields := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))
	id, _ := strconv.ParseInt(fields[0], 10, 64)
	return id
}

// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package locks

import (
	"context"
	"slices"
)

// armMutexKV is the instance of MutexKV for ARM resources
var (
	armMutexKV     = newMutexKV()
	armSemaphoreKV = newSemaphoreKV()
)

func ByID(id string) {
	armMutexKV.Lock(id)
}

// handle the case of using the same name for different kinds of resources
func ByName(name string, resourceType string) {
	armMutexKV.Lock(getLockName(name, resourceType))
}

func MultipleByID(ids *[]string) {
	newSlice := removeDuplicatesFromStringArray(*ids)

	slices.Sort(newSlice)

	for _, id := range newSlice {
		ByID(id)
	}
}

func MultipleByName(names *[]string, resourceType string) {
	newSlice := removeDuplicatesFromStringArray(*names)

	slices.Sort(newSlice)

	for _, name := range newSlice {
		ByName(name, resourceType)
	}
}

func UnlockByID(id string) {
	armMutexKV.Unlock(id)
}

func UnlockByName(name string, resourceType string) {
	armMutexKV.Unlock(getLockName(name, resourceType))
}

func UnlockMultipleByID(ids *[]string) {
	newSlice := removeDuplicatesFromStringArray(*ids)

	for _, id := range newSlice {
		UnlockByID(id)
	}
}

func UnlockMultipleByName(names *[]string, resourceType string) {
	newSlice := removeDuplicatesFromStringArray(*names)

	for _, name := range newSlice {
		UnlockByName(name, resourceType)
	}
}

func NewWeightedByName(name string, resourceType string, max int64) {
	armSemaphoreKV.NewWeighted(getLockName(name, resourceType), max)
}

func AcquireOneByName(ctx context.Context, name string, resourceType string) error {
	return armSemaphoreKV.Acquire(ctx, getLockName(name, resourceType), 1)
}

func ReleaseOneByName(name string, resourceType string) {
	armSemaphoreKV.Release(getLockName(name, resourceType), 1)
}

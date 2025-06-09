// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package locks

import "slices"

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = newMutexKV()

func ByID(id string) {
	armMutexKV.Lock(id)
}

// handle the case of using the same name for different kinds of resources
func ByName(name string, resourceType string) {
	updatedName := resourceType + "." + name
	armMutexKV.Lock(updatedName)
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
	updatedName := resourceType + "." + name
	armMutexKV.Unlock(updatedName)
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

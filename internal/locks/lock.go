// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package locks

import "slices"

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = newMutexKV()

func ByID(id string) {
	armMutexKV.Lock(id)
}

func MultipleByID(ids *[]string) {
	newSlice := removeDuplicatesFromStringArray(*ids)

	slices.Sort(newSlice)

	for _, id := range newSlice {
		ByID(id)
	}
}

func UnlockByID(id string) {
	armMutexKV.Unlock(id)
}

func UnlockMultipleByID(ids *[]string) {
	newSlice := removeDuplicatesFromStringArray(*ids)

	for _, id := range newSlice {
		UnlockByID(id)
	}
}

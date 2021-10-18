package locks

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = NewMutexKV()

func ByID(id string) {
	armMutexKV.Lock(id)
}

// handle the case of using the same name for different kinds of resources
func ByName(name string, resourceType string) {
	updatedName := resourceType + "." + name
	armMutexKV.Lock(updatedName)
}

func MultipleByName(names *[]string, resourceType string) {
	newSlice := removeDuplicatesFromStringArray(*names)

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

func UnlockMultipleByName(names *[]string, resourceType string) {
	newSlice := removeDuplicatesFromStringArray(*names)

	for _, name := range newSlice {
		UnlockByName(name, resourceType)
	}
}

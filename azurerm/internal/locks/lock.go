package locks

import "github.com/hashicorp/terraform/helper/mutexkv"

// armMutexKV is the instance of MutexKV for ARM resources
var armMutexKV = mutexkv.NewMutexKV()

func ByID(id string) {
	armMutexKV.Lock(id)
}

// handle the case of using the same name for different kinds of resources
func ByName(name string, resourceType string) {
	updatedName := resourceType + "." + name
	armMutexKV.Lock(updatedName)
}

func MultipleByName(names *[]string, resourceType string) {
	for _, name := range *names {
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
	for _, name := range *names {
		UnlockByName(name, resourceType)
	}
}

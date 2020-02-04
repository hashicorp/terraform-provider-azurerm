package locks

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/mutexkv"
)

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
	for i, name := range *names {
		// at the end of every array item add its index. This way we guarantee that this item will be unique (no duplicates are possible)
		uniqueValue := fmt.Sprintf("%s-arrIdx-%d", name, i)

		ByName(uniqueValue, resourceType)
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
	for i, name := range *names {
		// at the end of every array item add its index. We need to add this sufix because it was added during the lock process
		uniqueValue := fmt.Sprintf("%s-arrIdx-%d", name, i)

		UnlockByName(uniqueValue, resourceType)
	}
}

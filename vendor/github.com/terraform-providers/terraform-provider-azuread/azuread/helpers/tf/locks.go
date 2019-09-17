package tf

import "github.com/hashicorp/terraform/helper/mutexkv"

// mutex is the instance of MutexKV for ARM resources
var mutex = mutexkv.NewMutexKV()

// handles the case of using the same name for different kinds of resources
func LockByName(resourceType string, name string) {
	mutex.Lock(resourceType + "." + name)
}

func UnlockByName(resourceType string, name string) {
	mutex.Unlock(resourceType + "." + name)
}

package azuread

// handles the case of using the same name for different kinds of resources
func azureADLockByName(resourceType string, name string) {
	armMutexKV.Lock(resourceType + "." + name)
}

func azureADUnlockByName(resourceType string, name string) {
	armMutexKV.Unlock(resourceType + "." + name)
}

package azurerm

// handle the case of using the same name for different kinds of resources
func azureRMLockByName(name string, resourceType string) {
	updatedName := resourceType + "." + name
	armMutexKV.Lock(updatedName)
}

func azureRMLockMultipleByName(names *[]string, resourceType string) {
	for _, name := range *names {
		azureRMLockByName(name, resourceType)
	}
}

func azureRMUnlockByName(name string, resourceType string) {
	updatedName := resourceType + "." + name
	armMutexKV.Unlock(updatedName)
}

func azureRMUnlockMultipleByName(names *[]string, resourceType string) {
	for _, name := range *names {
		azureRMUnlockByName(name, resourceType)
	}
}

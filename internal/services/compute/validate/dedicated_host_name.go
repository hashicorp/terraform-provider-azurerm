package validate

func DedicatedHostName() func(i interface{}, k string) (warnings []string, errors []error) {
	return DedicatedHostGroupName()
}

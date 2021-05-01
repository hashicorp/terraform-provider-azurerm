package validate

import (
	"fmt"
	"regexp"
)

func InferenceClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// The portal says: The workspace name must be between 1 and 16 characters. The name may only include alphanumeric characters and '-'.
	// If you provide invalid name, the rest api will return an error with the following regex.
	if matched := regexp.MustCompile(`^[a-zA-Z0-9][\w-]{1,15}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 2 and 16 characters, and may only include alphanumeric characters and '-' character", k))
	}
	return
}

func KubernetesClusterResourceGroupName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// Azure portal (not https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules) says:
	// name for managed clusters have to be between 1 and 90 characters. The name may only include alphanumeric characters and '-, _'.
	// start and end has to be alphanumeric. If you provide invalid name, the rest api will return an error with the following regex.
	if matched := regexp.MustCompile(`^[a-zA-Z0-9][\w-_]{1,90}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 1 and 90 characters, and may only include alphanumeric characters and '-, :' character", k))
	}
	return
}

func NodePoolName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9][\w-_]{1,12}$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 1 and 12 characters, and may only include alphanumeric characters and '-, :' character", k))
	}
	return
}

func ClusterPurpose(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	// Azure portal (not https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules) says:
	// name for managed clusters have to be between 1 and 90 characters. The name may only include alphanumeric characters and '-, _'.
	// start and end has to be alphanumeric. If you provide invalid name, the rest api will return an error with the following regex.

	switch v {
	case
		"Prod",
		"Dev",
		"Test":
		return
	}
	errors = append(errors, fmt.Errorf("%s must be one of \"Prod\", \"Dev\", \"Test\" ", k))
	return
}

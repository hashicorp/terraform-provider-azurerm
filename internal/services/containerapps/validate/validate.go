// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
)

func InitTimeout(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^[0-9]+[smh]$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q should be a duration in whole intervals of seconds, minutes, or hours. e.g. 5s, 10m, 1h", k))
	}

	return
}

func DaprComponentName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([a-z])[a-z0-9-]{0,58}[a-z]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") || strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character and cannot have '--'. The length must not be more than 60 characters", k))
		return
	}

	return
}

func SecretName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") || strings.HasSuffix(v, ".") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character", k))
	}
	return
}

func CertificateName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([a-z0-9])[a-z0-9-.]*[a-z0-9]$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") || strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character", k))
		return
	}

	return
}

func ContainerAppName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([a-z])[a-z0-9-]{0,58}[a-z]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") || strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character and cannot have '--'. The length must not be more than 60 characters", k))
		return
	}

	return
}

func ContainerCpu(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be float", k))
		return
	}
	switch v {
	case 0.25, 0.5, 0.75, 1.0, 1.25, 1.5, 1.75, 2.0:
	default:
		errors = append(errors, fmt.Errorf("%s must be one of `0.25`, `0.5`, `0.75`, `1.0`, `1.25`, `1.5`, `1.75`, or `2.0`", k))
	}

	return
}

func ManagedEnvironmentStorageName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([a-z])[a-z0-9-]{0,30}[a-z]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") || strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character and cannot have '--'. The length must not be more than 32 characters", k))
		return
	}

	return
}

func ManagedEnvironmentName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([a-zA-Z])[a-zA-Z0-9-]{0,58}[a-z]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character. The length must not be more than 60 characters", k))
		return
	}
	return
}

func ContainerAppContainerName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if matched := regexp.MustCompile(`^([a-zA-Z0-9])[a-zA-Z0-9-.]{0,254}[a-z]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters, '-', or '.', start with an alphabetic character, and end with an alphanumeric character. The length must not be more than 60 characters", k))
		return
	}
	return
}

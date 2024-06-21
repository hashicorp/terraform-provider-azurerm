// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strconv"
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

	if matched := regexp.MustCompile(`^([a-z])[a-z0-9-]{0,58}[a-z0-9]?$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") || strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character and cannot have '--'. The length must not be more than 32 characters", k))
		return
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
	}

	return
}

func ContainerAppJobName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) == 1 {
		if matched := regexp.MustCompile(`^[a-z0-9]$`).Match([]byte(v)); !matched {
			errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters, '-', or '.', start and end with an alphanumeric character", k))
		}
	} else {
		if matched := regexp.MustCompile(`^([a-z0-9])[a-z0-9-]*[a-z0-9]$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") {
			errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters, or '-', start and end with an alphanumeric character", k))
		}
	}

	if len(v) > 32 {
		errors = append(errors, fmt.Errorf("%q must not exceed 32 characters", k))
	}

	if strings.Contains(v, "--") {
		errors = append(errors, fmt.Errorf("%q must not contain --", k))
	}

	return
}

func LowerCaseAlphaNumericWithHyphensAndPeriods(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) == 1 {
		if matched := regexp.MustCompile(`^[a-z0-9]$`).Match([]byte(v)); !matched {
			errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters, '-', or '.', start and end with an alphanumeric character", k))
		}
	} else {
		if matched := regexp.MustCompile(`^([a-z0-9])[a-z0-9-.]*[a-z0-9]$`).Match([]byte(v)); !matched || strings.HasSuffix(v, "-") {
			errors = append(errors, fmt.Errorf("%q must consist of lower case alphanumeric characters, '-', or '.', start and end with an alphanumeric character", k))
		}
	}

	return
}

func ContainerAppScaleRuleConcurrentRequests(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	c, err := strconv.Atoi(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s to be a string representation of an integer, got %+v", k, v))
		return
	}

	if c <= 0 {
		errors = append(errors, fmt.Errorf("value for %s must be at least `1`, got %d", k, c))
	}

	return
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
)

func CloudVMClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if len(v) < 1 || len(v) > 255 {
		errors = append(errors, fmt.Errorf("name must be %d to %d characters", 1, 255))
		return warnings, errors
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		errors = append(errors, fmt.Errorf("name must start with a letter or underscore (_)"))
		return warnings, errors
	}

	re := regexp.MustCompile("--")
	if re.MatchString(v) {
		errors = append(errors, fmt.Errorf("name must not contain any consecutive hyphers (--)"))
		return warnings, errors
	}

	return warnings, errors
}

func CpuCoreCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return warnings, errors
	}

	if v < 2 {
		errors = append(errors, fmt.Errorf("cpu_core_count must be at least %v", 2))
		return warnings, errors
	}

	return warnings, errors
}

func DataStorageSizeInTbs(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be float", k))
		return warnings, errors
	}

	if v < 2 || v > 192 {
		errors = append(errors, fmt.Errorf("%v must be between %v and %v", k, 2, 192))
		return warnings, errors
	}

	return warnings, errors
}

func LicenseModel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	if v != string(cloudvmclusters.LicenseModelBringYourOwnLicense) && v != string(cloudvmclusters.LicenseModelLicenseIncluded) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(cloudvmclusters.LicenseModelBringYourOwnLicense), string(cloudvmclusters.LicenseModelLicenseIncluded)))
		return warnings, errors
	}

	return warnings, errors
}

func DataStoragePercentage(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return warnings, errors
	}

	if v != 35 && v != 40 && v != 60 && v != 80 {
		errors = append(errors, fmt.Errorf("%v must 35, 40, 60 or 80", k))
		return warnings, errors
	}

	return warnings, errors
}

func SystemVersion(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return warnings, errors
	}

	pattern := "(?:19|22|23|24|25)\\.[0-9]+(\\.[0-9]+)*|[0-9]+(\\.[0-9]+)*"
	re := regexp.MustCompile(pattern)
	if !re.MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must match one of the following patterns: %v", k, pattern))
	}

	return warnings, errors
}

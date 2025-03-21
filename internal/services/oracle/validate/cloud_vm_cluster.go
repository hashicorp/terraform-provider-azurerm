// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
)

func CloudVMClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 || len(v) > 255 {
		errors = append(errors, fmt.Errorf("name must be %d to %d characters", 1, 255))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		errors = append(errors, fmt.Errorf("name must start with a letter or underscore (_)"))
		return
	}

	re := regexp.MustCompile("--")
	if re.MatchString(v) {
		errors = append(errors, fmt.Errorf("name must not contain any consecutive hyphers (--)"))
		return
	}

	return
}

func CpuCoreCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 2 {
		errors = append(errors, fmt.Errorf("cpu_core_count must be at least %v", 2))
		return
	}

	return
}

func DataStorageSizeInTbs(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(float64)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be float", k))
		return
	}

	if v < 2 || v > 192 {
		errors = append(errors, fmt.Errorf("%v must be between %v and %v", k, 2, 192))
		return
	}

	return
}

func LicenseModel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if v != string(cloudvmclusters.LicenseModelBringYourOwnLicense) && v != string(cloudvmclusters.LicenseModelLicenseIncluded) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(cloudvmclusters.LicenseModelBringYourOwnLicense), string(cloudvmclusters.LicenseModelLicenseIncluded)))
		return
	}

	return
}

func DataStoragePercentage(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v != 35 && v != 40 && v != 60 && v != 80 {
		errors = append(errors, fmt.Errorf("%v must 35, 40, 60 or 80", k))
		return
	}

	return
}

func SystemVersion(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	pattern := "(?:19|22|23|24|25)\\.[0-9]+(\\.[0-9]+)*|[0-9]+(\\.[0-9]+)*"
	re := regexp.MustCompile(pattern)
	if !re.MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must match one of the following patterns: %v", k, pattern))
	}

	return
}

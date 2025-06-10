// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-03-01/exadbvmclusters"
)

func ExadbVMClusterName(i interface{}, k string) (warnings []string, errors []error) {
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

func ClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) > 11 {
		errors = append(errors, fmt.Errorf("name can be no longer than %d characters", 11))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) {
		errors = append(errors, fmt.Errorf("name must start with an alphabetic character"))
		return
	}

	re := regexp.MustCompile("_")
	if re.MatchString(v) {
		errors = append(errors, fmt.Errorf("name must not contain any underscores (_)"))
		return
	}

	return
}

func EcpuCount(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(int)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be int", k))
		return
	}

	if v < 8 || v > 200 {
		errors = append(errors, fmt.Errorf("ecpu_count must be between %v and %v", 8, 200))
		return
	}

	if v%4 != 0 {
		errors = append(errors, fmt.Errorf("ecpu_count needs to be multiple of %v", 4))
		return
	}

	return
}

func ExadbLicenseModel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if v != string(exadbvmclusters.LicenseModelBringYourOwnLicense) && v != string(exadbvmclusters.LicenseModelLicenseIncluded) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(exadbvmclusters.LicenseModelBringYourOwnLicense), string(exadbvmclusters.LicenseModelLicenseIncluded)))
		return
	}

	return
}

func ExadbSystemVersion(i interface{}, k string) (warnings []string, errors []error) {
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

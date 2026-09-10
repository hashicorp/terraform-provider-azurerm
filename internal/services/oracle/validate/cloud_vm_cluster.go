// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/cloudvmclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CloudVMClusterName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(1, 255),
		validation.StringMatch(regexp.MustCompile(`^[\p{L}_]`), "must start with a letter or underscore (_)"),
		validation.StringDoesNotMatch(regexp.MustCompile("--"), "must not contain any consecutive hyphers (--)"),
	)(i, k)
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

func SystemVersion(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`(?:19|22|23|24|25)\.[0-9]+(\.[0-9]+)*|[0-9]+(\.[0-9]+)*`), "must match one of the following patterns: (?:19|22|23|24|25).[0-9]+(.[0-9]+)* or [0-9]+(.[0-9]+)*")(i, k)
}

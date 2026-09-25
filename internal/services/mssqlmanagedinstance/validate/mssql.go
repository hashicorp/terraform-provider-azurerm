// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Your server name can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.
func ValidateMsSqlManagedInstanceServerName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`), "can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters")(i, k)
}

// Your database name can't end with '.' or ' ', can't contain '<,>,*,%,&,:,\,/,?' or control characters, and can't have more than 128 characters.
func ValidateMsSqlManagedInstanceDatabaseName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^<>*%&:\\\/?]{0,127}[^\s.<>*%&:\\\/?]$`), `can't end with '.' or ' ', can't contain '<,>,*,%,&,:,\,/,?' or control characters, and can't have more than 128 characters`)(i, k)
}

func ValidateMsSqlManagedInstanceFailoverGroupName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[0-9a-z]([-0-9a-z]{0,61}[0-9a-z])?$`), "can contain only lowercase letters, numbers, and '-', but can't start or end with '-'")(i, k)
}

// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AutonomousDatabaseName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile(`^\p{L}`), "must start with a letter"),
		validation.StringMatch(regexp.MustCompile(`^[\p{L}\p{N}]*$`), "must contain only letters and numbers"),
		validation.StringLenBetween(0, 30),
	)(i, k)
}

// lintignore:V011 // the length check is combined with password complexity rules
func AutonomousDatabasePassword(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if len(v) < 12 || len(v) > 30 {
		return []string{}, append(errors, fmt.Errorf("%v must be 12 to 30 characters", k))
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasDoubleQuote := false
	for _, r := range v {
		if r == '"' {
			hasDoubleQuote = true
		}
		if unicode.IsUpper(r) {
			hasUpper = true
		}
		if unicode.IsLower(r) {
			hasLower = true
		}
		if unicode.IsNumber(r) {
			hasNumber = true
		}
	}
	if hasDoubleQuote {
		return []string{}, append(errors, fmt.Errorf("%v must not contain the double quote (\") character", k))
	}
	if !hasUpper {
		return []string{}, append(errors, fmt.Errorf("%v must contain at least one uppercase letter", k))
	}
	if !hasLower {
		return []string{}, append(errors, fmt.Errorf("%v must contain at least one lowercase letter", k))
	}
	if !hasNumber {
		return []string{}, append(errors, fmt.Errorf("%v must contain at least one number", k))
	}
	if strings.Contains(v, "admin") {
		return []string{}, append(errors, fmt.Errorf("%v must not contain the username \"admin\"", k))
	}

	return []string{}, []error{}
}

func AdbsComputeModel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if v != string(autonomousdatabases.ComputeModelECPU) && v != string(autonomousdatabases.ComputeModelOCPU) {
		return []string{}, append(errors, fmt.Errorf("%v must be %v or %v", k, string(autonomousdatabases.ComputeModelECPU), string(autonomousdatabases.ComputeModelOCPU)))
	}

	return []string{}, []error{}
}

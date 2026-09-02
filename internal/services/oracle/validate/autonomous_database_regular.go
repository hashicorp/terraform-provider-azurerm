// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

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

func AutonomousDatabasePassword(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(12, 30),
		validation.StringDoesNotContainAny(`"`),
		validation.StringMatch(regexp.MustCompile(`\p{Lu}`), "must contain at least one uppercase letter"),
		validation.StringMatch(regexp.MustCompile(`\p{Ll}`), "must contain at least one lowercase letter"),
		validation.StringMatch(regexp.MustCompile(`\p{N}`), "must contain at least one number"),
		validation.StringDoesNotMatch(regexp.MustCompile(`admin`), "must not contain the username \"admin\""),
	)(i, k)
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

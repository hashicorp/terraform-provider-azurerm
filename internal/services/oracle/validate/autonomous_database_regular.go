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

// AutonomousDatabasePassword checks the password rules one at a time and never puts the value
// itself in an error, since validation errors end up in logs.
func AutonomousDatabasePassword(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	rules := []struct {
		ok  bool
		msg string
	}{
		{len(v) >= 12 && len(v) <= 30, "must be 12 to 30 characters"},
		{!strings.Contains(v, `"`), `must not contain the double quote (") character`},
		{strings.ContainsFunc(v, unicode.IsUpper), "must contain at least one uppercase letter"},
		{strings.ContainsFunc(v, unicode.IsLower), "must contain at least one lowercase letter"},
		{strings.ContainsFunc(v, unicode.IsNumber), "must contain at least one number"},
		{!strings.Contains(v, "admin"), `must not contain the username "admin"`},
	}
	for _, r := range rules {
		if !r.ok {
			return nil, []error{fmt.Errorf("%q %s", k, r.msg)}
		}
	}

	return nil, nil
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

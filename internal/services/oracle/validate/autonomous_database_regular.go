// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
)

func AutonomousDatabaseName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) {
		return []string{}, append(errors, fmt.Errorf("%v must start with a letter", k))
	}

	for _, r := range v {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return []string{}, append(errors, fmt.Errorf("%v must contain only letters and numbers", k))
		}
	}

	if len(v) > 30 {
		return []string{}, append(errors, fmt.Errorf("%v must be 30 characters max", k))
	}

	return []string{}, []error{}
}

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

func CustomerContactEmail(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if _, err := mail.ParseAddress(v); err != nil {
		return []string{}, append(errors, fmt.Errorf("%v must be a valid email address", k))
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

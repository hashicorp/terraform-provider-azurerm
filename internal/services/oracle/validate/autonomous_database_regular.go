package validate

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
)

func AutonomousDatabaseName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) {
		errors = append(errors, fmt.Errorf("%v must start with a letter", k))
		return
	}

	for _, r := range v {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			errors = append(errors, fmt.Errorf("%v must contain only letters and numbers", k))
			return
		}
	}

	if len(v) > 30 {
		errors = append(errors, fmt.Errorf("%v must be 30 characers max", k))
		return
	}

	return
}

func AutonomousDatabasePassword(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 12 || len(v) > 30 {
		errors = append(errors, fmt.Errorf("%v must be 12 to 30 characters", k))
		return
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
		errors = append(errors, fmt.Errorf("%v must not contain the double quote (\") character", k))
		return
	}
	if !hasUpper {
		errors = append(errors, fmt.Errorf("%v must contain at least one uppercase letter", k))
		return
	}
	if !hasLower {
		errors = append(errors, fmt.Errorf("%v must contain at least one lowercase letter", k))
		return
	}
	if !hasNumber {
		errors = append(errors, fmt.Errorf("%v must contain at least one number", k))
		return
	}
	if strings.Contains(v, "admin") {
		errors = append(errors, fmt.Errorf("%v must not contain the username \"admin\"", k))
		return
	}

	return
}

func CustomerContactEmail(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	_, err := mail.ParseAddress(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%v must be a valid email address", k))
		return
	}

	return warnings, errors
}

func AdbsComputeModel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if v != string(autonomousdatabases.ComputeModelECPU) && v != string(autonomousdatabases.ComputeModelOCPU) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(autonomousdatabases.ComputeModelECPU), string(autonomousdatabases.ComputeModelOCPU)))
		return
	}

	return
}

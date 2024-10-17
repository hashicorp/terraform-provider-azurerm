package validate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
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
		} else if unicode.IsUpper(r) {
			hasUpper = true
		} else if unicode.IsLower(r) {
			hasLower = true
		} else if unicode.IsNumber(r) {
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

func LicenseType(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if v != string(autonomousdatabases.LicenseModelLicenseIncluded) && v != string(autonomousdatabases.LicenseModelBringYourOwnLicense) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(cloudexadatainfrastructures.PatchingModeRolling), string(cloudexadatainfrastructures.PatchingModeNonRolling)))
		return
	}

	return
}

func CustomerContactEmail(v interface{}, k string) (warnings []string, errors []error) {
	value, ok := v.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%q cannot be an empty string: %q", k, v))
		return warnings, errors
	}

	vSegments := strings.Split(value, ".")
	if len(vSegments) < 2 || len(vSegments) > 34 {
		errors = append(errors, fmt.Errorf("%q must be between 2 and 34 segments", k))
		return warnings, errors
	}

	for _, segment := range vSegments {
		if segment == "" {
			errors = append(errors, fmt.Errorf("%q cannot contain consecutive period", k))
			return warnings, errors
		}

		if len(segment) > 63 {
			errors = append(errors, fmt.Errorf("the each segment of the `email` must contain between 1 and 63 characters"))
			return warnings, errors
		}
	}

	if !regexp.MustCompile(`^[a-zA-Z\d._-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q only contains letters, numbers, underscores, dashes and periods", k))
		return warnings, errors
	}

	return warnings, errors
}

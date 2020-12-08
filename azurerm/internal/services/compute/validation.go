package compute

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ValidateVmName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("Expected %q to be a string but it wasn't!", k))
		return
	}

	// The value must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	const maxLength = 80
	// VM name can be 1-80 characters in length
	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, dots, dashes and underscores", k))
	}

	if matched := regexp.MustCompile(`^[a-zA-Z0-9]`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must begin with an alphanumeric character", k))
	}

	if matched := regexp.MustCompile(`\w$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q must end with an alphanumeric character or underscore", k))
	}

	// Portal: Virtual machine name cannot contain only numbers.
	if matched := regexp.MustCompile(`^\d+$`).Match([]byte(v)); matched {
		errors = append(errors, fmt.Errorf("%q cannot contain only numbers", k))
	}

	return warnings, errors
}

func ValidateLinuxComputerNameFull(i interface{}, k string) (warnings []string, errors []error) {
	// Linux host name cannot exceed 64 characters in length
	return ValidateLinuxComputerName(i, k, 64, false)
}

func ValidateLinuxComputerNamePrefix(i interface{}, k string) (warnings []string, errors []error) {
	// Linux host name prefix cannot exceed 58 characters in length
	return ValidateLinuxComputerName(i, k, 58, true)
}

func ValidateOrchestratedVMSSName(i interface{}, k string) (warnings []string, errors []error) {
	return ValidateVmName(i, k)
}

func ValidateLinuxComputerName(i interface{}, k string, maxLength int, allowDashSuffix bool) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't!", k))
		return
	}

	// The value must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
	}

	if strings.HasPrefix(v, "_") {
		errors = append(errors, fmt.Errorf("%q cannot begin with an underscore", k))
	}

	if strings.HasSuffix(v, ".") {
		errors = append(errors, fmt.Errorf("%q cannot end with a period", k))
	}

	if !allowDashSuffix && strings.HasSuffix(v, "-") {
		errors = append(errors, fmt.Errorf("%q cannot end with a dash", k))
	}

	// Linux host name cannot contain the following characters
	specialCharacters := `\/"[]:|<>+=;,?*@&~!#$%^()_{}'`
	if strings.ContainsAny(v, specialCharacters) {
		errors = append(errors, fmt.Errorf("%q cannot contain the special characters: `%s`", k, specialCharacters))
	}

	return warnings, errors
}

func ValidateWindowsComputerNameFull(i interface{}, k string) (warnings []string, errors []error) {
	// Windows computer name cannot be more than 15 characters long
	return ValidateWindowsComputerName(i, k, 15)
}

func ValidateWindowsComputerNamePrefix(i interface{}, k string) (warnings []string, errors []error) {
	// Windows computer name prefix cannot be more than 9 characters long
	return ValidateWindowsComputerName(i, k, 9)
}

func ValidateWindowsComputerName(i interface{}, k string, maxLength int) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't!", k))
		return
	}

	// The value must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	if len(v) > maxLength {
		errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
	}

	if strings.HasSuffix(v, "-") {
		errors = append(errors, fmt.Errorf("%q cannot end with dash", k))
	}

	// A windows computer name can only contain alphanumeric characters and hyphens
	if matched := regexp.MustCompile(`^[a-zA-Z0-9-]+$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters and dashes", k))
	}

	// Windows computer name cannot contain only numbers
	if matched := regexp.MustCompile(`^\d+$`).Match([]byte(v)); matched {
		errors = append(errors, fmt.Errorf("%q cannot contain only numbers", k))
	}

	return warnings, errors
}

func validateDiskEncryptionSetName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	// Swagger says: Supported characters for the name are a-z, A-Z, 0-9 and _. The maximum name length is 80 characters.
	// Confirmed with the service team, they gave me a regex: ^[^_\W][\w-._]{0,79}(?<![-.])$
	// This means the name can contain a-z, A-Z, 0-9, underscore, dot or hyphen, and must not starts with a underscore or any other non-word characters (underscore is considered as word character)
	// additionally, the name cannot end with hyphen or dot.
	// Golang regex does not support "negative look ahead" (aka `?<!`), therefore I transformed this regex to the following regular expression.
	if matched := regexp.MustCompile(`^[^_\W][\w-._]{0,78}\w$`).Match([]byte(v)); !matched {
		errors = append(errors, fmt.Errorf("%s must be between 1 - 80 characters long, and contains only a-z, A-Z, 0-9 and _", k))
	}
	return
}

func validateDiskSizeGB(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 32767 {
		errors = append(errors, fmt.Errorf(
			"The `disk_size_gb` can only be between 0 and 32767"))
	}
	return warnings, errors
}

func validateManagedDiskSizeGB(v interface{}, _ string) (warnings []string, errors []error) {
	value := v.(int)
	if value < 0 || value > 65536 {
		errors = append(errors, fmt.Errorf(
			"The `disk_size_gb` can only be between 0 and 65536"))
	}
	return warnings, errors
}

func validateDedicatedHostGroupName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^_\W][\w-.]{0,78}[\w]$`), "")
}

func validateDedicatedHostName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validateDedicatedHostGroupName()
}

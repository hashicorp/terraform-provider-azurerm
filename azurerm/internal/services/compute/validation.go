package compute

import (
	"fmt"
	"regexp"
	"strings"
)

func ValidateLinuxName(i interface{}, k string) (warnings []string, errors []error) {
	return validateName(64)(i, k)
}

func ValidateWindowsName(i interface{}, k string) (warnings []string, errors []error) {
	return validateName(16)(i, k)
}

func ValidateScaleSetResourceID(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	id, err := ParseVirtualMachineScaleSetID(v)
	if err != nil {
		es = append(es, fmt.Errorf("Error parsing %q as a VM Scale Set Resource ID: %s", v, err))
		return
	}

	if id.Name == "" {
		es = append(es, fmt.Errorf("Error parsing %q as a VM Scale Set Resource ID: `virtualMachineScaleSets` segment was empty", v))
		return
	}

	return
}

func validateName(maxLength int) func(i interface{}, k string) (warnings []string, errors []error) {
	return func(i interface{}, k string) (warnings []string, errors []error) {
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

		// The value must be between 1 and 64 (Linux) or 16 (Windows) characters long.
		if len(v) >= maxLength {
			errors = append(errors, fmt.Errorf("%q can be at most %d characters, got %d", k, maxLength, len(v)))
		}

		if strings.HasPrefix(v, "_") {
			errors = append(errors, fmt.Errorf("%q cannot begin with an underscore", k))
		}

		if strings.HasSuffix(v, ".") || strings.HasSuffix(v, "-") {
			errors = append(errors, fmt.Errorf("%q cannot end with an period or dash", k))
		}

		// Azure resource names cannot contain special characters \/""[]:|<>+=;,?*@& or begin with '_' or end with '.' or '-'
		specialCharacters := `\/""[]:|<>+=;,?*@&`
		if strings.ContainsAny(v, specialCharacters) {
			errors = append(errors, fmt.Errorf("%q cannot contain the special characters: `%s`", k, specialCharacters))
		}

		// The value can only contain alphanumeric characters and cannot start with a number.
		if matched := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`).Match([]byte(v)); !matched {
			errors = append(errors, fmt.Errorf("%q may only contain alphanumeric characters, dashes and underscores", k))
		}

		return
	}
}

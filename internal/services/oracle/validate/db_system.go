package validate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/dbsystems"
)

func DbSystemName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) < 1 || len(v) > 255 {
		errors = append(errors, fmt.Errorf("name must be %d to %d characters", 1, 255))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		errors = append(errors, fmt.Errorf("name must start with a letter or underscore (_)"))
		return
	}

	re := regexp.MustCompile("--")
	if re.MatchString(v) {
		errors = append(errors, fmt.Errorf("name must not contain any consecutive hyphens (--)"))
		return
	}

	return
}

func DbSystemLicenseModel(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if v != string(dbsystems.LicenseModelBringYourOwnLicense) && v != string(dbsystems.LicenseModelLicenseIncluded) {
		errors = append(errors, fmt.Errorf("%v must be %v or %v", k,
			string(dbsystems.LicenseModelBringYourOwnLicense), string(dbsystems.LicenseModelLicenseIncluded)))
		return
	}

	return
}

func DbSystemPassword(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if v == "" {
		errors = append(errors, fmt.Errorf("%v must not be an empty string", k))
		return
	}

	if len(v) < 9 || len(v) > 255 {
		return []string{}, append(errors, fmt.Errorf("%v must be 9 to 255 characters", k))
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	hasInvalid := false
	upperCount := 0
	lowerCount := 0
	numberCount := 0
	specialCount := 0

	// Allowed characters are letters, numbers, and _, #, or -
	var allowedCharsPattern = regexp.MustCompile(`^[A-Za-z0-9_#-]+$`)
	if !allowedCharsPattern.MatchString(v) {
		hasInvalid = true
	}

	for _, r := range v {
		if unicode.IsUpper(r) {
			upperCount++
			if upperCount >= 2 {
				hasUpper = true
			}
		}
		if unicode.IsLower(r) {
			lowerCount++
			if lowerCount >= 2 {
				hasLower = true
			}
		}
		if unicode.IsNumber(r) {
			numberCount++
			if numberCount >= 2 {
				hasNumber = true
			}
		}
		if strings.ContainsRune("_#-", r) {
			specialCount++
			if specialCount >= 2 {
				hasSpecial = true
			}
		}
	}

	if hasInvalid {
		return []string{}, append(errors, fmt.Errorf("%v must contain only the following special characters: _, #, or -", k))
	}
	if !hasUpper {
		return []string{}, append(errors, fmt.Errorf("%v must contain at least two uppercase letters", k))
	}
	if !hasLower {
		return []string{}, append(errors, fmt.Errorf("%v must contain at least two lowercase letters", k))
	}
	if !hasNumber {
		return []string{}, append(errors, fmt.Errorf("%v must contain at least two numbers", k))
	}
	if !hasSpecial {
		return []string{}, append(errors, fmt.Errorf("%v must contain at least two special characters. The special characters must be _, #, or -", k))
	}

	return []string{}, []error{}
}

func PluggableDatabaseName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) > 30 {
		errors = append(errors, fmt.Errorf("name must be no more than %d characters", 30))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) {
		errors = append(errors, fmt.Errorf("name must start with a letter"))
		return
	}

	re := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	if !re.MatchString(v) {
		errors = append(errors, fmt.Errorf("name must not contain any special characters"))
		return
	}

	return
}

func ClusterName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	if len(v) > 11 {
		errors = append(errors, fmt.Errorf("name must be no more than %d characters", 11))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) && firstChar != '-' {
		errors = append(errors, fmt.Errorf("name must start with a letter or hyphen (-)"))
		return
	}

	re := regexp.MustCompile("_")
	if re.MatchString(v) {
		errors = append(errors, fmt.Errorf("name must not contain any underscores (_)"))
		return
	}

	return
}

package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

func DatabaseSystemName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	minLength := 1
	maxLength := 255
	if len(v) < minLength || len(v) > maxLength {
		errors = append(errors, fmt.Errorf("`name` must be %d to %d characters", minLength, maxLength))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		errors = append(errors, fmt.Errorf("`name` must start with a letter or underscore (_)"))
		return
	}

	re := regexp.MustCompile("--")
	if re.MatchString(v) {
		errors = append(errors, fmt.Errorf("`name` must not contain any consecutive hyphens (--)"))
		return
	}

	return
}

func DatabaseSystemPassword(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return []string{}, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	if v == "" {
		errors = append(errors, fmt.Errorf("%v must not be an empty string", k))
		return
	}

	minLength := 9
	maxLength := 255
	if len(v) < minLength || len(v) > maxLength {
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

func PluggableDatabaseName(i interface{}, k string) (warnings []string, errorsList []error) {
	v, ok := i.(string)
	if !ok {
		errorsList = append(errorsList, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	maxLength := 30
	if len(v) > maxLength {
		errorsList = append(errorsList, fmt.Errorf("`name` must be no more than %d characters", maxLength))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) {
		errorsList = append(errorsList, errors.New("`name` must start with a letter"))
		return
	}

	re := regexp.MustCompile(`^[A-Za-z0-9]+$`)
	if !re.MatchString(v) {
		errorsList = append(errorsList, errors.New("`name` must not contain any special characters"))
		return
	}

	return
}

func ClusterName(i interface{}, k string) (warnings []string, errorsList []error) {
	v, ok := i.(string)
	if !ok {
		errorsList = append(errorsList, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	maxLength := 11
	if len(v) > maxLength {
		errorsList = append(errorsList, fmt.Errorf("`name` must be no more than %d characters", maxLength))
		return
	}

	firstChar, _ := utf8.DecodeRuneInString(v)
	if !unicode.IsLetter(firstChar) && firstChar != '-' {
		errorsList = append(errorsList, errors.New("`name` must start with a letter or hyphen (-)"))
		return
	}

	re := regexp.MustCompile("_")
	if re.MatchString(v) {
		errorsList = append(errorsList, errors.New("`name` must not contain any underscores (_)"))
		return
	}

	return
}

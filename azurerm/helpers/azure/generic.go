package azure

import (
	"fmt"
	"regexp"
)

const hyphenUnderscoreParenthesesPeriod string = "-_\\(\\)\\."
const hyphenUnderscorePeriod string = "-_\\."
const hyphenUnderscore string = "-_"
const hyphenPeriod string = "-\\."
const hyphen string = "-"
const none string = ""
const alphanumericLower string = "a-z0-9"
const alphanumericUpper string = "A-Z0-9"
const alphanumericBoth string = "a-zA-Z0-9"

// Generic method to validate Azure Resource names which enforce the Azure standards for naming conventions...
// 1. First two characters must be a number or a letter
// 2. Last characters must be a number or a letter
// 3. No consecutive hyphens, underscores, parentheses, or periods
// 4. Min and Max length for values vary depending on the resource
// 5. Value can not start or end with any special character
// NOTE: There is an absolute minimum length for all values of 3, it can not be lower than 3 because of the
//       first two chars and last char must be a number or letter rule, it can however be larger than 3 if desired.
func ValidateNameGeneric(i interface{}, attributeName string, pattern string, specialChars string, errorPrefix string, minLength int, maxLength int) (_ []string, errors []error) {
	value, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", attributeName))
		return nil, errors
	}

	if minLength < 3 {
		minLength = 3
	}

	if maxLength < minLength {
		maxLength = minLength + 1
	}

	if len(value) < minLength || len(value) > maxLength {
		errors = append(errors, fmt.Errorf("%s %q must be %d - %d characters in length", errorPrefix, attributeName, minLength, maxLength))
	}

	regEx := fmt.Sprintf("^[%s]{2}.*[%s]{1}$", pattern, pattern)
	r := regexp.MustCompile(regEx)
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s %q first, second, and last characters must be a %s", errorPrefix, attributeName, getOrTxt(pattern)))
	}

	regEx = fmt.Sprintf("^[%s%s]*$", pattern, specialChars)
	r = regexp.MustCompile(regEx)
	if !r.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s %q can only contain %s", errorPrefix, attributeName, getAndTxt(pattern, specialChars)))
	}

	r = regexp.MustCompile("(--|__|\\.\\.|\\(\\(|\\)\\))")
	if r.MatchString(value) {
		errors = append(errors, fmt.Errorf("%s %q must not contain any consecutive hyphens, underscores, parentheses, or periods", errorPrefix, attributeName))
	}

	return nil, errors
}

func getAndTxt(pattern string, specialChars string) (msg string) {

	if specialChars == none {
		switch pattern {
		case alphanumericBoth:
			msg = "letters and numbers"
		case alphanumericLower:
			msg = "lowercase letters and numbers"
		case alphanumericUpper:
			msg = "uppercase letters and numbers"
		default:
			msg = ""
		}
	} else {
		tmpMsg := ""

		switch pattern {
		case alphanumericBoth:
			tmpMsg = "letters, numbers"
		case alphanumericLower:
			tmpMsg = "lowercase letters, numbers"
		case alphanumericUpper:
			tmpMsg = "uppercase letters, numbers"
		}

		switch specialChars {
		case hyphenUnderscoreParenthesesPeriod:
			msg = fmt.Sprintf("%s, hyphens, underscores, parentheses, and periods", tmpMsg)
		case hyphenUnderscorePeriod:
			msg = fmt.Sprintf("%s, hyphens, underscores, and periods", tmpMsg)
		case hyphenUnderscore:
			msg = fmt.Sprintf("%s, hyphens and underscores", tmpMsg)
		case hyphenPeriod:
			msg = fmt.Sprintf("%s, hyphens and periods", tmpMsg)
		case hyphen:
			msg = fmt.Sprintf("%s and hyphens", tmpMsg)
		default:
			msg = tmpMsg
		}
	}

	return msg
}

func getOrTxt(expression string) (msg string) {

	switch expression {
	case alphanumericBoth:
		msg = "letter or number"
	case alphanumericLower:
		msg = "lowercase letter or number"
	case alphanumericUpper:
		msg = "uppercase letter or number"
	default:
		msg = ""
	}

	return msg
}

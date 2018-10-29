package validate

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func DevSpaceName() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		// Length should be between 3 and 31.
		if s, es = validation.StringLenBetween(3, 31)(i, k); len(es) > 0 {
			return s, es
		}

		// Naming rule.
		regexStr := "^[a-zA-Z0-9](-?[a-zA-Z0-9])*$"
		errMsg := "DevSpace name can only include alphanumeric characters, hyphens."
		if s, es = validation.StringMatch(regexp.MustCompile(regexStr), errMsg)(i, k); len(es) > 0 {
			return s, es
		}

		return s, es
	}
}

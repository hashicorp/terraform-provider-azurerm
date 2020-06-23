package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func KustoTableName(i interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringIsNotEmpty,
		validation.StringLenBetween(1, 1024),
		validation.StringMatch(regexp.MustCompile("^[0-9]*[a-zA-Z][0-9a-zA-Z]*"), "must container only letters and numbers. if first char is a number it has to be followed by a letter."),
	)(i, k)
}

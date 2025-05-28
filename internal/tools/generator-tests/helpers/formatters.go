package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

// TerraformResourceName generates a Terraform-compliant resource name by combining the provider and resource name.
func TerraformResourceName(provider, resourceName string) string {
	fmtStr := "%s_%s"
	return fmt.Sprintf(fmtStr, strings.ToLower(provider), strcase.ToSnake(resourceName))
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

// PrefixedDescriptionString returns a string prefixed with "a" or "an" based on whether the input starts with a vowel.
func PrefixedDescriptionString(input string) string {
	prefix := "a"
	first := input[0:1]
	vowel, _ := regexp.Match(first, []byte(`aeiouAEIOU`))

	if vowel {
		prefix = "an"
	}
	return fmt.Sprintf("%s %s", prefix, strings.Title(strcase.ToDelimited(input, ' ')))
}

// ToDelimTitle converts the input string to a title-cased string with words delimited by spaces.
func ToDelimTitle(input string) string {
	return strings.Title(strcase.ToDelimited(input, ' '))
}

// PrefixedLabelString determines whether a given label should use "A" or "An" as its prefix based on its starting letter.
func PrefixedLabelString(input string) string {
	prefix := "A"
	vowel, _ := regexp.Match(input[0:1], []byte(`aeiouAEIOU`))

	if vowel {
		prefix = "An"
	}

	return fmt.Sprintf("%s `%s`", prefix, input)
}

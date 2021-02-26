package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// ScheduleName validates Automation Account Schedule names
func ScheduleName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[^<>*%&:\\?.+/]{0,127}[^<>*%&:\\?.+/\s]$`),
		`The name length must be from 1 to 128 characters. The name cannot contain special characters < > * % & : \ ? . + / and cannot end with a whitespace character.`,
	)
}

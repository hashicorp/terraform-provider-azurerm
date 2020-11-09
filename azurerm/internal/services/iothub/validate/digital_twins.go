package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DigitaltwinsName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9-]{1,61}[A-Za-z0-9]$`),
		`Name contains invalid characters or exceeds allowed length.`,
	)
}

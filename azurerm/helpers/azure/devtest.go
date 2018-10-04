package azure

import (
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func ValidateDevTestLabName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		"Lab Name can only include alphanumeric characters, underscores, hyphens.")
}

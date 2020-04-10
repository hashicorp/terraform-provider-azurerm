package maps

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ValidateName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[A-Za-z0-9]{1}[A-Za-z0-9._-]{1,}$`),
		"First character must be alphanumeric. Subsequent character(s) must be any combination of alphanumeric, underscore (_), period (.), or hyphen (-).")
}

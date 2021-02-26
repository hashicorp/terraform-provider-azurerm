package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func AuthorizationRuleName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9][-._a-zA-Z0-9]{0,48}([a-zA-Z0-9])?$"),
		"The name can contain only letters, numbers, periods, hyphens and underscores. The name must start and end with a letter or number and be less the 50 characters long.",
	)
}

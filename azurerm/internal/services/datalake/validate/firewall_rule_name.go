package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func FirewallRuleName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`\A([-_a-zA-Z0-9]{3,50})\z`),
		"Name can only consist of letters, numbers, underscores and hyphens and must be between 3 and 50 characters long",
	)
}

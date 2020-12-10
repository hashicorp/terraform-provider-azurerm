package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func IntegrationAccountName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[\w-().]{1,80}$`), `Integration name can contain only letters, numbers, '_','-', '(', ')' or '.'`,
	)
}

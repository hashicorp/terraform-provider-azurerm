package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func AccountName() pluginsdk.SchemaValidateFunc {
	// store and analytic account names are the same
	return validation.StringMatch(
		regexp.MustCompile(`\A([a-z0-9]{3,24})\z`),
		"Name can only consist of lowercase letters and numbers and must be between 3 and 24 characters long",
	)
}

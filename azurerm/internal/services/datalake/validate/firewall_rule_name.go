package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func FirewallRuleName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`\A([-_a-zA-Z0-9]{3,50})\z`),
		"Name can only consist of letters, numbers, underscores and hyphens and must be between 3 and 50 characters long",
	)
}

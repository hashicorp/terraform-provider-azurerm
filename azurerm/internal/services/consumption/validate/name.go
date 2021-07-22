package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func ConsumptionBudgetName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[-_a-zA-Z0-9]{1,63}$"),
		"The consumption budget name can contain only letters, numbers, underscores, and hyphens. The consumption budget name be between 6 and 63 characters long.",
	)
}

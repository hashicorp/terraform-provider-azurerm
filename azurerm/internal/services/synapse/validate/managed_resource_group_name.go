package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func ManagedResourceGroupName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[-\w\._\(\)]{0,89}[-\w_\(\)]$`),
		"The resource group name must be no longer than 90 characters long, and must be alphanumeric characters and '-', '_', '(', ')' and'.'. Note that the name cannot end with '.'")
}

package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func IntegrationAccountName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[\w-().]{1,80}$`), `Integration name can contain only letters, numbers, '_','-', '(', ')' or '.'`,
	)
}

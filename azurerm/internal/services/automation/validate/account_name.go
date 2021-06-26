package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

// AutomationAccount validates Automation Account names
func AutomationAccount() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[0-9a-zA-Z][-0-9a-zA-Z]{4,48}[0-9a-zA-Z]$`),
		`The account name must start with a letter or number.  The account name can contain letters, numbers, and dashes. The final character must be a letter or a number. The account name length must be from 6 to 50 characters.`,
	)
}

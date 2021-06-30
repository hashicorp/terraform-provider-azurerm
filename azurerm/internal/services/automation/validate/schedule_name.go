package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

// ScheduleName validates Automation Account Schedule names
func ScheduleName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[^<>*%&:\\?.+/]{0,127}[^<>*%&:\\?.+/\s]$`),
		`The name length must be from 1 to 128 characters. The name cannot contain special characters < > * % & : \ ? . + / and cannot end with a whitespace character.`,
	)
}

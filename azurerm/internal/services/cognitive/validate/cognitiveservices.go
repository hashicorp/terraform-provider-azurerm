package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func CognitiveServicesAccountName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^([a-zA-Z0-9]{1}[a-zA-Z0-9_.-]{1,})$"),
		"The Cognitive Services Account Name can only start with an alphanumeric character, and must only contain alphanumeric characters, periods, dashes or underscores.",
	)
}

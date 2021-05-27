package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func QueueName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z0-9][\w-./~]{0,258}([a-zA-Z0-9])?$`),
		"The topic name can contain only letters, numbers, periods, hyphens, tildas, forward slashes and underscores. The namespace must start and end with a letter or number and be less then 260 characters long.",
	)
}

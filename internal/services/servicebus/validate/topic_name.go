package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func TopicName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9]$|^[A-Za-z0-9][\\w-\\.\\/\\~]*[A-Za-z0-9]$"),
		"The topic name can contain only letters, numbers, periods, hyphens, tildas, forward slashes and underscores. The namespace must start with a letter or number, and it must end with a letter or number and be less then 260 characters long.",
	)
}

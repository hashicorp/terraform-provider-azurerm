package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AccountConnectionName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_-]{2,32}$"),
		"`name` must be between 3 and 33 characters long, start with an alphanumeric character, and contain only alphanumeric characters, dashes(-) or underscores(_).",
	)
}

package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"regexp"
)

func ValidateWebpubsubName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{1,61}[a-zA-Z0-9]$"),
		"The web pubsub name can contain only letters, numbers and hyphens. The first character must be a letter. The last character must be a letter or number. The value must be between 3 and 63 characters long.",
	)
}

func ValidateWebPubsbHubName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z][A-Za-z0-9_`,.\\[\\]]{0,127}$"),
		"The web pubsub hub name can contain only letters, numbers and special characters including `,` , `_` , `.` , `[` ,` . The first character must be a letter. The value must be between 0 and 127 characters long.",
	)
}

package web

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ValidateAppServiceSourceControlTokenName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"BitBucket",
		"Dropbox",
		"GitHub",
		"OneDrive",
	}, false)
}

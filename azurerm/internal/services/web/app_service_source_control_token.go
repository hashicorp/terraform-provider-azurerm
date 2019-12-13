package web

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ValidateAppServiceSourceControlTokenName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"BitBucket",
		"Dropbox",
		"GitHub",
		"OneDrive",
	}, false)
}

package web

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ValidateAppServiceSourceControlTokenName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"Bitbucket",
		"Dropbox",
		"GitHub",
		"OneDrive",
	}, false)
}

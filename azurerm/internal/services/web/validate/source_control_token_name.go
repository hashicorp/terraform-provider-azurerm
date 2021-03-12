package validate

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func SourceControlTokenName() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		"BitBucket",
		"Dropbox",
		"GitHub",
		"OneDrive",
	}, false)
}

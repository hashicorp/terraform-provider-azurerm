package validate

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// deprecated: please use validation.StringIsBase64 instead
func Base64String() schema.SchemaValidateFunc {
	return validation.StringIsBase64
}

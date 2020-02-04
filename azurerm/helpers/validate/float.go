package validate

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// deprecated: please use validation.FloatAtLeast instead
func FloatAtLeast(min float64) schema.SchemaValidateFunc {
	return validation.FloatAtLeast(min)
}

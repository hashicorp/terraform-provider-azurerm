package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ManagedResourceGroupName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[-\\w\\._\\(\\)]{0,89}[-\\w_\\(\\)]$"),
		"The resource group name must be no longer than 90 characters long, and must be alphanumeric characters and '-', '_', '(', ')' and'.'. Note that the name cannot end with '.'")
}

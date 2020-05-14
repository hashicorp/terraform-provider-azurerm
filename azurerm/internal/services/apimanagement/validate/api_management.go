package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func ApimSkuName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^Consumption_0$|^Basic_(1|2)$|^Developer_1$|^Premium_([1-9]|10)$|^Standard_[1-4]$`),
		`This is not a valid Api Management sku name.`,
	)
}

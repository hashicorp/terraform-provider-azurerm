package validate

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
)

func FlexibleServerName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?$`),

		`This is not a valid collation.`,
	)
}

func FlexibleServerSkuName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^(Standard_E(2|4|8|16|32|48|64)s_v3)|(Standard_B(1m|2)s)|(Standard_D(2|4|8|16|32|48|64)s_v3)$`),

		`This is not a valid collation.`,
	)
}

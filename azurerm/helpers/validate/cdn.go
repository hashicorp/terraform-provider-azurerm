package validate

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"regexp"
)

func CdnEndpointDeliveryPolicyRuleName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*$"),
		"The Delivery Policy Rule Name must start with a letter any may only contain letters and numbers.",
	)
}

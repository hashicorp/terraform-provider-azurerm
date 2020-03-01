package cdn

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func RuleConditionRequestScheme() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Equal",
				ValidateFunc: validation.StringInSlice([]string{
					"Equal",
				}, false),
			},

			"negate_condition": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				ValidateFunc: validation.StringInSlice([]string{
					"HTTP",
					"HTTPS",
				}, false),
			},
		},
	}
}

package cdn

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func EndpointDeliveryRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.CdnEndpointDeliveryPolicyRuleName(),
			},

			"order": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"request_scheme_condition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     RuleConditionRequestScheme(),
			},

			"url_redirect_action": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     RuleActionUrlRedirect(),
			},
		},
	}
}

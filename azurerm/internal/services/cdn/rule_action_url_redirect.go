package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func RuleActionUrlRedirect() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"redirect_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.Found),
					string(cdn.Moved),
					string(cdn.PermanentRedirect),
					string(cdn.TemporaryRedirect),
				}, false),
			},

			"destination_protocol": {
				Type:     schema.TypeString,
				Required: false,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.MatchRequest),
					string(cdn.HTTP),
					string(cdn.HTTPS),
				}, false),
			},

			"hostname": {
				Type:     schema.TypeString,
				Required: false,
			},

			"path": {
				Type:         schema.TypeString,
				Required:     false,
				ValidateFunc: validate.RuleActionUrlRedirectPath(),
			},

			"query_string": {
				Type:         schema.TypeString,
				Required:     false,
				ValidateFunc: validate.RuleActionUrlRedirectPath(),
			},

			"fragment": {
				Type:         schema.TypeString,
				Required:     false,
				ValidateFunc: validate.RuleActionUrlRedirectFragment(),
			},
		},
	}
}

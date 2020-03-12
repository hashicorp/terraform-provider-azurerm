package deliveryruleactions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func URLRewrite() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_pattern": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.RuleActionUrlRewriteSourcePattern(),
			},

			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.RuleActionUrlRewriteDestination(),
			},

			"preserve_unmatched_path": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func ExpandArmCdnEndpointActionURLRewrite(ura map[string]interface{}) *cdn.URLRewriteAction {
	return &cdn.URLRewriteAction{
		Name: cdn.NameURLRewrite,
		Parameters: &cdn.URLRewriteActionParameters{
			OdataType:             utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlRewriteActionParameters"),
			SourcePattern:         utils.String(ura["source_pattern"].(string)),
			Destination:           utils.String(ura["destination"].(string)),
			PreserveUnmatchedPath: utils.Bool(ura["preserve_unmatched_path"].(bool)),
		},
	}
}

func FlattenArmCdnEndpointActionURLRewrite(ura *cdn.URLRewriteAction) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := ura.Parameters; params != nil {
		res["source_pattern"] = *params.SourcePattern
		res["destination"] = *params.Destination

		if params.PreserveUnmatchedPath != nil {
			res["preserve_unmatched_path"] = *params.PreserveUnmatchedPath
		} else {
			res["preserve_unmatched_path"] = true
		}
	}

	return res
}

package delivery_rule_actions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func URLRedirect() *schema.Resource {
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

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(cdn.MatchRequest),
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.MatchRequest),
					string(cdn.HTTP),
					string(cdn.HTTPS),
				}, false),
			},

			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionUrlRedirectPath(),
			},

			"query_string": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionUrlRedirectPath(),
			},

			"fragment": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionUrlRedirectFragment(),
			},
		},
	}
}

func ExpandArmCdnEndpointActionUrlRedirect(ura map[string]interface{}) *cdn.URLRedirectAction {
	urlRedirectAction := cdn.URLRedirectAction{
		Name: cdn.NameURLRedirect,
	}

	params := cdn.URLRedirectActionParameters{
		OdataType:    utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlRedirectActionParameters"),
		RedirectType: cdn.RedirectType(ura["redirect_type"].(string)),
	}

	if destProt := ura["protocol"]; destProt.(string) != "" {
		params.DestinationProtocol = cdn.DestinationProtocol(destProt.(string))
	}

	if hostname := ura["hostname"]; hostname.(string) != "" {
		params.CustomHostname = utils.String(hostname.(string))
	}

	if path := ura["path"]; path.(string) != "" {
		params.CustomPath = utils.String(path.(string))
	}

	if queryString := ura["query_string"]; queryString.(string) != "" {
		params.CustomQueryString = utils.String(queryString.(string))
	}

	if fragment := ura["fragment"]; fragment.(string) != "" {
		params.CustomFragment = utils.String(fragment.(string))
	}

	urlRedirectAction.Parameters = &params

	return &urlRedirectAction
}

func FlattenArmCdnEndpointActionUrlRedirect(ura *cdn.URLRedirectAction) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := ura.Parameters; params != nil {
		res["redirect_type"] = string(params.RedirectType)

		res["protocol"] = string(params.DestinationProtocol)

		if params.CustomHostname != nil {
			res["hostname"] = *params.CustomHostname
		}

		if params.CustomPath != nil {
			res["path"] = *params.CustomPath
		}

		if params.CustomQueryString != nil {
			res["query_string"] = *params.CustomQueryString
		}

		if params.CustomFragment != nil {
			res["fragment"] = *params.CustomFragment
		}
	}

	return res
}

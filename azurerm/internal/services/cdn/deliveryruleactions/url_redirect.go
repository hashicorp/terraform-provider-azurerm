package deliveryruleactions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func URLRedirect() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"redirect_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.Found),
					string(cdn.Moved),
					string(cdn.PermanentRedirect),
					string(cdn.TemporaryRedirect),
				}, false),
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(cdn.MatchRequest),
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.MatchRequest),
					string(cdn.HTTP),
					string(cdn.HTTPS),
				}, false),
			},

			"hostname": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"path": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionUrlRedirectPath(),
			},

			"query_string": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionUrlRedirectQueryString(),
			},

			"fragment": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionUrlRedirectFragment(),
			},
		},
	}
}

func ExpandArmCdnEndpointActionUrlRedirect(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		params := cdn.URLRedirectActionParameters{
			OdataType:    utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlRedirectActionParameters"),
			RedirectType: cdn.RedirectType(item["redirect_type"].(string)),
		}

		if destProt := item["protocol"]; destProt.(string) != "" {
			params.DestinationProtocol = cdn.DestinationProtocol(destProt.(string))
		}

		if hostname := item["hostname"]; hostname.(string) != "" {
			params.CustomHostname = utils.String(hostname.(string))
		}

		if path := item["path"]; path.(string) != "" {
			params.CustomPath = utils.String(path.(string))
		}

		if queryString := item["query_string"]; queryString.(string) != "" {
			params.CustomQueryString = utils.String(queryString.(string))
		}

		if fragment := item["fragment"]; fragment.(string) != "" {
			params.CustomFragment = utils.String(fragment.(string))
		}

		output = append(output, cdn.URLRedirectAction{
			Name:       cdn.NameURLRedirect,
			Parameters: &params,
		})
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionUrlRedirect(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsURLRedirectAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url redirect action!")
	}

	customHostname := ""
	customPath := ""
	fragment := ""
	queryString := ""
	protocol := ""
	redirectType := ""

	if params := action.Parameters; params != nil {
		redirectType = string(params.RedirectType)
		protocol = string(params.DestinationProtocol)

		if params.CustomHostname != nil {
			customHostname = *params.CustomHostname
		}

		if params.CustomPath != nil {
			customPath = *params.CustomPath
		}

		if params.CustomQueryString != nil {
			queryString = *params.CustomQueryString
		}

		if params.CustomFragment != nil {
			fragment = *params.CustomFragment
		}
	}

	return &map[string]interface{}{
		"fragment":      fragment,
		"hostname":      customHostname,
		"query_string":  queryString,
		"path":          customPath,
		"protocol":      protocol,
		"redirect_type": redirectType,
	}, nil
}

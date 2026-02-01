// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package deliveryruleactions

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func URLRedirect() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"redirect_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.RedirectTypeFound),
					string(cdn.RedirectTypeMoved),
					string(cdn.RedirectTypePermanentRedirect),
					string(cdn.RedirectTypeTemporaryRedirect),
				}, false),
			},

			"protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(cdn.DestinationProtocolMatchRequest),
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.DestinationProtocolMatchRequest),
					string(cdn.DestinationProtocolHTTP),
					string(cdn.DestinationProtocolHTTPS),
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
			OdataType:    pointer.To("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlRedirectActionParameters"),
			RedirectType: cdn.RedirectType(item["redirect_type"].(string)),
		}

		if destProt := item["protocol"]; destProt.(string) != "" {
			params.DestinationProtocol = cdn.DestinationProtocol(destProt.(string))
		}

		if hostname := item["hostname"]; hostname.(string) != "" {
			params.CustomHostname = pointer.To(hostname.(string))
		}

		if path := item["path"]; path.(string) != "" {
			params.CustomPath = pointer.To(path.(string))
		}

		if queryString := item["query_string"]; queryString.(string) != "" {
			params.CustomQueryString = pointer.To(queryString.(string))
		}

		if fragment := item["fragment"]; fragment.(string) != "" {
			params.CustomFragment = pointer.To(fragment.(string))
		}

		output = append(output, cdn.URLRedirectAction{
			Name:       cdn.NameBasicDeliveryRuleActionNameURLRedirect,
			Parameters: &params,
		})
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionUrlRedirect(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsURLRedirectAction()
	if !ok {
		return nil, errors.New("expected a delivery rule url redirect action")
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

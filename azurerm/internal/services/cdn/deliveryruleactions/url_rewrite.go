package deliveryruleactions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func URLRewrite() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"source_pattern": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.RuleActionUrlRewriteSourcePattern(),
			},

			"destination": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.RuleActionUrlRewriteDestination(),
			},

			"preserve_unmatched_path": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func ExpandArmCdnEndpointActionURLRewrite(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		output = append(output, cdn.URLRewriteAction{
			Name: cdn.NameURLRewrite,
			Parameters: &cdn.URLRewriteActionParameters{
				OdataType:             utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlRewriteActionParameters"),
				SourcePattern:         utils.String(item["source_pattern"].(string)),
				Destination:           utils.String(item["destination"].(string)),
				PreserveUnmatchedPath: utils.Bool(item["preserve_unmatched_path"].(bool)),
			},
		})
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionURLRewrite(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsURLRewriteAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url rewrite action!")
	}

	sourcePattern := ""
	destination := ""
	preserveUnmatchedPath := true
	if params := action.Parameters; params != nil {
		if params.Destination != nil {
			destination = *params.Destination
		}

		if params.SourcePattern != nil {
			sourcePattern = *params.SourcePattern
		}

		if params.PreserveUnmatchedPath != nil {
			preserveUnmatchedPath = *params.PreserveUnmatchedPath
		}
	}

	return &map[string]interface{}{
		"destination":             destination,
		"preserve_unmatched_path": preserveUnmatchedPath,
		"source_pattern":          sourcePattern,
	}, nil
}

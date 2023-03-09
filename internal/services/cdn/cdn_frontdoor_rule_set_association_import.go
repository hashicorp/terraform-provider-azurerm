package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importCdnFrontDoorRuleSetAssociation() pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := parse.FrontDoorRuleSetAssociationID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.RouteProperties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `AFDDomainProperties` was nil", id)
		}
		ruleSets := flattenRuleSetResourceArray(resp.RouteProperties.RuleSets)

		d.SetId(id.ID())
		d.Set("cdn_frontdoor_route_id", parse.NewFrontDoorRouteID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName).ID())
		d.Set("cdn_frontdoor_rule_set_ids", ruleSets)

		return []*pluginsdk.ResourceData{d}, nil
	}
}

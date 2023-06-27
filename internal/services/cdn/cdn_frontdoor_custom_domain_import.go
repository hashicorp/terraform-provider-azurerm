package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func importCdnFrontDoorCustomDomain() pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		insensitive, err := parse.FrontDoorCustomDomainIDInsensitively(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId

		id := parse.NewFrontDoorCustomDomainID(subscriptionId, insensitive.ResourceGroup, insensitive.ProfileName, insensitive.CustomDomainName)

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.AFDDomainProperties == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `AFDDomainProperties` was nil", id)
		}

		d.SetId(id.ID())

		return []*pluginsdk.ResourceData{d}, nil
	}
}

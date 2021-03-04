package subscription

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-11-01/subscriptions"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/subscription/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

func importSubscriptionByAlias() func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		aliasClient := meta.(*clients.Client).Subscription.AliasClient
		client := meta.(*clients.Client).Subscription.Client
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()
		aliasId, err := parse.SubscriptionAliasID(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, fmt.Errorf("failed parsing Subscription Alias ID for import")
		}
		alias, err := aliasClient.Get(ctx, aliasId.Name)
		if err != nil {
			return []*schema.ResourceData{}, fmt.Errorf("failed reading Subscription Alias: %+v", err)
		}
		if alias.Properties == nil || alias.Properties.SubscriptionID == nil {
			return []*schema.ResourceData{}, fmt.Errorf("failed reading Subscription Alias Properties, empty response or missing Subscription ID")
		}
		subscription, err := client.Get(ctx, *alias.Properties.SubscriptionID)
		if err != nil {
			return []*schema.ResourceData{}, fmt.Errorf("failed parsing Subscription details for import: %+v", err)
		}
		if subscription.State != subscriptions.Enabled {
			return []*schema.ResourceData{}, fmt.Errorf("cannot import a cancelled Subscription by Alias ID, please enable the subscription prior to import")
		}
		return []*schema.ResourceData{d}, nil
	}
}

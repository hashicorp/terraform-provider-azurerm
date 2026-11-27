// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_server

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	virtualwanClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/virtualwan/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteServerCreate(d *pluginsdk.ResourceData, meta any) error {
	client := virtualwanClient.From(meta).VirtualWANs
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := virtualwans.NewVirtualHubID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByName(id.VirtualHubName, "azurerm_route_server")
	defer locks.UnlockByName(id.VirtualHubName, "azurerm_route_server")

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.VirtualHubsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_route_server", id.ID())
		}
	}

	parameters := virtualwans.VirtualHub{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &virtualwans.VirtualHubProperties{
			Sku:                        pointer.To(d.Get("sku").(string)),
			AllowBranchToBranchTraffic: pointer.To(d.Get("branch_to_branch_traffic_enabled").(bool)),
			HubRoutingPreference:       pointer.ToEnum[virtualwans.HubRoutingPreference](d.Get("hub_routing_preference").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.VirtualHubsCreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	timeout, _ := ctx.Deadline()
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Provisioning", "Updating"},
		Target:                    []string{"Succeeded", "Provisioned"},
		Refresh:                   routeServerCreateRefreshFunc(ctx, client, id),
		PollInterval:              15 * time.Second,
		ContinuousTargetOccurence: 5,
		Timeout:                   time.Until(timeout),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	ipConfigName := "ipConfig1"
	ipConfigs := virtualwans.HubIPConfiguration{
		Name: pointer.To(ipConfigName),
		Properties: &virtualwans.HubIPConfigurationPropertiesFormat{
			PublicIPAddress: &virtualwans.PublicIPAddress{
				Id: pointer.To(d.Get("public_ip_address_id").(string)),
			},
			Subnet: &virtualwans.Subnet{
				Id: pointer.To(d.Get("subnet_id").(string)),
			},
		},
	}

	ipConfigId := commonids.NewVirtualHubIPConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, ipConfigName)

	if err := client.VirtualHubIPConfigurationCreateOrUpdateThenPoll(ctx, ipConfigId, ipConfigs); err != nil {
		return fmt.Errorf("creating %s: %+v", ipConfigId, err)
	}

	return resourceRouteServerRead(d, meta)
}

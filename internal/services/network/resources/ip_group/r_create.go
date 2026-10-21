// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ip_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipgroups"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIpGroupCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.IPGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	for _, fw := range d.Get("firewall_ids").([]any) {
		id, err := azurefirewalls.ParseAzureFirewallID(fw.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall ID %q: %+v", fw, err)
		}
		locks.ByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, firewall.AzureFirewallResourceName)
	}

	for _, fwpol := range d.Get("firewall_policy_ids").([]any) {
		id, err := firewallpolicies.ParseFirewallPolicyID(fwpol.(string))
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Policy ID %q: %+v", fwpol, err)
		}
		locks.ByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.FirewallPolicyName, firewall.AzureFirewallPolicyResourceName)
	}

	id := ipgroups.NewIPGroupID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, ipgroups.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_ip_group", id.ID())
		}
	}

	ipAddresses := d.Get("cidrs").(*pluginsdk.Set).List()

	sg := ipgroups.IPGroup{
		Name:     &id.IpGroupName,
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &ipgroups.IPGroupPropertiesFormat{
			IPAddresses: pluginsdk.ExpandStringSlice(ipAddresses),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, sg, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceIpGroupRead(d, meta)
}

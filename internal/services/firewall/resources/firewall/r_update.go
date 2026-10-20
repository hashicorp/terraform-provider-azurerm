// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFirewallCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := azurefirewalls.NewAzureFirewallID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if d.IsNewResource() && !response.WasNotFound(existing.HttpResponse) {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			return tf.ImportAsExistsError("azurerm_firewall", id.ID())
		}
	}

	if err := validateFirewallIPConfigurationSettings(d.Get("ip_configuration").([]any)); err != nil {
		return fmt.Errorf("validating %s: %+v", id, err)
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]any)
	i := d.Get("ip_configuration").([]any)
	ipConfigs, subnetToLock, vnetToLock, err := expandFirewallIPConfigurations(i)
	if err != nil {
		return fmt.Errorf("building list of Azure Firewall IP Configurations: %+v", err)
	}

	parameters := azurefirewalls.AzureFirewall{
		Location: &location,
		Properties: &azurefirewalls.AzureFirewallPropertiesFormat{
			IPConfigurations:     ipConfigs,
			ThreatIntelMode:      pointer.ToEnum[azurefirewalls.AzureFirewallThreatIntelMode](d.Get("threat_intel_mode").(string)),
			AdditionalProperties: pointer.To(make(map[string]string)),
		},
		Tags: tags.Expand(t),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = &zones
	}

	m := d.Get("management_ip_configuration").([]any)
	if len(m) == 1 {
		mgmtIPConfig, mgmtSubnetName, mgmtVirtualNetworkName, err := expandFirewallIPConfigurations(m)
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Management IP Configurations: %+v", err)
		}

		if !slices.Contains(*subnetToLock, (*mgmtSubnetName)[0]) {
			*subnetToLock = append(*subnetToLock, (*mgmtSubnetName)[0])
		}

		if !slices.Contains(*vnetToLock, (*mgmtVirtualNetworkName)[0]) {
			*vnetToLock = append(*vnetToLock, (*mgmtVirtualNetworkName)[0])
		}
		if *mgmtIPConfig != nil {
			if parameters.Properties.IPConfigurations != nil {
				for k, v := range *parameters.Properties.IPConfigurations {
					if v.Name != nil && (*mgmtIPConfig)[0].Name != nil && *v.Name == *(*mgmtIPConfig)[0].Name {
						return fmt.Errorf("`management_ip_configuration.0.name` must not be the same as `ip_configuration.%d.name`", k)
					}
				}
			}

			parameters.Properties.ManagementIPConfiguration = &(*mgmtIPConfig)[0]
		}
	}

	if threatIntelMode := d.Get("threat_intel_mode").(string); threatIntelMode != "" {
		parameters.Properties.ThreatIntelMode = pointer.ToEnum[azurefirewalls.AzureFirewallThreatIntelMode](threatIntelMode)
	}

	if policyId := d.Get("firewall_policy_id").(string); policyId != "" {
		parameters.Properties.FirewallPolicy = &azurefirewalls.SubResource{Id: &policyId}
	}

	vhub, hubIpAddresses, ok := expandFirewallVirtualHubSetting(existing.Model, d.Get("virtual_hub").([]any))
	if ok {
		parameters.Properties.VirtualHub = vhub
		parameters.Properties.HubIPAddresses = hubIpAddresses
	}

	if skuName := d.Get("sku_name").(string); skuName != "" {
		if parameters.Properties.Sku == nil {
			parameters.Properties.Sku = &azurefirewalls.AzureFirewallSku{}
		}
		parameters.Properties.Sku.Name = pointer.ToEnum[azurefirewalls.AzureFirewallSkuName](skuName)
	}

	if skuTier := d.Get("sku_tier").(string); skuTier != "" {
		if parameters.Properties.Sku == nil {
			parameters.Properties.Sku = &azurefirewalls.AzureFirewallSku{}
		}
		parameters.Properties.Sku.Tier = pointer.ToEnum[azurefirewalls.AzureFirewallSkuTier](skuTier)
	}

	if dnsServerSetting := expandFirewallAdditionalProperty(d); dnsServerSetting != nil {
		for k, v := range dnsServerSetting {
			attrs := *parameters.Properties.AdditionalProperties
			attrs[k] = v
		}
	}

	if privateIpRangeSetting := expandFirewallPrivateIpRange(d.Get("private_ip_ranges").(*pluginsdk.Set).List()); privateIpRangeSetting != nil {
		for k, v := range privateIpRangeSetting {
			attrs := *parameters.Properties.AdditionalProperties
			attrs[k] = v
		}
	}

	if policyId, ok := d.GetOk("firewall_policy_id"); ok {
		id, _ := firewallpolicies.ParseFirewallPolicyID(policyId.(string))
		locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
		defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
	}

	locks.ByName(id.AzureFirewallName, AzureFirewallResourceName)
	defer locks.UnlockByName(id.AzureFirewallName, AzureFirewallResourceName)

	locks.MultipleByName(vnetToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetToLock, VirtualNetworkResourceName)

	locks.MultipleByName(subnetToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(subnetToLock, SubnetResourceName)

	if !d.IsNewResource() {
		exists, err2 := client.Get(ctx, id)
		if err2 != nil {
			if response.WasNotFound(exists.HttpResponse) {
				return fmt.Errorf("retrieving existing %s: firewall not found in resource group", id)
			}
			return fmt.Errorf("retrieving existing %s: %+v", id, err2)
		}
		if exists.Model == nil {
			return fmt.Errorf("retrieving existing rules for %s: `model` was nil", id)
		}

		if exists.Model.Properties == nil {
			return fmt.Errorf("retrieving existing rules for %s: `props` was nil", id)
		}
		props := *exists.Model.Properties
		parameters.Properties.ApplicationRuleCollections = props.ApplicationRuleCollections
		parameters.Properties.NetworkRuleCollections = props.NetworkRuleCollections
		parameters.Properties.NatRuleCollections = props.NatRuleCollections
	}

	if d.IsNewResource() {
		if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, azurefirewalls.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		d.SetId(id.ID())
		if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
			return err
		}
	} else {
		if err := client.CreateOrUpdateThenPoll(ctx, id, parameters, azurefirewalls.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceFirewallRead(d, meta)
}

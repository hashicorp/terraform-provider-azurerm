// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var AzureFirewallResourceName = "azurerm_firewall"

func resourceFirewall() *pluginsdk.Resource {
	resource := pluginsdk.Resource{
		Create: resourceFirewallCreateUpdate,
		Read:   resourceFirewallRead,
		Update: resourceFirewallCreateUpdate,
		Delete: resourceFirewallDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := azurefirewalls.ParseAzureFirewallID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FirewallName,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			// lintignore:S013
			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(azurefirewalls.AzureFirewallSkuNameAZFWHub),
					string(azurefirewalls.AzureFirewallSkuNameAZFWVNet),
				}, false),
			},

			// lintignore:S013
			"sku_tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(azurefirewalls.AzureFirewallSkuTierPremium),
					string(azurefirewalls.AzureFirewallSkuTierStandard),
					string(azurefirewalls.AzureFirewallSkuTierBasic),
				}, false),
			},

			"firewall_policy_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: firewallpolicies.ValidateFirewallPolicyID,
			},

			"ip_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.FirewallSubnetName,
						},
						"public_ip_address_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidatePublicIPAddressID,
						},
						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"management_ip_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.FirewallManagementSubnetName,
						},
						"public_ip_address_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: commonids.ValidatePublicIPAddressID,
						},
						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"threat_intel_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(azurefirewalls.AzureFirewallThreatIntelModeOff),
					string(azurefirewalls.AzureFirewallThreatIntelModeAlert),
					string(azurefirewalls.AzureFirewallThreatIntelModeDeny),
				}, false),
			},

			"dns_servers": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
			},

			"dns_proxy_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"private_ip_ranges": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.Any(
						validation.IsCIDR,
						validation.StringInSlice([]string{"IANAPrivateRanges"}, false),
					),
				},
			},

			"virtual_hub": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"virtual_hub_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: virtualwans.ValidateVirtualHubID,
						},
						"public_ip_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Default:      1,
						},
						"public_ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
						},
						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"zones": commonschema.ZonesMultipleOptionalForceNew(),

			"tags": commonschema.Tags(),
		},
	}

	return &resource
}

func resourceFirewallCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Azure Firewall creation")

	id := azurefirewalls.NewAzureFirewallID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if d.IsNewResource() && !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_firewall", id.ID())
	}

	if err := validateFirewallIPConfigurationSettings(d.Get("ip_configuration").([]interface{})); err != nil {
		return fmt.Errorf("validating %s: %+v", id, err)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	i := d.Get("ip_configuration").([]interface{})
	ipConfigs, subnetToLock, vnetToLock, err := expandFirewallIPConfigurations(i)
	if err != nil {
		return fmt.Errorf("building list of Azure Firewall IP Configurations: %+v", err)
	}

	parameters := azurefirewalls.AzureFirewall{
		Location: &location,
		Properties: &azurefirewalls.AzureFirewallPropertiesFormat{
			IPConfigurations:     ipConfigs,
			ThreatIntelMode:      pointer.To(azurefirewalls.AzureFirewallThreatIntelMode(d.Get("threat_intel_mode").(string))),
			AdditionalProperties: pointer.To(make(map[string]string)),
		},
		Tags: tags.Expand(t),
	}

	zones := zones.ExpandUntyped(d.Get("zones").(*schema.Set).List())
	if len(zones) > 0 {
		parameters.Zones = &zones
	}

	m := d.Get("management_ip_configuration").([]interface{})
	if len(m) == 1 {
		mgmtIPConfig, mgmtSubnetName, mgmtVirtualNetworkName, err := expandFirewallIPConfigurations(m)
		if err != nil {
			return fmt.Errorf("parsing Azure Firewall Management IP Configurations: %+v", err)
		}

		if !utils.SliceContainsValue(*subnetToLock, (*mgmtSubnetName)[0]) {
			*subnetToLock = append(*subnetToLock, (*mgmtSubnetName)[0])
		}

		if !utils.SliceContainsValue(*vnetToLock, (*mgmtVirtualNetworkName)[0]) {
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
		parameters.Properties.ThreatIntelMode = pointer.To(azurefirewalls.AzureFirewallThreatIntelMode(threatIntelMode))
	}

	if policyId := d.Get("firewall_policy_id").(string); policyId != "" {
		parameters.Properties.FirewallPolicy = &azurefirewalls.SubResource{Id: &policyId}
	}

	vhub, hubIpAddresses, ok := expandFirewallVirtualHubSetting(existing.Model, d.Get("virtual_hub").([]interface{}))
	if ok {
		parameters.Properties.VirtualHub = vhub
		parameters.Properties.HubIPAddresses = hubIpAddresses
	}

	if skuName := d.Get("sku_name").(string); skuName != "" {
		if parameters.Properties.Sku == nil {
			parameters.Properties.Sku = &azurefirewalls.AzureFirewallSku{}
		}
		parameters.Properties.Sku.Name = pointer.To(azurefirewalls.AzureFirewallSkuName(skuName))
	}

	if skuTier := d.Get("sku_tier").(string); skuTier != "" {
		if parameters.Properties.Sku == nil {
			parameters.Properties.Sku = &azurefirewalls.AzureFirewallSku{}
		}
		parameters.Properties.Sku.Tier = pointer.To(azurefirewalls.AzureFirewallSkuTier(skuTier))
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

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFirewallRead(d, meta)
}

func resourceFirewallRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azurefirewalls.ParseAzureFirewallID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - removing from state!", id.AzureFirewallName, id.ResourceGroupName)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Firewall %s : %+v", *id, err)
	}

	d.Set("name", id.AzureFirewallName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := read.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("zones", zones.FlattenUntyped(model.Zones))

		if props := model.Properties; props != nil {
			if err := d.Set("ip_configuration", flattenFirewallIPConfigurations(props.IPConfigurations)); err != nil {
				return fmt.Errorf("setting `ip_configuration`: %+v", err)
			}
			managementIPConfigs := make([]interface{}, 0)
			if props.ManagementIPConfiguration != nil {
				managementIPConfigs = flattenFirewallIPConfigurations(&[]azurefirewalls.AzureFirewallIPConfiguration{
					*props.ManagementIPConfiguration,
				})
			}
			if err := d.Set("management_ip_configuration", managementIPConfigs); err != nil {
				return fmt.Errorf("setting `management_ip_configuration`: %+v", err)
			}

			d.Set("threat_intel_mode", string(pointer.From(props.ThreatIntelMode)))

			dnsProxyEnabled, dnsServers := flattenFirewallAdditionalProperty(props.AdditionalProperties)
			if err := d.Set("dns_proxy_enabled", dnsProxyEnabled); err != nil {
				return fmt.Errorf("setting `dns_proxy_enabled`: %+v", err)
			}
			if err := d.Set("dns_servers", dnsServers); err != nil {
				return fmt.Errorf("setting `dns_servers`: %+v", err)
			}

			if err := d.Set("private_ip_ranges", flattenFirewallPrivateIpRange(props.AdditionalProperties)); err != nil {
				return fmt.Errorf("setting `private_ip_ranges`: %+v", err)
			}

			firewallPolicyId := ""
			if props.FirewallPolicy != nil && props.FirewallPolicy.Id != nil {
				firewallPolicyId = *props.FirewallPolicy.Id
				if policyId, err := firewallpolicies.ParseFirewallPolicyIDInsensitively(firewallPolicyId); err == nil {
					firewallPolicyId = policyId.ID()
				}
			}
			d.Set("firewall_policy_id", firewallPolicyId)

			if sku := props.Sku; sku != nil {
				d.Set("sku_name", string(pointer.From(sku.Name)))
				d.Set("sku_tier", string(pointer.From(sku.Tier)))
			}

			if err := d.Set("virtual_hub", flattenFirewallVirtualHubSetting(props)); err != nil {
				return fmt.Errorf("setting `virtual_hub`: %+v", err)
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceFirewallDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.AzureFirewalls
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azurefirewalls.ParseAzureFirewallID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(read.HttpResponse) {
			// deleted outside of TF
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - assuming removed!", id.AzureFirewallName, id.ResourceGroupName)
			return nil
		}

		return fmt.Errorf("retrieving Firewall %s : %+v", *id, err)
	}

	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)
	if model := read.Model; model != nil {
		if props := model.Properties; props != nil {
			if configs := props.IPConfigurations; configs != nil {
				for _, config := range *configs {
					if config.Properties == nil || config.Properties.Subnet == nil || config.Properties.Subnet.Id == nil {
						continue
					}

					parsedSubnetID, err2 := commonids.ParseSubnetID(*config.Properties.Subnet.Id)
					if err2 != nil {
						return err2
					}

					if !utils.SliceContainsValue(subnetNamesToLock, parsedSubnetID.SubnetName) {
						subnetNamesToLock = append(subnetNamesToLock, parsedSubnetID.SubnetName)
					}

					if !utils.SliceContainsValue(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName) {
						virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName)
					}
				}
			}

			if mconfig := props.ManagementIPConfiguration; mconfig != nil {
				if mconfig.Properties != nil && mconfig.Properties.Subnet != nil && mconfig.Properties.Subnet.Id != nil {
					parsedSubnetID, err2 := commonids.ParseSubnetID(*mconfig.Properties.Subnet.Id)
					if err2 != nil {
						return err2
					}

					if !utils.SliceContainsValue(subnetNamesToLock, parsedSubnetID.SubnetName) {
						subnetNamesToLock = append(subnetNamesToLock, parsedSubnetID.SubnetName)
					}

					if !utils.SliceContainsValue(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName) {
						virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName)
					}
				}
			}
		}

		if read.Model.Properties != nil && read.Model.Properties.FirewallPolicy != nil && read.Model.Properties.FirewallPolicy.Id != nil {
			id, err := firewallpolicies.ParseFirewallPolicyIDInsensitively(*read.Model.Properties.FirewallPolicy.Id)
			if err != nil {
				return err
			}
			locks.ByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
			defer locks.UnlockByName(id.FirewallPolicyName, AzureFirewallPolicyResourceName)
		}

		locks.ByName(id.AzureFirewallName, AzureFirewallResourceName)
		defer locks.UnlockByName(id.AzureFirewallName, AzureFirewallResourceName)

		locks.MultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)
		defer locks.UnlockMultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)

		locks.MultipleByName(&subnetNamesToLock, SubnetResourceName)
		defer locks.UnlockMultipleByName(&subnetNamesToLock, SubnetResourceName)

		// todo see if this is still needed this way
		/*
			// Change this back to using the SDK method once https://github.com/Azure/azure-sdk-for-go/issues/17013 is addressed.
			future, err := azuresdkhacks.DeleteFirewall(ctx, client, id.ResourceGroup, id.AzureFirewallName)
			if err != nil {
				return fmt.Errorf("deleting Azure Firewall %s : %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the deletion of Azure Firewall %s : %+v", *id, err)
			}
		*/

		return client.DeleteThenPoll(ctx, *id)
	}

	return err
}

func expandFirewallIPConfigurations(configs []interface{}) (*[]azurefirewalls.AzureFirewallIPConfiguration, *[]string, *[]string, error) {
	ipConfigs := make([]azurefirewalls.AzureFirewallIPConfiguration, 0)
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		name := data["name"].(string)
		subnetId := data["subnet_id"].(string)
		pubID := data["public_ip_address_id"].(string)

		ipConfig := azurefirewalls.AzureFirewallIPConfiguration{
			Name:       utils.String(name),
			Properties: &azurefirewalls.AzureFirewallIPConfigurationPropertiesFormat{},
		}

		if pubID != "" {
			ipConfig.Properties.PublicIPAddress = &azurefirewalls.SubResource{
				Id: utils.String(pubID),
			}
		}

		if subnetId != "" {
			subnetID, err := commonids.ParseSubnetID(subnetId)
			if err != nil {
				return nil, nil, nil, err
			}

			if !utils.SliceContainsValue(subnetNamesToLock, subnetID.SubnetName) {
				subnetNamesToLock = append(subnetNamesToLock, subnetID.SubnetName)
			}

			if !utils.SliceContainsValue(virtualNetworkNamesToLock, subnetID.VirtualNetworkName) {
				virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, subnetID.VirtualNetworkName)
			}

			ipConfig.Properties.Subnet = &azurefirewalls.SubResource{
				Id: utils.String(subnetId),
			}
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return &ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func flattenFirewallIPConfigurations(input *[]azurefirewalls.AzureFirewallIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		afIPConfig := make(map[string]interface{})
		props := v.Properties
		if props == nil {
			continue
		}

		if name := v.Name; name != nil {
			afIPConfig["name"] = *name
		}

		if subnet := props.Subnet; subnet != nil {
			if id := subnet.Id; id != nil {
				afIPConfig["subnet_id"] = *id
			}
		}

		if ipAddress := props.PrivateIPAddress; ipAddress != nil {
			afIPConfig["private_ip_address"] = *ipAddress
		}

		if pip := props.PublicIPAddress; pip != nil {
			if id := pip.Id; id != nil {
				afIPConfig["public_ip_address_id"] = *id
			}
		}
		result = append(result, afIPConfig)
	}

	return result
}

func expandFirewallAdditionalProperty(d *pluginsdk.ResourceData) map[string]string {
	// Swagger issue asking finalize these properties: https://github.com/Azure/azure-rest-api-specs/issues/11278
	res := map[string]string{}
	if servers := d.Get("dns_servers").([]interface{}); len(servers) > 0 {
		var servs []string
		for _, server := range servers {
			servs = append(servs, server.(string))
		}
		res["Network.DNS.EnableProxy"] = "true"
		res["Network.DNS.Servers"] = strings.Join(servs, ",")
	}
	if enabled := d.Get("dns_proxy_enabled").(bool); enabled {
		res["Network.DNS.EnableProxy"] = "true"
	}
	return res
}

func flattenFirewallAdditionalProperty(input *map[string]string) (enabled interface{}, servers []interface{}) {
	if input == nil || len(*input) == 0 {
		return nil, nil
	}

	if enabledPtr, ok := (*input)["Network.DNS.EnableProxy"]; ok {
		enabled = enabledPtr == "true"
	}

	if serversPtr, ok := (*input)["Network.DNS.Servers"]; ok {
		for _, val := range strings.Split(serversPtr, ",") {
			servers = append(servers, val)
		}
	}
	return
}

func expandFirewallPrivateIpRange(input []interface{}) map[string]string {
	if len(input) == 0 {
		return nil
	}

	rangeSlice := *utils.ExpandStringSlice(input)
	if len(rangeSlice) == 0 {
		return nil
	}

	// Swagger issue asking finalize these properties: https://github.com/Azure/azure-rest-api-specs/issues/10015
	return map[string]string{
		"Network.SNAT.PrivateRanges": strings.Join(rangeSlice, ","),
	}
}

func flattenFirewallPrivateIpRange(input *map[string]string) []interface{} {
	if input == nil && len(*input) == 0 {
		return nil
	}

	attrs := *input
	rangeSlice := []string{}
	if privateIpRanges := attrs["Network.SNAT.PrivateRanges"]; privateIpRanges != "" {
		rangeSlice = strings.Split(attrs["Network.SNAT.PrivateRanges"], ",")
	}
	return utils.FlattenStringSlice(&rangeSlice)
}

func expandFirewallVirtualHubSetting(existing *azurefirewalls.AzureFirewall, input []interface{}) (vhub *azurefirewalls.SubResource, ipAddresses *azurefirewalls.HubIPAddresses, ok bool) {
	if len(input) == 0 {
		return nil, nil, false
	}

	b := input[0].(map[string]interface{})

	// The API requires both "Count" and "Addresses" for the "PublicIPs" setting.
	// The "Count" means how many PIP to provision.
	// The "Addresses" means differently in different cases:
	// - Create: only "Count" is needed, "Addresses" is not necessary
	// - Update: both "Count" and "Addresses" are needed:
	//     Scale up: "Addresses" should remain same as before scaling up
	//     Scale down: "Addresses" should indicate the addresses to be retained (in this case we retain the first new "Count" ones)
	newCount := b["public_ip_count"].(int)
	var addresses *[]azurefirewalls.AzureFirewallPublicIPAddress
	if existing != nil {
		if prop := existing.Properties; prop != nil {
			if ipaddress := prop.HubIPAddresses; ipaddress != nil {
				if pips := ipaddress.PublicIPs; pips != nil {
					if count := pips.Count; count != nil {
						oldCount := int(*count)
						addresses = pips.Addresses

						// In case of scale down, keep the first new "Count" addresses.
						if oldCount > newCount {
							keptAddresses := make([]azurefirewalls.AzureFirewallPublicIPAddress, newCount)
							for i := 0; i < newCount; i++ {
								keptAddresses[i] = (*addresses)[i]
							}
							addresses = &keptAddresses
						}
					}
				}
			}
		}
	}

	vhub = &azurefirewalls.SubResource{Id: utils.String(b["virtual_hub_id"].(string))}
	ipAddresses = &azurefirewalls.HubIPAddresses{
		PublicIPs: &azurefirewalls.HubPublicIPAddresses{
			Count:     utils.Int64(int64(b["public_ip_count"].(int))),
			Addresses: addresses,
		},
	}

	return vhub, ipAddresses, true
}

func flattenFirewallVirtualHubSetting(props *azurefirewalls.AzureFirewallPropertiesFormat) []interface{} {
	if props.VirtualHub == nil {
		return nil
	}

	var vhubId string
	if props.VirtualHub.Id != nil {
		vhubId = *props.VirtualHub.Id
	}

	var (
		publicIpCount int
		publicIps     []string
		privateIp     string
	)
	if hubIP := props.HubIPAddresses; hubIP != nil {
		if hubIP.PrivateIPAddress != nil {
			privateIp = *hubIP.PrivateIPAddress
		}
		if pubIPs := hubIP.PublicIPs; pubIPs != nil {
			if pubIPs.Count != nil {
				publicIpCount = int(*pubIPs.Count)
			}
			if pubIPs.Addresses != nil {
				for _, addr := range *pubIPs.Addresses {
					if addr.Address != nil {
						publicIps = append(publicIps, *addr.Address)
					}
				}
			}
		}
	}

	return []interface{}{
		map[string]interface{}{
			"virtual_hub_id":      vhubId,
			"public_ip_count":     publicIpCount,
			"public_ip_addresses": publicIps,
			"private_ip_address":  privateIp,
		},
	}
}

func validateFirewallIPConfigurationSettings(configs []interface{}) error {
	if len(configs) == 0 {
		return nil
	}

	subnetNumber := 0

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		if subnet, exist := data["subnet_id"].(string); exist && subnet != "" {
			subnetNumber++
		}
	}

	if subnetNumber != 1 {
		return fmt.Errorf(`The "ip_configuration" is invalid, %d "subnet_id" have been set, one "subnet_id" should be set among all "ip_configuration" blocks`, subnetNumber)
	}

	return nil
}

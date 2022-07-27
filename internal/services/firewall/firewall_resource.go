package firewall

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var azureFirewallResourceName = "azurerm_firewall"

func resourceFirewall() *pluginsdk.Resource {
	resource := pluginsdk.Resource{
		Create: resourceFirewallCreateUpdate,
		Read:   resourceFirewallRead,
		Update: resourceFirewallCreateUpdate,
		Delete: resourceFirewallDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FirewallID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			//lintignore:S013
			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallSkuNameAZFWHub),
					string(network.AzureFirewallSkuNameAZFWVNet),
				}, false),
			},

			//lintignore:S013
			"sku_tier": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.AzureFirewallSkuTierPremium),
					string(network.AzureFirewallSkuTierStandard),
					string(network.AzureFirewallSkuTierBasic),
				}, false),
			},

			"firewall_policy_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.FirewallPolicyID,
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
							Required:     true,
							ValidateFunc: networkValidate.PublicIpAddressID,
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
							ValidateFunc: networkValidate.PublicIpAddressID,
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
					string(network.AzureFirewallThreatIntelModeOff),
					string(network.AzureFirewallThreatIntelModeAlert),
					string(network.AzureFirewallThreatIntelModeDeny),
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
							ValidateFunc: networkValidate.VirtualHubID,
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

			"tags": tags.Schema(),
		},
	}

	return &resource
}

func resourceFirewallCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Azure Firewall creation")

	id := parse.NewFirewallID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if d.IsNewResource() && !utils.ResponseWasNotFound(existing.Response) {
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

	parameters := network.AzureFirewall{
		Location: &location,
		AzureFirewallPropertiesFormat: &network.AzureFirewallPropertiesFormat{
			IPConfigurations:     ipConfigs,
			ThreatIntelMode:      network.AzureFirewallThreatIntelMode(d.Get("threat_intel_mode").(string)),
			AdditionalProperties: make(map[string]*string),
		},
		Tags: tags.Expand(t),
	}

	zones := zones.Expand(d.Get("zones").(*schema.Set).List())
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
			parameters.ManagementIPConfiguration = &(*mgmtIPConfig)[0]
		}
	}

	if threatIntelMode := d.Get("threat_intel_mode").(string); threatIntelMode != "" {
		parameters.AzureFirewallPropertiesFormat.ThreatIntelMode = network.AzureFirewallThreatIntelMode(threatIntelMode)
	}

	if policyId := d.Get("firewall_policy_id").(string); policyId != "" {
		parameters.AzureFirewallPropertiesFormat.FirewallPolicy = &network.SubResource{ID: &policyId}
	}

	vhub, hubIpAddresses, ok := expandFirewallVirtualHubSetting(existing, d.Get("virtual_hub").([]interface{}))
	if ok {
		parameters.AzureFirewallPropertiesFormat.VirtualHub = vhub
		parameters.AzureFirewallPropertiesFormat.HubIPAddresses = hubIpAddresses
	}

	if skuName := d.Get("sku_name").(string); skuName != "" {
		if parameters.Sku == nil {
			parameters.Sku = &network.AzureFirewallSku{}
		}
		parameters.Sku.Name = network.AzureFirewallSkuName(skuName)
	}

	if skuTier := d.Get("sku_tier").(string); skuTier != "" {
		if parameters.Sku == nil {
			parameters.Sku = &network.AzureFirewallSku{}
		}
		parameters.Sku.Tier = network.AzureFirewallSkuTier(skuTier)
	}

	if dnsServerSetting := expandFirewallDNSServers(d.Get("dns_servers").([]interface{})); dnsServerSetting != nil {
		for k, v := range dnsServerSetting {
			parameters.AdditionalProperties[k] = v
		}
	}

	if privateIpRangeSetting := expandFirewallPrivateIpRange(d.Get("private_ip_ranges").(*pluginsdk.Set).List()); privateIpRangeSetting != nil {
		for k, v := range privateIpRangeSetting {
			parameters.AdditionalProperties[k] = v
		}
	}

	locks.ByName(id.AzureFirewallName, azureFirewallResourceName)
	defer locks.UnlockByName(id.AzureFirewallName, azureFirewallResourceName)

	locks.MultipleByName(vnetToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetToLock, VirtualNetworkResourceName)

	locks.MultipleByName(subnetToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(subnetToLock, SubnetResourceName)

	if !d.IsNewResource() {
		exists, err2 := client.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
		if err2 != nil {
			if utils.ResponseWasNotFound(exists.Response) {
				return fmt.Errorf("retrieving existing %s: firewall not found in resource group", id)
			}
			return fmt.Errorf("retrieving existing %s: %+v", id, err2)
		}
		if exists.AzureFirewallPropertiesFormat == nil {
			return fmt.Errorf("retrieving existing rules for %s: `props` was nil", id)
		}
		props := *exists.AzureFirewallPropertiesFormat
		parameters.AzureFirewallPropertiesFormat.ApplicationRuleCollections = props.ApplicationRuleCollections
		parameters.AzureFirewallPropertiesFormat.NetworkRuleCollections = props.NetworkRuleCollections
		parameters.AzureFirewallPropertiesFormat.NatRuleCollections = props.NatRuleCollections
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.AzureFirewallName, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFirewallRead(d, meta)
}

func resourceFirewallRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - removing from state!", id.AzureFirewallName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure Firewall %s : %+v", *id, err)
	}

	d.Set("name", id.AzureFirewallName)
	d.Set("resource_group_name", id.ResourceGroup)

	d.Set("location", location.NormalizeNilable(read.Location))
	d.Set("zones", zones.Flatten(read.Zones))

	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if err := d.Set("ip_configuration", flattenFirewallIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("setting `ip_configuration`: %+v", err)
		}
		managementIPConfigs := make([]interface{}, 0)
		if props.ManagementIPConfiguration != nil {
			managementIPConfigs = flattenFirewallIPConfigurations(&[]network.AzureFirewallIPConfiguration{
				*props.ManagementIPConfiguration,
			})
		}
		if err := d.Set("management_ip_configuration", managementIPConfigs); err != nil {
			return fmt.Errorf("setting `management_ip_configuration`: %+v", err)
		}

		d.Set("threat_intel_mode", string(props.ThreatIntelMode))

		if err := d.Set("dns_servers", flattenFirewallDNSServers(props.AdditionalProperties)); err != nil {
			return fmt.Errorf("setting `dns_servers`: %+v", err)
		}

		if err := d.Set("private_ip_ranges", flattenFirewallPrivateIpRange(props.AdditionalProperties)); err != nil {
			return fmt.Errorf("setting `private_ip_ranges`: %+v", err)
		}

		if policy := props.FirewallPolicy; policy != nil {
			d.Set("firewall_policy_id", policy.ID)
		}

		if sku := props.Sku; sku != nil {
			d.Set("sku_name", string(sku.Name))
			d.Set("sku_tier", string(sku.Tier))
		}

		if err := d.Set("virtual_hub", flattenFirewallVirtualHubSetting(props)); err != nil {
			return fmt.Errorf("setting `virtual_hub`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceFirewallDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Firewall.AzureFirewallsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FirewallID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] Firewall %q was not found in Resource Group %q - assuming removed!", id.AzureFirewallName, id.ResourceGroup)
			return nil
		}

		return fmt.Errorf("retrieving Firewall %s : %+v", *id, err)
	}

	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)
	if props := read.AzureFirewallPropertiesFormat; props != nil {
		if configs := props.IPConfigurations; configs != nil {
			for _, config := range *configs {
				if config.Subnet == nil || config.Subnet.ID == nil {
					continue
				}

				parsedSubnetID, err2 := networkParse.SubnetID(*config.Subnet.ID)
				if err2 != nil {
					return err2
				}

				if !utils.SliceContainsValue(subnetNamesToLock, parsedSubnetID.Name) {
					subnetNamesToLock = append(subnetNamesToLock, parsedSubnetID.Name)
				}

				if !utils.SliceContainsValue(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName) {
					virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName)
				}
			}
		}

		if mconfig := props.ManagementIPConfiguration; mconfig != nil {
			if mconfig.Subnet != nil && mconfig.Subnet.ID != nil {
				parsedSubnetID, err2 := networkParse.SubnetID(*mconfig.Subnet.ID)
				if err2 != nil {
					return err2
				}

				if !utils.SliceContainsValue(subnetNamesToLock, parsedSubnetID.Name) {
					subnetNamesToLock = append(subnetNamesToLock, parsedSubnetID.Name)
				}

				if !utils.SliceContainsValue(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName) {
					virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, parsedSubnetID.VirtualNetworkName)
				}
			}
		}
	}

	locks.ByName(id.AzureFirewallName, azureFirewallResourceName)
	defer locks.UnlockByName(id.AzureFirewallName, azureFirewallResourceName)

	locks.MultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNamesToLock, VirtualNetworkResourceName)

	locks.MultipleByName(&subnetNamesToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(&subnetNamesToLock, SubnetResourceName)

	// Change this back to using the SDK method once https://github.com/Azure/azure-sdk-for-go/issues/17013 is addressed.
	future, err := azuresdkhacks.DeleteFirewall(ctx, client, id.ResourceGroup, id.AzureFirewallName)
	if err != nil {
		return fmt.Errorf("deleting Azure Firewall %s : %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of Azure Firewall %s : %+v", *id, err)
	}

	return err
}

func expandFirewallIPConfigurations(configs []interface{}) (*[]network.AzureFirewallIPConfiguration, *[]string, *[]string, error) {
	ipConfigs := make([]network.AzureFirewallIPConfiguration, 0)
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})
		name := data["name"].(string)
		subnetId := data["subnet_id"].(string)
		pubID := data["public_ip_address_id"].(string)

		ipConfig := network.AzureFirewallIPConfiguration{
			Name: utils.String(name),
			AzureFirewallIPConfigurationPropertiesFormat: &network.AzureFirewallIPConfigurationPropertiesFormat{
				PublicIPAddress: &network.SubResource{
					ID: utils.String(pubID),
				},
			},
		}

		if subnetId != "" {
			subnetID, err := networkParse.SubnetID(subnetId)
			if err != nil {
				return nil, nil, nil, err
			}

			if !utils.SliceContainsValue(subnetNamesToLock, subnetID.Name) {
				subnetNamesToLock = append(subnetNamesToLock, subnetID.Name)
			}

			if !utils.SliceContainsValue(virtualNetworkNamesToLock, subnetID.VirtualNetworkName) {
				virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, subnetID.VirtualNetworkName)
			}

			ipConfig.AzureFirewallIPConfigurationPropertiesFormat.Subnet = &network.SubResource{
				ID: utils.String(subnetId),
			}
		}
		ipConfigs = append(ipConfigs, ipConfig)
	}
	return &ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func flattenFirewallIPConfigurations(input *[]network.AzureFirewallIPConfiguration) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		afIPConfig := make(map[string]interface{})
		props := v.AzureFirewallIPConfigurationPropertiesFormat
		if props == nil {
			continue
		}

		if name := v.Name; name != nil {
			afIPConfig["name"] = *name
		}

		if subnet := props.Subnet; subnet != nil {
			if id := subnet.ID; id != nil {
				afIPConfig["subnet_id"] = *id
			}
		}

		if ipAddress := props.PrivateIPAddress; ipAddress != nil {
			afIPConfig["private_ip_address"] = *ipAddress
		}

		if pip := props.PublicIPAddress; pip != nil {
			if id := pip.ID; id != nil {
				afIPConfig["public_ip_address_id"] = *id
			}
		}
		result = append(result, afIPConfig)
	}

	return result
}

func expandFirewallDNSServers(input []interface{}) map[string]*string {
	if len(input) == 0 {
		return nil
	}

	var servers []string
	for _, server := range input {
		servers = append(servers, server.(string))
	}

	// Swagger issue asking finalize these properties: https://github.com/Azure/azure-rest-api-specs/issues/11278
	return map[string]*string{
		"Network.DNS.EnableProxy": utils.String("true"),
		"Network.DNS.Servers":     utils.String(strings.Join(servers, ",")),
	}
}

func flattenFirewallDNSServers(input map[string]*string) []interface{} {
	if len(input) == 0 {
		return nil
	}

	enabled := false
	if enabledPtr := input["Network.DNS.EnableProxy"]; enabledPtr != nil {
		enabled = *enabledPtr == "true"
	}

	if !enabled {
		return nil
	}

	servers := []string{}
	if serversPtr := input["Network.DNS.Servers"]; serversPtr != nil {
		servers = strings.Split(*serversPtr, ",")
	}
	return utils.FlattenStringSlice(&servers)
}

func expandFirewallPrivateIpRange(input []interface{}) map[string]*string {
	if len(input) == 0 {
		return nil
	}

	rangeSlice := *utils.ExpandStringSlice(input)
	if len(rangeSlice) == 0 {
		return nil
	}

	// Swagger issue asking finalize these properties: https://github.com/Azure/azure-rest-api-specs/issues/10015
	return map[string]*string{
		"Network.SNAT.PrivateRanges": utils.String(strings.Join(rangeSlice, ",")),
	}
}

func flattenFirewallPrivateIpRange(input map[string]*string) []interface{} {
	if len(input) == 0 {
		return nil
	}

	rangeSlice := []string{}
	if privateIpRanges := input["Network.SNAT.PrivateRanges"]; privateIpRanges != nil {
		rangeSlice = strings.Split(*privateIpRanges, ",")
	}
	return utils.FlattenStringSlice(&rangeSlice)
}

func expandFirewallVirtualHubSetting(existing network.AzureFirewall, input []interface{}) (vhub *network.SubResource, ipAddresses *network.HubIPAddresses, ok bool) {
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
	var addresses *[]network.AzureFirewallPublicIPAddress
	if prop := existing.AzureFirewallPropertiesFormat; prop != nil {
		if ipaddress := prop.HubIPAddresses; ipaddress != nil {
			if pips := ipaddress.PublicIPs; pips != nil {
				if count := pips.Count; count != nil {
					oldCount := int(*count)
					addresses = pips.Addresses

					// In case of scale down, keep the first new "Count" addresses.
					if oldCount > newCount {
						keptAddresses := make([]network.AzureFirewallPublicIPAddress, newCount)
						for i := 0; i < newCount; i++ {
							keptAddresses[i] = (*addresses)[i]
						}
						addresses = &keptAddresses
					}
				}
			}
		}
	}

	vhub = &network.SubResource{ID: utils.String(b["virtual_hub_id"].(string))}
	ipAddresses = &network.HubIPAddresses{
		PublicIPs: &network.HubPublicIPAddresses{
			Count:     utils.Int32(int32(b["public_ip_count"].(int))),
			Addresses: addresses,
		},
	}

	return vhub, ipAddresses, true
}

func flattenFirewallVirtualHubSetting(props *network.AzureFirewallPropertiesFormat) []interface{} {
	if props.VirtualHub == nil {
		return nil
	}

	var vhubId string
	if props.VirtualHub.ID != nil {
		vhubId = *props.VirtualHub.ID
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

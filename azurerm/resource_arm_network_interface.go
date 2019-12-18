package azurerm

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/state"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var networkInterfaceResourceName = "azurerm_network_interface"

func resourceArmNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkInterfaceCreateUpdate,
		Read:   resourceArmNetworkInterfaceRead,
		Update: resourceArmNetworkInterfaceCreateUpdate,
		Delete: resourceArmNetworkInterfaceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"network_security_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"mac_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.MACAddress,
			},

			"virtual_machine_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},

						"subnet_id": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     azure.ValidateResourceID,
						},

						"private_ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_version": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(network.IPv4),
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPv4),
								string(network.IPv6),
							}, false),
						},

						//TODO: should this be renamed to private_ip_address_allocation_method or private_ip_allocation_method ?
						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Dynamic),
								string(network.Static),
							}, true),
							StateFunc:        state.IgnoreCase,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"public_ip_address_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},

						"application_gateway_backend_address_pools_ids": {
							Type:       schema.TypeSet,
							Optional:   true,
							Computed:   true,
							Deprecated: "This field has been deprecated in favour of the `azurerm_network_interface_application_gateway_backend_address_pool_association` resource.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
							Set: schema.HashString,
						},

						"load_balancer_backend_address_pools_ids": {
							Type:       schema.TypeSet,
							Optional:   true,
							Computed:   true,
							Deprecated: "This field has been deprecated in favour of the `azurerm_network_interface_backend_address_pool_association` resource.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
							Set: schema.HashString,
						},

						"load_balancer_inbound_nat_rules_ids": {
							Type:       schema.TypeSet,
							Optional:   true,
							Computed:   true,
							Deprecated: "This field has been deprecated in favour of the `azurerm_network_interface_nat_rule_association` resource.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
							Set: schema.HashString,
						},

						"application_security_group_ids": {
							Type:       schema.TypeSet,
							Optional:   true,
							Computed:   true,
							Deprecated: "This field has been deprecated in favour of the `azurerm_network_interface_application_security_group_association` resource.",
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
							Set: schema.HashString,
						},

						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
				Set: schema.HashString,
			},

			"internal_dns_name_label": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"applied_dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.NoEmptyStrings,
				},
				Set: schema.HashString,
			},

			"internal_fqdn": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "This field has been removed by Azure",
			},

			/**
			 * As of 2018-01-06: AN (aka. SR-IOV) on Azure is GA on Windows and Linux.
			 *
			 * Refer to: https://azure.microsoft.com/en-us/blog/maximize-your-vm-s-performance-with-accelerated-networking-now-generally-available-for-both-windows-and-linux/
			 *
			 * Refer to: https://docs.microsoft.com/en-us/azure/virtual-network/create-vm-accelerated-networking-cli
			 * For details, VM configuration and caveats.
			 */
			"enable_accelerated_networking": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enable_ip_forwarding": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// todo consider removing this one day as it is exposed in `private_ip_addresses.0`
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetworkInterfaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.InterfacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Network Interface creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Network Interface %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_interface", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	enableIpForwarding := d.Get("enable_ip_forwarding").(bool)
	enableAcceleratedNetworking := d.Get("enable_accelerated_networking").(bool)
	t := d.Get("tags").(map[string]interface{})

	properties := network.InterfacePropertiesFormat{
		EnableIPForwarding:          &enableIpForwarding,
		EnableAcceleratedNetworking: &enableAcceleratedNetworking,
	}

	locks.ByName(name, networkInterfaceResourceName)
	defer locks.UnlockByName(name, networkInterfaceResourceName)

	if v, ok := d.GetOk("network_security_group_id"); ok {
		nsgId := v.(string)
		properties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: &nsgId,
		}

		parsedNsgID, err := azure.ParseAzureResourceID(nsgId)
		if err != nil {
			return fmt.Errorf("Error parsing Network Security Group ID %q: %+v", nsgId, err)
		}

		networkSecurityGroupName := parsedNsgID.Path["networkSecurityGroups"]

		locks.ByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer locks.UnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	dns, hasDns := d.GetOk("dns_servers")
	nameLabel, hasNameLabel := d.GetOk("internal_dns_name_label")
	if hasDns || hasNameLabel {
		ifaceDnsSettings := network.InterfaceDNSSettings{}

		if hasDns {
			var dnsServers []string
			dns := dns.(*schema.Set).List()
			for _, v := range dns {
				str := v.(string)
				dnsServers = append(dnsServers, str)
			}
			ifaceDnsSettings.DNSServers = &dnsServers
		}

		if hasNameLabel {
			name_label := nameLabel.(string)
			ifaceDnsSettings.InternalDNSNameLabel = &name_label
		}

		properties.DNSSettings = &ifaceDnsSettings
	}

	ipConfigs, subnetnToLock, vnnToLock, sgErr := expandAzureRmNetworkInterfaceIpConfigurations(d)
	if sgErr != nil {
		return fmt.Errorf("Error Building list of Network Interface IP Configurations: %+v", sgErr)
	}

	locks.MultipleByName(subnetnToLock, subnetResourceName)
	defer locks.UnlockMultipleByName(subnetnToLock, subnetResourceName)

	locks.MultipleByName(vnnToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnnToLock, virtualNetworkResourceName)

	if len(ipConfigs) > 0 {
		properties.IPConfigurations = &ipConfigs
	}

	iface := network.Interface{
		Name:                      &name,
		Location:                  &location,
		InterfacePropertiesFormat: &properties,
		Tags:                      tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, iface)
	if err != nil {
		return err
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read NIC %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNetworkInterfaceRead(d, meta)
}

func resourceArmNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	resp, err := client.Get(ctx, resGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.InterfacePropertiesFormat; props != nil {
		d.Set("mac_address", props.MacAddress)
		addresses := make([]interface{}, 0)
		if configs := props.IPConfigurations; configs != nil {
			for i, config := range *props.IPConfigurations {
				if ipProps := config.InterfaceIPConfigurationPropertiesFormat; ipProps != nil {
					if v := ipProps.PrivateIPAddress; v != nil {
						if i == 0 {
							d.Set("private_ip_address", v)
						}
						addresses = append(addresses, *v)
					}
				}
			}
		}
		if err := d.Set("private_ip_addresses", addresses); err != nil {
			return err
		}

		if props.IPConfigurations != nil {
			configs := flattenNetworkInterfaceIPConfigurations(props.IPConfigurations)
			if err := d.Set("ip_configuration", configs); err != nil {
				return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
			}
		}

		if vm := props.VirtualMachine; vm != nil {
			d.Set("virtual_machine_id", vm.ID)
		}

		var appliedDNSServers []string
		var dnsServers []string
		if dnsSettings := props.DNSSettings; dnsSettings != nil {
			if s := dnsSettings.AppliedDNSServers; s != nil {
				appliedDNSServers = *s
			}

			if s := dnsSettings.DNSServers; s != nil {
				dnsServers = *s
			}

			d.Set("internal_fqdn", dnsSettings.InternalFqdn)
			d.Set("internal_dns_name_label", dnsSettings.InternalDNSNameLabel)
		}

		d.Set("applied_dns_servers", appliedDNSServers)
		d.Set("dns_servers", dnsServers)

		if nsg := props.NetworkSecurityGroup; nsg != nil {
			d.Set("network_security_group_id", nsg.ID)
		} else {
			d.Set("network_security_group_id", "")
		}

		d.Set("enable_ip_forwarding", resp.EnableIPForwarding)
		d.Set("enable_accelerated_networking", resp.EnableAcceleratedNetworking)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.InterfacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	locks.ByName(name, networkInterfaceResourceName)
	defer locks.UnlockByName(name, networkInterfaceResourceName)

	if v, ok := d.GetOk("network_security_group_id"); ok {
		networkSecurityGroupId := v.(string)
		parsedNsgID, err := azure.ParseAzureResourceID(networkSecurityGroupId)
		if err != nil {
			return fmt.Errorf("Error parsing Network Security Group ID %q: %+v", networkSecurityGroupId, err)
		}

		networkSecurityGroupName := parsedNsgID.Path["networkSecurityGroups"]

		locks.ByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer locks.UnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	configs := d.Get("ip_configuration").([]interface{})
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		subnet_id := data["subnet_id"].(string)
		if subnet_id != "" {
			subnetId, err2 := azure.ParseAzureResourceID(subnet_id)
			if err2 != nil {
				return err2
			}
			subnetName := subnetId.Path["subnets"]
			if !sliceContainsValue(subnetNamesToLock, subnetName) {
				subnetNamesToLock = append(subnetNamesToLock, subnetName)
			}

			virtualNetworkName := subnetId.Path["virtualNetworks"]
			if !sliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
				virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
			}
		}
	}

	locks.MultipleByName(&subnetNamesToLock, subnetResourceName)
	defer locks.UnlockMultipleByName(&subnetNamesToLock, subnetResourceName)

	locks.MultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNamesToLock, virtualNetworkResourceName)

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Network Interface %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return err
}

func flattenNetworkInterfaceIPConfigurations(ipConfigs *[]network.InterfaceIPConfiguration) []interface{} {
	result := make([]interface{}, 0, len(*ipConfigs))
	for _, ipConfig := range *ipConfigs {
		niIPConfig := make(map[string]interface{})

		props := ipConfig.InterfaceIPConfigurationPropertiesFormat

		niIPConfig["name"] = *ipConfig.Name

		if props.Subnet != nil && props.Subnet.ID != nil {
			niIPConfig["subnet_id"] = *props.Subnet.ID
		}

		niIPConfig["private_ip_address_allocation"] = strings.ToLower(string(props.PrivateIPAllocationMethod))

		if props.PrivateIPAddress != nil {
			niIPConfig["private_ip_address"] = *props.PrivateIPAddress
		}

		if props.PrivateIPAddressVersion != "" {
			niIPConfig["private_ip_address_version"] = string(props.PrivateIPAddressVersion)
		}

		if props.PublicIPAddress != nil {
			niIPConfig["public_ip_address_id"] = *props.PublicIPAddress.ID
		}

		if props.Primary != nil {
			niIPConfig["primary"] = *props.Primary
		}

		var poolsAG []interface{}
		if props.ApplicationGatewayBackendAddressPools != nil {
			for _, pool := range *props.ApplicationGatewayBackendAddressPools {
				poolsAG = append(poolsAG, *pool.ID)
			}
		}
		niIPConfig["application_gateway_backend_address_pools_ids"] = schema.NewSet(schema.HashString, poolsAG)

		var pools []interface{}
		if props.LoadBalancerBackendAddressPools != nil {
			for _, pool := range *props.LoadBalancerBackendAddressPools {
				pools = append(pools, *pool.ID)
			}
		}
		niIPConfig["load_balancer_backend_address_pools_ids"] = schema.NewSet(schema.HashString, pools)

		var rules []interface{}
		if props.LoadBalancerInboundNatRules != nil {
			for _, rule := range *props.LoadBalancerInboundNatRules {
				rules = append(rules, *rule.ID)
			}
		}
		niIPConfig["load_balancer_inbound_nat_rules_ids"] = schema.NewSet(schema.HashString, rules)

		securityGroups := make([]interface{}, 0)
		if sgs := props.ApplicationSecurityGroups; sgs != nil {
			for _, sg := range *sgs {
				securityGroups = append(securityGroups, *sg.ID)
			}
		}
		niIPConfig["application_security_group_ids"] = schema.NewSet(schema.HashString, securityGroups)

		result = append(result, niIPConfig)
	}
	return result
}

func expandAzureRmNetworkInterfaceIpConfigurations(d *schema.ResourceData) ([]network.InterfaceIPConfiguration, *[]string, *[]string, error) {
	configs := d.Get("ip_configuration").([]interface{})
	ipConfigs := make([]network.InterfaceIPConfiguration, 0, len(configs))
	subnetNamesToLock := make([]string, 0)
	virtualNetworkNamesToLock := make([]string, 0)

	for _, configRaw := range configs {
		data := configRaw.(map[string]interface{})

		subnet_id := data["subnet_id"].(string)
		private_ip_allocation_method := data["private_ip_address_allocation"].(string)
		private_ip_address_version := network.IPVersion(data["private_ip_address_version"].(string))

		allocationMethod := network.IPAllocationMethod(private_ip_allocation_method)
		properties := network.InterfaceIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: allocationMethod,
			PrivateIPAddressVersion:   private_ip_address_version,
		}

		if private_ip_address_version == network.IPv4 && subnet_id == "" {
			return nil, nil, nil, fmt.Errorf("A Subnet ID must be specified for an IPv4 Network Interface.")
		}

		if subnet_id != "" {
			properties.Subnet = &network.Subnet{
				ID: &subnet_id,
			}

			subnetId, err := azure.ParseAzureResourceID(subnet_id)
			if err != nil {
				return []network.InterfaceIPConfiguration{}, nil, nil, err
			}

			subnetName := subnetId.Path["subnets"]
			virtualNetworkName := subnetId.Path["virtualNetworks"]

			if !sliceContainsValue(subnetNamesToLock, subnetName) {
				subnetNamesToLock = append(subnetNamesToLock, subnetName)
			}

			if !sliceContainsValue(virtualNetworkNamesToLock, virtualNetworkName) {
				virtualNetworkNamesToLock = append(virtualNetworkNamesToLock, virtualNetworkName)
			}
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &network.PublicIPAddress{
				ID: &v,
			}
		}

		if v, ok := data["primary"]; ok {
			b := v.(bool)
			properties.Primary = &b
		}

		if v, ok := data["application_gateway_backend_address_pools_ids"]; ok {
			var ids []network.ApplicationGatewayBackendAddressPool
			pools := v.(*schema.Set).List()
			for _, p := range pools {
				pool_id := p.(string)
				id := network.ApplicationGatewayBackendAddressPool{
					ID: &pool_id,
				}

				ids = append(ids, id)
			}

			properties.ApplicationGatewayBackendAddressPools = &ids
		}

		if v, ok := data["load_balancer_backend_address_pools_ids"]; ok {
			var ids []network.BackendAddressPool
			pools := v.(*schema.Set).List()
			for _, p := range pools {
				pool_id := p.(string)
				id := network.BackendAddressPool{
					ID: &pool_id,
				}

				ids = append(ids, id)
			}

			properties.LoadBalancerBackendAddressPools = &ids
		}

		if v, ok := data["load_balancer_inbound_nat_rules_ids"]; ok {
			var natRules []network.InboundNatRule
			rules := v.(*schema.Set).List()
			for _, r := range rules {
				rule_id := r.(string)
				rule := network.InboundNatRule{
					ID: &rule_id,
				}

				natRules = append(natRules, rule)
			}

			properties.LoadBalancerInboundNatRules = &natRules
		}

		if v, ok := data["application_security_group_ids"]; ok {
			var securityGroups []network.ApplicationSecurityGroup
			rules := v.(*schema.Set).List()
			for _, r := range rules {
				groupId := r.(string)
				group := network.ApplicationSecurityGroup{
					ID: &groupId,
				}

				securityGroups = append(securityGroups, group)
			}

			properties.ApplicationSecurityGroups = &securityGroups
		}

		name := data["name"].(string)
		ipConfig := network.InterfaceIPConfiguration{
			Name:                                     &name,
			InterfaceIPConfigurationPropertiesFormat: &properties,
		}

		ipConfigs = append(ipConfigs, ipConfig)
	}

	// if we've got multiple IP Configurations - one must be designated Primary
	if len(ipConfigs) > 1 {
		hasPrimary := false
		for _, config := range ipConfigs {
			if config.Primary != nil && *config.Primary {
				hasPrimary = true
				break
			}
		}

		if !hasPrimary {
			return nil, nil, nil, fmt.Errorf("If multiple `ip_configurations` are specified - one must be designated as `primary`.")
		}
	}

	return ipConfigs, &subnetNamesToLock, &virtualNetworkNamesToLock, nil
}

func sliceContainsValue(input []string, value string) bool {
	for _, v := range input {
		if v == value {
			return true
		}
	}

	return false
}

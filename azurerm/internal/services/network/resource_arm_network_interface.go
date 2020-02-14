package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
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
		Create: resourceArmNetworkInterfaceCreate,
		Read:   resourceArmNetworkInterfaceRead,
		Update: resourceArmNetworkInterfaceUpdate,
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

			// NOTE: does this want it's own association resource?
			"network_security_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Set: schema.HashString,
			},

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

			"internal_dns_name_label": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"tags": tags.Schema(),

			// Computed
			"applied_dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

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

			"virtual_machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Network Interface %q (Resource Group %q): %s", name, resourceGroup, err)
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
		dnsSettings := network.InterfaceDNSSettings{}

		if hasDns {
			dnsRaw := dns.(*schema.Set).List()
			dns := expandNetworkInterfaceDnsServers(dnsRaw)
			dnsSettings.DNSServers = &dns
		}

		if hasNameLabel {
			name_label := nameLabel.(string)
			dnsSettings.InternalDNSNameLabel = &name_label
		}

		properties.DNSSettings = &dnsSettings
	}

	ipConfigsRaw := d.Get("ip_configuration").([]interface{})
	ipConfigs, err := expandNetworkInterfaceIPConfigurations(ipConfigsRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `ip_configuration`: %+v", err)
	}
	lockingDetails, err := determineResourcesToLockFromIPConfiguration(&ipConfigs)
	if err != nil {
		return fmt.Errorf("Error determing locking details: %+v", err)
	}

	lockingDetails.lock()
	defer lockingDetails.unlock()

	if len(ipConfigs) > 0 {
		properties.IPConfigurations = &ipConfigs
	}

	iface := network.Interface{
		Name:                      utils.String(name),
		Location:                  utils.String(location),
		InterfacePropertiesFormat: &properties,
		Tags:                      tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, iface)
	if err != nil {
		return fmt.Errorf("Error creating Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): ID was nil", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmNetworkInterfaceRead(d, meta)
}

func resourceArmNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	locks.ByName(name, networkInterfaceResourceName)
	defer locks.UnlockByName(name, networkInterfaceResourceName)

	location := azure.NormalizeLocation(d.Get("location").(string))
	update := network.Interface{
		Name:     utils.String(name),
		Location: utils.String(location),
	}

	if d.HasChange("dns_servers") {
		if update.InterfacePropertiesFormat == nil {
			update.InterfacePropertiesFormat = &network.InterfacePropertiesFormat{}
		}
		if update.InterfacePropertiesFormat.DNSSettings == nil {
			update.InterfacePropertiesFormat.DNSSettings = &network.InterfaceDNSSettings{}
		}

		dnsServersRaw := d.Get("dns_servers").(*schema.Set).List()
		dnsServers := expandNetworkInterfaceDnsServers(dnsServersRaw)

		update.InterfacePropertiesFormat.DNSSettings.DNSServers = &dnsServers
	}

	if d.HasChange("enable_accelerated_networking") {
		if update.InterfacePropertiesFormat == nil {
			update.InterfacePropertiesFormat = &network.InterfacePropertiesFormat{}
		}

		update.InterfacePropertiesFormat.EnableAcceleratedNetworking = utils.Bool(d.Get("enable_accelerated_networking").(bool))
	}

	if d.HasChange("enable_ip_forwarding") {
		if update.InterfacePropertiesFormat == nil {
			update.InterfacePropertiesFormat = &network.InterfacePropertiesFormat{}
		}

		update.InterfacePropertiesFormat.EnableIPForwarding = utils.Bool(d.Get("enable_ip_forwarding").(bool))
	}

	if d.HasChange("internal_dns_name_label") {
		if update.InterfacePropertiesFormat == nil {
			update.InterfacePropertiesFormat = &network.InterfacePropertiesFormat{}
		}
		if update.InterfacePropertiesFormat.DNSSettings == nil {
			update.InterfacePropertiesFormat.DNSSettings = &network.InterfaceDNSSettings{}
		}

		update.InterfacePropertiesFormat.DNSSettings.InternalDNSNameLabel = utils.String(d.Get("internal_dns_name_label").(string))
	}

	if d.HasChange("ip_configuration") {
		if update.InterfacePropertiesFormat == nil {
			update.InterfacePropertiesFormat = &network.InterfacePropertiesFormat{}
		}

		ipConfigsRaw := d.Get("ip_configuration").([]interface{})
		ipConfigs, err := expandNetworkInterfaceIPConfigurations(ipConfigsRaw)
		if err != nil {
			return fmt.Errorf("Error expanding `ip_configuration`: %+v", err)
		}
		lockingDetails, err := determineResourcesToLockFromIPConfiguration(&ipConfigs)
		if err != nil {
			return fmt.Errorf("Error determing locking details: %+v", err)
		}

		lockingDetails.lock()
		defer lockingDetails.unlock()

		update.InterfacePropertiesFormat.IPConfigurations = &ipConfigs
	}

	if d.HasChange("network_security_group_id") {
		if update.InterfacePropertiesFormat == nil {
			update.InterfacePropertiesFormat = &network.InterfacePropertiesFormat{}
		}

		update.InterfacePropertiesFormat.NetworkSecurityGroup = &network.SecurityGroup{
			ID: utils.String(d.Get("network_security_group_id").(string)),
		}
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, update)
	if err != nil {
		return fmt.Errorf("Error updating Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func resourceArmNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.InterfacePropertiesFormat; props != nil {
		primaryPrivateIPAddress := ""
		privateIPAddresses := make([]interface{}, 0)
		if configs := props.IPConfigurations; configs != nil {
			for i, config := range *props.IPConfigurations {
				if ipProps := config.InterfaceIPConfigurationPropertiesFormat; ipProps != nil {
					v := ipProps.PrivateIPAddress
					if v == nil {
						continue
					}

					if i == 0 {
						primaryPrivateIPAddress = *v
					}

					privateIPAddresses = append(privateIPAddresses, *v)
				}
			}
		}

		appliedDNSServers := make([]string, 0)
		dnsServers := make([]string, 0)
		internalDnsNameLabel := ""
		if dnsSettings := props.DNSSettings; dnsSettings != nil {
			appliedDNSServers = flattenNetworkInterfaceDnsServers(dnsSettings.AppliedDNSServers)
			dnsServers = flattenNetworkInterfaceDnsServers(dnsSettings.DNSServers)

			if dnsSettings.InternalDNSNameLabel != nil {
				internalDnsNameLabel = *dnsSettings.InternalDNSNameLabel
			}
		}

		networkSecurityGroupId := ""
		if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.ID != nil {
			networkSecurityGroupId = *props.NetworkSecurityGroup.ID
		}
		virtualMachineId := ""
		if props.VirtualMachine != nil && props.VirtualMachine.ID != nil {
			virtualMachineId = *props.VirtualMachine.ID
		}

		if err := d.Set("applied_dns_servers", appliedDNSServers); err != nil {
			return fmt.Errorf("Error setting `applied_dns_servers`: %+v", err)
		}

		if err := d.Set("dns_servers", dnsServers); err != nil {
			return fmt.Errorf("Error setting `applied_dns_servers`: %+v", err)
		}

		d.Set("enable_ip_forwarding", resp.EnableIPForwarding)
		d.Set("enable_accelerated_networking", resp.EnableAcceleratedNetworking)
		d.Set("internal_dns_name_label", internalDnsNameLabel)
		d.Set("mac_address", props.MacAddress)
		d.Set("network_security_group_id", networkSecurityGroupId)
		d.Set("private_ip_address", primaryPrivateIPAddress)
		d.Set("virtual_machine_id", virtualMachineId)

		if err := d.Set("ip_configuration", flattenNetworkInterfaceIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}

		if err := d.Set("private_ip_addresses", privateIPAddresses); err != nil {
			return fmt.Errorf("Error setting `private_ip_addresses`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkInterfaces"]

	locks.ByName(name, networkInterfaceResourceName)
	defer locks.UnlockByName(name, networkInterfaceResourceName)

	existing, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			log.Printf("[DEBUG] Network Interface %q was not found in Resource Group %q - removing from state", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if existing.InterfacePropertiesFormat == nil {
		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): `properties` was nil", name, resourceGroup)
	}
	props := *existing.InterfacePropertiesFormat

	if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.ID != nil {
		networkSecurityGroupId := *props.NetworkSecurityGroup.ID
		parsedNsgID, err := azure.ParseAzureResourceID(networkSecurityGroupId)
		if err != nil {
			return fmt.Errorf("Error parsing Network Security Group ID %q: %+v", networkSecurityGroupId, err)
		}

		networkSecurityGroupName := parsedNsgID.Path["networkSecurityGroups"]

		locks.ByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer locks.UnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	lockingDetails, err := determineResourcesToLockFromIPConfiguration(props.IPConfigurations)
	if err != nil {
		return fmt.Errorf("Error determing locking details: %+v", err)
	}

	lockingDetails.lock()
	defer lockingDetails.unlock()

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandNetworkInterfaceIPConfigurations(input []interface{}) ([]network.InterfaceIPConfiguration, error) {
	ipConfigs := make([]network.InterfaceIPConfiguration, 0)

	for _, configRaw := range input {
		data := configRaw.(map[string]interface{})

		subnetId := data["subnet_id"].(string)
		privateIpAllocationMethod := data["private_ip_address_allocation"].(string)
		privateIpAddressVersion := network.IPVersion(data["private_ip_address_version"].(string))

		allocationMethod := network.IPAllocationMethod(privateIpAllocationMethod)
		properties := network.InterfaceIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: allocationMethod,
			PrivateIPAddressVersion:   privateIpAddressVersion,
		}

		if privateIpAddressVersion == network.IPv4 && subnetId == "" {
			return nil, fmt.Errorf("A Subnet ID must be specified for an IPv4 Network Interface.")
		}

		if subnetId != "" {
			properties.Subnet = &network.Subnet{
				ID: &subnetId,
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
			properties.Primary = utils.Bool(v.(bool))
		}

		if v, ok := data["load_balancer_backend_address_pools_ids"]; ok {
			var ids []network.BackendAddressPool
			pools := v.(*schema.Set).List()
			for _, p := range pools {
				poolId := p.(string)
				id := network.BackendAddressPool{
					ID: &poolId,
				}

				ids = append(ids, id)
			}

			properties.LoadBalancerBackendAddressPools = &ids
		}

		if v, ok := data["load_balancer_inbound_nat_rules_ids"]; ok {
			var natRules []network.InboundNatRule
			rules := v.(*schema.Set).List()
			for _, r := range rules {
				ruleId := r.(string)
				rule := network.InboundNatRule{
					ID: &ruleId,
				}

				natRules = append(natRules, rule)
			}

			properties.LoadBalancerInboundNatRules = &natRules
		}

		name := data["name"].(string)
		ipConfigs = append(ipConfigs, network.InterfaceIPConfiguration{
			Name:                                     &name,
			InterfaceIPConfigurationPropertiesFormat: &properties,
		})
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
			return nil, fmt.Errorf("If multiple `ip_configurations` are specified - one must be designated as `primary`.")
		}
	}

	return ipConfigs, nil
}

func flattenNetworkInterfaceIPConfigurations(input *[]network.InterfaceIPConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, ipConfig := range *input {
		props := ipConfig.InterfaceIPConfigurationPropertiesFormat

		name := ""
		if ipConfig.Name != nil {
			name = *ipConfig.Name
		}

		subnetId := ""
		if props.Subnet != nil && props.Subnet.ID != nil {
			subnetId = *props.Subnet.ID
		}

		privateIPAddress := ""
		if props.PrivateIPAddress != nil {
			privateIPAddress = *props.PrivateIPAddress
		}

		privateIPAddressVersion := ""
		if props.PrivateIPAddressVersion != "" {
			privateIPAddressVersion = string(props.PrivateIPAddressVersion)
		}

		publicIPAddressId := ""
		if props.PublicIPAddress != nil && props.PublicIPAddress.ID != nil {
			publicIPAddressId = *props.PublicIPAddress.ID
		}

		var lbBackendAddressPools []interface{}
		if props.LoadBalancerBackendAddressPools != nil {
			for _, pool := range *props.LoadBalancerBackendAddressPools {
				lbBackendAddressPools = append(lbBackendAddressPools, *pool.ID)
			}
		}

		var lbInboundNatRules []interface{}
		if props.LoadBalancerInboundNatRules != nil {
			for _, rule := range *props.LoadBalancerInboundNatRules {
				lbInboundNatRules = append(lbInboundNatRules, *rule.ID)
			}
		}

		primary := false
		if props.Primary != nil {
			primary = *props.Primary
		}

		result = append(result, map[string]interface{}{
			"load_balancer_backend_address_pools_ids": schema.NewSet(schema.HashString, lbBackendAddressPools),
			"load_balancer_inbound_nat_rules_ids":     schema.NewSet(schema.HashString, lbInboundNatRules),
			"name":                                    name,
			"primary":                                 primary,
			"private_ip_address":                      privateIPAddress,
			"private_ip_address_allocation":           string(props.PrivateIPAllocationMethod),
			"private_ip_address_version":              privateIPAddressVersion,
			"public_ip_address_id":                    publicIPAddressId,
			"subnet_id":                               subnetId,
		})
	}
	return result
}

func expandNetworkInterfaceDnsServers(input []interface{}) []string {
	dnsServers := make([]string, 0)
	for _, v := range input {
		dnsServers = append(dnsServers, v.(string))
	}
	return dnsServers
}

func flattenNetworkInterfaceDnsServers(input *[]string) []string {
	output := make([]string, 0)
	if input == nil {
		return output
	}

	for _, v := range *input {
		output = append(output, v)
	}
	return output
}

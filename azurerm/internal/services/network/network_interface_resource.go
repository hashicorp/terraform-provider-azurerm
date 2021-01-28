package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/state"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var networkInterfaceResourceName = "azurerm_network_interface"

func resourceNetworkInterface() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkInterfaceCreate,
		Read:   resourceNetworkInterfaceRead,
		Update: resourceNetworkInterfaceUpdate,
		Delete: resourceNetworkInterfaceDelete,

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

						"primary": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
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

			"internal_domain_name_suffix": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceNetworkInterfaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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

	dns, hasDns := d.GetOk("dns_servers")
	nameLabel, hasNameLabel := d.GetOk("internal_dns_name_label")
	if hasDns || hasNameLabel {
		dnsSettings := network.InterfaceDNSSettings{}

		if hasDns {
			dnsRaw := dns.([]interface{})
			dns := expandNetworkInterfaceDnsServers(dnsRaw)
			dnsSettings.DNSServers = &dns
		}

		if hasNameLabel {
			dnsSettings.InternalDNSNameLabel = utils.String(nameLabel.(string))
		}

		properties.DNSSettings = &dnsSettings
	}

	ipConfigsRaw := d.Get("ip_configuration").([]interface{})
	ipConfigs, err := expandNetworkInterfaceIPConfigurations(ipConfigsRaw)
	if err != nil {
		return fmt.Errorf("Error expanding `ip_configuration`: %+v", err)
	}
	lockingDetails, err := determineResourcesToLockFromIPConfiguration(ipConfigs)
	if err != nil {
		return fmt.Errorf("Error determining locking details: %+v", err)
	}

	lockingDetails.lock()
	defer lockingDetails.unlock()

	if len(*ipConfigs) > 0 {
		properties.IPConfigurations = ipConfigs
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

	return resourceNetworkInterfaceRead(d, meta)
}

func resourceNetworkInterfaceUpdate(d *schema.ResourceData, meta interface{}) error {
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

	// first get the existing one so that we can pull things as needed
	existing, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if existing.InterfacePropertiesFormat == nil {
		return fmt.Errorf("Error retrieving Network Interface %q (Resource Group %q): `properties` was nil", name, resourceGroup)
	}

	// then pull out things we need to lock on
	info := parseFieldsFromNetworkInterface(*existing.InterfacePropertiesFormat)

	location := azure.NormalizeLocation(d.Get("location").(string))
	update := network.Interface{
		Name:     utils.String(name),
		Location: utils.String(location),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			EnableAcceleratedNetworking: utils.Bool(d.Get("enable_accelerated_networking").(bool)),
			DNSSettings:                 &network.InterfaceDNSSettings{},
		},
	}

	if d.HasChange("dns_servers") {
		dnsServersRaw := d.Get("dns_servers").([]interface{})
		dnsServers := expandNetworkInterfaceDnsServers(dnsServersRaw)

		update.InterfacePropertiesFormat.DNSSettings.DNSServers = &dnsServers
	} else {
		update.InterfacePropertiesFormat.DNSSettings.DNSServers = existing.InterfacePropertiesFormat.DNSSettings.DNSServers
	}

	if d.HasChange("enable_ip_forwarding") {
		update.InterfacePropertiesFormat.EnableIPForwarding = utils.Bool(d.Get("enable_ip_forwarding").(bool))
	} else {
		update.InterfacePropertiesFormat.EnableIPForwarding = existing.InterfacePropertiesFormat.EnableIPForwarding
	}

	if d.HasChange("internal_dns_name_label") {
		update.InterfacePropertiesFormat.DNSSettings.InternalDNSNameLabel = utils.String(d.Get("internal_dns_name_label").(string))
	} else {
		update.InterfacePropertiesFormat.DNSSettings.InternalDNSNameLabel = existing.InterfacePropertiesFormat.DNSSettings.InternalDNSNameLabel
	}

	if d.HasChange("ip_configuration") {
		ipConfigsRaw := d.Get("ip_configuration").([]interface{})
		ipConfigs, err := expandNetworkInterfaceIPConfigurations(ipConfigsRaw)
		if err != nil {
			return fmt.Errorf("Error expanding `ip_configuration`: %+v", err)
		}
		lockingDetails, err := determineResourcesToLockFromIPConfiguration(ipConfigs)
		if err != nil {
			return fmt.Errorf("Error determining locking details: %+v", err)
		}

		lockingDetails.lock()
		defer lockingDetails.unlock()

		// then map the fields managed in other resources back
		ipConfigs = mapFieldsToNetworkInterface(ipConfigs, info)

		update.InterfacePropertiesFormat.IPConfigurations = ipConfigs
	} else {
		update.InterfacePropertiesFormat.IPConfigurations = existing.InterfacePropertiesFormat.IPConfigurations
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	} else {
		update.Tags = existing.Tags
	}

	// this can be managed in another resource, so just port it over
	update.InterfacePropertiesFormat.NetworkSecurityGroup = existing.InterfacePropertiesFormat.NetworkSecurityGroup

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, update)
	if err != nil {
		return fmt.Errorf("Error updating Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Network Interface %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func resourceNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
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
		internalDomainNameSuffix := ""
		if dnsSettings := props.DNSSettings; dnsSettings != nil {
			appliedDNSServers = flattenNetworkInterfaceDnsServers(dnsSettings.AppliedDNSServers)
			dnsServers = flattenNetworkInterfaceDnsServers(dnsSettings.DNSServers)

			if dnsSettings.InternalDNSNameLabel != nil {
				internalDnsNameLabel = *dnsSettings.InternalDNSNameLabel
			}

			if dnsSettings.InternalDomainNameSuffix != nil {
				internalDomainNameSuffix = *dnsSettings.InternalDomainNameSuffix
			}
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
		d.Set("internal_domain_name_suffix", internalDomainNameSuffix)
		d.Set("mac_address", props.MacAddress)
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

func resourceNetworkInterfaceDelete(d *schema.ResourceData, meta interface{}) error {
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

	lockingDetails, err := determineResourcesToLockFromIPConfiguration(props.IPConfigurations)
	if err != nil {
		return fmt.Errorf("Error determining locking details: %+v", err)
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

func expandNetworkInterfaceIPConfigurations(input []interface{}) (*[]network.InterfaceIPConfiguration, error) {
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

	return &ipConfigs, nil
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

		primary := false
		if props.Primary != nil {
			primary = *props.Primary
		}

		result = append(result, map[string]interface{}{
			"name":                          name,
			"primary":                       primary,
			"private_ip_address":            privateIPAddress,
			"private_ip_address_allocation": string(props.PrivateIPAllocationMethod),
			"private_ip_address_version":    privateIPAddressVersion,
			"public_ip_address_id":          publicIPAddressId,
			"subnet_id":                     subnetId,
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
	if input == nil {
		return make([]string, 0)
	}

	return *input
}

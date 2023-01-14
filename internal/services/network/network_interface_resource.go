package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	lbvalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-05-01/network"
)

var networkInterfaceResourceName = "azurerm_network_interface"

func resourceNetworkInterface() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceCreate,
		Read:   resourceNetworkInterfaceRead,
		Update: resourceNetworkInterfaceUpdate,
		Delete: resourceNetworkInterfaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NetworkInterfaceID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": commonschema.LocationOptional(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"ip_configuration": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"subnet_id": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc:     azure.ValidateResourceID,
						},

						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(network.IPVersionIPv4),
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPVersionIPv4),
								string(network.IPVersionIPv6),
							}, false),
						},

						"private_ip_address_allocation": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPAllocationMethodDynamic),
								string(network.IPAllocationMethodStatic),
							}, false),
						},

						"public_ip_address_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},

						"primary": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Computed: true,
						},

						"gateway_load_balancer_frontend_ip_configuration_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: lbvalidate.LoadBalancerFrontendIpConfigurationID,
						},
					},
				},
			},

			// Optional
			"dns_servers": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_accelerated_networking": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_ip_forwarding": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"internal_dns_name_label": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"internal_domain_name_suffix": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),

			// Computed
			"applied_dns_servers": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"mac_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"private_ip_addresses": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"virtual_machine_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkInterfaceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewNetworkInterfaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_network_interface", id.ID())
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

	locks.ByName(id.Name, networkInterfaceResourceName)
	defer locks.UnlockByName(id.Name, networkInterfaceResourceName)

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
		return fmt.Errorf("expanding `ip_configuration`: %+v", err)
	}
	lockingDetails, err := determineResourcesToLockFromIPConfiguration(ipConfigs)
	if err != nil {
		return fmt.Errorf("determining locking details: %+v", err)
	}

	lockingDetails.lock()
	defer lockingDetails.unlock()

	if len(*ipConfigs) > 0 {
		properties.IPConfigurations = ipConfigs
	}

	iface := network.Interface{
		Name:                      utils.String(id.Name),
		ExtendedLocation:          expandEdgeZone(d.Get("edge_zone").(string)),
		Location:                  utils.String(location),
		InterfacePropertiesFormat: &properties,
		Tags:                      tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, iface)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNetworkInterfaceRead(d, meta)
}

func resourceNetworkInterfaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, networkInterfaceResourceName)
	defer locks.UnlockByName(id.Name, networkInterfaceResourceName)

	// first get the existing one so that we can pull things as needed
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.InterfacePropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	// then pull out things we need to lock on
	info := parseFieldsFromNetworkInterface(*existing.InterfacePropertiesFormat)

	location := azure.NormalizeLocation(d.Get("location").(string))
	update := network.Interface{
		Name:             utils.String(id.Name),
		ExtendedLocation: expandEdgeZone(d.Get("edge_zone").(string)),
		Location:         utils.String(location),
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
			return fmt.Errorf("expanding `ip_configuration`: %+v", err)
		}
		lockingDetails, err := determineResourcesToLockFromIPConfiguration(ipConfigs)
		if err != nil {
			return fmt.Errorf("determining locking details: %+v", err)
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, update)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	return nil
}

func resourceNetworkInterfaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("edge_zone", flattenEdgeZone(resp.ExtendedLocation))

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
			return fmt.Errorf("setting `applied_dns_servers`: %+v", err)
		}

		if err := d.Set("dns_servers", dnsServers); err != nil {
			return fmt.Errorf("setting `applied_dns_servers`: %+v", err)
		}

		d.Set("enable_ip_forwarding", resp.EnableIPForwarding)
		d.Set("enable_accelerated_networking", resp.EnableAcceleratedNetworking)
		d.Set("internal_dns_name_label", internalDnsNameLabel)
		d.Set("internal_domain_name_suffix", internalDomainNameSuffix)
		d.Set("mac_address", props.MacAddress)
		d.Set("private_ip_address", primaryPrivateIPAddress)
		d.Set("virtual_machine_id", virtualMachineId)

		if err := d.Set("ip_configuration", flattenNetworkInterfaceIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("setting `ip_configuration`: %+v", err)
		}

		if err := d.Set("private_ip_addresses", privateIPAddresses); err != nil {
			return fmt.Errorf("setting `private_ip_addresses`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceNetworkInterfaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.InterfacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.Name, networkInterfaceResourceName)
	defer locks.UnlockByName(id.Name, networkInterfaceResourceName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.InterfacePropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}
	props := *existing.InterfacePropertiesFormat

	lockingDetails, err := determineResourcesToLockFromIPConfiguration(props.IPConfigurations)
	if err != nil {
		return fmt.Errorf("determining locking details: %+v", err)
	}

	lockingDetails.lock()
	defer lockingDetails.unlock()

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
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

		if privateIpAddressVersion == network.IPVersionIPv4 && subnetId == "" {
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

		if v := data["gateway_load_balancer_frontend_ip_configuration_id"].(string); v != "" {
			properties.GatewayLoadBalancer = &network.SubResource{ID: &v}
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

		gatewayLBFrontendIPConfigId := ""
		if props.GatewayLoadBalancer != nil && props.GatewayLoadBalancer.ID != nil {
			gatewayLBFrontendIPConfigId = *props.GatewayLoadBalancer.ID
		}

		result = append(result, map[string]interface{}{
			"name":                          name,
			"primary":                       primary,
			"private_ip_address":            privateIPAddress,
			"private_ip_address_allocation": string(props.PrivateIPAllocationMethod),
			"private_ip_address_version":    privateIPAddressVersion,
			"public_ip_address_id":          publicIPAddressId,
			"subnet_id":                     subnetId,
			"gateway_load_balancer_frontend_ip_configuration_id": gatewayLBFrontendIPConfigId,
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

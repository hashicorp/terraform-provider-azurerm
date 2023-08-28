// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	lbvalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var networkInterfaceResourceName = "azurerm_network_interface"

func resourceNetworkInterface() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkInterfaceCreate,
		Read:   resourceNetworkInterfaceRead,
		Update: resourceNetworkInterfaceUpdate,
		Delete: resourceNetworkInterfaceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseNetworkInterfaceID(id)
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

			"location": commonschema.Location(),

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
							ValidateFunc:     commonids.ValidateSubnetID,
						},

						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Computed: true,
						},

						"private_ip_address_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(networkinterfaces.IPVersionIPvFour),
							ValidateFunc: validation.StringInSlice([]string{
								string(networkinterfaces.IPVersionIPvFour),
								string(networkinterfaces.IPVersionIPvSix),
							}, false),
						},

						"private_ip_address_allocation": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(networkinterfaces.IPAllocationMethodDynamic),
								string(networkinterfaces.IPAllocationMethodStatic),
							}, false),
						},

						"public_ip_address_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.PublicIpAddressID,
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

			"tags": commonschema.Tags(),

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
	client := meta.(*clients.Client).Network.NetworkInterfaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewNetworkInterfaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_network_interface", id.ID())
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	enableIpForwarding := d.Get("enable_ip_forwarding").(bool)
	enableAcceleratedNetworking := d.Get("enable_accelerated_networking").(bool)
	t := d.Get("tags").(map[string]interface{})

	properties := networkinterfaces.NetworkInterfacePropertiesFormat{
		EnableIPForwarding:          &enableIpForwarding,
		EnableAcceleratedNetworking: &enableAcceleratedNetworking,
	}

	locks.ByName(id.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.NetworkInterfaceName, networkInterfaceResourceName)

	dns, hasDns := d.GetOk("dns_servers")
	nameLabel, hasNameLabel := d.GetOk("internal_dns_name_label")
	if hasDns || hasNameLabel {
		dnsSettings := networkinterfaces.NetworkInterfaceDnsSettings{}

		if hasDns {
			dnsRaw := dns.([]interface{})
			dns := expandNetworkInterfaceDnsServers(dnsRaw)
			dnsSettings.DnsServers = &dns
		}

		if hasNameLabel {
			dnsSettings.InternalDnsNameLabel = utils.String(nameLabel.(string))
		}

		properties.DnsSettings = &dnsSettings
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

	iface := networkinterfaces.NetworkInterface{
		Name:             pointer.To(id.NetworkInterfaceName),
		ExtendedLocation: expandEdgeZoneModel(d.Get("edge_zone").(string)),
		Location:         utils.String(location),
		Properties:       &properties,
		Tags:             tags.Expand(t),
	}

	err = client.CreateOrUpdateThenPoll(ctx, id, iface)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNetworkInterfaceRead(d, meta)
}

func resourceNetworkInterfaceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseNetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.NetworkInterfaceName, networkInterfaceResourceName)

	// first get the existing one so that we can pull things as needed
	existing, err := client.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	// then pull out things we need to lock on
	info := parseFieldsFromNetworkInterface(*existing.Model.Properties)

	location := azure.NormalizeLocation(d.Get("location").(string))
	update := networkinterfaces.NetworkInterface{
		Name:             utils.String(id.NetworkInterfaceName),
		ExtendedLocation: expandEdgeZoneModel(d.Get("edge_zone").(string)),
		Location:         utils.String(location),
		Properties: &networkinterfaces.NetworkInterfacePropertiesFormat{
			EnableAcceleratedNetworking: utils.Bool(d.Get("enable_accelerated_networking").(bool)),
			DnsSettings:                 &networkinterfaces.NetworkInterfaceDnsSettings{},
		},
	}

	if d.HasChange("dns_servers") {
		dnsServersRaw := d.Get("dns_servers").([]interface{})
		dnsServers := expandNetworkInterfaceDnsServers(dnsServersRaw)

		update.Properties.DnsSettings.DnsServers = &dnsServers
	} else {
		update.Properties.DnsSettings.DnsServers = existing.Model.Properties.DnsSettings.DnsServers
	}

	if d.HasChange("enable_ip_forwarding") {
		update.Properties.EnableIPForwarding = utils.Bool(d.Get("enable_ip_forwarding").(bool))
	} else {
		update.Properties.EnableIPForwarding = existing.Model.Properties.EnableIPForwarding
	}

	if d.HasChange("internal_dns_name_label") {
		update.Properties.DnsSettings.InternalDnsNameLabel = utils.String(d.Get("internal_dns_name_label").(string))
	} else {
		update.Properties.DnsSettings.InternalDnsNameLabel = existing.Model.Properties.DnsSettings.InternalDnsNameLabel
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

		update.Properties.IPConfigurations = ipConfigs
	} else {
		update.Properties.IPConfigurations = existing.Model.Properties.IPConfigurations
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	} else {
		update.Tags = existing.Model.Tags
	}

	// this can be managed in another resource, so just port it over
	update.Properties.NetworkSecurityGroup = existing.Model.Properties.NetworkSecurityGroup

	err = client.CreateOrUpdateThenPoll(ctx, *id, update)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return nil
}

func resourceNetworkInterfaceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseNetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NetworkInterfaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZoneModel(model.ExtendedLocation))

		if props := model.Properties; props != nil {
			primaryPrivateIPAddress := ""
			privateIPAddresses := make([]interface{}, 0)
			if configs := props.IPConfigurations; configs != nil {
				for i, config := range *props.IPConfigurations {
					if ipProps := config.Properties; ipProps != nil {
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
			if dnsSettings := props.DnsSettings; dnsSettings != nil {
				appliedDNSServers = flattenNetworkInterfaceDnsServers(dnsSettings.AppliedDnsServers)
				dnsServers = flattenNetworkInterfaceDnsServers(dnsSettings.DnsServers)

				if dnsSettings.InternalDnsNameLabel != nil {
					internalDnsNameLabel = *dnsSettings.InternalDnsNameLabel
				}

				if dnsSettings.InternalDomainNameSuffix != nil {
					internalDomainNameSuffix = *dnsSettings.InternalDomainNameSuffix
				}
			}

			virtualMachineId := ""
			if props.VirtualMachine != nil && props.VirtualMachine.Id != nil {
				virtualMachineId = *props.VirtualMachine.Id
			}

			if err := d.Set("applied_dns_servers", appliedDNSServers); err != nil {
				return fmt.Errorf("setting `applied_dns_servers`: %+v", err)
			}

			if err := d.Set("dns_servers", dnsServers); err != nil {
				return fmt.Errorf("setting `applied_dns_servers`: %+v", err)
			}

			d.Set("enable_ip_forwarding", props.EnableIPForwarding)
			d.Set("enable_accelerated_networking", props.EnableAcceleratedNetworking)
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

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetworkInterfaceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseNetworkInterfaceID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NetworkInterfaceName, networkInterfaceResourceName)
	defer locks.UnlockByName(id.NetworkInterfaceName, networkInterfaceResourceName)

	existing, err := client.Get(ctx, *id, networkinterfaces.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			log.Printf("[DEBUG] %q was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	props := *existing.Model.Properties

	lockingDetails, err := determineResourcesToLockFromIPConfiguration(props.IPConfigurations)
	if err != nil {
		return fmt.Errorf("determining locking details: %+v", err)
	}

	lockingDetails.lock()
	defer lockingDetails.unlock()

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandNetworkInterfaceIPConfigurations(input []interface{}) (*[]networkinterfaces.NetworkInterfaceIPConfiguration, error) {
	ipConfigs := make([]networkinterfaces.NetworkInterfaceIPConfiguration, 0)

	for _, configRaw := range input {
		data := configRaw.(map[string]interface{})

		subnetId := data["subnet_id"].(string)
		privateIpAllocationMethod := data["private_ip_address_allocation"].(string)
		privateIpAddressVersion := networkinterfaces.IPVersion(data["private_ip_address_version"].(string))

		allocationMethod := networkinterfaces.IPAllocationMethod(privateIpAllocationMethod)
		properties := networkinterfaces.NetworkInterfaceIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: &allocationMethod,
			PrivateIPAddressVersion:   &privateIpAddressVersion,
		}

		if privateIpAddressVersion == networkinterfaces.IPVersionIPvFour && subnetId == "" {
			return nil, fmt.Errorf("A Subnet ID must be specified for an IPv4 Network Interface.")
		}

		if subnetId != "" {
			properties.Subnet = &networkinterfaces.Subnet{
				Id: &subnetId,
			}
		}

		if v := data["private_ip_address"].(string); v != "" {
			properties.PrivateIPAddress = &v
		}

		if v := data["public_ip_address_id"].(string); v != "" {
			properties.PublicIPAddress = &networkinterfaces.PublicIPAddress{
				Id: &v,
			}
		}

		if v, ok := data["primary"]; ok {
			properties.Primary = utils.Bool(v.(bool))
		}

		if v := data["gateway_load_balancer_frontend_ip_configuration_id"].(string); v != "" {
			properties.GatewayLoadBalancer = &networkinterfaces.SubResource{Id: &v}
		}

		name := data["name"].(string)
		ipConfigs = append(ipConfigs, networkinterfaces.NetworkInterfaceIPConfiguration{
			Name:       &name,
			Properties: &properties,
		})
	}

	// if we've got multiple IP Configurations - one must be designated Primary
	if len(ipConfigs) > 1 {
		hasPrimary := false
		for _, config := range ipConfigs {
			if config.Properties != nil && config.Properties.Primary != nil && *config.Properties.Primary {
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

func flattenNetworkInterfaceIPConfigurations(input *[]networkinterfaces.NetworkInterfaceIPConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, ipConfig := range *input {
		props := ipConfig.Properties

		name := ""
		if ipConfig.Name != nil {
			name = *ipConfig.Name
		}

		subnetId := ""
		if props.Subnet != nil && props.Subnet.Id != nil {
			subnetId = *props.Subnet.Id
		}

		privateIPAddress := ""
		if props.PrivateIPAddress != nil {
			privateIPAddress = *props.PrivateIPAddress
		}

		privateIPAllocationMethod := ""
		if props.PrivateIPAllocationMethod != nil {
			privateIPAllocationMethod = string(*props.PrivateIPAllocationMethod)
		}

		privateIPAddressVersion := ""
		if props.PrivateIPAddressVersion != nil {
			privateIPAddressVersion = string(*props.PrivateIPAddressVersion)
		}

		publicIPAddressId := ""
		if props.PublicIPAddress != nil && props.PublicIPAddress.Id != nil {
			publicIPAddressId = *props.PublicIPAddress.Id
		}

		primary := false
		if props.Primary != nil {
			primary = *props.Primary
		}

		gatewayLBFrontendIPConfigId := ""
		if props.GatewayLoadBalancer != nil && props.GatewayLoadBalancer.Id != nil {
			gatewayLBFrontendIPConfigId = *props.GatewayLoadBalancer.Id
		}

		result = append(result, map[string]interface{}{
			"name":                          name,
			"primary":                       primary,
			"private_ip_address":            privateIPAddress,
			"private_ip_address_allocation": privateIPAllocationMethod,
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

// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkInterfaceCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.NetworkInterfaces
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := commonids.NewNetworkInterfaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id, networkinterfaces.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_network_interface", id.ID())
		}
	}

	var enableIpForwarding, enableAcceleratedNetworking bool

	enableIpForwarding = d.Get("ip_forwarding_enabled").(bool)
	enableAcceleratedNetworking = d.Get("accelerated_networking_enabled").(bool)

	properties := networkinterfaces.NetworkInterfacePropertiesFormat{
		EnableIPForwarding:          &enableIpForwarding,
		EnableAcceleratedNetworking: &enableAcceleratedNetworking,
	}

	locks.ByID(strings.ToLower(id.ID()))
	defer locks.UnlockByID(strings.ToLower(id.ID()))

	if auxiliaryMode, hasAuxiliaryMode := d.GetOk("auxiliary_mode"); hasAuxiliaryMode {
		properties.AuxiliaryMode = pointer.ToEnum[networkinterfaces.NetworkInterfaceAuxiliaryMode](auxiliaryMode.(string))
	}

	if auxiliarySku, hasAuxiliarySku := d.GetOk("auxiliary_sku"); hasAuxiliarySku {
		properties.AuxiliarySku = pointer.ToEnum[networkinterfaces.NetworkInterfaceAuxiliarySku](auxiliarySku.(string))
	}

	dns, hasDns := d.GetOk("dns_servers")
	nameLabel, hasNameLabel := d.GetOk("internal_dns_name_label")
	if hasDns || hasNameLabel {
		dnsSettings := networkinterfaces.NetworkInterfaceDnsSettings{}

		if hasDns {
			dnsRaw := dns.([]any)
			dns := expandNetworkInterfaceDnsServers(dnsRaw)
			dnsSettings.DnsServers = &dns
		}

		if hasNameLabel {
			dnsSettings.InternalDnsNameLabel = pointer.To(nameLabel.(string))
		}

		properties.DnsSettings = &dnsSettings
	}

	ipConfigsRaw := d.Get("ip_configuration").([]any)
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
		ExtendedLocation: helpers.ExpandEdgeZoneModel(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Properties:       &properties,
		Tags:             tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, iface, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceNetworkInterfaceRead(d, meta)
}

func expandNetworkInterfaceIPConfigurations(input []any) (*[]networkinterfaces.NetworkInterfaceIPConfiguration, error) {
	ipConfigs := make([]networkinterfaces.NetworkInterfaceIPConfiguration, 0)

	for _, configRaw := range input {
		data := configRaw.(map[string]any)

		subnetId := data["subnet_id"].(string)
		privateIpAllocationMethod := data["private_ip_address_allocation"].(string)
		privateIpAddressVersion := networkinterfaces.IPVersion(data["private_ip_address_version"].(string))

		properties := networkinterfaces.NetworkInterfaceIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: pointer.ToEnum[networkinterfaces.IPAllocationMethod](privateIpAllocationMethod),
			PrivateIPAddressVersion:   &privateIpAddressVersion,
		}

		if privateIpAddressVersion == networkinterfaces.IPVersionIPvFour && subnetId == "" {
			return nil, errors.New("a Subnet ID must be specified for an IPv4 Network Interface")
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
			properties.Primary = pointer.To(v.(bool))
		}

		if v := data["gateway_load_balancer_frontend_ip_configuration_id"].(string); v != "" {
			properties.GatewayLoadBalancer = &networkinterfaces.SubResource{Id: &v}
		}

		ipConfigs = append(ipConfigs, networkinterfaces.NetworkInterfaceIPConfiguration{
			Name:       pointer.To(data["name"].(string)),
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
			return nil, errors.New("if multiple `ip_configurations` are specified - one must be designated as `primary`")
		}
	}

	return &ipConfigs, nil
}

func expandNetworkInterfaceDnsServers(input []any) []string {
	dnsServers := make([]string, 0)
	for _, v := range input {
		dnsServers = append(dnsServers, v.(string))
	}
	return dnsServers
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-03-01/virtualmachinescalesets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type OrchestratedVirtualMachineScaleSetDataSource struct{}

var _ sdk.DataSource = OrchestratedVirtualMachineScaleSetDataSource{}

type OrchestratedVirtualMachineScaleSetDataSourceModel struct {
	Name             string                                     `tfschema:"name"`
	ResourceGroup    string                                     `tfschema:"resource_group_name"`
	Location         string                                     `tfschema:"location"`
	NetworkInterface []VirtualMachineScaleSetNetworkInterface   `tfschema:"network_interface"`
	Identity         []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
}

type VirtualMachineScaleSetNetworkInterface struct {
	Name                         string                                                  `tfschema:"name"`
	IPConfiguration              []VirtualMachineScaleSetNetworkInterfaceIPConfiguration `tfschema:"ip_configuration"`
	DNSServers                   []string                                                `tfschema:"dns_servers"`
	AcceleratedNetworkingEnabled bool                                                    `tfschema:"accelerated_networking_enabled"`
	IPForwardingEnabled          bool                                                    `tfschema:"ip_forwarding_enabled"`
	NetworkSecurityGroupId       string                                                  `tfschema:"network_security_group_id"`
	Primary                      bool                                                    `tfschema:"primary"`
}

type VirtualMachineScaleSetNetworkInterfaceIPConfiguration struct {
	Name                                    string                                                                 `tfschema:"name"`
	ApplicationGatewayBackendAddressPoolIds []string                                                               `tfschema:"application_gateway_backend_address_pool_ids"`
	ApplicationSecurityGroupIds             []string                                                               `tfschema:"application_security_group_ids"`
	LoadBalancerBackendAddressPoolIds       []string                                                               `tfschema:"load_balancer_backend_address_pool_ids"`
	Primary                                 bool                                                                   `tfschema:"primary"`
	PublicIPAddress                         []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress `tfschema:"public_ip_address"`
	SubnetId                                string                                                                 `tfschema:"subnet_id"`
	Version                                 string                                                                 `tfschema:"version"`
}

type VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress struct {
	Name                 string                                                                      `tfschema:"name"`
	DomainNameLabel      string                                                                      `tfschema:"domain_name_label"`
	IdleTimeoutInMinutes int                                                                         `tfschema:"idle_timeout_in_minutes"`
	IPTag                []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag `tfschema:"ip_tag"`
	PublicIpPrefixId     string                                                                      `tfschema:"public_ip_prefix_id"`
	Version              string                                                                      `tfschema:"version"`
}

type VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag struct {
	Tag  string `tfschema:"tag"`
	Type string `tfschema:"type"`
}

func (r OrchestratedVirtualMachineScaleSetDataSource) ModelObject() interface{} {
	return &OrchestratedVirtualMachineScaleSetDataSourceModel{}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) ResourceType() string {
	return "azurerm_orchestrated_virtual_machine_scale_set"
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: computeValidate.VirtualMachineName,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"network_interface": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"ip_configuration": virtualMachineScaleSetIPConfigurationSchemaForDataSource(),

					"dns_servers": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"accelerated_networking_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"ip_forwarding_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"network_security_group_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"primary": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"identity": commonschema.UserAssignedIdentityComputed(),
	}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VMScaleSetClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var orchestratedVMSS OrchestratedVirtualMachineScaleSetDataSourceModel
			if err := metadata.Decode(&orchestratedVMSS); err != nil {
				return err
			}

			id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, orchestratedVMSS.ResourceGroup, orchestratedVMSS.Name)

			existing, err := client.Get(ctx, id, virtualmachinescalesets.GetOperationOptions{Expand: pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)})
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := existing.Model; model != nil {
				orchestratedVMSS.Location = location.Normalize(model.Location)

				if props := model.Properties; props != nil {
					if profile := props.VirtualMachineProfile; profile != nil {
						if nwProfile := profile.NetworkProfile; nwProfile != nil {
							orchestratedVMSS.NetworkInterface = flattenVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)
						}
					}
				}

				userIdentity, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				orchestratedVMSS.Identity = *userIdentity
			}

			metadata.SetID(id)

			return metadata.Encode(&orchestratedVMSS)
		},
	}
}

func flattenVirtualMachineScaleSetNetworkInterface(input *[]virtualmachinescalesets.VirtualMachineScaleSetNetworkConfiguration) []VirtualMachineScaleSetNetworkInterface {
	if input == nil {
		return []VirtualMachineScaleSetNetworkInterface{}
	}

	networkInterfaces := make([]VirtualMachineScaleSetNetworkInterface, 0)
	for _, v := range *input {
		if v.Properties == nil {
			continue
		}
		var networkSecurityGroupId string
		if v.Properties.NetworkSecurityGroup != nil && v.Properties.NetworkSecurityGroup.Id != nil {
			networkSecurityGroupId = *v.Properties.NetworkSecurityGroup.Id
		}
		var acceleratedNetworkingEnabled, ipForwardingEnabled, primary bool
		if v.Properties.EnableAcceleratedNetworking != nil {
			acceleratedNetworkingEnabled = *v.Properties.EnableAcceleratedNetworking
		}
		if v.Properties.EnableIPForwarding != nil {
			ipForwardingEnabled = *v.Properties.EnableIPForwarding
		}
		if v.Properties.Primary != nil {
			primary = *v.Properties.Primary
		}

		var dnsServers []string
		if settings := v.Properties.DnsSettings; settings != nil {
			dnsServers = *v.Properties.DnsSettings.DnsServers
		}

		networkInterfaces = append(networkInterfaces, VirtualMachineScaleSetNetworkInterface{
			Name:                         v.Name,
			NetworkSecurityGroupId:       networkSecurityGroupId,
			AcceleratedNetworkingEnabled: acceleratedNetworkingEnabled,
			IPForwardingEnabled:          ipForwardingEnabled,
			Primary:                      primary,
			DNSServers:                   dnsServers,
			IPConfiguration:              flattenOrchestratedVirtualMachineScaleSetNetworkInterfaceIPConfiguration(v.Properties.IPConfigurations),
		})
	}

	return networkInterfaces
}

func flattenOrchestratedVirtualMachineScaleSetNetworkInterfaceIPConfiguration(input []virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration) []VirtualMachineScaleSetNetworkInterfaceIPConfiguration {
	ipConfigurations := make([]VirtualMachineScaleSetNetworkInterfaceIPConfiguration, 0)
	for _, v := range input {
		if v.Properties == nil {
			continue
		}
		var subnetId string
		if v.Properties.Subnet != nil && v.Properties.Subnet.Id != nil {
			subnetId = *v.Properties.Subnet.Id
		}

		var primary bool
		if v.Properties.Primary != nil {
			primary = *v.Properties.Primary
		}

		ipConfigurations = append(ipConfigurations, VirtualMachineScaleSetNetworkInterfaceIPConfiguration{
			Name:                                    v.Name,
			SubnetId:                                subnetId,
			Primary:                                 primary,
			PublicIPAddress:                         flattenOrchestratedVirtualMachineScaleSetPublicIPAddress(v.Properties.PublicIPAddressConfiguration),
			ApplicationGatewayBackendAddressPoolIds: flattenSubResourcesToStringIDs(v.Properties.ApplicationGatewayBackendAddressPools),
			ApplicationSecurityGroupIds:             flattenSubResourcesToStringIDs(v.Properties.ApplicationSecurityGroups),
			LoadBalancerBackendAddressPoolIds:       flattenSubResourcesToStringIDs(v.Properties.LoadBalancerBackendAddressPools),
		})
	}

	return ipConfigurations
}

func flattenOrchestratedVirtualMachineScaleSetPublicIPAddress(input *virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration) []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress {
	if input == nil || input.Properties == nil {
		return []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress{}
	}

	ipTags := make([]VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag, 0)
	if input.Properties.IPTags != nil {
		for _, rawTag := range *input.Properties.IPTags {
			var tag, tagType string

			if rawTag.IPTagType != nil {
				tagType = *rawTag.IPTagType
			}

			if rawTag.Tag != nil {
				tag = *rawTag.Tag
			}

			ipTags = append(ipTags, VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag{
				Tag:  tag,
				Type: tagType,
			})
		}
	}

	var domainNameLabel, publicIPPrefixId, version string
	if input.Properties.DnsSettings != nil {
		domainNameLabel = input.Properties.DnsSettings.DomainNameLabel
	}

	if input.Properties.PublicIPPrefix != nil && input.Properties.PublicIPPrefix.Id != nil {
		publicIPPrefixId = *input.Properties.PublicIPPrefix.Id
	}

	if pointer.From(input.Properties.PublicIPAddressVersion) != "" {
		version = string(pointer.From(input.Properties.PublicIPAddressVersion))
	}

	var idleTimeoutInMinutes int
	if input.Properties.IdleTimeoutInMinutes != nil {
		idleTimeoutInMinutes = int(*input.Properties.IdleTimeoutInMinutes)
	}

	return []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress{{
		Name:                 input.Name,
		DomainNameLabel:      domainNameLabel,
		IdleTimeoutInMinutes: idleTimeoutInMinutes,
		IPTag:                ipTags,
		PublicIpPrefixId:     publicIPPrefixId,
		Version:              version,
	}}
}

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
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
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
	IdleTimeoutInMinutes int64                                                                       `tfschema:"idle_timeout_in_minutes"`
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

		"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),
	}
}

func (r OrchestratedVirtualMachineScaleSetDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachineScaleSetsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var orchestratedVMSS OrchestratedVirtualMachineScaleSetDataSourceModel
			if err := metadata.Decode(&orchestratedVMSS); err != nil {
				return err
			}

			id := virtualmachinescalesets.NewVirtualMachineScaleSetID(subscriptionId, orchestratedVMSS.ResourceGroup, orchestratedVMSS.Name)

			options := virtualmachinescalesets.DefaultGetOperationOptions()
			options.Expand = pointer.To(virtualmachinescalesets.ExpandTypesForGetVMScaleSetsUserData)
			existing, err := client.Get(ctx, id, options)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if model := existing.Model; model != nil {
				orchestratedVMSS.Location = location.Normalize(model.Location)

				identityFlattened, err := identity.FlattenSystemAndUserAssignedMapToModel(model.Identity)
				if err != nil {
					return err
				}
				orchestratedVMSS.Identity = pointer.From(identityFlattened)
				if props := model.Properties; props != nil {
					if profile := props.VirtualMachineProfile; profile != nil {
						if nwProfile := profile.NetworkProfile; nwProfile != nil {
							orchestratedVMSS.NetworkInterface = flattenVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)
						}
					}

				}
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
		var networkSecurityGroupId string
		var acceleratedNetworkingEnabled, ipForwardingEnabled, primary bool
		var dnsServers []string

		if props := v.Properties; props != nil {
			if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.Id != nil {
				networkSecurityGroupId = *props.NetworkSecurityGroup.Id
			}
			if props.EnableAcceleratedNetworking != nil {
				acceleratedNetworkingEnabled = *props.EnableAcceleratedNetworking
			}
			if props.EnableIPForwarding != nil {
				ipForwardingEnabled = *props.EnableIPForwarding
			}
			if props.Primary != nil {
				primary = *props.Primary
			}

			if settings := props.DnsSettings; settings != nil {
				dnsServers = *props.DnsSettings.DnsServers
			}

			networkInterfaces = append(networkInterfaces, VirtualMachineScaleSetNetworkInterface{
				Name:                         v.Name,
				NetworkSecurityGroupId:       networkSecurityGroupId,
				AcceleratedNetworkingEnabled: acceleratedNetworkingEnabled,
				IPForwardingEnabled:          ipForwardingEnabled,
				Primary:                      primary,
				DNSServers:                   dnsServers,
				IPConfiguration:              flattenOrchestratedVirtualMachineScaleSetNetworkInterfaceIPConfiguration(&props.IPConfigurations),
			})
		}
	}

	return networkInterfaces
}

func flattenOrchestratedVirtualMachineScaleSetNetworkInterfaceIPConfiguration(input *[]virtualmachinescalesets.VirtualMachineScaleSetIPConfiguration) []VirtualMachineScaleSetNetworkInterfaceIPConfiguration {
	if input == nil {
		return []VirtualMachineScaleSetNetworkInterfaceIPConfiguration{}
	}

	ipConfigurations := make([]VirtualMachineScaleSetNetworkInterfaceIPConfiguration, 0)
	for _, v := range *input {
		var subnetId string
		var primary bool
		if props := v.Properties; props != nil {
			if props.Subnet != nil && props.Subnet.Id != nil {
				subnetId = *props.Subnet.Id
			}

			if props.Primary != nil {
				primary = *props.Primary
			}

			ipConfigurations = append(ipConfigurations, VirtualMachineScaleSetNetworkInterfaceIPConfiguration{
				Name:                                    v.Name,
				SubnetId:                                subnetId,
				Primary:                                 primary,
				PublicIPAddress:                         flattenOrchestratedVirtualMachineScaleSetPublicIPAddress(props.PublicIPAddressConfiguration),
				ApplicationGatewayBackendAddressPoolIds: flattenSubResourcesToStringIDs(props.ApplicationGatewayBackendAddressPools),
				ApplicationSecurityGroupIds:             flattenSubResourcesToStringIDs(props.ApplicationSecurityGroups),
				LoadBalancerBackendAddressPoolIds:       flattenSubResourcesToStringIDs(props.LoadBalancerBackendAddressPools),
			})
		}
	}

	return ipConfigurations
}

func flattenOrchestratedVirtualMachineScaleSetPublicIPAddress(input *virtualmachinescalesets.VirtualMachineScaleSetPublicIPAddressConfiguration) []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress {
	if input == nil {
		return []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress{}
	}

	ipTags := make([]VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag, 0)
	var domainNameLabel, publicIPPrefixId, version string
	var idleTimeoutInMinutes int64
	if props := input.Properties; props != nil && props.IPTags != nil {
		for _, rawTag := range *props.IPTags {
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

		if props.DnsSettings != nil {
			domainNameLabel = props.DnsSettings.DomainNameLabel
		}

		if props.PublicIPPrefix != nil && props.PublicIPPrefix.Id != nil {
			publicIPPrefixId = *props.PublicIPPrefix.Id
		}

		if props.PublicIPAddressVersion != nil {
			version = string(pointer.From(props.PublicIPAddressVersion))
		}

		if props.IdleTimeoutInMinutes != nil {
			idleTimeoutInMinutes = *props.IdleTimeoutInMinutes
		}
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

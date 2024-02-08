// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	computeValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/compute/2023-03-01/compute"
)

type OrchestratedVirtualMachineScaleSetDataSource struct{}

var _ sdk.DataSource = OrchestratedVirtualMachineScaleSetDataSource{}

type OrchestratedVirtualMachineScaleSetDataSourceModel struct {
	Name             string                                   `tfschema:"name"`
	ResourceGroup    string                                   `tfschema:"resource_group_name"`
	Location         string                                   `tfschema:"location"`
	NetworkInterface []VirtualMachineScaleSetNetworkInterface `tfschema:"network_interface"`
	Identity         []identity.ModelUserAssigned             `tfschema:"identity"`
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

			id := commonids.NewVirtualMachineScaleSetID(subscriptionId, orchestratedVMSS.ResourceGroup, orchestratedVMSS.Name)

			existing, err := client.Get(ctx, id.ResourceGroupName, id.VirtualMachineScaleSetName, compute.ExpandTypesForGetVMScaleSetsUserData)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return fmt.Errorf("%s not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			orchestratedVMSS.Location = location.NormalizeNilable(existing.Location)

			if profile := existing.VirtualMachineProfile; profile != nil {
				if nwProfile := profile.NetworkProfile; nwProfile != nil {
					orchestratedVMSS.NetworkInterface = flattenVirtualMachineScaleSetNetworkInterface(nwProfile.NetworkInterfaceConfigurations)
				}
			}

			userIdentity, err := flattenOrchestratedVirtualMachineScaleSetIdentityToModel(existing.Identity)
			if err != nil {
				return err
			}
			orchestratedVMSS.Identity = userIdentity

			metadata.SetID(id)

			return metadata.Encode(&orchestratedVMSS)
		},
	}
}

func flattenVirtualMachineScaleSetNetworkInterface(input *[]compute.VirtualMachineScaleSetNetworkConfiguration) []VirtualMachineScaleSetNetworkInterface {
	if input == nil {
		return []VirtualMachineScaleSetNetworkInterface{}
	}

	networkInterfaces := make([]VirtualMachineScaleSetNetworkInterface, 0)
	for _, v := range *input {
		var name, networkSecurityGroupId string
		if v.Name != nil {
			name = *v.Name
		}
		if v.NetworkSecurityGroup != nil && v.NetworkSecurityGroup.ID != nil {
			networkSecurityGroupId = *v.NetworkSecurityGroup.ID
		}
		var acceleratedNetworkingEnabled, ipForwardingEnabled, primary bool
		if v.EnableAcceleratedNetworking != nil {
			acceleratedNetworkingEnabled = *v.EnableAcceleratedNetworking
		}
		if v.EnableIPForwarding != nil {
			ipForwardingEnabled = *v.EnableIPForwarding
		}
		if v.Primary != nil {
			primary = *v.Primary
		}

		var dnsServers []string
		if settings := v.DNSSettings; settings != nil {
			dnsServers = *v.DNSSettings.DNSServers
		}

		networkInterfaces = append(networkInterfaces, VirtualMachineScaleSetNetworkInterface{
			Name:                         name,
			NetworkSecurityGroupId:       networkSecurityGroupId,
			AcceleratedNetworkingEnabled: acceleratedNetworkingEnabled,
			IPForwardingEnabled:          ipForwardingEnabled,
			Primary:                      primary,
			DNSServers:                   dnsServers,
			IPConfiguration:              flattenOrchestratedVirtualMachineScaleSetNetworkInterfaceIPConfiguration(v.IPConfigurations),
		})
	}

	return networkInterfaces
}

func flattenOrchestratedVirtualMachineScaleSetNetworkInterfaceIPConfiguration(input *[]compute.VirtualMachineScaleSetIPConfiguration) []VirtualMachineScaleSetNetworkInterfaceIPConfiguration {
	if input == nil {
		return []VirtualMachineScaleSetNetworkInterfaceIPConfiguration{}
	}

	ipConfigurations := make([]VirtualMachineScaleSetNetworkInterfaceIPConfiguration, 0)
	for _, v := range *input {
		var name, subnetId string
		if v.Name != nil {
			name = *v.Name
		}
		if v.Subnet != nil && v.Subnet.ID != nil {
			subnetId = *v.Subnet.ID
		}

		var primary bool
		if v.Primary != nil {
			primary = *v.Primary
		}

		ipConfigurations = append(ipConfigurations, VirtualMachineScaleSetNetworkInterfaceIPConfiguration{
			Name:                                    name,
			SubnetId:                                subnetId,
			Primary:                                 primary,
			PublicIPAddress:                         flattenOrchestratedVirtualMachineScaleSetPublicIPAddress(v.PublicIPAddressConfiguration),
			ApplicationGatewayBackendAddressPoolIds: flattenSubResourcesToStringIDs(v.ApplicationGatewayBackendAddressPools),
			ApplicationSecurityGroupIds:             flattenSubResourcesToStringIDs(v.ApplicationSecurityGroups),
			LoadBalancerBackendAddressPoolIds:       flattenSubResourcesToStringIDs(v.LoadBalancerBackendAddressPools),
		})
	}

	return ipConfigurations
}

func flattenOrchestratedVirtualMachineScaleSetPublicIPAddress(input *compute.VirtualMachineScaleSetPublicIPAddressConfiguration) []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress {
	if input == nil {
		return []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress{}
	}

	ipTags := make([]VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddressIPTag, 0)
	if input.IPTags != nil {
		for _, rawTag := range *input.IPTags {
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

	var domainNameLabel, name, publicIPPrefixId, version string
	if input.DNSSettings != nil && input.DNSSettings.DomainNameLabel != nil {
		domainNameLabel = *input.DNSSettings.DomainNameLabel
	}

	if input.Name != nil {
		name = *input.Name
	}

	if input.PublicIPPrefix != nil && input.PublicIPPrefix.ID != nil {
		publicIPPrefixId = *input.PublicIPPrefix.ID
	}

	if input.PublicIPAddressVersion != "" {
		version = string(input.PublicIPAddressVersion)
	}

	var idleTimeoutInMinutes int
	if input.IdleTimeoutInMinutes != nil {
		idleTimeoutInMinutes = int(*input.IdleTimeoutInMinutes)
	}

	return []VirtualMachineScaleSetNetworkInterfaceIPConfigurationPublicIPAddress{{
		Name:                 name,
		DomainNameLabel:      domainNameLabel,
		IdleTimeoutInMinutes: idleTimeoutInMinutes,
		IPTag:                ipTags,
		PublicIpPrefixId:     publicIPPrefixId,
		Version:              version,
	}}
}

func flattenOrchestratedVirtualMachineScaleSetIdentityToModel(input *compute.VirtualMachineScaleSetIdentity) ([]identity.ModelUserAssigned, error) {
	if input == nil {
		return nil, nil
	}

	identityIds := make(map[string]identity.UserAssignedIdentityDetails, 0)
	for k, v := range input.UserAssignedIdentities {
		if v != nil {
			identityIds[k] = identity.UserAssignedIdentityDetails{
				ClientId:    v.ClientID,
				PrincipalId: v.PrincipalID,
			}
		}
	}

	tmp := identity.UserAssignedMap{
		Type:        identity.Type(input.Type),
		IdentityIds: identityIds,
	}

	output, err := identity.FlattenUserAssignedMapToModel(&tmp)
	if err != nil {
		return nil, fmt.Errorf("expanding `identity`: %+v", err)
	}

	return *output, nil
}

// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkappliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -test-sequential

type VirtualNetworkRoutingApplianceResource struct{}

var (
	_ sdk.ResourceWithUpdate   = VirtualNetworkRoutingApplianceResource{}
	_ sdk.ResourceWithIdentity = VirtualNetworkRoutingApplianceResource{}
)

type VirtualNetworkRoutingApplianceResourceModel struct {
	Name                    string                                               `tfschema:"name"`
	ResourceGroupName       string                                               `tfschema:"resource_group_name"`
	Location                string                                               `tfschema:"location"`
	BandwidthInGbps         float64                                              `tfschema:"bandwidth_in_gbps"`
	SubnetId                string                                               `tfschema:"subnet_id"`
	PrivateIPAddressVersion string                                               `tfschema:"private_ip_address_version"`
	Tags                    map[string]string                                    `tfschema:"tags"`
	IPConfigurations        []VirtualNetworkRoutingApplianceIPConfigurationModel `tfschema:"ip_configuration"`
	ResourceGuid            string                                               `tfschema:"resource_guid"`
}

type VirtualNetworkRoutingApplianceIPConfigurationModel struct {
	Name                      string `tfschema:"name"`
	Primary                   bool   `tfschema:"primary"`
	PrivateIPAddress          string `tfschema:"private_ip_address"`
	PrivateIPAllocationMethod string `tfschema:"private_ip_address_allocation"`
}

func (VirtualNetworkRoutingApplianceResource) ResourceType() string {
	return "azurerm_virtual_network_routing_appliance"
}

func (VirtualNetworkRoutingApplianceResource) ModelObject() any {
	return &VirtualNetworkRoutingApplianceResourceModel{}
}

func (VirtualNetworkRoutingApplianceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualnetworkappliances.ValidateVirtualNetworkApplianceID
}

func (VirtualNetworkRoutingApplianceResource) Identity() resourceids.ResourceId {
	return &virtualnetworkappliances.VirtualNetworkApplianceId{}
}

func (VirtualNetworkRoutingApplianceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[0-9a-zA-Z]([0-9a-zA-Z_.-]{0,62}[0-9a-zA-Z_])?$`),
				"`name` must be between 1 and 64 characters, start with a letter or number, end with a letter, number, or underscore, and contain only letters, numbers, underscores, periods, and hyphens",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"bandwidth_in_gbps": {
			Type:         pluginsdk.TypeFloat,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.FloatInSlice([]float64{10, 50, 100, 200}),
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				commonids.ValidateSubnetID,
				validation.StringMatch(regexp.MustCompile(`(?i)/subnets/VirtualNetworkApplianceSubnet$`), "`subnet_id` must reference a subnet named `VirtualNetworkApplianceSubnet`"),
			),
		},

		"private_ip_address_version": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(virtualnetworkappliances.VirtualNetworkApplianceIPVersionTypeIPvFour),
			ValidateFunc: validation.StringInSlice(
				virtualnetworkappliances.PossibleValuesForVirtualNetworkApplianceIPVersionType(),
				false,
			),
		},

		"tags": commonschema.Tags(),
	}
}

func (VirtualNetworkRoutingApplianceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"primary": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"private_ip_address_allocation": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"resource_guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r VirtualNetworkRoutingApplianceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkAppliances

			var config VirtualNetworkRoutingApplianceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := virtualnetworkappliances.NewVirtualNetworkApplianceID(metadata.Client.Account.SubscriptionId, config.ResourceGroupName, config.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := virtualnetworkappliances.VirtualNetworkAppliance{
				Location: pointer.To(location.Normalize(config.Location)),
				Properties: &virtualnetworkappliances.VirtualNetworkAppliancePropertiesFormat{
					BandwidthInGbps:         pointer.To(config.BandwidthInGbps),
					PrivateIPAddressVersion: pointer.ToEnum[virtualnetworkappliances.VirtualNetworkApplianceIPVersionType](config.PrivateIPAddressVersion),
					Subnet: &virtualnetworkappliances.Subnet{
						Id: pointer.To(config.SubnetId),
					},
				},
				Tags: pointer.To(config.Tags),
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(&id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r VirtualNetworkRoutingApplianceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkAppliances

			id, err := virtualnetworkappliances.ParseVirtualNetworkApplianceIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (VirtualNetworkRoutingApplianceResource) flatten(metadata sdk.ResourceMetaData, id *virtualnetworkappliances.VirtualNetworkApplianceId, model *virtualnetworkappliances.VirtualNetworkAppliance) error {
	state := VirtualNetworkRoutingApplianceResourceModel{
		Name:              id.VirtualNetworkApplianceName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.NormalizeNilable(model.Location)
		state.Tags = pointer.From(model.Tags)

		if props := model.Properties; props != nil {
			state.BandwidthInGbps = pointer.From(props.BandwidthInGbps)
			state.PrivateIPAddressVersion = pointer.FromEnum(props.PrivateIPAddressVersion)
			state.ResourceGuid = pointer.From(props.ResourceGuid)

			if props.Subnet != nil && props.Subnet.Id != nil {
				subnetId, err := commonids.ParseSubnetIDInsensitively(pointer.From(props.Subnet.Id))
				if err != nil {
					return err
				}

				state.SubnetId = subnetId.ID()
			}

			state.IPConfigurations = make([]VirtualNetworkRoutingApplianceIPConfigurationModel, 0)
			for _, config := range pointer.From(props.IPConfigurations) {
				ipConfiguration := VirtualNetworkRoutingApplianceIPConfigurationModel{
					Name: pointer.From(config.Name),
				}

				if configProps := config.Properties; configProps != nil {
					ipConfiguration.Primary = pointer.From(configProps.Primary)
					ipConfiguration.PrivateIPAddress = pointer.From(configProps.PrivateIPAddress)
					ipConfiguration.PrivateIPAllocationMethod = pointer.FromEnum(configProps.PrivateIPAllocationMethod)
				}

				state.IPConfigurations = append(state.IPConfigurations, ipConfiguration)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (VirtualNetworkRoutingApplianceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkAppliances

			id, err := virtualnetworkappliances.ParseVirtualNetworkApplianceIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config VirtualNetworkRoutingApplianceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pointer.To(config.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (VirtualNetworkRoutingApplianceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.VirtualNetworkAppliances

			id, err := virtualnetworkappliances.ParseVirtualNetworkApplianceIDInsensitively(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

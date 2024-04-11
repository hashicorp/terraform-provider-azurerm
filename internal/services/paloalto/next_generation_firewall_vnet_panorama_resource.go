// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package paloalto

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls"
	helpersValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallVNetPanoramaResource struct{}

type NextGenerationFirewallVnetPanoramaModel struct {
	Name                 string                      `tfschema:"name"`
	ResourceGroupName    string                      `tfschema:"resource_group_name"`
	Location             string                      `tfschema:"location"`
	PanoramaBase64Config string                      `tfschema:"panorama_base64_config"`
	NetworkProfile       []schema.NetworkProfileVnet `tfschema:"network_profile"`
	DNSSettings          []schema.DNSSettings        `tfschema:"dns_settings"`
	FrontEnd             []schema.DestinationNAT     `tfschema:"destination_nat"`
	PanoramaConfig       []schema.Panorama           `tfschema:"panorama"`
	Tags                 map[string]interface{}      `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallVNetPanoramaResource{}

func (r NextGenerationFirewallVNetPanoramaResource) ModelObject() interface{} {
	return &NextGenerationFirewallVnetPanoramaModel{}
}

func (r NextGenerationFirewallVNetPanoramaResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"panorama_base64_config": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: helpersValidate.Base64EncodedString,
		},

		"network_profile": schema.VnetNetworkProfileSchema(),

		// Optional
		"dns_settings": schema.DNSSettingsSchema(),

		"destination_nat": schema.DestinationNATSchema(),

		"tags": commonschema.Tags(),
	}
}

func (r NextGenerationFirewallVNetPanoramaResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"panorama": schema.PanoramaSchema(),
	}
}

func (r NextGenerationFirewallVNetPanoramaResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama"
}

func (r NextGenerationFirewallVNetPanoramaResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			var model NextGenerationFirewallVnetPanoramaModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := firewalls.NewFirewallID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			firewall := firewalls.FirewallResource{
				Location: location.Normalize(model.Location),
				Properties: firewalls.FirewallDeploymentProperties{
					PanoramaConfig: &firewalls.PanoramaConfig{
						ConfigString: model.PanoramaBase64Config,
					},
					IsPanoramaManaged: pointer.To(firewalls.BooleanEnumTRUE),
					DnsSettings:       schema.ExpandDNSSettings(model.DNSSettings),
					MarketplaceDetails: firewalls.MarketplaceDetails{
						OfferId:     "pan_swfw_cloud_ngfw", // TODO - Will just supplying the offer ID `panw-cloud-ngfw-payg` work?
						PublisherId: "paloaltonetworks",
					},
					NetworkProfile: schema.ExpandNetworkProfileVnet(model.NetworkProfile),
					PlanData: firewalls.PlanData{
						BillingCycle: firewalls.BillingCycleMONTHLY,
						PlanId:       "panw-cloud-ngfw-payg",
					},
					FrontEndSettings: schema.ExpandDestinationNAT(model.FrontEnd),
				},
				Tags: tags.Expand(model.Tags),
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, firewall); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NextGenerationFirewallVNetPanoramaResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVnetPanoramaModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state.Name = id.FirewallName
			state.ResourceGroupName = id.ResourceGroupName
			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				props := model.Properties

				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)

				state.NetworkProfile = schema.FlattenNetworkProfileVnet(props.NetworkProfile)

				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)

				if panoramaConfig := props.PanoramaConfig; panoramaConfig != nil {
					state.PanoramaBase64Config = panoramaConfig.ConfigString
					state.PanoramaConfig = []schema.Panorama{{
						Name:            pointer.From(panoramaConfig.CgName),
						DeviceGroupName: pointer.From(panoramaConfig.DgName),
						HostName:        pointer.From(panoramaConfig.HostName),
						PanoramaServer:  pointer.From(panoramaConfig.PanoramaServer),
						PanoramaServer2: pointer.From(panoramaConfig.PanoramaServer2),
						TplName:         pointer.From(panoramaConfig.TplName),
						VMAuthKey:       pointer.From(panoramaConfig.VMAuthKey),
					}}
				}

				state.Tags = tags.Flatten(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallVNetPanoramaResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err = client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r NextGenerationFirewallVNetPanoramaResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewalls.ValidateFirewallID
}

func (r NextGenerationFirewallVNetPanoramaResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := NextGenerationFirewallVnetPanoramaModel{}

			if err = metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retreiving %s: %+v", *id, err)
			}

			firewall := *existing.Model
			props := firewall.Properties

			if metadata.ResourceData.HasChange("panorama_base64_config") {
				props.PanoramaConfig.ConfigString = model.PanoramaBase64Config
			}

			if metadata.ResourceData.HasChange("network_profile") {
				props.NetworkProfile = schema.ExpandNetworkProfileVnet(model.NetworkProfile)
			}

			if metadata.ResourceData.HasChange("dns_settings") {
				props.DnsSettings = schema.ExpandDNSSettings(model.DNSSettings)
			}

			if metadata.ResourceData.HasChange("destination_nat") {
				props.FrontEndSettings = schema.ExpandDestinationNAT(model.FrontEnd)
			}

			firewall.Properties = props

			if metadata.ResourceData.HasChange("tags") {
				firewall.Tags = tags.Expand(model.Tags)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, firewall); err != nil {
				return err
			}

			return nil
		},
	}
}

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NextGenerationFirewallVHubPanoramaResource struct{}

type NextGenerationFirewallVHubPanoramaResourceModel struct {
	Name                 string                      `tfschema:"name"`
	ResourceGroupName    string                      `tfschema:"resource_group_name"`
	PanoramaBase64Config string                      `tfschema:"panorama_base64_config"`
	Location             string                      `tfschema:"location"`
	NetworkProfile       []schema.NetworkProfileVHub `tfschema:"network_profile"`
	DNSSettings          []schema.DNSSettings        `tfschema:"dns_settings"`
	FrontEnd             []schema.DestinationNAT     `tfschema:"destination_nat"`
	MarketplaceOfferId   string                      `tfschema:"marketplace_offer_id"`
	PlanId               string                      `tfschema:"plan_id"`
	Tags                 map[string]interface{}      `tfschema:"tags"`

	// Computed
	PanoramaConfig []schema.Panorama `tfschema:"panorama"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallVHubPanoramaResource{}

func (r NextGenerationFirewallVHubPanoramaResource) ModelObject() interface{} {
	return &NextGenerationFirewallVHubPanoramaResourceModel{}
}

func (r NextGenerationFirewallVHubPanoramaResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewalls.ValidateFirewallID
}

func (r NextGenerationFirewallVHubPanoramaResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama"
}

func (r NextGenerationFirewallVHubPanoramaResource) Arguments() map[string]*pluginsdk.Schema {
	args := map[string]*pluginsdk.Schema{
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

		"network_profile": schema.VHubNetworkProfileSchema(),

		"dns_settings": schema.DNSSettingsSchema(),

		"destination_nat": schema.DestinationNATSchema(),

		"marketplace_offer_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "pan_swfw_cloud_ngfw",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "panw-cngfw-payg",
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"tags": commonschema.Tags(),
	}

	if !features.FivePointOh() {
		args["plan_id"].Default = "panw-cloud-ngfw-payg"
	}

	return args
}

func (r NextGenerationFirewallVHubPanoramaResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"panorama": schema.PanoramaSchema(),
	}
}

func (r NextGenerationFirewallVHubPanoramaResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			var model NextGenerationFirewallVHubPanoramaResourceModel

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
						OfferId:     model.MarketplaceOfferId,
						PublisherId: "paloaltonetworks",
					},
					NetworkProfile: schema.ExpandNetworkProfileVHub(model.NetworkProfile),
					PlanData: firewalls.PlanData{
						BillingCycle: firewalls.BillingCycleMONTHLY,
						PlanId:       model.PlanId,
					},
					FrontEndSettings: schema.ExpandDestinationNAT(model.FrontEnd),
				},
				Tags: tags.Expand(model.Tags),
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, firewall); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NextGenerationFirewallVHubPanoramaResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVHubPanoramaResourceModel

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

				netProfile, err := schema.FlattenNetworkProfileVHub(props.NetworkProfile)
				if err != nil {
					return fmt.Errorf("flattening Network Profile for %s: %+v", *id, err)
				}

				state.NetworkProfile = []schema.NetworkProfileVHub{*netProfile}

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

				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId

				state.PlanId = props.PlanData.PlanId

				state.Tags = tags.Flatten(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallVHubPanoramaResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
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

func (r NextGenerationFirewallVHubPanoramaResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := NextGenerationFirewallVHubPanoramaResourceModel{}

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
				props.NetworkProfile = schema.ExpandNetworkProfileVHub(model.NetworkProfile)
			}

			if metadata.ResourceData.HasChange("dns_settings") {
				props.DnsSettings = schema.ExpandDNSSettings(model.DNSSettings)
			}

			if metadata.ResourceData.HasChange("destination_nat") {
				props.FrontEndSettings = schema.ExpandDestinationNAT(model.FrontEnd)
			}

			if metadata.ResourceData.HasChange("plan_id") {
				props.PlanData.PlanId = model.PlanId
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

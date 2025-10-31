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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-05-23/firewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NextGenerationFirewallVNetStrataCloudManagerResource struct{}

type NextGenerationFirewallVNetStrataCloudManagerModel struct {
	Name                         string                       `tfschema:"name"`
	ResourceGroupName            string                       `tfschema:"resource_group_name"`
	Location                     string                       `tfschema:"location"`
	NetworkProfile               []schema.NetworkProfileVnet  `tfschema:"network_profile"`
	StrataCloudManagerTenantName string                       `tfschema:"strata_cloud_manager_tenant_name"`
	DNSSettings                  []schema.DNSSettings         `tfschema:"dns_settings"`
	FrontEnd                     []schema.DestinationNAT      `tfschema:"destination_nat"`
	MarketplaceOfferId           string                       `tfschema:"marketplace_offer_id"`
	PlanId                       string                       `tfschema:"plan_id"`
	Identity                     []identity.ModelUserAssigned `tfschema:"identity"`
	Tags                         map[string]interface{}       `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallVNetStrataCloudManagerResource{}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) ModelObject() interface{} {
	return &NextGenerationFirewallVNetStrataCloudManagerModel{}
}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) Arguments() map[string]*pluginsdk.Schema {
	args := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"network_profile": schema.VnetNetworkProfileSchema(),

		"strata_cloud_manager_tenant_name": {
			Type:        pluginsdk.TypeString,
			Required:    true,
			Description: "Strata Cloud Manager name which is intended to manage the policy for this firewall.",
		},

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

		"identity": commonschema.UserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}

	return args
}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager"
}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2025_05_23.Firewalls

			var model NextGenerationFirewallVNetStrataCloudManagerModel

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

			expandedIdentity, err := expandPaloAltoLegacyToUserAssignedIdentity(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			firewall := firewalls.FirewallResource{
				Location: location.Normalize(model.Location),
				Properties: firewalls.FirewallDeploymentProperties{
					IsStrataCloudManaged: pointer.To(firewalls.BooleanEnumTRUE),
					StrataCloudManagerConfig: &firewalls.StrataCloudManagerConfig{
						CloudManagerName: model.StrataCloudManagerTenantName,
					},
					DnsSettings: schema.ExpandDNSSettings(model.DNSSettings),
					MarketplaceDetails: firewalls.MarketplaceDetails{
						OfferId:     model.MarketplaceOfferId,
						PublisherId: "paloaltonetworks",
					},
					NetworkProfile: schema.ExpandNetworkProfileVnet(model.NetworkProfile),
					PlanData: firewalls.PlanData{
						BillingCycle: firewalls.BillingCycleMONTHLY,
						PlanId:       model.PlanId,
					},
					FrontEndSettings: schema.ExpandDestinationNAT(model.FrontEnd),
				},
				Identity: expandedIdentity,
				Tags:     tags.Expand(model.Tags),
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, firewall); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2025_05_23.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVNetStrataCloudManagerModel

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
				props := model.Properties

				state.Location = location.Normalize(model.Location)

				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)

				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)

				state.NetworkProfile = schema.FlattenNetworkProfileVnet(props.NetworkProfile)

				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId

				state.PlanId = props.PlanData.PlanId

				flattenedIdentity, err := flattenPaloAltoUserAssignedToLegacyIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}
				state.Identity = flattenedIdentity

				state.Tags = tags.Flatten(existing.Model.Tags)

				if strataCloudManagerConfig := props.StrataCloudManagerConfig; strataCloudManagerConfig != nil {
					state.StrataCloudManagerTenantName = strataCloudManagerConfig.CloudManagerName
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2025_05_23.Firewalls

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

func (r NextGenerationFirewallVNetStrataCloudManagerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewalls.ValidateFirewallID
}

func (r NextGenerationFirewallVNetStrataCloudManagerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2025_05_23.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model NextGenerationFirewallVNetStrataCloudManagerModel
			if err = metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding model: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			firewall := *existing.Model
			props := firewall.Properties

			if metadata.ResourceData.HasChange("strata_cloud_manager_tenant_name") {
				props.StrataCloudManagerConfig.CloudManagerName = model.StrataCloudManagerTenantName
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

			if metadata.ResourceData.HasChange("plan_id") {
				props.PlanData.PlanId = model.PlanId
			}

			firewall.Properties = props

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := expandPaloAltoLegacyToUserAssignedIdentity(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding `identity`: %+v", err)
				}
				firewall.Identity = expandedIdentity
			}

			if metadata.ResourceData.HasChange("tags") {
				firewall.Tags = tags.Expand(model.Tags)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, firewall); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

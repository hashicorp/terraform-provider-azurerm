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
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2023-09-01/firewalls"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type NextGenerationFirewallVNetLocalRulestackResource struct{}

type NextGenerationFirewallVnetLocalRulestackModel struct {
	Name               string                      `tfschema:"name"`
	ResourceGroupName  string                      `tfschema:"resource_group_name"`
	NetworkProfile     []schema.NetworkProfileVnet `tfschema:"network_profile"`
	RuleStackId        string                      `tfschema:"rulestack_id"`
	DNSSettings        []schema.DNSSettings        `tfschema:"dns_settings"`
	FrontEnd           []schema.DestinationNAT     `tfschema:"destination_nat"`
	MarketplaceOfferId string                      `tfschema:"marketplace_offer_id"`
	PlanId             string                      `tfschema:"plan_id"`
	Tags               map[string]interface{}      `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallVNetLocalRulestackResource{}

func (r NextGenerationFirewallVNetLocalRulestackResource) ModelObject() interface{} {
	return &NextGenerationFirewallVnetLocalRulestackModel{}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Arguments() map[string]*pluginsdk.Schema {
	args := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: localrulestacks.ValidateLocalRulestackID,
		},

		"network_profile": schema.VnetNetworkProfileSchema(),

		// Optional
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

func (r NextGenerationFirewallVNetLocalRulestackResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack"
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls
			localRulestackClient := metadata.Client.PaloAlto.Client.LocalRulestacks

			var model NextGenerationFirewallVnetLocalRulestackModel

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

			ruleStackID, err := localrulestacks.ParseLocalRulestackID(model.RuleStackId)
			if err != nil {
				return err
			}

			ruleStack, err := localRulestackClient.Get(ctx, *ruleStackID)
			if err != nil {
				return fmt.Errorf("reading %s for %s: %+v", ruleStackID, id, err)
			}

			loc := location.Normalize(ruleStack.Model.Location)

			firewall := firewalls.FirewallResource{
				Location: loc,
				Properties: firewalls.FirewallDeploymentProperties{
					AssociatedRulestack: &firewalls.RulestackDetails{
						ResourceId: pointer.To(ruleStackID.ID()),
						Location:   pointer.To(location.Normalize(ruleStack.Model.Location)),
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
				Tags: tags.Expand(model.Tags),
			}

			locks.ByID(ruleStackID.ID())
			defer locks.UnlockByID(ruleStackID.ID())

			if err = client.CreateOrUpdateThenPoll(ctx, id, firewall); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state NextGenerationFirewallVnetLocalRulestackModel

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

				state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)

				state.NetworkProfile = schema.FlattenNetworkProfileVnet(props.NetworkProfile)

				state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)

				state.RuleStackId = pointer.From(props.AssociatedRulestack.ResourceId)

				state.MarketplaceOfferId = props.MarketplaceDetails.OfferId

				state.PlanId = props.PlanData.PlanId

				state.Tags = tags.Flatten(existing.Model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Delete() sdk.ResourceFunc {
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

func (r NextGenerationFirewallVNetLocalRulestackResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return firewalls.ValidateFirewallID
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.PaloAltoClient_v2023_09_01.Firewalls

			id, err := firewalls.ParseFirewallID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := NextGenerationFirewallVnetLocalRulestackModel{}

			if err = metadata.Decode(&model); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			firewall := *existing.Model
			props := firewall.Properties

			if metadata.ResourceData.HasChange("rulestack_id") {
				ruleStackID, err := localrulestacks.ParseLocalRulestackID(model.RuleStackId)
				if err != nil {
					return err
				}

				ruleStack := &firewalls.RulestackDetails{
					Location:    props.AssociatedRulestack.Location,
					ResourceId:  nil,
					RulestackId: pointer.To(ruleStackID.ID()),
				}

				props.AssociatedRulestack = ruleStack
				locks.ByID(ruleStackID.ID())
				defer locks.UnlockByID(ruleStackID.ID())
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

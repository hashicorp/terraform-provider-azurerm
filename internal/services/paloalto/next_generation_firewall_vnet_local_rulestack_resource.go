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
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NextGenerationFirewallVNetLocalRulestackResource struct{}

type NextGenerationFirewallVnetLocalRulestackModel struct {
	Name              string                      `tfschema:"name"`
	ResourceGroupName string                      `tfschema:"resource_group_name"`
	NetworkProfile    []schema.NetworkProfileVnet `tfschema:"network_profile"`
	RuleStackId       string                      `tfschema:"rule_stack_id"`
	DNSSettings       []schema.DNSSettings        `tfschema:"dns_settings"`
	FrontEnd          []schema.DestinationNAT     `tfschema:"destination_nat"`

	// Computed
	PlanData []schema.Plan `tfschema:"plan"`
	PanEtag  string        `tfschema:"pan_etag"`

	Tags map[string]interface{} `tfschema:"tags"`
}

var _ sdk.ResourceWithUpdate = NextGenerationFirewallVNetLocalRulestackResource{}

func (r NextGenerationFirewallVNetLocalRulestackResource) ModelObject() interface{} {
	return &NextGenerationFirewallVnetLocalRulestackModel{}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"rule_stack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: localrulestacks.ValidateLocalRulestackID,
		},

		"network_profile": schema.VnetNetworkProfileSchema(),

		// Optional
		"dns_settings": schema.DNSSettingsSchema(),

		"destination_nat": schema.DestinationNATSchema(),

		"tags": commonschema.Tags(),
	}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"plan": schema.PlanSchema(),

		"pan_etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) ResourceType() string {
	return "azurerm_palo_alto_next_generation_firewall_vnet_local_rulestack"
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient
			localRulestackClient := metadata.Client.PaloAlto.LocalRulestacksClient

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

func (r NextGenerationFirewallVNetLocalRulestackResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient

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

			props := existing.Model.Properties

			state.DNSSettings = schema.FlattenDNSSettings(props.DnsSettings)

			netProfile := schema.FlattenNetworkProfileVnet(props.NetworkProfile)

			state.NetworkProfile = []schema.NetworkProfileVnet{netProfile}

			state.FrontEnd = schema.FlattenDestinationNAT(props.FrontEndSettings)

			state.RuleStackId = pointer.From(props.AssociatedRulestack.ResourceId)

			state.PanEtag = pointer.From(props.PanEtag)

			state.PlanData = schema.FlattenPlanData(props.PlanData)

			state.Tags = tags.Flatten(existing.Model.Tags)

			return metadata.Encode(&state)
		},
	}
}

func (r NextGenerationFirewallVNetLocalRulestackResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient

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
		Timeout: 2 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.FirewallClient

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

			if metadata.ResourceData.HasChange("rule_stack_id") {
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

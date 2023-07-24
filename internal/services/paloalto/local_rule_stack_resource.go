package paloalto

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRuleStack struct{}

var _ sdk.ResourceWithUpdate = LocalRuleStack{}

const (
	RuleStackSecurityServicesCustom       string = "Custom"
	RuleStackSecurityServicesNone         string = "None"
	RuleStackSecurityServicesBestPractice string = "BestPractice"
)

type LocalRuleStackModel struct {
	Name                 string `tfschema:"name"`
	ResourceGroupName    string `tfschema:"resource_group_name"`
	Location             string `tfschema:"location"`
	AntiSpywareProfile   string `tfschema:"anti_spyware_profile"`
	AntiVirusProfile     string `tfschema:"anti_virus_profile"`
	DNSSubscription      string `tfschema:"dns_subscription"`
	FileBlockingProfile  string `tfschema:"file_blocking_profile"`
	URLFilteringProfile  string `tfschema:"url_filtering_profile"`
	VulnerabilityProfile string `tfschema:"vulnerability_profile"`
	Description          string `tfschema:"description"`
}

func (r LocalRuleStack) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return localrulestacks.ValidateLocalRulestackID
}

func (r LocalRuleStack) ResourceType() string {
	return "azurerm_palo_alto_local_rule_stack"
}

func (r LocalRuleStack) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LocalRuleStackName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"vulnerability_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RuleStackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesNone,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"anti_spyware_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RuleStackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesNone,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"anti_virus_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RuleStackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesNone,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"url_filtering_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RuleStackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesNone,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"file_blocking_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RuleStackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesNone,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"dns_subscription": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  RuleStackSecurityServicesNone,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesNone,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r LocalRuleStack) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRuleStack) ModelObject() interface{} {
	return &LocalRuleStackModel{}
}

func (r LocalRuleStack) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			model := LocalRuleStackModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := localrulestacks.NewLocalRulestackID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			localRuleStack := localrulestacks.LocalRulestackResource{
				Location: model.Location,
				Properties: localrulestacks.RulestackProperties{
					DefaultMode: pointer.To(localrulestacks.DefaultModeNONE),
					Description: pointer.To(model.Description),
					Scope:       pointer.To(localrulestacks.ScopeTypeLOCAL),
					SecurityServices: &localrulestacks.SecurityServices{
						AntiSpywareProfile:   pointer.To(model.AntiSpywareProfile),
						AntiVirusProfile:     pointer.To(model.AntiVirusProfile),
						DnsSubscription:      pointer.To(model.DNSSubscription),
						FileBlockingProfile:  pointer.To(model.FileBlockingProfile),
						UrlFilteringProfile:  pointer.To(model.URLFilteringProfile),
						VulnerabilityProfile: pointer.To(model.VulnerabilityProfile),
					},
				},
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, localRuleStack); err != nil {
				return err
			}

			metadata.SetID(id)

			if err = client.CommitThenPoll(ctx, id); err != nil {
				return fmt.Errorf("committing config for %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r LocalRuleStack) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRuleStackModel

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			props := existing.Model.Properties

			state.Name = id.LocalRulestackName
			state.ResourceGroupName = id.ResourceGroupName
			state.Description = pointer.From(props.Description)
			state.Location = location.Normalize(existing.Model.Location)

			if secServices := props.SecurityServices; secServices != nil {
				state.VulnerabilityProfile = pointer.From(secServices.VulnerabilityProfile)
				state.AntiSpywareProfile = pointer.From(secServices.AntiSpywareProfile)
				state.AntiVirusProfile = pointer.From(secServices.AntiVirusProfile)
				state.FileBlockingProfile = pointer.From(secServices.FileBlockingProfile)
				state.URLFilteringProfile = pointer.From(secServices.UrlFilteringProfile)
				state.DNSSubscription = pointer.From(secServices.DnsSubscription)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRuleStack) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient
			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
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

func (r LocalRuleStack) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := LocalRuleStackModel{}

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

			localRuleStack := *existing.Model
			update := localRuleStack.Properties

			if metadata.ResourceData.HasChange("description") {
				update.Description = pointer.To(model.Description)
			}

			secServices := pointer.From(update.SecurityServices)

			if metadata.ResourceData.HasChange("dns_subscription") {
				secServices.DnsSubscription = pointer.To(model.DNSSubscription)
			}

			if metadata.ResourceData.HasChange("vulnerability_profile") {
				secServices.VulnerabilityProfile = pointer.To(model.VulnerabilityProfile)
			}

			if metadata.ResourceData.HasChange("anti_spyware_profile") {
				secServices.AntiSpywareProfile = pointer.To(model.AntiSpywareProfile)
			}

			if metadata.ResourceData.HasChange("anti_virus_profile") {
				secServices.AntiVirusProfile = pointer.To(model.AntiVirusProfile)
			}

			if metadata.ResourceData.HasChange("url_filtering_profile") {
				secServices.UrlFilteringProfile = pointer.To(model.URLFilteringProfile)
			}

			if metadata.ResourceData.HasChange("file_blocking_profile") {
				secServices.FileBlockingProfile = pointer.To(model.FileBlockingProfile)
			}

			update.SecurityServices = pointer.To(secServices)

			localRuleStack.Properties = update

			if err = client.CreateOrUpdateThenPoll(ctx, *id, localRuleStack); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			if err = client.CommitThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("committing config for %s: %+v", *id, err)
			}

			return nil
		},
	}
}

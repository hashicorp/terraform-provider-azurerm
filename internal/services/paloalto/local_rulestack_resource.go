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
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
	return "azurerm_palo_alto_local_rulestack"
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
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"anti_spyware_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"anti_virus_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"url_filtering_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"file_blocking_profile": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
				RuleStackSecurityServicesBestPractice,
			}, false),
		},

		"dns_subscription": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				RuleStackSecurityServicesCustom,
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
			client := metadata.Client.PaloAlto.Client.LocalRulestacks

			model := LocalRuleStackModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := localrulestacks.NewLocalRulestackID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)
			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			secServices := localrulestacks.SecurityServices{
				AntiSpywareProfile:   pointer.To(RuleStackSecurityServicesNone),
				AntiVirusProfile:     pointer.To(RuleStackSecurityServicesNone),
				DnsSubscription:      pointer.To(RuleStackSecurityServicesNone),
				FileBlockingProfile:  pointer.To(RuleStackSecurityServicesNone),
				UrlFilteringProfile:  pointer.To(RuleStackSecurityServicesNone),
				VulnerabilityProfile: pointer.To(RuleStackSecurityServicesNone),
			}

			if model.AntiSpywareProfile != "" {
				secServices.AntiSpywareProfile = pointer.To(model.AntiSpywareProfile)
			}
			if model.AntiVirusProfile != "" {
				secServices.AntiVirusProfile = pointer.To(model.AntiVirusProfile)
			}
			if model.DNSSubscription != "" {
				secServices.DnsSubscription = pointer.To(model.DNSSubscription)
			}
			if model.FileBlockingProfile != "" {
				secServices.FileBlockingProfile = pointer.To(model.FileBlockingProfile)
			}
			if model.URLFilteringProfile != "" {
				secServices.UrlFilteringProfile = pointer.To(model.URLFilteringProfile)
			}
			if model.VulnerabilityProfile != "" {
				secServices.VulnerabilityProfile = pointer.To(model.VulnerabilityProfile)
			}

			localRuleStack := localrulestacks.LocalRulestackResource{
				Location: location.Normalize(model.Location),
				Properties: localrulestacks.RulestackProperties{
					DefaultMode:      pointer.To(localrulestacks.DefaultModeNONE),
					Description:      pointer.To(model.Description),
					Scope:            pointer.To(localrulestacks.ScopeTypeLOCAL),
					SecurityServices: pointer.To(secServices),
				},
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, localRuleStack); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r LocalRuleStack) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRulestacks

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

			state.Name = id.LocalRulestackName
			state.ResourceGroupName = id.ResourceGroupName
			if model := existing.Model; model != nil {
				props := model.Properties

				state.Description = pointer.From(props.Description)
				state.Location = location.Normalize(existing.Model.Location)

				if secServices := props.SecurityServices; secServices != nil {
					if v := pointer.From(secServices.VulnerabilityProfile); v != RuleStackSecurityServicesNone {
						state.VulnerabilityProfile = v
					}
					if v := pointer.From(secServices.AntiSpywareProfile); v != RuleStackSecurityServicesNone {
						state.AntiSpywareProfile = v
					}
					if v := pointer.From(secServices.AntiVirusProfile); v != RuleStackSecurityServicesNone {
						state.AntiVirusProfile = v
					}
					if v := pointer.From(secServices.FileBlockingProfile); v != RuleStackSecurityServicesNone {
						state.FileBlockingProfile = v
					}
					if v := pointer.From(secServices.UrlFilteringProfile); v != RuleStackSecurityServicesNone {
						state.URLFilteringProfile = v
					}
					if v := pointer.From(secServices.DnsSubscription); v != RuleStackSecurityServicesNone {
						state.DNSSubscription = v
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRuleStack) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.Client.LocalRulestacks
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
			client := metadata.Client.PaloAlto.Client.LocalRulestacks

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByID(id.ID())
			defer locks.UnlockByID(id.ID())

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
				if model.DNSSubscription != "" {
					secServices.DnsSubscription = pointer.To(model.DNSSubscription)
				} else {
					secServices.DnsSubscription = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("vulnerability_profile") {
				if model.VulnerabilityProfile != "" {
					secServices.VulnerabilityProfile = pointer.To(model.VulnerabilityProfile)
				} else {
					secServices.VulnerabilityProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("anti_spyware_profile") {
				if model.AntiSpywareProfile != "" {
					secServices.AntiSpywareProfile = pointer.To(model.AntiSpywareProfile)
				} else {
					secServices.AntiSpywareProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("anti_virus_profile") {
				if model.AntiVirusProfile != "" {
					secServices.AntiVirusProfile = pointer.To(model.AntiVirusProfile)
				} else {
					secServices.AntiVirusProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("url_filtering_profile") {
				if model.URLFilteringProfile != "" {
					secServices.UrlFilteringProfile = pointer.To(model.URLFilteringProfile)
				} else {
					secServices.UrlFilteringProfile = pointer.To(RuleStackSecurityServicesNone)
				}
			}

			if metadata.ResourceData.HasChange("file_blocking_profile") {
				if model.FileBlockingProfile != "" {
					secServices.FileBlockingProfile = pointer.To(model.FileBlockingProfile)
				} else {
					secServices.FileBlockingProfile = pointer.To(RuleStackSecurityServicesNone)
				}
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

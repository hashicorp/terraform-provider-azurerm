package cognitive

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/commitmentplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CognitiveCommitmentPlanModel struct {
	Name               string            `tfschema:"name"`
	AutoRenewEnabled   bool              `tfschema:"auto_renew_enabled"`
	CognitiveAccountId string            `tfschema:"cognitive_account_id"`
	CurrentTier        string            `tfschema:"current_tier"`
	HostingModel       string            `tfschema:"hosting_model"`
	RenewalTier        string            `tfschema:"renewal_tier"`
	PlanType           string            `tfschema:"plan_type"`
	Tags               map[string]string `tfschema:"tags"`
}

type CognitiveCommitmentPlanResource struct{}

var _ sdk.ResourceWithUpdate = CognitiveCommitmentPlanResource{}

func (r CognitiveCommitmentPlanResource) ResourceType() string {
	return "azurerm_cognitive_commitment_plan"
}

func (r CognitiveCommitmentPlanResource) ModelObject() interface{} {
	return &CognitiveCommitmentPlanModel{}
}

func (r CognitiveCommitmentPlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commitmentplans.ValidateAccountCommitmentPlanID
}

func (r CognitiveCommitmentPlanResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9_-]{2,64}$"),
				"The name can only include alphanumeric characters, underscores and hyphens, must contain between 2 and 64 characters.",
			),
		},

		"cognitive_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: cognitiveservicesaccounts.ValidateAccountID,
		},

		"hosting_model": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(commitmentplans.PossibleValuesForHostingModel(), false),
		},

		"current_tier": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"auto_renew_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"renewal_tier": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"tags": commonschema.Tags(),
	}
}

func (r CognitiveCommitmentPlanResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CognitiveCommitmentPlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.CommitmentPlansClient

			var model CognitiveCommitmentPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := cognitiveservicesaccounts.ParseAccountID(model.CognitiveAccountId)
			if err != nil {
				return err
			}

			id := commitmentplans.NewAccountCommitmentPlanID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.AccountName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := commitmentplans.CommitmentPlan{
				Properties: &commitmentplans.CommitmentPlanProperties{
					AutoRenew:    pointer.To(model.AutoRenewEnabled),
					Current:      pointer.To(commitmentplans.CommitmentPeriod{Tier: pointer.To(model.CurrentTier)}),
					HostingModel: pointer.To(commitmentplans.HostingModel(model.HostingModel)),
					PlanType:     pointer.To(model.PlanType),
				},
			}

			if model.RenewalTier != "" {
				properties.Properties.Next = pointer.To(commitmentplans.CommitmentPeriod{Tier: pointer.To(model.RenewalTier)})
			}

			if model.Tags != nil {
				properties.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r CognitiveCommitmentPlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.CommitmentPlansClient

			id, err := commitmentplans.ParseAccountCommitmentPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model CognitiveCommitmentPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			properties := resp.Model

			// Set startDate and endDate to nil as API returns error "Both 'startDate' and 'endDate' are read only fields and cannot be specified"
			if properties.Properties.Current != nil {
				properties.Properties.Current.StartDate = nil
				properties.Properties.Current.EndDate = nil
			}
			if properties.Properties.Next != nil {
				properties.Properties.Next.StartDate = nil
				properties.Properties.Next.EndDate = nil
			}

			if metadata.ResourceData.HasChange("auto_renew_enabled") {
				properties.Properties.AutoRenew = &model.AutoRenewEnabled
			}

			if metadata.ResourceData.HasChange("renewal_tier") {
				if model.RenewalTier != "" {
					properties.Properties.Next = pointer.To(commitmentplans.CommitmentPeriod{Tier: pointer.To(model.RenewalTier)})
				} else {
					properties.Properties.Next = nil
				}
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = pointer.To(model.Tags)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CognitiveCommitmentPlanResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.CommitmentPlansClient

			id, err := commitmentplans.ParseAccountCommitmentPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := CognitiveCommitmentPlanModel{
				Name:               id.CommitmentPlanName,
				CognitiveAccountId: cognitiveservicesaccounts.NewAccountID(id.SubscriptionId, id.ResourceGroupName, id.AccountName).ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.AutoRenewEnabled = pointer.From(properties.AutoRenew)
					state.HostingModel = string(pointer.From(properties.HostingModel))
					state.PlanType = pointer.From(properties.PlanType)

					if properties.Current != nil {
						state.CurrentTier = pointer.From(properties.Current.Tier)
					}

					if properties.Next != nil {
						state.RenewalTier = pointer.From(properties.Next.Tier)
					}
				}

				state.Tags = pointer.From(model.Tags)
			}
			return metadata.Encode(&state)
		},
	}
}

func (r CognitiveCommitmentPlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Cognitive.CommitmentPlansClient

			id, err := commitmentplans.ParseAccountCommitmentPlanID(metadata.ResourceData.Id())
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

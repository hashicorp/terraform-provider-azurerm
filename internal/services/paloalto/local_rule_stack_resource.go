package paloalto

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrules"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/localrulestacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type LocalRuleStack struct{}

var _ sdk.ResourceWithUpdate = LocalRuleStack{}

type LocalRuleStackModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`
	DefaultMode       string `tfschema:"default_mode"`
	Description       string `tfschema:"description"`
}

func (r LocalRuleStack) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return localrulestacks.ValidateLocalRuleStackID
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
			ValidateFunc: validate.LocalRuleStackName, // TODO - Need validation rules either 30 or 128 alphanumeric, no `-` at beginning or end.
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"default_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(localrulestacks.DefaultModeFIREWALL),
			ValidateFunc: validation.StringInSlice(localrulestacks.PossibleValuesForDefaultMode(), false),
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
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRuleStacksClient

			model := LocalRuleStackModel{}

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := localrulestacks.NewLocalRuleStackID(metadata.Client.Account.SubscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			localRuleSet := localrulestacks.LocalRulestackResource{
				Location: model.Location,
				Properties: localrulestacks.RulestackProperties{
					DefaultMode: pointer.To(localrulestacks.DefaultMode(model.DefaultMode)),
					Description: pointer.To(model.Description),
					Scope:       pointer.To(localrulestacks.ScopeTypeLOCAL),
				},
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, localRuleSet); err != nil {
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
			client := metadata.Client.PaloAlto.LocalRuleStacksClient

			id, err := localrulestacks.ParseLocalRuleStackID(metadata.ResourceData.Id())
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

			state.Name = id.LocalRuleStackName
			state.ResourceGroupName = id.ResourceGroupName
			state.DefaultMode = string(pointer.From(props.DefaultMode))
			state.Description = pointer.From(props.Description)
			state.Location = location.Normalize(existing.Model.Location)

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRuleStack) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRuleStacksClient

			id, err := localrulestacks.ParseLocalRuleStackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// The API won't delete the RuleStack if the default rule it creates exists, so to continue development, we'll hack around that for now.
			// This must not ship.
			if err = hackForCheckForDefaultRuleAndRemoveIt(ctx, metadata, *id); err != nil {
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

			return nil
		},
	}
}

func hackForCheckForDefaultRuleAndRemoveIt(ctx context.Context, metadata sdk.ResourceMetaData, localRuleStackID localrulestacks.LocalRuleStackId) error {
	client := metadata.Client.PaloAlto.LocalRulesClient

	id := localrules.NewLocalRuleStackID(localRuleStackID.SubscriptionId, localRuleStackID.ResourceGroupName, localRuleStackID.LocalRuleStackName)

	resp, err := client.ListByLocalRulestacksComplete(ctx, id)
	if err != nil {
		return fmt.Errorf("listing rules for %s: %+v", id, err)
	}

	if len(resp.Items) > 0 {
		if len(resp.Items) == 1 && resp.Items[0].Properties.Priority != nil && *resp.Items[0].Properties.Priority == 1000000 {
			ruleId := localrules.NewLocalRuleID(id.SubscriptionId, id.ResourceGroupName, id.LocalRuleStackName, "1000000")
			if err != nil {
				return err
			}
			if err := client.DeleteThenPoll(ctx, ruleId); err != nil {
				return fmt.Errorf("deleting default rule for %s: %+v", id, err)
			}
		}
	}

	return nil
}

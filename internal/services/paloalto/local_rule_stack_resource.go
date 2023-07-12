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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type LocalRulestack struct{}

var _ sdk.ResourceWithUpdate = LocalRulestack{}

type LocalRulestackModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	Location          string `tfschema:"location"`
	Description       string `tfschema:"description"`
}

func (r LocalRulestack) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return localrulestacks.ValidateLocalRulestackID
}

func (r LocalRulestack) ResourceType() string {
	return "azurerm_palo_alto_local_rule_stack"
}

func (r LocalRulestack) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LocalRulestackName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (r LocalRulestack) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LocalRulestack) ModelObject() interface{} {
	return &LocalRulestackModel{}
}

func (r LocalRulestack) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			model := LocalRulestackModel{}

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

			localRulestack := localrulestacks.LocalRulestackResource{
				Location: model.Location,
				Properties: localrulestacks.RulestackProperties{
					DefaultMode: pointer.To(localrulestacks.DefaultModeNONE),
					Description: pointer.To(model.Description),
					Scope:       pointer.To(localrulestacks.ScopeTypeLOCAL),
				},
			}

			if err = client.CreateOrUpdateThenPoll(ctx, id, localRulestack); err != nil {
				return err
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r LocalRulestack) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LocalRulestackModel

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

			return metadata.Encode(&state)
		},
	}
}

func (r LocalRulestack) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 3 * time.Hour,
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

func (r LocalRulestack) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PaloAlto.LocalRulestacksClient

			id, err := localrulestacks.ParseLocalRulestackID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			model := LocalRulestackModel{}

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

			localRulestack := *existing.Model

			if metadata.ResourceData.HasChange("description") {
				localRulestack.Properties.Description = pointer.To(model.Description)
			}

			if err = client.CreateOrUpdateThenPoll(ctx, *id, localRulestack); err != nil {
				return err
			}

			return nil
		},
	}
}

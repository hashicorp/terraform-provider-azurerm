// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_actions_task_definition -service-package-name storage -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type StorageActionsTaskDefinitionResource struct{}

var (
	_ sdk.ResourceWithIdentity      = StorageActionsTaskDefinitionResource{}
	_ sdk.ResourceWithUpdate        = StorageActionsTaskDefinitionResource{}
	_ sdk.ResourceWithCustomizeDiff = StorageActionsTaskDefinitionResource{}
)

type StorageActionsTaskDefinitionModel struct {
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Location          string                                     `tfschema:"location"`
	Description       string                                     `tfschema:"description"`
	Action            []StorageActionsTaskDefinitionActionModel  `tfschema:"action"`
	Enabled           bool                                       `tfschema:"enabled"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags              map[string]string                          `tfschema:"tags"`
}

type StorageActionsTaskDefinitionActionModel struct {
	If   []StorageActionsTaskDefinitionIfModel   `tfschema:"if"`
	Else []StorageActionsTaskDefinitionElseModel `tfschema:"else"`
}

type StorageActionsTaskDefinitionIfModel struct {
	Condition string                                       `tfschema:"condition"`
	Operation []StorageActionsTaskDefinitionOperationModel `tfschema:"operation"`
}

type StorageActionsTaskDefinitionElseModel struct {
	Operation []StorageActionsTaskDefinitionOperationModel `tfschema:"operation"`
}

type StorageActionsTaskDefinitionOperationModel struct {
	Name       string            `tfschema:"name"`
	OnFailure  string            `tfschema:"on_failure"`
	OnSuccess  string            `tfschema:"on_success"`
	Parameters map[string]string `tfschema:"parameters"`
}

func (r StorageActionsTaskDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagetasks.ValidateStorageTaskID
}

func (r StorageActionsTaskDefinitionResource) Identity() resourceids.ResourceId {
	return &storagetasks.StorageTaskId{}
}

func (r StorageActionsTaskDefinitionResource) ModelObject() interface{} {
	return &StorageActionsTaskDefinitionModel{}
}

func (r StorageActionsTaskDefinitionResource) ResourceType() string {
	return "azurerm_storage_actions_task_definition"
}

func (StorageActionsTaskDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z0-9]{3,18}$`),
				"The `name` must be between 3 and 18 characters in length and use numbers and lower-case letters only.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"action": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"if": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"condition": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"operation": storageActionsTaskDefinitionOperationSchema(),
							},
						},
					},
					"else": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{ // kept as a block to mirror the `if` block structure for HCL readability
							Schema: map[string]*pluginsdk.Schema{
								"operation": storageActionsTaskDefinitionOperationSchema(),
							},
						},
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityRequired(),

		"tags": commonschema.Tags(),
	}
}

func (StorageActionsTaskDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StorageActionsTaskDefinitionResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageActionsTaskDefinitionModel
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(model.Action) == 0 || len(model.Action[0].If) == 0 {
				return nil
			}

			if err := validateNoDeleteBlobWithOtherOperations(model.Action[0].If[0].Operation, "if"); err != nil {
				return err
			}

			if len(model.Action[0].Else) > 0 {
				if err := validateNoDeleteBlobWithOtherOperations(model.Action[0].Else[0].Operation, "else"); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func validateNoDeleteBlobWithOtherOperations(ops []StorageActionsTaskDefinitionOperationModel, blockName string) error {
	if len(ops) <= 1 {
		return nil
	}

	// A `DeleteBlob` operation must be the only operation within its `if` or `else` block.
	for _, op := range ops {
		if op.Name == string(storagetasks.StorageTaskOperationNameDeleteBlob) {
			return fmt.Errorf("`action.0.%s.0.operation`: a `DeleteBlob` operation cannot be combined with other operations in the same block", blockName)
		}
	}
	return nil
}

func storageActionsTaskDefinitionOperationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(storagetasks.PossibleValuesForStorageTaskOperationName(), false),
				},
				"on_failure": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(storagetasks.PossibleValuesForOnFailure(), false),
				},
				"on_success": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(storagetasks.PossibleValuesForOnSuccess(), false),
				},
				"parameters": {
					Type:     pluginsdk.TypeMap,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func (r StorageActionsTaskDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageActionsTaskDefinitionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Storage.StorageTasksClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := storagetasks.NewStorageTaskID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding identity: %+v", err)
			}

			properties := storagetasks.StorageTask{
				Location: location.Normalize(model.Location),
				Identity: pointer.From(expandedIdentity),
				Properties: storagetasks.StorageTaskProperties{
					Description: model.Description,
					Enabled:     model.Enabled,
					Action:      expandStorageActionsTaskDefinitionAction(model.Action),
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateThenPoll(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func (r StorageActionsTaskDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageTasksClient

			id, err := storagetasks.ParseStorageTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			return r.flatten(metadata, id, model)
		},
	}
}

func (r StorageActionsTaskDefinitionResource) flatten(metadata sdk.ResourceMetaData, id *storagetasks.StorageTaskId, model *storagetasks.StorageTask) error {
	state := StorageActionsTaskDefinitionModel{
		Name:              id.StorageTaskName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.Tags = pointer.From(model.Tags)

		flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(pointer.To(model.Identity))
		if err != nil {
			return fmt.Errorf("flattening identity: %+v", err)
		}
		state.Identity = flattenedIdentity

		state.Description = model.Properties.Description
		state.Enabled = model.Properties.Enabled
		state.Action = flattenStorageActionsTaskDefinitionAction(model.Properties.Action)
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}
	return metadata.Encode(&state)
}

func (r StorageActionsTaskDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageTasksClient

			// parse the existing Resource ID from the State
			id, err := storagetasks.ParseStorageTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageActionsTaskDefinitionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = model.Description
			}

			if metadata.ResourceData.HasChange("enabled") {
				payload.Properties.Enabled = model.Enabled
			}

			if metadata.ResourceData.HasChange("action") {
				payload.Properties.Action = expandStorageActionsTaskDefinitionAction(model.Action)
			}

			if metadata.ResourceData.HasChange("identity") {
				expandedIdentity, err := identity.ExpandLegacySystemAndUserAssignedMapFromModel(model.Identity)
				if err != nil {
					return fmt.Errorf("expanding identity: %+v", err)
				}
				payload.Identity = pointer.From(expandedIdentity)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r StorageActionsTaskDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.StorageTasksClient

			id, err := storagetasks.ParseStorageTaskID(metadata.ResourceData.Id())
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

func expandStorageActionsTaskDefinitionAction(input []StorageActionsTaskDefinitionActionModel) storagetasks.StorageTaskAction {
	action := input[0]

	result := storagetasks.StorageTaskAction{
		If: storagetasks.IfCondition{
			Condition:  action.If[0].Condition,
			Operations: expandStorageActionsTaskDefinitionOperations(action.If[0].Operation),
		},
	}

	if len(action.Else) > 0 {
		result.Else = &storagetasks.ElseCondition{
			Operations: expandStorageActionsTaskDefinitionOperations(action.Else[0].Operation),
		}
	}

	return result
}

func expandStorageActionsTaskDefinitionOperations(input []StorageActionsTaskDefinitionOperationModel) []storagetasks.StorageTaskOperation {
	result := make([]storagetasks.StorageTaskOperation, 0, len(input))
	for _, op := range input {
		operation := storagetasks.StorageTaskOperation{
			Name:      storagetasks.StorageTaskOperationName(op.Name),
			OnFailure: pointer.ToEnum[storagetasks.OnFailure](op.OnFailure),
			OnSuccess: pointer.ToEnum[storagetasks.OnSuccess](op.OnSuccess),
		}

		if len(op.Parameters) > 0 {
			operation.Parameters = &op.Parameters
		}

		result = append(result, operation)
	}
	return result
}

func flattenStorageActionsTaskDefinitionAction(input storagetasks.StorageTaskAction) []StorageActionsTaskDefinitionActionModel {
	action := StorageActionsTaskDefinitionActionModel{
		If: []StorageActionsTaskDefinitionIfModel{
			{
				Condition: input.If.Condition,
				Operation: flattenStorageActionsTaskDefinitionOperations(input.If.Operations),
			},
		},
	}

	if input.Else != nil {
		action.Else = []StorageActionsTaskDefinitionElseModel{
			{
				Operation: flattenStorageActionsTaskDefinitionOperations(input.Else.Operations),
			},
		}
	}

	return []StorageActionsTaskDefinitionActionModel{action}
}

func flattenStorageActionsTaskDefinitionOperations(input []storagetasks.StorageTaskOperation) []StorageActionsTaskDefinitionOperationModel {
	result := make([]StorageActionsTaskDefinitionOperationModel, 0, len(input))
	for _, op := range input {
		operation := StorageActionsTaskDefinitionOperationModel{
			Name:       string(op.Name),
			OnFailure:  string(pointer.From(op.OnFailure)),
			OnSuccess:  string(pointer.From(op.OnSuccess)),
			Parameters: pointer.From(op.Parameters),
		}
		result = append(result, operation)
	}

	return result
}

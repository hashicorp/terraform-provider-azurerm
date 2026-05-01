// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package storageactions

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

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_actions_task -service-package-name storage -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary"

type StorageActionsTaskResource struct{}

var (
	_ sdk.ResourceWithIdentity = StorageActionsTaskResource{}
	_ sdk.ResourceWithUpdate   = StorageActionsTaskResource{}
)

type StorageActionsTaskModel struct {
	Name              string                                     `tfschema:"name"`
	ResourceGroupName string                                     `tfschema:"resource_group_name"`
	Location          string                                     `tfschema:"location"`
	Description       string                                     `tfschema:"description"`
	Action            []StorageActionsTaskActionModel            `tfschema:"action"`
	Enabled           bool                                       `tfschema:"enabled"`
	Identity          []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	Tags              map[string]string                          `tfschema:"tags"`
}

type StorageActionsTaskActionModel struct {
	If   []StorageActionsTaskIfModel   `tfschema:"if"`
	Else []StorageActionsTaskElseModel `tfschema:"else"`
}

type StorageActionsTaskIfModel struct {
	Condition string                             `tfschema:"condition"`
	Operation []StorageActionsTaskOperationModel `tfschema:"operation"`
}

type StorageActionsTaskElseModel struct {
	Operation []StorageActionsTaskOperationModel `tfschema:"operation"`
}

type StorageActionsTaskOperationModel struct {
	Name       string            `tfschema:"name"`
	OnFailure  string            `tfschema:"on_failure"`
	OnSuccess  string            `tfschema:"on_success"`
	Parameters map[string]string `tfschema:"parameters"`
}

func (r StorageActionsTaskResource) ModelObject() interface{} {
	return &StorageActionsTaskModel{}
}

func (r StorageActionsTaskResource) ResourceType() string {
	return "azurerm_storage_actions_task"
}

func (StorageActionsTaskResource) Arguments() map[string]*pluginsdk.Schema {
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
								"operation": storageActionsTaskOperationSchema(),
							},
						},
					},
					"else": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{ // kept as a block to mirror the `if` block structure for HCL readability
							Schema: map[string]*pluginsdk.Schema{
								"operation": storageActionsTaskOperationSchema(),
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

func (StorageActionsTaskResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func storageActionsTaskOperationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(storagetasks.PossibleValuesForStorageTaskOperationName(), false),
				},
				"on_failure": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(storagetasks.PossibleValuesForOnFailure(), false),
				},
				"on_success": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
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

func (r StorageActionsTaskResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model StorageActionsTaskModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.StorageActions.StorageTasksClient
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
					Action:      expandStorageActionsTaskAction(model.Action),
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

func (r StorageActionsTaskResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageActions.StorageTasksClient

			// parse the existing Resource ID from the State
			id, err := storagetasks.ParseStorageTaskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StorageActionsTaskModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Retrieve the existing resource and patch it. We use Create (PUT)
			// rather than Update (PATCH) because the PATCH model treats absent
			// optional sub-fields (e.g. `action.else`) as "no change", which
			// would prevent users from clearing them.
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
				payload.Properties.Action = expandStorageActionsTaskAction(model.Action)
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

func (StorageActionsTaskResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageActions.StorageTasksClient

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

			state := StorageActionsTaskModel{
				Name:              id.StorageTaskName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
				Tags:              pointer.From(model.Tags),
			}

			flattenedIdentity, err := identity.FlattenLegacySystemAndUserAssignedMapToModel(pointer.To(model.Identity))
			if err != nil {
				return fmt.Errorf("flattening identity: %+v", err)
			}
			state.Identity = flattenedIdentity

			state.Description = model.Properties.Description
			state.Enabled = model.Properties.Enabled
			state.Action = flattenStorageActionsTaskAction(model.Properties.Action)

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func (r StorageActionsTaskResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.StorageActions.StorageTasksClient

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

func (r StorageActionsTaskResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagetasks.ValidateStorageTaskID
}

func (r StorageActionsTaskResource) Identity() resourceids.ResourceId {
	return &storagetasks.StorageTaskId{}
}

func expandStorageActionsTaskAction(input []StorageActionsTaskActionModel) storagetasks.StorageTaskAction {
	action := input[0]

	result := storagetasks.StorageTaskAction{
		If: storagetasks.IfCondition{
			Condition:  action.If[0].Condition,
			Operations: expandStorageActionsTaskOperations(action.If[0].Operation),
		},
	}

	if len(action.Else) > 0 {
		result.Else = &storagetasks.ElseCondition{
			Operations: expandStorageActionsTaskOperations(action.Else[0].Operation),
		}
	}

	return result
}

func expandStorageActionsTaskOperations(input []StorageActionsTaskOperationModel) []storagetasks.StorageTaskOperation {
	result := make([]storagetasks.StorageTaskOperation, 0, len(input))
	for _, op := range input {
		operation := storagetasks.StorageTaskOperation{
			Name: storagetasks.StorageTaskOperationName(op.Name),
		}

		if op.OnFailure != "" {
			onFailure := storagetasks.OnFailure(op.OnFailure)
			operation.OnFailure = &onFailure
		}

		if op.OnSuccess != "" {
			onSuccess := storagetasks.OnSuccess(op.OnSuccess)
			operation.OnSuccess = &onSuccess
		}

		if len(op.Parameters) > 0 {
			operation.Parameters = &op.Parameters
		}

		result = append(result, operation)
	}
	return result
}

func flattenStorageActionsTaskAction(input storagetasks.StorageTaskAction) []StorageActionsTaskActionModel {
	action := StorageActionsTaskActionModel{
		If: []StorageActionsTaskIfModel{
			{
				Condition: input.If.Condition,
				Operation: flattenStorageActionsTaskOperations(input.If.Operations),
			},
		},
	}

	if input.Else != nil {
		action.Else = []StorageActionsTaskElseModel{
			{
				Operation: flattenStorageActionsTaskOperations(input.Else.Operations),
			},
		}
	}

	return []StorageActionsTaskActionModel{action}
}

func flattenStorageActionsTaskOperations(input []storagetasks.StorageTaskOperation) []StorageActionsTaskOperationModel {
	result := make([]StorageActionsTaskOperationModel, 0, len(input))
	for _, op := range input {
		operation := StorageActionsTaskOperationModel{
			Name:       string(op.Name),
			OnFailure:  string(pointer.From(op.OnFailure)),
			OnSuccess:  string(pointer.From(op.OnSuccess)),
			Parameters: pointer.From(op.Parameters),
		}
		result = append(result, operation)
	}

	return result
}

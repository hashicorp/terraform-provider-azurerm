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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storagetaskassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storageactions/2023-01-01/storagetasks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name storage_actions_task_assignment -service-package-name storage -properties "name" -compare-values "subscription_id:storage_account_id,resource_group_name:storage_account_id,storage_account_name:storage_account_id"

type StorageActionsTaskAssignmentResource struct{}

var (
	_ sdk.ResourceWithUpdate        = StorageActionsTaskAssignmentResource{}
	_ sdk.ResourceWithCustomizeDiff = StorageActionsTaskAssignmentResource{}
	_ sdk.ResourceWithIdentity      = StorageActionsTaskAssignmentResource{}
)

type StorageActionsTaskAssignmentModel struct {
	Name             string                                         `tfschema:"name"`
	StorageAccountId string                                         `tfschema:"storage_account_id"`
	Description      string                                         `tfschema:"description"`
	Enabled          bool                                           `tfschema:"enabled"`
	ExecutionContext []StorageActionsTaskAssignmentExecutionContext `tfschema:"execution_context"`
	ReportPrefix     string                                         `tfschema:"report_prefix"`
	TaskId           string                                         `tfschema:"task_id"`
}

type StorageActionsTaskAssignmentExecutionContext struct {
	Trigger []StorageActionsTaskAssignmentTrigger `tfschema:"trigger"`
	Target  []StorageActionsTaskAssignmentTarget  `tfschema:"target"`
}

type StorageActionsTaskAssignmentTrigger struct {
	Type      string `tfschema:"type"`
	EndBy     string `tfschema:"end_by"`
	Interval  int64  `tfschema:"interval"`
	StartFrom string `tfschema:"start_from"`
	StartOn   string `tfschema:"start_on"`
}

type StorageActionsTaskAssignmentTarget struct {
	ExcludePrefix []string `tfschema:"exclude_prefix"`
	Prefix        []string `tfschema:"prefix"`
}

func (StorageActionsTaskAssignmentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-z][a-z0-9]{2,23}$`),
				"The name of the storage task assignment within the specified resource group. Storage task assignment names must be between 3 and 24 characters in length and use numbers and lower-case letters only.",
			),
		},
		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
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
		"execution_context": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"trigger": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(storagetaskassignments.PossibleValuesForTriggerType(), false),
								},
								"end_by": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
								"interval": {
									Type:         pluginsdk.TypeInt,
									Optional:     true,
									ValidateFunc: validation.IntAtLeast(1),
								},
								"start_from": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
								"start_on": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
							},
						},
					},
					"target": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"prefix": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"exclude_prefix": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},
		"report_prefix": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"task_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storagetasks.ValidateStorageTaskID,
		},
	}
}

func (StorageActionsTaskAssignmentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (StorageActionsTaskAssignmentResource) ModelObject() interface{} {
	return &StorageActionsTaskAssignmentModel{}
}

func (StorageActionsTaskAssignmentResource) ResourceType() string {
	return "azurerm_storage_actions_task_assignment"
}

func (r StorageActionsTaskAssignmentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.StorageTaskAssignments

			var plan StorageActionsTaskAssignmentModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountId, err := commonids.ParseStorageAccountID(plan.StorageAccountId)
			if err != nil {
				return err
			}

			id := storagetaskassignments.NewStorageTaskAssignmentID(accountId.SubscriptionId, accountId.ResourceGroupName, accountId.StorageAccountName, plan.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := storagetaskassignments.StorageTaskAssignment{
				Properties: &storagetaskassignments.StorageTaskAssignmentProperties{
					Description:      plan.Description,
					Enabled:          plan.Enabled,
					TaskId:           plan.TaskId,
					ExecutionContext: expandStorageActionsTaskAssignmentExecutionContext(plan.ExecutionContext),
					Report: storagetaskassignments.StorageTaskAssignmentReport{
						Prefix: plan.ReportPrefix,
					},
				},
			}

			if err := client.CreateThenPoll(ctx, id, payload); err != nil {
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

func (r StorageActionsTaskAssignmentResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.StorageTaskAssignments

			id, err := storagetaskassignments.ParseStorageTaskAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var plan StorageActionsTaskAssignmentModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Use Create (PUT) so absent optional fields (e.g. `target`) are cleared.
			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}

			payload := *existing.Model

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.Description = plan.Description
			}
			if metadata.ResourceData.HasChange("enabled") {
				payload.Properties.Enabled = plan.Enabled
			}
			if metadata.ResourceData.HasChange("task_id") {
				payload.Properties.TaskId = plan.TaskId
			}
			if metadata.ResourceData.HasChange("execution_context") {
				payload.Properties.ExecutionContext = expandStorageActionsTaskAssignmentExecutionContext(plan.ExecutionContext)
			}
			if metadata.ResourceData.HasChange("report_prefix") {
				payload.Properties.Report = storagetaskassignments.StorageTaskAssignmentReport{
					Prefix: plan.ReportPrefix,
				}
			}

			if err := client.CreateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
	}
}

func (r StorageActionsTaskAssignmentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.StorageTaskAssignments

			id, err := storagetaskassignments.ParseStorageTaskAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := StorageActionsTaskAssignmentModel{
				Name:             id.StorageTaskAssignmentName,
				StorageAccountId: commonids.NewStorageAccountID(id.SubscriptionId, id.ResourceGroupName, id.StorageAccountName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.Description = props.Description
					state.Enabled = props.Enabled
					state.TaskId = props.TaskId
					state.ExecutionContext = flattenStorageActionsTaskAssignmentExecutionContext(props.ExecutionContext)
					state.ReportPrefix = props.Report.Prefix
				}
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}
			return metadata.Encode(&state)
		},
	}
}

func (r StorageActionsTaskAssignmentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Storage.ResourceManager.StorageTaskAssignments

			id, err := storagetaskassignments.ParseStorageTaskAssignmentID(metadata.ResourceData.Id())
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

func (StorageActionsTaskAssignmentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return storagetaskassignments.ValidateStorageTaskAssignmentID
}

func (StorageActionsTaskAssignmentResource) Identity() resourceids.ResourceId {
	return &storagetaskassignments.StorageTaskAssignmentId{}
}

func (StorageActionsTaskAssignmentResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var config StorageActionsTaskAssignmentModel
			if err := metadata.DecodeDiff(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if len(config.ExecutionContext) == 0 || len(config.ExecutionContext[0].Trigger) == 0 {
				return nil
			}

			t := config.ExecutionContext[0].Trigger[0]

			// API rule: Terminal trigger types (RunOnce, OnSchedule) cannot be changed once set. Only MockRun assignments can transition to RunOnce or OnSchedule.
			const typePath = "execution_context.0.trigger.0.type"
			if metadata.ResourceDiff.HasChange(typePath) {
				old, _ := metadata.ResourceDiff.GetChange(typePath)
				oldType := old.(string)
				if oldType != "" && oldType != string(storagetaskassignments.TriggerTypeMockRun) {
					if err := metadata.ResourceDiff.ForceNew(typePath); err != nil {
						return err
					}
				}
			}

			switch t.Type {
			case string(storagetaskassignments.TriggerTypeOnSchedule):
				if t.StartFrom == "" {
					return fmt.Errorf("`execution_context.0.trigger.0.start_from` is required when `type` is %q", t.Type)
				}
				if t.EndBy == "" {
					return fmt.Errorf("`execution_context.0.trigger.0.end_by` is required when `type` is %q", t.Type)
				}
				if t.Interval == 0 {
					return fmt.Errorf("`execution_context.0.trigger.0.interval` is required when `type` is %q", t.Type)
				}
				if t.StartOn != "" {
					return fmt.Errorf("`execution_context.0.trigger.0.start_on` must not be set when `type` is %q", t.Type)
				}

			case string(storagetaskassignments.TriggerTypeRunOnce), string(storagetaskassignments.TriggerTypeMockRun):
				if t.StartOn == "" {
					return fmt.Errorf("`execution_context.0.trigger.0.start_on` is required when `type` is %q", t.Type)
				}
				if t.StartFrom != "" {
					return fmt.Errorf("`execution_context.0.trigger.0.start_from` must not be set when `type` is %q", t.Type)
				}
				if t.EndBy != "" {
					return fmt.Errorf("`execution_context.0.trigger.0.end_by` must not be set when `type` is %q", t.Type)
				}
				if t.Interval != 0 {
					return fmt.Errorf("`execution_context.0.trigger.0.interval` must not be set when `type` is %q", t.Type)
				}
			}

			return nil
		},
	}
}

func expandStorageActionsTaskAssignmentExecutionContext(input []StorageActionsTaskAssignmentExecutionContext) storagetaskassignments.StorageTaskAssignmentExecutionContext {
	// schema marks the block Required + MaxItems 1, so input always has exactly one element
	ec := input[0]
	return storagetaskassignments.StorageTaskAssignmentExecutionContext{
		Trigger: expandStorageActionsTaskAssignmentTrigger(ec.Trigger),
		Target:  expandStorageActionsTaskAssignmentTarget(ec.Target),
	}
}

func expandStorageActionsTaskAssignmentTrigger(input []StorageActionsTaskAssignmentTrigger) storagetaskassignments.ExecutionTrigger {
	// schema marks the block Required + MaxItems 1, so input always has exactly one element
	t := input[0]

	params := storagetaskassignments.TriggerParameters{}
	if t.EndBy != "" {
		params.EndBy = pointer.To(t.EndBy)
	}
	if t.StartFrom != "" {
		params.StartFrom = pointer.To(t.StartFrom)
	}
	if t.StartOn != "" {
		params.StartOn = pointer.To(t.StartOn)
	}
	if t.Interval != 0 {
		params.Interval = pointer.To(t.Interval)
	}
	// `Days` is the only supported IntervalUnit and is required when type is OnSchedule.
	if t.Type == string(storagetaskassignments.TriggerTypeOnSchedule) {
		params.IntervalUnit = pointer.To(storagetaskassignments.IntervalUnitDays)
	}

	return storagetaskassignments.ExecutionTrigger{
		Type:       storagetaskassignments.TriggerType(t.Type),
		Parameters: params,
	}
}

func expandStorageActionsTaskAssignmentTarget(input []StorageActionsTaskAssignmentTarget) *storagetaskassignments.ExecutionTarget {
	if len(input) == 0 {
		return nil
	}
	t := input[0]
	out := storagetaskassignments.ExecutionTarget{
		Prefix: pointer.To(t.Prefix),
	}
	if len(t.ExcludePrefix) > 0 {
		out.ExcludePrefix = pointer.To(t.ExcludePrefix)
	}
	return &out
}

func flattenStorageActionsTaskAssignmentExecutionContext(input storagetaskassignments.StorageTaskAssignmentExecutionContext) []StorageActionsTaskAssignmentExecutionContext {
	return []StorageActionsTaskAssignmentExecutionContext{
		{
			Trigger: flattenStorageActionsTaskAssignmentTrigger(input.Trigger),
			Target:  flattenStorageActionsTaskAssignmentTarget(input.Target),
		},
	}
}

func flattenStorageActionsTaskAssignmentTrigger(input storagetaskassignments.ExecutionTrigger) []StorageActionsTaskAssignmentTrigger {
	out := StorageActionsTaskAssignmentTrigger{
		Type:      string(input.Type),
		EndBy:     pointer.From(input.Parameters.EndBy),
		StartFrom: pointer.From(input.Parameters.StartFrom),
		StartOn:   pointer.From(input.Parameters.StartOn),
		Interval:  pointer.From(input.Parameters.Interval),
	}
	return []StorageActionsTaskAssignmentTrigger{out}
}

func flattenStorageActionsTaskAssignmentTarget(input *storagetaskassignments.ExecutionTarget) []StorageActionsTaskAssignmentTarget {
	if input == nil {
		return []StorageActionsTaskAssignmentTarget{}
	}
	return []StorageActionsTaskAssignmentTarget{
		{
			Prefix:        pointer.From(input.Prefix),
			ExcludePrefix: pointer.From(input.ExcludePrefix),
		},
	}
}

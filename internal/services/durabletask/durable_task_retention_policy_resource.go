// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package durabletask

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name durable_task_retention_policy -service-package-name durabletask -compare-values "subscription_id:durable_task_scheduler_id,resource_group_name:durable_task_scheduler_id,scheduler_name:durable_task_scheduler_id"

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RetentionPolicyResourceModel struct {
	DurableTaskSchedulerId string                 `tfschema:"durable_task_scheduler_id"`
	RetentionPolicy        []RetentionPolicyModel `tfschema:"retention_policy"`
}

type RetentionPolicyModel struct {
	CanceledRetentionPeriodInDays   int64 `tfschema:"canceled_retention_period_in_days"`
	CompletedRetentionPeriodInDays  int64 `tfschema:"completed_retention_period_in_days"`
	DefaultRetentionPeriodInDays    int64 `tfschema:"default_retention_period_in_days"`
	FailedRetentionPeriodInDays     int64 `tfschema:"failed_retention_period_in_days"`
	TerminatedRetentionPeriodInDays int64 `tfschema:"terminated_retention_period_in_days"`
}

type RetentionPolicyResource struct{}

var (
	_ sdk.Resource             = RetentionPolicyResource{}
	_ sdk.ResourceWithUpdate   = RetentionPolicyResource{}
	_ sdk.ResourceWithIdentity = RetentionPolicyResource{}
)

func (r RetentionPolicyResource) Identity() resourceids.ResourceId {
	return &RetentionPolicyID{}
}

func (r RetentionPolicyResource) ResourceType() string {
	return "azurerm_durable_task_retention_policy"
}

func (r RetentionPolicyResource) ModelObject() interface{} {
	return &RetentionPolicyResourceModel{}
}

func (r RetentionPolicyResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return ValidateRetentionPolicyID
}

func (r RetentionPolicyResource) Arguments() map[string]*pluginsdk.Schema {
	retentionPolicyAtLeastOneOf := []string{
		"retention_policy.0.canceled_retention_period_in_days",
		"retention_policy.0.completed_retention_period_in_days",
		"retention_policy.0.default_retention_period_in_days",
		"retention_policy.0.failed_retention_period_in_days",
		"retention_policy.0.terminated_retention_period_in_days",
	}
	defaultRetentionConflictsWith := []string{
		"retention_policy.0.canceled_retention_period_in_days",
		"retention_policy.0.completed_retention_period_in_days",
		"retention_policy.0.failed_retention_period_in_days",
		"retention_policy.0.terminated_retention_period_in_days",
	}

	return map[string]*pluginsdk.Schema{
		"durable_task_scheduler_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: schedulers.ValidateSchedulerID,
		},

		"retention_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"canceled_retention_period_in_days": {
						Type:          pluginsdk.TypeInt,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(1, 90),
						AtLeastOneOf:  retentionPolicyAtLeastOneOf,
						ConflictsWith: []string{"retention_policy.0.default_retention_period_in_days"},
					},

					"completed_retention_period_in_days": {
						Type:          pluginsdk.TypeInt,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(1, 90),
						AtLeastOneOf:  retentionPolicyAtLeastOneOf,
						ConflictsWith: []string{"retention_policy.0.default_retention_period_in_days"},
					},

					"default_retention_period_in_days": {
						Type:          pluginsdk.TypeInt,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(1, 90),
						AtLeastOneOf:  retentionPolicyAtLeastOneOf,
						ConflictsWith: defaultRetentionConflictsWith,
					},

					"failed_retention_period_in_days": {
						Type:          pluginsdk.TypeInt,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(1, 90),
						AtLeastOneOf:  retentionPolicyAtLeastOneOf,
						ConflictsWith: []string{"retention_policy.0.default_retention_period_in_days"},
					},

					"terminated_retention_period_in_days": {
						Type:          pluginsdk.TypeInt,
						Optional:      true,
						ValidateFunc:  validation.IntBetween(1, 90),
						AtLeastOneOf:  retentionPolicyAtLeastOneOf,
						ConflictsWith: []string{"retention_policy.0.default_retention_period_in_days"},
					},
				},
			},
		},
	}
}

func (r RetentionPolicyResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r RetentionPolicyResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.RetentionPoliciesClient

			var model RetentionPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parsedId, err := schedulers.ParseSchedulerID(model.DurableTaskSchedulerId)
			if err != nil {
				return fmt.Errorf("parsing scheduler ID: %+v", err)
			}

			schedulerId := retentionpolicies.NewSchedulerID(parsedId.SubscriptionId, parsedId.ResourceGroupName, parsedId.SchedulerName)
			// Custom ID type needed because the retention policy is a singleton child resource with a
			// fixed path `/retentionPolicies/default` under the parent scheduler. The SDK only provides
			// scheduler ID helpers, but the Terraform resource's state/import ID must be the full child
			// resource path (ending in `/retentionPolicies/default`) to uniquely identify this singleton.
			id := NewRetentionPolicyID(parsedId.SubscriptionId, parsedId.ResourceGroupName, parsedId.SchedulerName)

			metadata.Logger.Infof("Import check for retention policy on %s", schedulerId.ID())
			existing, err := client.Get(ctx, schedulerId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing retention policy on %s: %+v", schedulerId.ID(), err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			metadata.Logger.Infof("Creating retention policy on %s", schedulerId.ID())

			properties := retentionpolicies.RetentionPolicy{
				Properties: &retentionpolicies.RetentionPolicyProperties{
					RetentionPolicies: expandRetentionPolicyDetails(model.RetentionPolicy),
				},
			}

			if err := client.CreateOrReplaceThenPoll(ctx, schedulerId, properties); err != nil {
				return fmt.Errorf("creating retention policy on %s: %+v", schedulerId.ID(), err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func (r RetentionPolicyResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.RetentionPoliciesClient

			id, err := ParseRetentionPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			schedulerId := retentionpolicies.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

			metadata.Logger.Infof("Reading retention policy on %s", schedulerId.ID())
			resp, err := client.Get(ctx, schedulerId)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving retention policy on %s: %+v", schedulerId.ID(), err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving retention policy on %s: model was nil", schedulerId.ID())
			}

			state := RetentionPolicyResourceModel{
				DurableTaskSchedulerId: schedulerId.ID(),
				RetentionPolicy:        make([]RetentionPolicyModel, 0),
			}

			if props := model.Properties; props != nil && props.RetentionPolicies != nil {
				state.RetentionPolicy = flattenRetentionPolicyDetails(props.RetentionPolicies)
			}

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&state)
		},
	}
}

func (r RetentionPolicyResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.RetentionPoliciesClient

			id, err := ParseRetentionPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model RetentionPolicyResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			schedulerId := retentionpolicies.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

			metadata.Logger.Infof("Updating retention policy on %s", schedulerId.ID())

			if !metadata.ResourceData.HasChanges("retention_policy") {
				metadata.Logger.Infof("No changes detected for retention_policy on %s", schedulerId.ID())
				return nil
			}

			properties := retentionpolicies.RetentionPolicyUpdate{
				Properties: &retentionpolicies.RetentionPolicyProperties{
					RetentionPolicies: expandRetentionPolicyDetails(model.RetentionPolicy),
				},
			}

			if err := client.UpdateThenPoll(ctx, schedulerId, properties); err != nil {
				return fmt.Errorf("updating retention policy on %s: %+v", schedulerId.ID(), err)
			}

			return nil
		},
	}
}

func expandRetentionPolicyDetails(input []RetentionPolicyModel) *[]retentionpolicies.RetentionPolicyDetails {
	policies := make([]retentionpolicies.RetentionPolicyDetails, 0)
	if len(input) == 0 {
		return &policies
	}

	policy := input[0]
	appendPolicy := func(retentionPeriodInDays int64, orchestrationState *retentionpolicies.PurgeableOrchestrationState) {
		if retentionPeriodInDays == 0 {
			return
		}

		detail := retentionpolicies.RetentionPolicyDetails{
			RetentionPeriodInDays: retentionPeriodInDays,
		}
		if orchestrationState != nil {
			detail.OrchestrationState = orchestrationState
		}

		policies = append(policies, detail)
	}

	stateCanceled := retentionpolicies.PurgeableOrchestrationStateCanceled
	stateCompleted := retentionpolicies.PurgeableOrchestrationStateCompleted
	stateFailed := retentionpolicies.PurgeableOrchestrationStateFailed
	stateTerminated := retentionpolicies.PurgeableOrchestrationStateTerminated

	appendPolicy(policy.CanceledRetentionPeriodInDays, &stateCanceled)
	appendPolicy(policy.CompletedRetentionPeriodInDays, &stateCompleted)
	appendPolicy(policy.DefaultRetentionPeriodInDays, nil)
	appendPolicy(policy.FailedRetentionPeriodInDays, &stateFailed)
	appendPolicy(policy.TerminatedRetentionPeriodInDays, &stateTerminated)

	return &policies
}

func flattenRetentionPolicyDetails(input *[]retentionpolicies.RetentionPolicyDetails) []RetentionPolicyModel {
	if input == nil {
		return make([]RetentionPolicyModel, 0)
	}

	policy := RetentionPolicyModel{}
	for _, item := range *input {
		if item.OrchestrationState == nil {
			policy.DefaultRetentionPeriodInDays = item.RetentionPeriodInDays
			continue
		}

		switch *item.OrchestrationState {
		case retentionpolicies.PurgeableOrchestrationStateCanceled:
			policy.CanceledRetentionPeriodInDays = item.RetentionPeriodInDays
		case retentionpolicies.PurgeableOrchestrationStateCompleted:
			policy.CompletedRetentionPeriodInDays = item.RetentionPeriodInDays
		case retentionpolicies.PurgeableOrchestrationStateFailed:
			policy.FailedRetentionPeriodInDays = item.RetentionPeriodInDays
		case retentionpolicies.PurgeableOrchestrationStateTerminated:
			policy.TerminatedRetentionPeriodInDays = item.RetentionPeriodInDays
		}
	}

	if policy == (RetentionPolicyModel{}) {
		return make([]RetentionPolicyModel, 0)
	}

	return []RetentionPolicyModel{policy}
}

func (r RetentionPolicyResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DurableTask.RetentionPoliciesClient

			id, err := ParseRetentionPolicyID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			schedulerId := retentionpolicies.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

			metadata.Logger.Infof("Deleting retention policy on %s", schedulerId.ID())

			if err := client.DeleteThenPoll(ctx, schedulerId); err != nil {
				return fmt.Errorf("deleting retention policy on %s: %+v", schedulerId.ID(), err)
			}

			return nil
		},
	}
}

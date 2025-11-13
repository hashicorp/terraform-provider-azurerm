// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package durabletask

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/retentionpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/durabletask/2025-11-01/schedulers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type RetentionPolicyResourceModel struct {
	SchedulerId     string                     `tfschema:"scheduler_id"`
	RetentionPolicy []RetentionPolicyItemModel `tfschema:"retention_policy"`
}

type RetentionPolicyItemModel struct {
	RetentionPeriodInDays int64  `tfschema:"retention_period_in_days"`
	OrchestrationState    string `tfschema:"orchestration_state"`
}

type RetentionPolicyResource struct{}

var (
	_ sdk.Resource           = RetentionPolicyResource{}
	_ sdk.ResourceWithUpdate = RetentionPolicyResource{}
)

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
	return map[string]*pluginsdk.Schema{
		"scheduler_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: schedulers.ValidateSchedulerID,
		},

		"retention_policy": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"retention_period_in_days": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},

					"orchestration_state": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							"Completed",
							"Failed",
							"Terminated",
							"Canceled",
						}, false),
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

			schedulerId, err := schedulers.ParseSchedulerID(model.SchedulerId)
			if err != nil {
				return fmt.Errorf("parsing scheduler ID: %+v", err)
			}

			id := NewRetentionPolicyID(schedulerId.SubscriptionId, schedulerId.ResourceGroupName, schedulerId.SchedulerName)

			metadata.Logger.Infof("Import check for retention policy on %s", schedulerId.ID())
			existing, err := client.Get(ctx, *schedulerId)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing retention policy on %s: %+v", schedulerId.ID(), err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			metadata.Logger.Infof("Creating retention policy on %s", schedulerId.ID())

			policies := make([]retentionpolicies.RetentionPolicy, 0)
			for _, item := range model.RetentionPolicy {
				policy := retentionpolicies.RetentionPolicy{
					RetentionPeriodInDays: pointer.To(item.RetentionPeriodInDays),
				}

				if item.OrchestrationState != "" {
					policy.OrchestrationState = pointer.To(item.OrchestrationState)
				}

				policies = append(policies, policy)
			}

			properties := retentionpolicies.RetentionPolicyResource{
				Properties: &retentionpolicies.RetentionPolicyResourceProperties{
					Policies: &policies,
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *schedulerId, properties); err != nil {
				return fmt.Errorf("creating retention policy on %s: %+v", schedulerId.ID(), err)
			}

			metadata.SetID(id)
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

			schedulerId := schedulers.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

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
				SchedulerId:     schedulerId.ID(),
				RetentionPolicy: make([]RetentionPolicyItemModel, 0),
			}

			if props := model.Properties; props != nil && props.Policies != nil {
				for _, policy := range *props.Policies {
					item := RetentionPolicyItemModel{
						RetentionPeriodInDays: pointer.From(policy.RetentionPeriodInDays),
						OrchestrationState:    pointer.From(policy.OrchestrationState),
					}
					state.RetentionPolicy = append(state.RetentionPolicy, item)
				}
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

			schedulerId := schedulers.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

			metadata.Logger.Infof("Updating retention policy on %s", schedulerId.ID())

			policies := make([]retentionpolicies.RetentionPolicy, 0)
			for _, item := range model.RetentionPolicy {
				policy := retentionpolicies.RetentionPolicy{
					RetentionPeriodInDays: pointer.To(item.RetentionPeriodInDays),
				}

				if item.OrchestrationState != "" {
					policy.OrchestrationState = pointer.To(item.OrchestrationState)
				}

				policies = append(policies, policy)
			}

			properties := retentionpolicies.RetentionPolicyUpdate{
				Properties: &retentionpolicies.RetentionPolicyResourceProperties{
					Policies: &policies,
				},
			}

			if err := client.UpdateThenPoll(ctx, schedulerId, properties); err != nil {
				return fmt.Errorf("updating retention policy on %s: %+v", schedulerId.ID(), err)
			}

			return nil
		},
	}
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

			schedulerId := schedulers.NewSchedulerID(id.SubscriptionId, id.ResourceGroupName, id.SchedulerName)

			metadata.Logger.Infof("Deleting retention policy on %s", schedulerId.ID())

			if err := client.DeleteThenPoll(ctx, schedulerId); err != nil {
				return fmt.Errorf("deleting retention policy on %s: %+v", schedulerId.ID(), err)
			}

			return nil
		},
	}
}

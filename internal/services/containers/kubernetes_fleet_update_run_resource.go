// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2024-04-01/updateruns"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = KubernetesFleetUpdateRunResource{}
var _ sdk.ResourceWithUpdate = KubernetesFleetUpdateRunResource{}

type KubernetesFleetUpdateRunResource struct{}

func (r KubernetesFleetUpdateRunResource) ModelObject() interface{} {
	return &KubernetesFleetUpdateRunResourceSchema{}
}

type KubernetesFleetUpdateRunResourceSchema struct {
	KubernetesFleetManagerId string                                                       `tfschema:"kubernetes_fleet_manager_id"`
	Name                     string                                                       `tfschema:"name"`
	ManagedClusterUpdate     []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema `tfschema:"managed_cluster_update"`
	FleetUpdateStrategyId    string                                                       `tfschema:"fleet_update_strategy_id"`
	Stage                    []KubernetesFleetUpdateRunResourceUpdateStageSchema          `tfschema:"stage"`
}

type KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema struct {
	NodeImageSelection []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema `tfschema:"node_image_selection"`
	Upgrade            []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema            `tfschema:"upgrade"`
}

type KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema struct {
	Type string `tfschema:"type"`
}

type KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema struct {
	Type              string `tfschema:"type"`
	KubernetesVersion string `tfschema:"kubernetes_version"`
}

type KubernetesFleetUpdateRunResourceUpdateGroupSchema struct {
	Name string `tfschema:"name"`
}

type KubernetesFleetUpdateRunResourceUpdateStageSchema struct {
	AfterStageWaitInSeconds int64                                               `tfschema:"after_stage_wait_in_seconds"`
	Group                   []KubernetesFleetUpdateRunResourceUpdateGroupSchema `tfschema:"group"`
	Name                    string                                              `tfschema:"name"`
}

func (r KubernetesFleetUpdateRunResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return updateruns.ValidateUpdateRunID
}

func (r KubernetesFleetUpdateRunResource) ResourceType() string {
	return "azurerm_kubernetes_fleet_update_run"
}

func (r KubernetesFleetUpdateRunResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"kubernetes_fleet_manager_id": commonschema.ResourceIDReferenceRequiredForceNew(&commonids.KubernetesFleetId{}),

		"managed_cluster_update": {
			Required: true,
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"upgrade": {
						Required: true,
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Required: true,
									Type:     pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										string(updateruns.ManagedClusterUpgradeTypeFull),
										string(updateruns.ManagedClusterUpgradeTypeNodeImageOnly),
									}, false),
								},

								"kubernetes_version": {
									Optional:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"node_image_selection": {
						Optional: true,
						Type:     pluginsdk.TypeList,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"type": {
									Required: true,
									Type:     pluginsdk.TypeString,
									ValidateFunc: validation.StringInSlice([]string{
										string(updateruns.NodeImageSelectionTypeConsistent),
										string(updateruns.NodeImageSelectionTypeLatest),
									}, false),
								},
							},
						},
					},
				},
			},
		},

		"fleet_update_strategy_id": {
			Optional:      true,
			Type:          pluginsdk.TypeString,
			ConflictsWith: []string{"stage"},
		},

		"stage": {
			Optional:      true,
			Type:          pluginsdk.TypeList,
			ConflictsWith: []string{"fleet_update_strategy_id"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Required:     true,
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"group": {
						Required: true,
						Type:     pluginsdk.TypeList,
						MinItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Required:     true,
									Type:         pluginsdk.TypeString,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"after_stage_wait_in_seconds": {
						Optional: true,
						Type:     pluginsdk.TypeInt,
					},
				},
			},
		},
	}
}

func (r KubernetesFleetUpdateRunResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r KubernetesFleetUpdateRunResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateRunsClient

			var config KubernetesFleetUpdateRunResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			fleetId, err := commonids.ParseKubernetesFleetID(config.KubernetesFleetManagerId)
			if err != nil {
				return err
			}

			id := updateruns.NewUpdateRunID(fleetId.SubscriptionId, fleetId.ResourceGroupName, fleetId.FleetName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := updateruns.UpdateRun{
				Properties: &updateruns.UpdateRunProperties{
					ManagedClusterUpdate: expandKubernetesFleetUpdateRunManagedClusterUpdate(config.ManagedClusterUpdate),
					Strategy:             expandKubernetesFleetUpdateRunStrategy(config.Stage),
					UpdateStrategyId:     pointer.To(config.FleetUpdateStrategyId),
				},
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload, updateruns.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r KubernetesFleetUpdateRunResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateRunsClient

			id, err := updateruns.ParseUpdateRunID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config KubernetesFleetUpdateRunResourceSchema
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving existing %s: properties was nil", *id)
			}
			payload := *existing.Model

			if metadata.ResourceData.HasChange("managed_cluster_update") {
				payload.Properties.ManagedClusterUpdate = expandKubernetesFleetUpdateRunManagedClusterUpdate(config.ManagedClusterUpdate)
			}

			if metadata.ResourceData.HasChange("fleet_update_strategy_id") {
				payload.Properties.UpdateStrategyId = pointer.To(config.FleetUpdateStrategyId)
			}

			if metadata.ResourceData.HasChange("stage") {
				payload.Properties.Strategy.Stages = expandKubernetesFleetUpdateRunStage(config.Stage)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload, updateruns.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r KubernetesFleetUpdateRunResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateRunsClient
			schema := KubernetesFleetUpdateRunResourceSchema{}

			id, err := updateruns.ParseUpdateRunID(metadata.ResourceData.Id())
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

			if model := resp.Model; model != nil {
				schema.Name = id.UpdateRunName
				schema.KubernetesFleetManagerId = commonids.NewKubernetesFleetID(id.SubscriptionId, id.ResourceGroupName, id.FleetName).ID()
				if model.Properties != nil {
					schema.ManagedClusterUpdate = flattenKubernetesFleetUpdateRunManagedClusterUpdate(model.Properties.ManagedClusterUpdate)
					schema.FleetUpdateStrategyId = pointer.From(model.Properties.UpdateStrategyId)
					schema.Stage = flattenKubernetesFleetUpdateRunStage(model.Properties.Strategy.Stages)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r KubernetesFleetUpdateRunResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Containers.FleetUpdateRunsClient

			id, err := updateruns.ParseUpdateRunID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, updateruns.DefaultDeleteOperationOptions()); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandKubernetesFleetUpdateRunManagedClusterUpdate(input []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema) updateruns.ManagedClusterUpdate {
	if len(input) == 0 {
		return updateruns.ManagedClusterUpdate{}
	}
	return updateruns.ManagedClusterUpdate{
		NodeImageSelection: expandKubernetesFleetUpdateRunNodeImageSelection(input[0].NodeImageSelection),
		Upgrade:            expandKubernetesFleetUpdateRunUpgrade(input[0].Upgrade),
	}
}

func expandKubernetesFleetUpdateRunUpgrade(input []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema) updateruns.ManagedClusterUpgradeSpec {
	if len(input) == 0 {
		return updateruns.ManagedClusterUpgradeSpec{}
	}
	return updateruns.ManagedClusterUpgradeSpec{
		Type:              updateruns.ManagedClusterUpgradeType(input[0].Type),
		KubernetesVersion: pointer.To(input[0].KubernetesVersion),
	}
}

func expandKubernetesFleetUpdateRunNodeImageSelection(input []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema) *updateruns.NodeImageSelection {
	if len(input) == 0 {
		return nil
	}
	return &updateruns.NodeImageSelection{
		Type: updateruns.NodeImageSelectionType(input[0].Type),
	}
}

func expandKubernetesFleetUpdateRunStrategy(input []KubernetesFleetUpdateRunResourceUpdateStageSchema) *updateruns.UpdateRunStrategy {
	if len(input) == 0 {
		return nil
	}

	return &updateruns.UpdateRunStrategy{
		Stages: expandKubernetesFleetUpdateRunStage(input),
	}
}

func expandKubernetesFleetUpdateRunStage(input []KubernetesFleetUpdateRunResourceUpdateStageSchema) []updateruns.UpdateStage {
	output := make([]updateruns.UpdateStage, 0)
	for _, stage := range input {
		output = append(output, updateruns.UpdateStage{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.FromInt64(stage.AfterStageWaitInSeconds),
			Groups:                  expandKubernetesFleetUpdateRunGroup(stage.Group),
		})
	}
	return output
}

func expandKubernetesFleetUpdateRunGroup(input []KubernetesFleetUpdateRunResourceUpdateGroupSchema) *[]updateruns.UpdateGroup {
	output := make([]updateruns.UpdateGroup, 0)
	for _, group := range input {
		output = append(output, updateruns.UpdateGroup{
			Name: group.Name,
		})
	}
	return &output
}

func flattenKubernetesFleetUpdateRunManagedClusterUpdate(input updateruns.ManagedClusterUpdate) []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema {
	return []KubernetesFleetUpdateRunResourceManagedClusterUpdateSchema{
		{
			NodeImageSelection: flattenKubernetesFleetUpdateRunNodeImageSelection(input.NodeImageSelection),
			Upgrade:            flattenKubernetesFleetUpdateRunUpgrade(input.Upgrade),
		},
	}
}

func flattenKubernetesFleetUpdateRunNodeImageSelection(input *updateruns.NodeImageSelection) []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema {
	if input == nil {
		return []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema{}
	}
	return []KubernetesFleetUpdateRunResourceManagedClusterUpdateNodeImageSelectionSchema{
		{
			Type: string(input.Type),
		},
	}
}

func flattenKubernetesFleetUpdateRunUpgrade(input updateruns.ManagedClusterUpgradeSpec) []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema {
	return []KubernetesFleetUpdateRunResourceManagedClusterUpdateUpgradeSchema{
		{
			Type:              string(input.Type),
			KubernetesVersion: pointer.From(input.KubernetesVersion),
		},
	}
}

func flattenKubernetesFleetUpdateRunStage(input []updateruns.UpdateStage) []KubernetesFleetUpdateRunResourceUpdateStageSchema {
	output := make([]KubernetesFleetUpdateRunResourceUpdateStageSchema, 0)
	for _, stage := range input {
		output = append(output, KubernetesFleetUpdateRunResourceUpdateStageSchema{
			Name:                    stage.Name,
			AfterStageWaitInSeconds: pointer.ToInt64(stage.AfterStageWaitInSeconds),
			Group:                   flattenKubernetesFleetUpdateRunGroup(stage.Groups),
		})
	}
	return output

}

func flattenKubernetesFleetUpdateRunGroup(input *[]updateruns.UpdateGroup) []KubernetesFleetUpdateRunResourceUpdateGroupSchema {
	output := make([]KubernetesFleetUpdateRunResourceUpdateGroupSchema, 0)
	for _, group := range *input {
		output = append(output, KubernetesFleetUpdateRunResourceUpdateGroupSchema{
			Name: group.Name,
		})
	}
	return output
}

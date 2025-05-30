// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteRecoveryReplicationRecoveryPlanModel struct {
	Name                   string                                         `tfschema:"name"`
	ShutdownRecoveryGroup  []GenericRecoveryGroupModel                    `tfschema:"shutdown_recovery_group"`
	FailoverRecoveryGroup  []GenericRecoveryGroupModel                    `tfschema:"failover_recovery_group"`
	BootRecoveryGroup      []BootRecoveryGroupModel                       `tfschema:"boot_recovery_group"`
	RecoveryVaultId        string                                         `tfschema:"recovery_vault_id"`
	SourceRecoveryFabricId string                                         `tfschema:"source_recovery_fabric_id"`
	TargetRecoveryFabricId string                                         `tfschema:"target_recovery_fabric_id"`
	A2ASettings            []ReplicationRecoveryPlanA2ASpecificInputModel `tfschema:"azure_to_azure_settings"`
}

type GenericRecoveryGroupModel struct {
	PostAction []ActionModel `tfschema:"post_action"`
	PreAction  []ActionModel `tfschema:"pre_action"`
}

type BootRecoveryGroupModel struct {
	PostAction               []ActionModel `tfschema:"post_action"`
	PreAction                []ActionModel `tfschema:"pre_action"`
	ReplicatedProtectedItems []string      `tfschema:"replicated_protected_items"`
}

type RecoveryGroupModel struct {
	GroupType                string        `tfschema:"type"`
	PostAction               []ActionModel `tfschema:"post_action"`
	PreAction                []ActionModel `tfschema:"pre_action"`
	ReplicatedProtectedItems []string      `tfschema:"replicated_protected_items"`
}

type ActionModel struct {
	ActionDetailType        string   `tfschema:"type"`
	FabricLocation          string   `tfschema:"fabric_location"`
	FailOverDirections      []string `tfschema:"fail_over_directions"`
	FailOverTypes           []string `tfschema:"fail_over_types"`
	ManualActionInstruction string   `tfschema:"manual_action_instruction"`
	Name                    string   `tfschema:"name"`
	RunbookId               string   `tfschema:"runbook_id"`
	ScriptPath              string   `tfschema:"script_path"`
}

type ReplicationRecoveryPlanA2ASpecificInputModel struct {
	PrimaryZone      string `tfschema:"primary_zone"`
	RecoveryZone     string `tfschema:"recovery_zone"`
	PrimaryEdgeZone  string `tfschema:"primary_edge_zone"`
	RecoveryEdgeZone string `tfschema:"recovery_edge_zone"`
}

type SiteRecoveryReplicationRecoveryPlanResource struct{}

var _ sdk.ResourceWithUpdate = SiteRecoveryReplicationRecoveryPlanResource{}

func (r SiteRecoveryReplicationRecoveryPlanResource) ResourceType() string {
	return "azurerm_site_recovery_replication_recovery_plan"
}

func (r SiteRecoveryReplicationRecoveryPlanResource) ModelObject() interface{} {
	return &SiteRecoveryReplicationRecoveryPlanModel{}
}

func (r SiteRecoveryReplicationRecoveryPlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationrecoveryplans.ValidateReplicationRecoveryPlanID
}

func (r SiteRecoveryReplicationRecoveryPlanResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9-]{1,63}[a-zA-Z0-9]$`), "The name can contain only letters, numbers, and hyphens. It should start with a letter and end with a letter or a number."),
		},

		"recovery_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateVaultID,
		},

		"source_recovery_fabric_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		"target_recovery_fabric_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		// lintignore:S013
		"shutdown_recovery_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pre_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},

					"post_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},
				},
			},
		},

		// lintignore:S013
		"failover_recovery_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"pre_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},

					"post_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},
				},
			},
		},

		// lintignore:S013
		"boot_recovery_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"replicated_protected_items": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},

					"pre_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},

					"post_action": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem:     replicationRecoveryPlanActionSchema(),
					},
				},
			},
		},

		"azure_to_azure_settings": replicationRecoveryPlanA2ASchema(),
	}
}

// we do not split action into three different schema because all actions should keep the order from user.
func replicationRecoveryPlanActionSchema() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AutomationRunbookActionDetails",
					"ManualActionDetails",
					"ScriptActionDetails",
				}, false),
			},

			"fail_over_directions": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(replicationrecoveryplans.PossibleOperationsDirectionsPrimaryToRecovery),
						string(replicationrecoveryplans.PossibleOperationsDirectionsRecoveryToPrimary),
					}, false),
				},
			},

			"fail_over_types": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(replicationrecoveryplans.ReplicationProtectedItemOperationPlannedFailover),
						string(replicationrecoveryplans.ReplicationProtectedItemOperationTestFailover),
						string(replicationrecoveryplans.ReplicationProtectedItemOperationUnplannedFailover),
					}, false),
				},
			},

			"runbook_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"fabric_location": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(replicationrecoveryplans.RecoveryPlanActionLocationPrimary),
					string(replicationrecoveryplans.RecoveryPlanActionLocationRecovery),
				}, false),
			},

			"manual_action_instruction": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"script_path": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func replicationRecoveryPlanA2ASchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"primary_zone": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{"azure_to_azure_settings.0.recovery_zone"},
				},

				"recovery_zone": {
					Type:         schema.TypeString,
					Optional:     true,
					ForceNew:     true,
					ValidateFunc: validation.StringIsNotEmpty,
					RequiredWith: []string{"azure_to_azure_settings.0.primary_zone"},
				},

				"primary_edge_zone": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					ValidateFunc:     validation.StringIsNotEmpty,
					StateFunc:        edgezones.StateFunc,
					DiffSuppressFunc: edgezones.DiffSuppressFunc,
					RequiredWith:     []string{"azure_to_azure_settings.0.recovery_edge_zone"},
				},

				"recovery_edge_zone": {
					Type:             schema.TypeString,
					Optional:         true,
					ForceNew:         true,
					ValidateFunc:     validation.StringIsNotEmpty,
					StateFunc:        edgezones.StateFunc,
					DiffSuppressFunc: edgezones.DiffSuppressFunc,
					RequiredWith:     []string{"azure_to_azure_settings.0.primary_edge_zone"},
				},
			},
		},
	}
}

func (r SiteRecoveryReplicationRecoveryPlanResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (r SiteRecoveryReplicationRecoveryPlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SiteRecoveryReplicationRecoveryPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.RecoveryServices.ReplicationRecoveryPlansClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vaultId, err := replicationrecoveryplans.ParseVaultID(model.RecoveryVaultId)
			if err != nil {
				return err
			}

			id := replicationrecoveryplans.NewReplicationRecoveryPlanID(subscriptionId, vaultId.ResourceGroupName, vaultId.VaultName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
				if !response.WasNotFound(existing.HttpResponse) && !response.WasBadRequest(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing site recovery plan %q: %+v", id, err)
				}
			}

			if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
				return tf.ImportAsExistsError("azurerm_site_recovery_replication_recovery_plan", *existing.Model.Id)
			}

			// FailoverDeploymentModelClassic is used for other cloud service back up to Azure.
			deploymentModel := replicationrecoveryplans.FailoverDeploymentModelResourceManager

			groupValue, err := expandRecoveryGroup(model.ShutdownRecoveryGroup, model.FailoverRecoveryGroup, model.BootRecoveryGroup)
			if err != nil {
				return fmt.Errorf("expanding recovery group: %+v", err)
			}

			parameters := replicationrecoveryplans.CreateRecoveryPlanInput{
				Properties: replicationrecoveryplans.CreateRecoveryPlanInputProperties{
					PrimaryFabricId:         model.SourceRecoveryFabricId,
					RecoveryFabricId:        model.TargetRecoveryFabricId,
					FailoverDeploymentModel: &deploymentModel,
					Groups:                  groupValue,
				},
			}

			if len(model.A2ASettings) == 1 {
				parameters.Properties.ProviderSpecificInput = expandA2ASettings(model.A2ASettings[0])
			}

			err = client.CreateThenPoll(ctx, id, parameters)
			if err != nil {
				return fmt.Errorf("creating site recovery replication plan %q: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r SiteRecoveryReplicationRecoveryPlanResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ReplicationRecoveryPlansClient

			id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			vaultId := replicationrecoveryplans.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("making Read request on site recovery replication plan %s : %+v", id.String(), err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("making Read request on site recovery replication plan %s : model is nil", id.String())
			}

			state := SiteRecoveryReplicationRecoveryPlanModel{
				Name:            id.ReplicationRecoveryPlanName,
				RecoveryVaultId: vaultId.ID(),
			}

			if prop := model.Properties; prop != nil {
				if prop.PrimaryFabricId != nil {
					state.SourceRecoveryFabricId = handleAzureSdkForGoBug2824(*prop.PrimaryFabricId)
				}
				if prop.RecoveryFabricId != nil {
					state.TargetRecoveryFabricId = handleAzureSdkForGoBug2824(*prop.RecoveryFabricId)
				}
				if group := prop.Groups; group != nil {
					state.ShutdownRecoveryGroup, state.FailoverRecoveryGroup, state.BootRecoveryGroup = flattenRecoveryGroups(*group)
				}
				if details := prop.ProviderSpecificDetails; details != nil && len(*details) > 0 {
					state.A2ASettings = flattenRecoveryPlanProviderSpecificInput(details)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SiteRecoveryReplicationRecoveryPlanResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SiteRecoveryReplicationRecoveryPlanModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}
			client := metadata.Client.RecoveryServices.ReplicationRecoveryPlansClient

			id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parse Site reocvery replication plan id: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var groupValue []replicationrecoveryplans.RecoveryPlanGroup
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model is nil", *id)
			}

			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: properties is nil", *id)
			}

			if resp.Model.Properties.Groups == nil {
				return fmt.Errorf("retrieving %s: groups is nil", *id)
			}

			groupValue = *resp.Model.Properties.Groups

			if metadata.ResourceData.HasChange("boot_recovery_group") ||
				metadata.ResourceData.HasChange("failover_recovery_group") ||
				metadata.ResourceData.HasChange("shutdown_recovery_group") {
				groupValue, err = expandRecoveryGroup(model.ShutdownRecoveryGroup, model.FailoverRecoveryGroup, model.BootRecoveryGroup)
			}

			if err != nil {
				return fmt.Errorf("expanding recovery group: %+v", err)
			}

			parameters := replicationrecoveryplans.UpdateRecoveryPlanInput{
				Properties: &replicationrecoveryplans.UpdateRecoveryPlanInputProperties{
					Groups: &groupValue,
				},
			}

			err = client.UpdateThenPoll(ctx, *id, parameters)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SiteRecoveryReplicationRecoveryPlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ReplicationRecoveryPlansClient

			id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting site recovery protection replication plan %q : %+v", id, err)
			}

			return nil
		},
	}
}

func expandRecoveryGroup(shutdown []GenericRecoveryGroupModel, failover []GenericRecoveryGroupModel, boot []BootRecoveryGroupModel) ([]replicationrecoveryplans.RecoveryPlanGroup, error) {
	output := make([]replicationrecoveryplans.RecoveryPlanGroup, 0)

	for _, group := range shutdown {
		preActions, err := expandAction(group.PreAction)
		if err != nil {
			return output, err
		}
		postActions, err := expandAction(group.PostAction)
		if err != nil {
			return output, err
		}

		output = append(output, replicationrecoveryplans.RecoveryPlanGroup{
			GroupType:         replicationrecoveryplans.RecoveryPlanGroupTypeShutdown,
			StartGroupActions: &preActions,
			EndGroupActions:   &postActions,
		})
	}

	for _, group := range failover {
		preActions, err := expandAction(group.PreAction)
		if err != nil {
			return output, err
		}
		postActions, err := expandAction(group.PostAction)
		if err != nil {
			return output, err
		}

		output = append(output, replicationrecoveryplans.RecoveryPlanGroup{
			GroupType:         replicationrecoveryplans.RecoveryPlanGroupTypeFailover,
			StartGroupActions: &preActions,
			EndGroupActions:   &postActions,
		})
	}

	for _, group := range boot {
		protectedItems := make([]replicationrecoveryplans.RecoveryPlanProtectedItem, 0)
		for _, protectedItem := range group.ReplicatedProtectedItems {
			protectedItems = append(protectedItems, replicationrecoveryplans.RecoveryPlanProtectedItem{
				Id: pointer.To(protectedItem),
			})
		}

		preActions, err := expandAction(group.PreAction)
		if err != nil {
			return output, err
		}
		postActions, err := expandAction(group.PostAction)
		if err != nil {
			return output, err
		}

		output = append(output, replicationrecoveryplans.RecoveryPlanGroup{
			GroupType:                 replicationrecoveryplans.RecoveryPlanGroupTypeBoot,
			ReplicationProtectedItems: &protectedItems,
			StartGroupActions:         &preActions,
			EndGroupActions:           &postActions,
		})
	}

	return output, nil
}

func expandAction(input []ActionModel) ([]replicationrecoveryplans.RecoveryPlanAction, error) {
	output := make([]replicationrecoveryplans.RecoveryPlanAction, 0)
	for _, action := range input {
		failoverDirections := make([]replicationrecoveryplans.PossibleOperationsDirections, 0)
		for _, direction := range action.FailOverDirections {
			failoverDirections = append(failoverDirections, replicationrecoveryplans.PossibleOperationsDirections(direction))
		}

		failoverTypes := make([]replicationrecoveryplans.ReplicationProtectedItemOperation, 0)
		for _, failoverType := range action.FailOverTypes {
			failoverTypes = append(failoverTypes, replicationrecoveryplans.ReplicationProtectedItemOperation(failoverType))
		}

		if action.ActionDetailType == "ManualActionDetails" && action.FabricLocation != "" {
			return nil, fmt.Errorf("`fabric_location` must not be specified for `recovery_group` with `ManualActionDetails` type")
		}

		output = append(output, replicationrecoveryplans.RecoveryPlanAction{
			ActionName:         action.Name,
			FailoverDirections: failoverDirections,
			FailoverTypes:      failoverTypes,
			CustomDetails:      expandActionDetail(action),
		})
	}

	return output, nil
}

func expandA2ASettings(input ReplicationRecoveryPlanA2ASpecificInputModel) *[]replicationrecoveryplans.RecoveryPlanProviderSpecificInput {
	return &[]replicationrecoveryplans.RecoveryPlanProviderSpecificInput{
		replicationrecoveryplans.RecoveryPlanA2AInput{
			PrimaryZone:              pointer.To(input.PrimaryZone),
			RecoveryZone:             pointer.To(input.RecoveryZone),
			PrimaryExtendedLocation:  expandEdgeZone(input.PrimaryEdgeZone),
			RecoveryExtendedLocation: expandEdgeZone(input.RecoveryEdgeZone),
		},
	}
}

func flattenRecoveryGroups(input []replicationrecoveryplans.RecoveryPlanGroup) (shutdown []GenericRecoveryGroupModel, failover []GenericRecoveryGroupModel, boot []BootRecoveryGroupModel) {
	shutdown = make([]GenericRecoveryGroupModel, 0)
	failover = make([]GenericRecoveryGroupModel, 0)
	boot = make([]BootRecoveryGroupModel, 0)
	for _, groupItem := range input {
		switch groupItem.GroupType {
		case replicationrecoveryplans.RecoveryPlanGroupTypeShutdown:
			o := GenericRecoveryGroupModel{}
			if groupItem.StartGroupActions != nil {
				o.PreAction = flattenRecoveryPlanActions(groupItem.StartGroupActions)
			}
			if groupItem.EndGroupActions != nil {
				o.PostAction = flattenRecoveryPlanActions(groupItem.EndGroupActions)
			}
			shutdown = append(shutdown, o)
		case replicationrecoveryplans.RecoveryPlanGroupTypeFailover:
			o := GenericRecoveryGroupModel{}
			if groupItem.StartGroupActions != nil {
				o.PreAction = flattenRecoveryPlanActions(groupItem.StartGroupActions)
			}
			if groupItem.EndGroupActions != nil {
				o.PostAction = flattenRecoveryPlanActions(groupItem.EndGroupActions)
			}
			failover = append(failover, o)
		case replicationrecoveryplans.RecoveryPlanGroupTypeBoot:
			o := BootRecoveryGroupModel{}
			o.ReplicatedProtectedItems = flattenRecoveryPlanProtectedItems(groupItem.ReplicationProtectedItems)
			if groupItem.StartGroupActions != nil {
				o.PreAction = flattenRecoveryPlanActions(groupItem.StartGroupActions)
			}
			if groupItem.EndGroupActions != nil {
				o.PostAction = flattenRecoveryPlanActions(groupItem.EndGroupActions)
			}
			boot = append(boot, o)
		}
	}

	return shutdown, failover, boot
}

func expandActionDetail(input ActionModel) (output replicationrecoveryplans.RecoveryPlanActionDetails) {
	switch input.ActionDetailType {
	case "AutomationRunbookActionDetails":
		output = replicationrecoveryplans.RecoveryPlanAutomationRunbookActionDetails{
			RunbookId:      pointer.To(input.RunbookId),
			FabricLocation: replicationrecoveryplans.RecoveryPlanActionLocation(input.FabricLocation),
		}
	case "ManualActionDetails":
		output = replicationrecoveryplans.RecoveryPlanManualActionDetails{
			Description: pointer.To(input.ManualActionInstruction),
		}
	case "ScriptActionDetails":
		output = replicationrecoveryplans.RecoveryPlanScriptActionDetails{
			Path:           input.ScriptPath,
			FabricLocation: replicationrecoveryplans.RecoveryPlanActionLocation(input.FabricLocation),
		}
	}
	return
}

func flattenRecoveryPlanProtectedItems(input *[]replicationrecoveryplans.RecoveryPlanProtectedItem) []string {
	protectedItemOutputs := make([]string, 0)
	if input != nil {
		for _, protectedItem := range *input {
			protectedItemOutputs = append(protectedItemOutputs, handleAzureSdkForGoBug2824(*protectedItem.Id))
		}
	}
	return protectedItemOutputs
}

func flattenRecoveryPlanActions(input *[]replicationrecoveryplans.RecoveryPlanAction) []ActionModel {
	actionOutputs := make([]ActionModel, 0)
	for _, action := range *input {
		actionOutput := ActionModel{
			Name: action.ActionName,
		}
		switch detail := action.CustomDetails.(type) {
		case replicationrecoveryplans.RecoveryPlanAutomationRunbookActionDetails:
			actionOutput.ActionDetailType = "AutomationRunbookActionDetails"
			actionOutput.FabricLocation = string(detail.FabricLocation)
			if detail.RunbookId != nil {
				actionOutput.RunbookId = *detail.RunbookId
			}
		case replicationrecoveryplans.RecoveryPlanManualActionDetails:
			actionOutput.ActionDetailType = "ManualActionDetails"
			if detail.Description != nil {
				actionOutput.ManualActionInstruction = *detail.Description
			}
		case replicationrecoveryplans.RecoveryPlanScriptActionDetails:
			actionOutput.ActionDetailType = "ScriptActionDetails"
			actionOutput.ScriptPath = detail.Path
			actionOutput.FabricLocation = string(detail.FabricLocation)
		}

		failoverDirections := make([]string, 0)
		for _, failoverDirection := range action.FailoverDirections {
			failoverDirections = append(failoverDirections, string(failoverDirection))
		}

		failoverTypes := make([]string, 0)
		for _, failoverType := range action.FailoverTypes {
			failoverTypes = append(failoverTypes, string(failoverType))
		}
		actionOutput.FailOverDirections = failoverDirections
		actionOutput.FailOverTypes = failoverTypes
		actionOutputs = append(actionOutputs, actionOutput)
	}
	return actionOutputs
}

func flattenRecoveryPlanProviderSpecificInput(input *[]replicationrecoveryplans.RecoveryPlanProviderSpecificDetails) []ReplicationRecoveryPlanA2ASpecificInputModel {
	output := make([]ReplicationRecoveryPlanA2ASpecificInputModel, 0)
	for _, providerSpecificInput := range *input {
		if a2aInput, ok := providerSpecificInput.(replicationrecoveryplans.RecoveryPlanA2ADetails); ok {
			o := ReplicationRecoveryPlanA2ASpecificInputModel{
				PrimaryZone:      pointer.From(a2aInput.PrimaryZone),
				RecoveryZone:     pointer.From(a2aInput.RecoveryZone),
				PrimaryEdgeZone:  flattenEdgeZone(a2aInput.PrimaryExtendedLocation),
				RecoveryEdgeZone: flattenEdgeZone(a2aInput.RecoveryExtendedLocation),
			}
			output = append(output, o)
		}
	}
	return output
}

package recoveryservices

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SiteRecoveryReplicationRecoveryPlanModel struct {
	Name                   string               `tfschema:"name"`
	RecoveryGroup          []RecoveryGroupModel `tfschema:"recovery_group"`
	RecoveryVaultId        string               `tfschema:"recovery_vault_id"`
	SourceRecoveryFabricId string               `tfschema:"source_recovery_fabric_id"`
	TargetRecoveryFabricId string               `tfschema:"target_recovery_fabric_id"`
}

type RecoveryGroupModel struct {
	GroupType                replicationrecoveryplans.RecoveryPlanGroupType `tfschema:"type"`
	PostAction               []ActionModel                                  `tfschema:"post_action"`
	PreAction                []ActionModel                                  `tfschema:"pre_action"`
	ReplicatedProtectedItems []string                                       `tfschema:"replicated_protected_items"`
}

type ActionModel struct {
	ActionDetailType        string                                              `tfschema:"type"`
	FabricLocation          replicationrecoveryplans.RecoveryPlanActionLocation `tfschema:"fabric_location"`
	FailOverDirections      []string                                            `tfschema:"fail_over_directions"`
	FailOverTypes           []string                                            `tfschema:"fail_over_types"`
	ManualActionInstruction string                                              `tfschema:"manual_action_instruction"`
	Name                    string                                              `tfschema:"name"`
	RunbookId               string                                              `tfschema:"runbook_id"`
	ScriptPath              string                                              `tfschema:"script_path"`
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

		"recovery_group": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MinItems: 3,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(replicationrecoveryplans.RecoveryPlanGroupTypeBoot),
							string(replicationrecoveryplans.RecoveryPlanGroupTypeShutdown),
							string(replicationrecoveryplans.RecoveryPlanGroupTypeFailover),
						}, false),
					},
					"replicated_protected_items": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
					"pre_action": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     schemaAction(),
					},
					"post_action": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem:     schemaAction(),
					},
				},
			},
		},
	}
}

func schemaAction() *pluginsdk.Resource {
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
			groupValue, err := expandRecoverGroup(model.RecoveryGroup)
			if err != nil {
				return fmt.Errorf("when expanding recovery group: %s", err)
			}

			parameters := replicationrecoveryplans.CreateRecoveryPlanInput{
				Properties: replicationrecoveryplans.CreateRecoveryPlanInputProperties{
					PrimaryFabricId:         model.SourceRecoveryFabricId,
					RecoveryFabricId:        model.TargetRecoveryFabricId,
					FailoverDeploymentModel: &deploymentModel,
					Groups:                  groupValue,
				},
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
					state.RecoveryGroup = flattenRecoveryGroups(*group)
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

			recoveryPlanGroup, err := expandRecoverGroup(model.RecoveryGroup)
			if err != nil {
				return fmt.Errorf("when expanding recovery group: %s", err)
			}

			parameters := replicationrecoveryplans.UpdateRecoveryPlanInput{
				Properties: &replicationrecoveryplans.UpdateRecoveryPlanInputProperties{
					Groups: &recoveryPlanGroup,
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

func expandRecoverGroup(input []RecoveryGroupModel) ([]replicationrecoveryplans.RecoveryPlanGroup, error) {
	output := make([]replicationrecoveryplans.RecoveryPlanGroup, 0)
	if pass, err := validateRecoverGroup(input); !pass {
		return output, err
	}

	for _, group := range input {

		protectedItems := make([]replicationrecoveryplans.RecoveryPlanProtectedItem, 0)
		for _, protectedItem := range group.ReplicatedProtectedItems {
			protectedItems = append(protectedItems, replicationrecoveryplans.RecoveryPlanProtectedItem{
				Id: utils.String(protectedItem),
			})
		}

		preActions := make([]replicationrecoveryplans.RecoveryPlanAction, 0)
		for _, preActionInput := range group.PreAction {

			failoverDirections := make([]replicationrecoveryplans.PossibleOperationsDirections, 0)
			for _, direction := range preActionInput.FailOverDirections {
				failoverDirections = append(failoverDirections, replicationrecoveryplans.PossibleOperationsDirections(direction))
			}

			failoverTypes := make([]replicationrecoveryplans.ReplicationProtectedItemOperation, 0)
			for _, failoveType := range preActionInput.FailOverTypes {
				failoverTypes = append(failoverTypes, replicationrecoveryplans.ReplicationProtectedItemOperation(failoveType))
			}

			preActions = append(preActions, replicationrecoveryplans.RecoveryPlanAction{
				ActionName:         preActionInput.Name,
				FailoverDirections: failoverDirections,
				FailoverTypes:      failoverTypes,
				CustomDetails:      expandActionDetail(preActionInput),
			})
		}

		postActions := make([]replicationrecoveryplans.RecoveryPlanAction, 0)
		for _, postActionInput := range group.PostAction {

			failoverDirections := make([]replicationrecoveryplans.PossibleOperationsDirections, 0)
			for _, direction := range postActionInput.FailOverDirections {
				failoverDirections = append(failoverDirections, replicationrecoveryplans.PossibleOperationsDirections(direction))
			}

			failoverTypes := make([]replicationrecoveryplans.ReplicationProtectedItemOperation, 0)
			for _, failoveType := range postActionInput.FailOverTypes {
				failoverTypes = append(failoverTypes, replicationrecoveryplans.ReplicationProtectedItemOperation(failoveType))
			}

			postActions = append(postActions, replicationrecoveryplans.RecoveryPlanAction{
				ActionName:         postActionInput.Name,
				FailoverDirections: failoverDirections,
				FailoverTypes:      failoverTypes,
				CustomDetails:      expandActionDetail(postActionInput),
			})
		}

		output = append(output, replicationrecoveryplans.RecoveryPlanGroup{
			GroupType:                 group.GroupType,
			ReplicationProtectedItems: &protectedItems,
			StartGroupActions:         &preActions,
			EndGroupActions:           &postActions,
		})

	}
	return output, nil
}

func validateRecoverGroup(input []RecoveryGroupModel) (bool, error) {
	bootCount := 0
	shutdownCount := 0
	failoverCount := 0
	for _, group := range input {
		if group.GroupType == replicationrecoveryplans.RecoveryPlanGroupTypeBoot {
			bootCount += 1
		}
		if group.GroupType == replicationrecoveryplans.RecoveryPlanGroupTypeFailover {
			failoverCount += 1
		}
		if group.GroupType == replicationrecoveryplans.RecoveryPlanGroupTypeShutdown {
			shutdownCount += 1
		}
	}
	if bootCount == 0 || shutdownCount == 0 || failoverCount == 0 {
		return false, fmt.Errorf("every group type needs at least one recovery group")
	}
	return true, nil
}

func flattenRecoveryGroups(input []replicationrecoveryplans.RecoveryPlanGroup) []RecoveryGroupModel {
	output := make([]RecoveryGroupModel, 0)
	for _, groupItem := range input {
		recoveryGroupOutput := RecoveryGroupModel{}
		recoveryGroupOutput.GroupType = groupItem.GroupType
		if groupItem.ReplicationProtectedItems != nil {
			recoveryGroupOutput.ReplicatedProtectedItems = flattenRecoveryPlanProtectedItems(groupItem.ReplicationProtectedItems)
		}
		if groupItem.StartGroupActions != nil {
			recoveryGroupOutput.PreAction = flattenRecoveryPlanActions(groupItem.StartGroupActions)
		}
		if groupItem.EndGroupActions != nil {
			recoveryGroupOutput.PostAction = flattenRecoveryPlanActions(groupItem.EndGroupActions)
		}
		output = append(output, recoveryGroupOutput)
	}
	return output
}

func expandActionDetail(input ActionModel) (output replicationrecoveryplans.RecoveryPlanActionDetails) {
	switch input.ActionDetailType {
	case "AutomationRunbookActionDetails":
		output = replicationrecoveryplans.RecoveryPlanAutomationRunbookActionDetails{
			RunbookId:      utils.String(input.RunbookId),
			FabricLocation: input.FabricLocation,
		}
	case "ManualActionDetails":
		output = replicationrecoveryplans.RecoveryPlanManualActionDetails{
			Description: utils.String(input.ManualActionInstruction),
		}
	case "ScriptActionDetails":
		output = replicationrecoveryplans.RecoveryPlanScriptActionDetails{
			Path:           input.ScriptPath,
			FabricLocation: input.FabricLocation,
		}
	}
	return
}

func flattenRecoveryPlanProtectedItems(input *[]replicationrecoveryplans.RecoveryPlanProtectedItem) []string {
	protectedItemOutputs := make([]string, 0)
	for _, protectedItem := range *input {
		protectedItemOutputs = append(protectedItemOutputs, handleAzureSdkForGoBug2824(*protectedItem.Id))
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
			actionOutput.FabricLocation = detail.FabricLocation
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
			actionOutput.FabricLocation = detail.FabricLocation
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

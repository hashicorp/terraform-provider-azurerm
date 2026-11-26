// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_recovery_plan

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
)

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

func flattenRecoveryPlanProtectedItems(input *[]replicationrecoveryplans.RecoveryPlanProtectedItem) []string {
	protectedItemOutputs := make([]string, 0)
	if input != nil {
		for _, protectedItem := range *input {
			protectedItemOutputs = append(protectedItemOutputs, helpers.HandleAzureSdkForGoBug2824(*protectedItem.Id))
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
				PrimaryEdgeZone:  helpers.FlattenEdgeZone(a2aInput.PrimaryExtendedLocation),
				RecoveryEdgeZone: helpers.FlattenEdgeZone(a2aInput.RecoveryExtendedLocation),
			}
			output = append(output, o)
		}
	}
	return output
}

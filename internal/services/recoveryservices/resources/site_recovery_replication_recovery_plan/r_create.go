// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_recovery_plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
)

func (r Resource) Create() sdk.ResourceFunc {
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
			}

			// FailoverDeploymentModelClassic is used for other cloud service back up to Azure.

			groupValue, err := expandRecoveryGroup(model.ShutdownRecoveryGroup, model.FailoverRecoveryGroup, model.BootRecoveryGroup)
			if err != nil {
				return fmt.Errorf("expanding recovery group: %+v", err)
			}

			parameters := replicationrecoveryplans.CreateRecoveryPlanInput{
				Properties: replicationrecoveryplans.CreateRecoveryPlanInputProperties{
					PrimaryFabricId:         model.SourceRecoveryFabricId,
					RecoveryFabricId:        model.TargetRecoveryFabricId,
					FailoverDeploymentModel: pointer.To(replicationrecoveryplans.FailoverDeploymentModelResourceManager),
					Groups:                  groupValue,
				},
			}

			if len(model.A2ASettings) == 1 {
				parameters.Properties.ProviderSpecificInput = expandA2ASettings(model.A2ASettings[0])
			}

			if err := client.CreateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}
			metadata.SetID(id)

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
			PrimaryExtendedLocation:  helpers.ExpandEdgeZone(input.PrimaryEdgeZone),
			RecoveryExtendedLocation: helpers.ExpandEdgeZone(input.RecoveryEdgeZone),
		},
	}
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

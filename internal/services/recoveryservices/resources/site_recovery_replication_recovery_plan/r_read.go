// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package site_recovery_replication_recovery_plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Read() sdk.ResourceFunc {
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
					state.SourceRecoveryFabricId = helpers.HandleAzureSdkForGoBug2824(*prop.PrimaryFabricId)
				}
				if prop.RecoveryFabricId != nil {
					state.TargetRecoveryFabricId = helpers.HandleAzureSdkForGoBug2824(*prop.RecoveryFabricId)
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
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice(replicationrecoveryplans.PossibleValuesForPossibleOperationsDirections(), false),
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
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(replicationrecoveryplans.PossibleValuesForRecoveryPlanActionLocation(), false),
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

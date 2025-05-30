// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteRecoveryReplicationRecoveryPlanDataSource struct{}

type SiteRecoveryReplicationRecoveryPlanDataSourceModel struct {
	Name                   string                                                   `tfschema:"name"`
	RecoveryGroup          []RecoveryGroupDataSourceModel                           `tfschema:"recovery_group"`
	RecoveryVaultId        string                                                   `tfschema:"recovery_vault_id"`
	SourceRecoveryFabricId string                                                   `tfschema:"source_recovery_fabric_id"`
	TargetRecoveryFabricId string                                                   `tfschema:"target_recovery_fabric_id"`
	A2ASettings            []ReplicationRecoveryPlanA2ASpecificInputDataSourceModel `tfschema:"azure_to_azure_settings"`
}

type RecoveryGroupDataSourceModel struct {
	GroupType                string                  `tfschema:"type"`
	PostAction               []ActionDataSourceModel `tfschema:"post_action"`
	PreAction                []ActionDataSourceModel `tfschema:"pre_action"`
	ReplicatedProtectedItems []string                `tfschema:"replicated_protected_items"`
}

type ActionDataSourceModel struct {
	ActionDetailType        string   `tfschema:"type"`
	FabricLocation          string   `tfschema:"fabric_location"`
	FailOverDirections      []string `tfschema:"fail_over_directions"`
	FailOverTypes           []string `tfschema:"fail_over_types"`
	ManualActionInstruction string   `tfschema:"manual_action_instruction"`
	Name                    string   `tfschema:"name"`
	RunbookId               string   `tfschema:"runbook_id"`
	ScriptPath              string   `tfschema:"script_path"`
}

type ReplicationRecoveryPlanA2ASpecificInputDataSourceModel struct {
	PrimaryZone      string `tfschema:"primary_zone"`
	RecoveryZone     string `tfschema:"recovery_zone"`
	PrimaryEdgeZone  string `tfschema:"primary_edge_zone"`
	RecoveryEdgeZone string `tfschema:"recovery_edge_zone"`
}

var _ sdk.DataSource = SiteRecoveryReplicationRecoveryPlanDataSource{}

func (r SiteRecoveryReplicationRecoveryPlanDataSource) ResourceType() string {
	return "azurerm_site_recovery_replication_recovery_plan"
}

func (r SiteRecoveryReplicationRecoveryPlanDataSource) ModelObject() interface{} {
	return &SiteRecoveryReplicationRecoveryPlanModel{}
}

func (r SiteRecoveryReplicationRecoveryPlanDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationrecoveryplans.ValidateReplicationRecoveryPlanID
}

func (r SiteRecoveryReplicationRecoveryPlanDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel SiteRecoveryReplicationRecoveryPlanModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.RecoveryServices.ReplicationRecoveryPlansClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			vaultId, err := replicationrecoveryplans.ParseVaultID(metaModel.RecoveryVaultId)
			if err != nil {
				return err
			}

			id := replicationrecoveryplans.NewReplicationRecoveryPlanID(subscriptionId, vaultId.ResourceGroupName, vaultId.VaultName, metaModel.Name)
			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s was not found", id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("making Read request on site recovery replication plan %s : model is nil", id.String())
			}

			state := SiteRecoveryReplicationRecoveryPlanDataSourceModel{
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
					state.RecoveryGroup = flattenDataSourceRecoveryGroups(*group)
				}

				if details := prop.ProviderSpecificDetails; details != nil && len(*details) > 0 {
					state.A2ASettings = flattenDataSourceRecoveryPlanProviderSpecficInput(details)
				}
			}

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}
}

func flattenDataSourceRecoveryGroups(input []replicationrecoveryplans.RecoveryPlanGroup) []RecoveryGroupDataSourceModel {
	output := make([]RecoveryGroupDataSourceModel, 0)
	for _, groupItem := range input {
		recoveryGroupOutput := RecoveryGroupDataSourceModel{}
		recoveryGroupOutput.GroupType = string(groupItem.GroupType)
		if groupItem.ReplicationProtectedItems != nil {
			recoveryGroupOutput.ReplicatedProtectedItems = flattenRecoveryPlanProtectedItems(groupItem.ReplicationProtectedItems)
		}
		if groupItem.StartGroupActions != nil {
			recoveryGroupOutput.PreAction = flattenDataSourceRecoveryPlanActions(groupItem.StartGroupActions)
		}
		if groupItem.EndGroupActions != nil {
			recoveryGroupOutput.PostAction = flattenDataSourceRecoveryPlanActions(groupItem.EndGroupActions)
		}
		output = append(output, recoveryGroupOutput)
	}
	return output
}

func flattenDataSourceRecoveryPlanActions(input *[]replicationrecoveryplans.RecoveryPlanAction) []ActionDataSourceModel {
	actionOutputs := make([]ActionDataSourceModel, 0)
	for _, action := range *input {
		actionOutput := ActionDataSourceModel{
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

func (r SiteRecoveryReplicationRecoveryPlanDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"recovery_vault_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: replicationrecoveryplans.ValidateVaultID,
		},
	}
}

func flattenDataSourceRecoveryPlanProviderSpecficInput(input *[]replicationrecoveryplans.RecoveryPlanProviderSpecificDetails) []ReplicationRecoveryPlanA2ASpecificInputDataSourceModel {
	output := make([]ReplicationRecoveryPlanA2ASpecificInputDataSourceModel, 0)
	for _, providerSpecificInput := range *input {
		if a2aInput, ok := providerSpecificInput.(replicationrecoveryplans.RecoveryPlanA2ADetails); ok {
			o := ReplicationRecoveryPlanA2ASpecificInputDataSourceModel{
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

func (r SiteRecoveryReplicationRecoveryPlanDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"source_recovery_fabric_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"target_recovery_fabric_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"failover_deployment_model": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"recovery_group": {
			Type:     pluginsdk.TypeSet,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"replicated_protected_items": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"pre_action": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem:     dataSourceSiteRecoveryReplicationPlanActions(),
					},

					"post_action": {
						Type:     pluginsdk.TypeSet,
						Computed: true,
						Elem:     dataSourceSiteRecoveryReplicationPlanActions(),
					},
				},
			},
		},

		"azure_to_azure_settings": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"primary_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"recovery_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"primary_edge_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},

					"recovery_edge_zone": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func dataSourceSiteRecoveryReplicationPlanActions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type: pluginsdk.TypeList,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"type": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"fail_over_directions": {
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"fail_over_types": {
					Type:     pluginsdk.TypeSet,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
				"runbook_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"fabric_location": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"manual_action_instruction": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"script_path": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SiteRecoveryReplicationRecoveryPlanDataSource struct{}

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

			metadata.SetID(id)
			return metadata.Encode(&state)
		},
	}

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

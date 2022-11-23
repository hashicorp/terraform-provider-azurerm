package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-09-10/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceSiteRecoveryRecoverPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{

		Read: dataSourceSiteRecoveryRecoverPlanRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
			},

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

			"recovery_groups": {
				Type:     pluginsdk.TypeSet,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"group_type": {
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
							Elem:     dataSourceSiteRecoveryRecoverPlanActions(),
						},
						"post_action": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem:     dataSourceSiteRecoveryRecoverPlanActions(),
						},
					},
				},
			},
		},
	}
}

func dataSourceSiteRecoveryRecoverPlanActions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type: pluginsdk.TypeList,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"action_detail_type": {
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

func dataSourceSiteRecoveryRecoverPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoveryPlansClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery replication plan %s : %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("making Read request on site recovery replication plan %s : model is nil", id.String())
	}

	d.Set("name", id.RecoveryPlanName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.ResourceName)

	if prop := model.Properties; prop != nil {
		d.Set("source_recovery_fabric_id", prop.PrimaryFabricId)
		d.Set("target_recovery_fabric_id", prop.RecoveryFabricId)
		d.Set("failover_deployment_model", prop.FailoverDeploymentModel)

		protectedItemsOutput := make([]interface{}, 0)

		if group := prop.Groups; group != nil {
			for _, groupItem := range *group {
				if groupItem.GroupType == replicationrecoveryplans.RecoveryPlanGroupTypeBoot {
					for _, protectedItem := range *groupItem.ReplicationProtectedItems {
						protectedItemOutput := make(map[string]interface{})
						protectedItemOutput["id"] = protectedItem.Id
						protectedItemsOutput = append(protectedItemsOutput, protectedItemOutput)
					}
				}
			}
		}
		d.Set("replicated_protected_items", protectedItemsOutput)
	}

	return nil
}

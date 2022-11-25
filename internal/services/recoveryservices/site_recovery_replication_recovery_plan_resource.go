package recoveryservices

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-09-10/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-09-10/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSiteRecoveryReplicationRecoveryPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryReplicationRecoveryPlanCreate,
		Read:   resourceSiteRecoveryReplicationRecoveryPlanRead,
		Update: resourceSiteRecoveryReplicationRecoveryPlanUpdate,
		Delete: resourceSiteRecoveryReplicationRecoveryPlanDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": commonschema.ResourceGroupName(),

			"recovery_vault_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
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

			"failover_deployment_model": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(siterecovery.Classic),
					string(siterecovery.ResourceManager),
				}, false),
			},

			"recovery_groups": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"group_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(siterecovery.Boot),
								string(siterecovery.Shutdown),
								string(siterecovery.Failover),
							}, false),
						},
						"replicated_protected_items": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},
						"pre_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     schemaAction(),
						},
						"post_action": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem:     schemaAction(),
						},
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

			"action_detail_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(siterecovery.InstanceTypeAutomationRunbookActionDetails),
					string(siterecovery.InstanceTypeManualActionDetails),
					string(siterecovery.InstanceTypeScriptActionDetails),
				}, false),
			},

			"fail_over_directions": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(siterecovery.PrimaryToRecovery),
						string(siterecovery.RecoveryToPrimary),
					}, false),
				},
			},

			"fail_over_types": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(siterecovery.ReplicationProtectedItemOperationPlannedFailover),
						string(siterecovery.ReplicationProtectedItemOperationTestFailover),
						string(siterecovery.ReplicationProtectedItemOperationUnplannedFailover),
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
					string(siterecovery.Primary),
					string(siterecovery.Recovery),
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

func resourceSiteRecoveryReplicationRecoveryPlanCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoveryPlansClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationrecoveryplans.NewReplicationRecoveryPlanID(subscriptionId, resGroup, vaultName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
			if !response.WasNotFound(existing.HttpResponse) && !response.WasBadRequest(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing site recovery plan %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replication_recovery_plan", *existing.Model.Id)
		}
	}

	deploymentModel := replicationrecoveryplans.FailoverDeploymentModel(d.Get("failover_deployment_model").(string))

	parameters := replicationrecoveryplans.CreateRecoveryPlanInput{
		Properties: replicationrecoveryplans.CreateRecoveryPlanInputProperties{
			PrimaryFabricId:         d.Get("source_recovery_fabric_id").(string),
			RecoveryFabricId:        d.Get("target_recovery_fabric_id").(string),
			FailoverDeploymentModel: &deploymentModel,
			Groups:                  expandRecoverGroup(d.Get("recovery_groups").([]interface{})),
		},
	}

	err := client.CreateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery replication plan %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(id.ID())

	return resourceSiteRecoveryReplicationRecoveryPlanRead(d, meta)
}

func resourceSiteRecoveryReplicationRecoveryPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		if group := prop.Groups; group != nil {
			d.Set("recovery_groups", flattenRecoveryGroups(*group))
		}
	}
	return nil
}

func resourceSiteRecoveryReplicationRecoveryPlanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoveryPlansClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(d.Id())
	if err != nil {
		return fmt.Errorf("parse Site reocvery replication plan id: %+v", err)
	}

	parameters := replicationrecoveryplans.UpdateRecoveryPlanInput{
		Properties: &replicationrecoveryplans.UpdateRecoveryPlanInputProperties{},
	}
	err = client.UpdateThenPoll(ctx, *id, parameters)
	if err != nil {
		return fmt.Errorf("updating site recovery replication plan %s (vault %s): %+v", id.RecoveryPlanName, id.ResourceName, err)
	}

	return resourceSiteRecoveryReplicationRecoveryPlanRead(d, meta)
}

func resourceSiteRecoveryReplicationRecoveryPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoveryPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting site recovery protection replication plan %s : %+v", id.String(), err)
	}

	return nil
}

func expandRecoverGroup(input []interface{}) []replicationrecoveryplans.RecoveryPlanGroup {
	var output []replicationrecoveryplans.RecoveryPlanGroup
	for _, groupRaw := range input {
		groupInput := groupRaw.(map[string]interface{})

		var protectedItems []replicationrecoveryplans.RecoveryPlanProtectedItem
		for _, protectedItem := range groupInput["replicated_protected_items"].([]interface{}) {
			protectedItems = append(protectedItems, replicationrecoveryplans.RecoveryPlanProtectedItem{
				VirtualMachineId: utils.String(protectedItem.(string)),
			})
		}

		var startActions []replicationrecoveryplans.RecoveryPlanAction
		for _, startActionRaw := range groupInput["pre_action"].([]interface{}) {
			startActionInput := startActionRaw.(map[string]interface{})

			var failOverTypes []replicationrecoveryplans.ReplicationProtectedItemOperation
			for _, failOverType := range startActionInput["fail_over_types"].(*pluginsdk.Set).List() {
				failOverTypes = append(failOverTypes, replicationrecoveryplans.ReplicationProtectedItemOperation(failOverType.(string)))
			}

			startActions = append(startActions, replicationrecoveryplans.RecoveryPlanAction{
				ActionName: startActionInput["name"].(string),
				FailoverDirections: []replicationrecoveryplans.PossibleOperationsDirections{
					replicationrecoveryplans.PossibleOperationsDirectionsPrimaryToRecovery,
					replicationrecoveryplans.PossibleOperationsDirectionsRecoveryToPrimary,
				},
				FailoverTypes: failOverTypes,
				CustomDetails: expandActionDetail(startActionInput),
			})
		}

		output = append(output, replicationrecoveryplans.RecoveryPlanGroup{
			GroupType:                 replicationrecoveryplans.RecoveryPlanGroupType(groupInput["group_type"].(string)),
			ReplicationProtectedItems: &protectedItems,
		})

	}
	return output
}

func flattenRecoveryGroups(input []replicationrecoveryplans.RecoveryPlanGroup) []interface{} {
	output := make([]interface{}, 0)
	for _, groupItem := range input {
		recoveryGroupOutput := make(map[string]interface{})
		recoveryGroupOutput["group_type"] = groupItem.GroupType
		if groupItem.ReplicationProtectedItems != nil {
			recoveryGroupOutput["replicated_protected_items"] = flattenRecoveryPlanProtectedItems(groupItem.ReplicationProtectedItems)
		}
		if groupItem.StartGroupActions != nil {
			recoveryGroupOutput["pre_action"] = flattenRecoveryPlanActions(groupItem.StartGroupActions)
		}
		if groupItem.EndGroupActions != nil {
			recoveryGroupOutput["post_action"] = flattenRecoveryPlanActions(groupItem.StartGroupActions)
		}
		output = append(output, recoveryGroupOutput)
	}
	return output
}

func expandActionDetail(input map[string]interface{}) (output replicationrecoveryplans.RecoveryPlanActionDetails) {
	switch input["action_detail_type"].(string) {
	case "AutomationRunbookActionDetails":
		output = replicationrecoveryplans.RecoveryPlanAutomationRunbookActionDetails{
			RunbookId:      utils.String(input["runbook_id"].(string)),
			FabricLocation: replicationrecoveryplans.RecoveryPlanActionLocation(input["fabric_location"].(string)),
		}
	case "ManualActionDetails":
		output = replicationrecoveryplans.RecoveryPlanManualActionDetails{
			Description: utils.String(input["manual_action_instruction"].(string)),
		}
	case "ScriptActionDetails":
		output = replicationrecoveryplans.RecoveryPlanScriptActionDetails{
			Path:           input["script_path"].(string),
			FabricLocation: replicationrecoveryplans.RecoveryPlanActionLocation(input["fabric_location"].(string)),
		}
	}
	return
}

func flattenRecoveryPlanProtectedItems(input *[]replicationrecoveryplans.RecoveryPlanProtectedItem) []interface{} {
	protectedItemOutputs := make([]interface{}, 0)
	for _, protectedItem := range *input {
		protectedItemOutputs = append(protectedItemOutputs, protectedItem.VirtualMachineId)
	}
	return protectedItemOutputs
}

func flattenRecoveryPlanActions(input *[]replicationrecoveryplans.RecoveryPlanAction) []interface{} {
	actionOutputs := make([]interface{}, 0)
	for _, action := range *input {
		actionOutput := make(map[string]interface{}, 0)
		actionOutput["name"] = action.ActionName
		switch detail := action.CustomDetails.(type) {
		case replicationrecoveryplans.RecoveryPlanAutomationRunbookActionDetails:
			actionOutput["action_detail_type"] = "AutomationRunbookActionDetails"
			actionOutput["runbook_id"] = detail.RunbookId
			actionOutput["fabric_location"] = detail.FabricLocation
		case replicationrecoveryplans.RecoveryPlanManualActionDetails:
			actionOutput["action_detail_type"] = "ManualActionDetails"
			actionOutput["manual_action_instruction"] = detail.Description
		case replicationrecoveryplans.RecoveryPlanScriptActionDetails:
			actionOutput["action_detail_type"] = "ScriptActionDetails"
			actionOutput["script_path"] = detail.Path
			actionOutput["fabric_location"] = detail.FabricLocation
		}
		directions := make([]interface{}, 0)
		for _, direction := range action.FailoverDirections {
			directions = append(directions, string(direction))
		}
		actionOutput["fail_over_directions"] = directions
		failOverTypes := make([]interface{}, 0)
		for _, failOverType := range action.FailoverTypes {
			failOverTypes = append(failOverTypes, string(failOverType))
		}
		actionOutput["fail_over_types"] = failOverTypes
		actionOutputs = append(actionOutputs, actionOutput)
	}
	return actionOutputs
}

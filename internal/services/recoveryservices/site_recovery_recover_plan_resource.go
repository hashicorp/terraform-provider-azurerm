package recoveryservices

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-05-01/replicationrecoveryplans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSiteRecoveryRecoverPlan() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryRecoverPlanCreate,
		Read:   resourceSiteRecoveryRecoverPlanRead,
		Update: resourceSiteRecoveryRecoverPlanUpdate,
		Delete: resourceSiteRecoveryRecoverPlanDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationRecoverPlanID(id)
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
			"resource_group_name": azure.SchemaResourceGroupName(),

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
				ValidateFunc: azure.ValidateResourceID,
			},

			"target_recovery_fabric_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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
				Type:     pluginsdk.TypeSet,
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
						"pre_actions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     schemaAction(),
						},
						"post_actions": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem:     schemaAction(),
						},
					},
				},
			},
		},
	}
}

func schemaAction() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type: pluginsdk.TypeList,
		Elem: &pluginsdk.Resource{
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
						Type:     pluginsdk.TypeString,
						Required: true,
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
						Type:     pluginsdk.TypeString,
						Required: true,
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
		},
	}
}

func resourceSiteRecoveryRecoverPlanCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient
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
			return tf.ImportAsExistsError("azurerm_site_recovery_recover_plan", *existing.Model.Id)
		}
	}

	var recoverGroups []replicationrecoveryplans.RecoveryPlanGroup

	for _, groupRaw := range d.Get("recovery_groups").(*pluginsdk.Set).List() {
		groupInput := groupRaw.(map[string]interface{})

		var protectedItems []replicationrecoveryplans.RecoveryPlanProtectedItem
		for _, protectedItem := range groupInput["replicated_protected_items"].(*pluginsdk.Set).List() {
			protectedItems = append(protectedItems, replicationrecoveryplans.RecoveryPlanProtectedItem{
				VirtualMachineId: utils.String(protectedItem.(string)),
			})
		}

		var startActions []replicationrecoveryplans.RecoveryPlanAction
		for _, startActionRaw := range groupInput["pre_actions"].(*pluginsdk.Set).List() {
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

		recoverGroups = append(recoverGroups, replicationrecoveryplans.RecoveryPlanGroup{
			GroupType:                 replicationrecoveryplans.RecoveryPlanGroupType(groupInput["group_type"].(string)),
			ReplicationProtectedItems: &protectedItems,
		})

	}

	deploymentModel := replicationrecoveryplans.FailoverDeploymentModel(d.Get("failover_deployment_model").(string))

	parameters := replicationrecoveryplans.CreateRecoveryPlanInput{
		Properties: replicationrecoveryplans.CreateRecoveryPlanInputProperties{
			PrimaryFabricId:         d.Get("source_recovery_fabric_id").(string),
			RecoveryFabricId:        d.Get("target_recovery_fabric_id").(string),
			FailoverDeploymentModel: &deploymentModel,
			Groups:                  recoverGroups,
		},
	}

	err := client.CreateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery recover plan %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving site recovery recover plan %s (vault %s): %+v", name, vaultName, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("ertrieving site reocvery recover plan %s (vault: %s): model is nil", name, vaultName)
	}

	d.SetId(*model.Id)

	return resourceSiteRecoveryRecoverPlanRead(d, meta)
}

func resourceSiteRecoveryRecoverPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery recover plan %s : %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("making Read request on site recovery recover plan %s : model is nil", id.String())
	}

	d.Set("name", id.RecoveryPlanName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.ResourceName)

	if prop := model.Properties; prop != nil {
		d.Set("source_recovery_fabric_id", prop.PrimaryFabricId)
		d.Set("target_recovery_fabric_id", prop.RecoveryFabricId)
		d.Set("failover_deployment_model", prop.FailoverDeploymentModel)

		recoveryGroupOutputs := make([]interface{}, 0)

		if group := prop.Groups; group != nil {
			for _, groupItem := range *group {
				recoveryGroupOutput := make(map[string]interface{})
				recoveryGroupOutput["group_type"] = groupItem.GroupType
				if groupItem.ReplicationProtectedItems != nil {
					recoveryGroupOutput["replicated_protected_items"] = flattenRecoveryPlanProtectedItems(groupItem.ReplicationProtectedItems)
				}
				if groupItem.StartGroupActions != nil {
					recoveryGroupOutput["pre_actions"] = flattenRecoveryPlanActions(groupItem.StartGroupActions)
				}
				if groupItem.EndGroupActions != nil {
					recoveryGroupOutput["post_actions"] = flattenRecoveryPlanActions(groupItem.StartGroupActions)
				}
				recoveryGroupOutputs = append(recoveryGroupOutputs, recoveryGroupOutput)
			}
		}
		d.Set("recovery_groups", recoveryGroupOutputs)
	}
	return nil
}

func resourceSiteRecoveryRecoverPlanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(d.Id())
	if err != nil {
		return fmt.Errorf("parse Site reocvery recover plan id: %+v", err)
	}

	parameters := replicationrecoveryplans.UpdateRecoveryPlanInput{
		Properties: &replicationrecoveryplans.UpdateRecoveryPlanInputProperties{},
	}
	err = client.UpdateThenPoll(ctx, *id, parameters)
	if err != nil {
		return fmt.Errorf("updating site recovery recover plan %s (vault %s): %+v", id.RecoveryPlanName, id.ResourceName, err)
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving site recovery recover plan %s (vault %s): %+v", id.RecoveryPlanName, id.ResourceName, err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("ertrieving site reocvery recover plan %s (vault: %s): model is nil", id.RecoveryPlanName, id.ResourceName)
	}

	d.SetId(*model.Id)

	return resourceSiteRecoveryRecoverPlanRead(d, meta)
}

func resourceSiteRecoveryRecoverPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationrecoveryplans.ParseReplicationRecoveryPlanID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting site recovery protection recover plan %s : %+v", id.String(), err)
	}

	return nil
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

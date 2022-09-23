package recoveryservices

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
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
	"time"
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
							Type:     pluginsdk.TypeSet,
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
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			// NOTE: Bad Request due to https://github.com/Azure/azure-rest-api-specs/issues/12759
			if !utils.ResponseWasNotFound(existing.Response) && !utils.ResponseWasBadRequest(existing.Response) {
				return fmt.Errorf("checking for presence of existing site recovery plan %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_recover_plan", *existing.ID)
		}
	}

	var recoverGroups []siterecovery.RecoveryPlanGroup

	for _, groupRaw := range d.Get("recovery_groups").(*pluginsdk.Set).List() {
		groupInput := groupRaw.(map[string]interface{})

		var protectedItems []siterecovery.RecoveryPlanProtectedItem
		for _, protectedItem := range groupInput["replicated_protected_items"].(*pluginsdk.Set).List() {
			protectedItems = append(protectedItems, siterecovery.RecoveryPlanProtectedItem{
				ID: utils.String(protectedItem.(string)),
			})
		}

		var startActions []siterecovery.RecoveryPlanAction
		for _, startActionRaw := range groupInput["pre_actions"].(*pluginsdk.Set).List() {
			startActionInput := startActionRaw.(map[string]interface{})

			var failOverTypes []siterecovery.ReplicationProtectedItemOperation
			for _, failOverType := range startActionInput["fail_over_types"].(*pluginsdk.Set).List() {
				failOverTypes = append(failOverTypes, siterecovery.ReplicationProtectedItemOperation(failOverType.(string)))
			}

			startActions = append(startActions, siterecovery.RecoveryPlanAction{
				ActionName: utils.String(startActionInput["name"].(string)),
				FailoverDirections: &[]siterecovery.PossibleOperationsDirections{
					siterecovery.PrimaryToRecovery,
					siterecovery.RecoveryToPrimary,
				},
				FailoverTypes: &failOverTypes,
				CustomDetails: expandActionDetail(startActionInput),
			})
		}

		recoverGroups = append(recoverGroups, siterecovery.RecoveryPlanGroup{
			GroupType:                 siterecovery.RecoveryPlanGroupType(groupInput["group_type"].(string)),
			ReplicationProtectedItems: &protectedItems,
		})

	}

	parameters := siterecovery.CreateRecoveryPlanInput{
		Properties: &siterecovery.CreateRecoveryPlanInputProperties{
			PrimaryFabricID:         utils.String(d.Get("source_recovery_fabric_id").(string)),
			RecoveryFabricID:        utils.String(d.Get("target_recovery_fabric_id").(string)),
			FailoverDeploymentModel: siterecovery.FailoverDeploymentModel(d.Get("failover_deployment_model").(string)),
			Groups:                  &recoverGroups,
		},
	}

	future, err := client.Create(ctx, name, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery recover plan %s (vault %s): %+v", name, vaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of site recovery plan %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving site recovery recover plan %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryRecoverPlanRead(d, meta)
}

func resourceSiteRecoveryRecoverPlanRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationRecoverPlanID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id.ReplicationRecoveryPlanName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery recover plan %s : %+v", id.String(), err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)

	if prop := resp.Properties; prop != nil {
		d.Set("source_recovery_fabric_id", prop.PrimaryFabricID)
		d.Set("target_recovery_fabric_id", prop.RecoveryFabricID)
		d.Set("failover_deployment_model", prop.FailoverDeploymentModel)

		protectedItemsOutput := make([]interface{}, 0)

		if group := prop.Groups; group != nil {
			for _, groupItem := range *group {
				if groupItem.GroupType == siterecovery.Boot {
					for _, protectedItem := range *groupItem.ReplicationProtectedItems {
						protectedItemOutput := make(map[string]interface{})
						protectedItemOutput["id"] = protectedItem.ID
						protectedItemsOutput = append(protectedItemsOutput, protectedItemOutput)
					}
				}
			}
		}
		d.Set("replicated_protected_items", protectedItemsOutput)
	}
	return nil
}

func resourceSiteRecoveryRecoverPlanUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	parameters := siterecovery.UpdateRecoveryPlanInput{}
	future, err := client.Update(ctx, name, parameters)
	if err != nil {
		return fmt.Errorf("updating site recovery recover plan %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("updating site recovery recover plan %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving site recovery recover plan %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func resourceSiteRecoveryRecoverPlanDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationRecoverPlanID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationRecoverPlanClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, id.ReplicationRecoveryPlanName)
	if err != nil {
		return fmt.Errorf("deleting site recovery protection recover plan %s : %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of site recovery recover plan %s : %+v", id.String(), err)
	}

	return nil
}

func expandActionDetail(input map[string]interface{}) (output siterecovery.BasicRecoveryPlanActionDetails) {
	instanceType := siterecovery.InstanceTypeBasicRecoveryPlanActionDetails(input["action_detail_type"].(string))
	switch instanceType {
	case siterecovery.InstanceTypeAutomationRunbookActionDetails:
		output = siterecovery.RecoveryPlanAutomationRunbookActionDetails{
			InstanceType:   instanceType,
			RunbookID:      utils.String(input["runbook_id"].(string)),
			FabricLocation: siterecovery.RecoveryPlanActionLocation(input["fabric_location"].(string)),
		}
	case siterecovery.InstanceTypeManualActionDetails:
		output = siterecovery.RecoveryPlanManualActionDetails{
			InstanceType: instanceType,
			Description:  utils.String(input["manual_action_instruction"].(string)),
		}
	case siterecovery.InstanceTypeScriptActionDetails:
		output = siterecovery.RecoveryPlanScriptActionDetails{
			Path:           utils.String(input["script_path"].(string)),
			FabricLocation: siterecovery.RecoveryPlanActionLocation(input["fabric_location"].(string)),
		}
	}
	return
}

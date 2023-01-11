package recoveryservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectioncontainermappings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSiteRecoveryProtectionContainerMapping() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryContainerMappingCreate,
		Read:   resourceSiteRecoveryContainerMappingRead,
		Update: resourceSiteRecoveryContainerMappingUpdate,
		Delete: resourceSiteRecoveryServicesContainerMappingDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationProtectionContainerMappingsID(id)
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
			"recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_replication_policy_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"recovery_source_protection_container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_target_protection_container_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"automatic_update_extension": {
				Type:     pluginsdk.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Optional: true,
				Computed: true, // set it to computed because the service will return it no matter if we have passed it.
				Elem: &pluginsdk.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
						"automation_account_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},
		},
	}
}

func resourceSiteRecoveryContainerMappingCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("recovery_fabric_name").(string)
	policyId := d.Get("recovery_replication_policy_id").(string)
	protectionContainerName := d.Get("recovery_source_protection_container_name").(string)
	targetContainerId := d.Get("recovery_target_protection_container_id").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationprotectioncontainermappings.NewReplicationProtectionContainerMappingID(subscriptionId, resGroup, vaultName, fabricName, protectionContainerName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing site recovery protection container mapping %s (fabric %s, container %s): %+v", name, fabricName, protectionContainerName, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_protection_container_mapping", *existing.Model.Id)
		}
	}

	parameters := replicationprotectioncontainermappings.CreateProtectionContainerMappingInput{
		Properties: &replicationprotectioncontainermappings.CreateProtectionContainerMappingInputProperties{
			TargetProtectionContainerId: &targetContainerId,
			PolicyId:                    &policyId,
			//siterecovery.ReplicationProviderSpecificContainerMappingInput{}
			ProviderSpecificInput:       replicationprotectioncontainermappings.A2AContainerMappingInput{},
		},
	}

	autoUpdateEnabledValue, automationAccountArmId := expandAutoUpdateSettings(d.Get("automatic_update_extension").([]interface{}))

	parameters.Properties.ProviderSpecificInput = siterecovery.ReplicationProviderSpecificContainerMappingInput{}

	if autoUpdateEnabledValue == siterecovery.Enabled {
		parameters.Properties.ProviderSpecificInput = siterecovery.A2AContainerMappingInput{
			InstanceType:           siterecovery.InstanceTypeBasicReplicationProviderSpecificContainerMappingInputInstanceTypeA2A,
			AgentAutoUpdateStatus:  autoUpdateEnabledValue,
			AutomationAccountArmID: automationAccountArmId,
		}
	}

	err := client.CreateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating site recovery protection container mapping %s (vault %s): %+v", name, vaultName, err)
	}
	d.SetId(id.ID())

	return resourceSiteRecoveryContainerMappingRead(d, meta)
}

func resourceSiteRecoveryContainerMappingUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ReplicationProtectionContainerMappingsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectionContainerMappingName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : %+v", id.String(), err)
	}

	if resp.Properties == nil {
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : `Properties` is nil", id.String())
	}

	detail, ok := resp.Properties.ProviderSpecificDetails.AsA2AProtectionContainerMappingDetails()
	if !ok {
		return fmt.Errorf("update %s: type mismatch", id)
	}

	updateInput := siterecovery.A2AUpdateContainerMappingInput{
		AgentAutoUpdateStatus:  detail.AgentAutoUpdateStatus,
		AutomationAccountArmID: detail.AutomationAccountArmID,
	}

	if d.HasChange("automatic_update_extension") {
		autoUpdateEnabledValue, automationAccountArmId := expandAutoUpdateSettings(d.Get("automatic_update_extension").([]interface{}))
		updateInput.AgentAutoUpdateStatus = autoUpdateEnabledValue
		updateInput.AutomationAccountArmID = automationAccountArmId
	}

	update := siterecovery.UpdateProtectionContainerMappingInput{
		Properties: &siterecovery.UpdateProtectionContainerMappingInputProperties{
			ProviderSpecificInput: updateInput,
		},
	}

	future, err := client.Update(ctx, id.ReplicationFabricName, id.ReplicationProtectionContainerName, id.ReplicationProtectionContainerMappingName, update)
	if err != nil {
		return fmt.Errorf("update %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("update %s: %+v", id, err)
	}

	return resourceSiteRecoveryContainerMappingRead(d, meta)
}

func resourceSiteRecoveryContainerMappingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery protection container mapping %s : %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("retrieving site recovery protection container mapping %s (vault %s): model is nil", id.MappingName, id.ResourceName)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.ResourceName)
	d.Set("recovery_fabric_name", id.FabricName)
	d.Set("name", model.Name)
	if prop := model.Properties; prop != nil {
		d.Set("recovery_source_protection_container_name", model.Properties.SourceProtectionContainerFriendlyName)
		d.Set("recovery_replication_policy_id", model.Properties.PolicyId)
		d.Set("recovery_target_protection_container_id", model.Properties.TargetProtectionContainerId)
	}

	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("recovery_fabric_name", id.ReplicationFabricName)
	d.Set("recovery_source_protection_container_name", resp.Properties.SourceProtectionContainerFriendlyName)
	d.Set("name", resp.Name)
	d.Set("recovery_replication_policy_id", resp.Properties.PolicyID)
	d.Set("recovery_target_protection_container_id", resp.Properties.TargetProtectionContainerID)

	if detail, ok := resp.Properties.ProviderSpecificDetails.AsA2AProtectionContainerMappingDetails(); ok {
		d.Set("automatic_update_extension", flattenAutoUpdateSettings(detail))
	} else {
		d.Set("automatic_update_extension", flattenAutoUpdateSettings(nil))
	}

	return nil
}

func resourceSiteRecoveryServicesContainerMappingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationprotectioncontainermappings.ParseReplicationProtectionContainerMappingID(d.Id())
	if err != nil {
		return err
	}

	instanceType := "A2A"

	client := meta.(*clients.Client).RecoveryServices.ContainerMappingClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	input := replicationprotectioncontainermappings.RemoveProtectionContainerMappingInput{
		Properties: &replicationprotectioncontainermappings.RemoveProtectionContainerMappingInputProperties{
			ProviderSpecificInput: &replicationprotectioncontainermappings.ReplicationProviderContainerUnmappingInput{
				InstanceType: &instanceType,
			},
		},
	}

	err = client.DeleteThenPoll(ctx, *id, input)
	if err != nil {
		return fmt.Errorf("deleting site recovery protection container mapping %s : %+v", id.String(), err)
	}

	return nil
}

func expandAutoUpdateSettings(input []interface{}) (enabled siterecovery.AgentAutoUpdateStatus, automationAccountId *string) {
	if len(input) == 0 {
		return siterecovery.Disabled, nil
	}
	autoUpdateSettingMap := input[0].(map[string]interface{})

	autoUpdateEnabledValue := siterecovery.Disabled
	if autoUpdateSettingMap["enabled"].(bool) {
		autoUpdateEnabledValue = siterecovery.Enabled
	}

	var accountIdOutput *string
	accountId := autoUpdateSettingMap["automation_account_id"].(string)
	if accountId == "" {
		accountIdOutput = nil
	} else {
		accountIdOutput = &accountId
	}

	return autoUpdateEnabledValue, accountIdOutput
}

func flattenAutoUpdateSettings(input *siterecovery.A2AProtectionContainerMappingDetails) []interface{} {
	output := map[string]interface{}{}
	if input == nil {
		output["enabled"] = false
		return []interface{}{output}
	}
	output["enabled"] = input.AgentAutoUpdateStatus == siterecovery.Enabled
	output["automation_account_id"] = input.AutomationAccountArmID
	return []interface{}{output}
}

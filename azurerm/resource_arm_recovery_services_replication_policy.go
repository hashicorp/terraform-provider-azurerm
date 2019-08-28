package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesReplicationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryServicesReplicationPolicyCreate,
		Read:   resourceArmRecoveryServicesReplicationPolicyRead,
		Update: resourceArmRecoveryServicesReplicationPolicyUpdate,
		Delete: resourceArmSiteRecoveryReplicationPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
			},
			"recovery_point_retention_in_minutes": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(1, 365*24*60),
			},
			"application_consistent_snapshot_frequency_in_minutes": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(1, 365*24*60),
			},
		},
	}
}

func resourceArmRecoveryServicesReplicationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*ArmClient).recoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing recovery services replication policy %s: %+v", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_recovery_services_replication_policy", azure.HandleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	recoveryPoint := int32(d.Get("recovery_point_retention_in_minutes").(int))
	appConsitency := int32(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	var parameters = siterecovery.CreatePolicyInput{
		Properties: &siterecovery.CreatePolicyInputProperties{
			ProviderSpecificInput: &siterecovery.A2APolicyCreationInput{
				RecoveryPointHistory:            &recoveryPoint,
				AppConsistentFrequencyInMinutes: &appConsitency,
				MultiVMSyncStatus:               siterecovery.Enable,
				InstanceType:                    siterecovery.InstanceTypeBasicPolicyProviderSpecificInputInstanceTypeA2A,
			},
		},
	}
	future, err := client.Create(ctx, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating recovery services replication policy %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating recovery services replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(azure.HandleAzureSdkForGoBug2824(*resp.ID))

	return resourceArmRecoveryServicesReplicationPolicyRead(d, meta)
}

func resourceArmRecoveryServicesReplicationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*ArmClient).recoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	recoveryPoint := int32(d.Get("recovery_point_retention_in_minutes").(int))
	appConsitency := int32(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	var parameters = siterecovery.UpdatePolicyInput{
		Properties: &siterecovery.UpdatePolicyInputProperties{
			ReplicationProviderSettings: &siterecovery.A2APolicyCreationInput{
				RecoveryPointHistory:            &recoveryPoint,
				AppConsistentFrequencyInMinutes: &appConsitency,
				MultiVMSyncStatus:               siterecovery.Enable,
				InstanceType:                    siterecovery.InstanceTypeBasicPolicyProviderSpecificInputInstanceTypeA2A,
			},
		},
	}
	future, err := client.Update(ctx, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating recovery services replication policy %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error updating recovery services replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(azure.HandleAzureSdkForGoBug2824(*resp.ID))

	return resourceArmRecoveryServicesReplicationPolicyRead(d, meta)
}

func resourceArmRecoveryServicesReplicationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationPolicies"]

	client := meta.(*ArmClient).recoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on recovery services replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	if a2APolicyDetails, isA2A := resp.Properties.ProviderSpecificDetails.AsA2APolicyDetails(); isA2A {
		d.Set("recovery_point_retention_in_minutes", a2APolicyDetails.RecoveryPointHistory)
		d.Set("application_consistent_snapshot_frequency_in_minutes", a2APolicyDetails.AppConsistentFrequencyInMinutes)
	}
	return nil
}

func resourceArmSiteRecoveryReplicationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationPolicies"]

	client := meta.(*ArmClient).recoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	future, err := client.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting recovery services replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of recovery services replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	return nil
}

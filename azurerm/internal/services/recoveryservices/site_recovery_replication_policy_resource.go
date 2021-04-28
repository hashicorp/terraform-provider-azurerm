package recoveryservices

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/recoveryservices/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSiteRecoveryReplicationPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRecoveryReplicationPolicyCreate,
		Read:   resourceSiteRecoveryReplicationPolicyRead,
		Update: resourceSiteRecoveryReplicationPolicyUpdate,
		Delete: resourceSiteRecoveryReplicationPolicyDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.RecoveryServicesVaultName,
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

func resourceSiteRecoveryReplicationPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing site recovery replication policy %s: %+v", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replication_policy", handleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	recoveryPoint := int32(d.Get("recovery_point_retention_in_minutes").(int))
	appConsitency := int32(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	parameters := siterecovery.CreatePolicyInput{
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
		return fmt.Errorf("Error creating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func resourceSiteRecoveryReplicationPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	recoveryPoint := int32(d.Get("recovery_point_retention_in_minutes").(int))
	appConsitency := int32(d.Get("application_consistent_snapshot_frequency_in_minutes").(int))
	parameters := siterecovery.UpdatePolicyInput{
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
		return fmt.Errorf("Error updating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error updating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("Error retrieving site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func resourceSiteRecoveryReplicationPolicyRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationPolicies"]

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
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

func resourceSiteRecoveryReplicationPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	name := id.Path["replicationPolicies"]

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	return nil
}

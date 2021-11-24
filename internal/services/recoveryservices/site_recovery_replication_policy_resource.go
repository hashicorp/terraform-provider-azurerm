package recoveryservices

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
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

func resourceSiteRecoveryReplicationPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryReplicationPolicyCreate,
		Read:   resourceSiteRecoveryReplicationPolicyRead,
		Update: resourceSiteRecoveryReplicationPolicyUpdate,
		Delete: resourceSiteRecoveryReplicationPolicyDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationPolicyID(id)
			return err
		}),
		CustomizeDiff: resourceSiteRecoveryReplicationPolicyCustomDiff,

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
			"recovery_point_retention_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(0, 365*24*60),
			},
			"application_consistent_snapshot_frequency_in_minutes": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validation.IntBetween(0, 365*24*60),
			},
		},
	}
}

func resourceSiteRecoveryReplicationPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
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
				return fmt.Errorf("checking for presence of existing site recovery replication policy %s: %+v", name, err)
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
		return fmt.Errorf("creating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("creating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func resourceSiteRecoveryReplicationPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("updating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("updating site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return fmt.Errorf("retrieving site recovery replication policy %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	return resourceSiteRecoveryReplicationPolicyRead(d, meta)
}

func resourceSiteRecoveryReplicationPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationPolicyID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery replication policy %s : %+v", id.String(), err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("recovery_vault_name", id.VaultName)
	if a2APolicyDetails, isA2A := resp.Properties.ProviderSpecificDetails.AsA2APolicyDetails(); isA2A {
		d.Set("recovery_point_retention_in_minutes", a2APolicyDetails.RecoveryPointHistory)
		d.Set("application_consistent_snapshot_frequency_in_minutes", a2APolicyDetails.AppConsistentFrequencyInMinutes)
	}
	return nil
}

func resourceSiteRecoveryReplicationPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.ReplicationPolicyID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationPoliciesClient(id.ResourceGroup, id.VaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	future, err := client.Delete(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("deleting site recovery replication policy %s : %+v", id.String(), err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of site recovery replication policy %s : %+v", id.String(), err)
	}

	return nil
}

func resourceSiteRecoveryReplicationPolicyCustomDiff(ctx context.Context, d *pluginsdk.ResourceDiff, i interface{}) error {
	retention := d.Get("recovery_point_retention_in_minutes").(int)
	frequency := d.Get("application_consistent_snapshot_frequency_in_minutes").(int)

	if retention == 0 && frequency > 0 {
		return fmt.Errorf("application_consistent_snapshot_frequency_in_minutes cannot be greater than zero when recovery_point_retention_in_minutes is set to zero")
	}

	return nil
}

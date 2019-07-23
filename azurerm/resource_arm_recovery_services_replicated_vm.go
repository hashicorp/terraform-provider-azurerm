package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRecoveryServicesReplicatedVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRecoveryReplicatedItemCreate,
		Read:   resourceArmRecoveryReplicatedItemRead,
		Update: nil,
		Delete: resourceArmRecoveryReplicatedItemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupName(),

			"recovery_vault_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"source_recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"source_vm_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_recovery_fabric_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"recovery_replication_policy_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"source_recovery_protection_container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},
			"target_recovery_protection_container_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_resource_group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validate.NoEmptyStrings,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"managed_disk": {
				Type:       schema.TypeSet,
				ConfigMode: schema.SchemaConfigModeAttr,
				Optional:   true,
				ForceNew:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_id": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     validate.NoEmptyStrings,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"staging_storage_account_id": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     validate.NoEmptyStrings,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"target_resource_group_id": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     validate.NoEmptyStrings,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"targert_disk_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.StandardLRS),
								string(compute.PremiumLRS),
								string(compute.StandardSSDLRS),
								string(compute.UltraSSDLRS),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"targert_replica_disk_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(compute.StandardLRS),
								string(compute.PremiumLRS),
								string(compute.StandardSSDLRS),
								string(compute.UltraSSDLRS),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},
		},
	}
}

func resourceArmRecoveryReplicatedItemCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	fabricId := d.Get("source_vm_id").(string)
	policyId := d.Get("recovery_replication_policy_id").(string)
	sourceProtectionContainerName := d.Get("source_recovery_protection_container_name").(string)
	targetProtectionContainerId := d.Get("target_recovery_protection_container_id").(string)
	targetResourceGroupId := d.Get("target_resource_group_id").(string)

	client := meta.(*ArmClient).recoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	managedDisks := []siterecovery.A2AVMManagedDiskInputDetails{}

	for _, raw := range d.Get("managed_disk").(*schema.Set).List() {
		diskInput := raw.(map[string]interface{})
		diskId := diskInput["disk_id"].(string)
		primaryStagingAzureStorageAccountID := diskInput["staging_storage_account_id"].(string)
		recoveryResourceGroupId := diskInput["target_resource_group_id"].(string)
		targertReplicaDiskType := diskInput["targert_replica_disk_type"].(string)
		targertDiskType := diskInput["targert_disk_type"].(string)

		managedDisks = append(managedDisks, siterecovery.A2AVMManagedDiskInputDetails{
			DiskID:                              &diskId,
			PrimaryStagingAzureStorageAccountID: &primaryStagingAzureStorageAccountID,
			RecoveryResourceGroupID:             &recoveryResourceGroupId,
			RecoveryReplicaDiskAccountType:      &targertReplicaDiskType,
			RecoveryTargetDiskAccountType:       &targertDiskType,
		})
	}

	var parameters = siterecovery.EnableProtectionInput{
		Properties: &siterecovery.EnableProtectionInputProperties{
			PolicyID: &policyId,
			ProviderSpecificDetails: siterecovery.A2AEnableProtectionInput{
				FabricObjectID:          &fabricId,
				RecoveryContainerID:     &targetProtectionContainerId,
				RecoveryResourceGroupID: &targetResourceGroupId,
				VMManagedDisks:          &managedDisks,
			},
		},
	}
	future, err := client.Create(ctx, fabricName, sourceProtectionContainerName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating recovery network mapping: %+v", err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating recovery network mapping: %+v", err)
	}

	resp, err := client.Get(ctx, fabricName, sourceProtectionContainerName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving recovery network mapping: %+v", err)
	}

	d.SetId(*resp.ID)

	return resourceArmRecoveryReplicatedItemRead(d, meta)
}

func resourceArmRecoveryReplicatedItemRead(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectedItems"]

	client := meta.(*ArmClient).recoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext

	resp, err := client.Get(ctx, fabricName, protectionContainerName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on recovery services protection container mapping %q: %+v", name, err)
	}

	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	d.Set("source_recovery_fabric_name", fabricName)
	d.Set("target_recovery_fabric_id", resp.Properties.RecoveryFabricID)
	d.Set("recovery_replication_policy_id", resp.Properties.PolicyID)
	d.Set("source_recovery_protection_container_name", protectionContainerName)
	d.Set("target_recovery_protection_container_id", resp.Properties.RecoveryContainerID)

	if a2aDetails, isA2a := resp.Properties.ProviderSpecificDetails.AsA2AReplicationDetails(); isA2a {
		d.Set("source_vm_id", unifyIdCasing(a2aDetails.FabricObjectID))
		d.Set("target_resource_group_id", a2aDetails.RecoveryAzureResourceGroupID)
		if a2aDetails.ProtectedManagedDisks != nil {
			disksOutput := make([]map[string]interface{}, 0)
			for _, disk := range *a2aDetails.ProtectedManagedDisks {
				diskOutput := make(map[string]interface{})
				diskOutput["disk_id"] = disk.DiskID
				diskOutput["staging_storage_account_id"] = disk.PrimaryStagingAzureStorageAccountID
				diskOutput["target_resource_group_id"] = disk.RecoveryResourceGroupID
				diskOutput["targert_replica_disk_type"] = disk.RecoveryReplicaDiskAccountType
				diskOutput["targert_disk_type"] = disk.RecoveryTargetDiskAccountType

				disksOutput = append(disksOutput, diskOutput)
			}
			d.Set("managed_disk", disksOutput)
		}
	}

	return nil
}

func resourceArmRecoveryReplicatedItemDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectedItems"]

	disableProtectionInput := siterecovery.DisableProtectionInput{
		Properties: &siterecovery.DisableProtectionInputProperties{
			DisableProtectionReason:  siterecovery.NotSpecified,
			ReplicationProviderInput: siterecovery.DisableProtectionProviderSpecificInput{},
		},
	}

	client := meta.(*ArmClient).recoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	ctx := meta.(*ArmClient).StopContext
	future, err := client.Delete(ctx, fabricName, protectionContainerName, name, disableProtectionInput)
	if err != nil {
		return fmt.Errorf("Error deleting recovery services protection container mapping %q: %+v", name, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of recovery services protection container mapping %q: %+v", name, err)
	}
	return nil
}

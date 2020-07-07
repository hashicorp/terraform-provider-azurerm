package recoveryservices

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-01-10/siterecovery"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSiteRecoveryReplicatedVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSiteRecoveryReplicatedItemCreate,
		Read:   resourceArmSiteRecoveryReplicatedItemRead,
		Delete: resourceArmSiteRecoveryReplicatedItemDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(80 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(80 * time.Minute),
			Delete: schema.DefaultTimeout(80 * time.Minute),
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
				ValidateFunc: azure.ValidateRecoveryServicesVaultName,
			},
			"source_recovery_fabric_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"source_vm_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_recovery_fabric_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"recovery_replication_policy_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"source_recovery_protection_container_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"target_recovery_protection_container_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_resource_group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_availability_set_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"managed_disk": {
				Type:       schema.TypeSet,
				ConfigMode: schema.SchemaConfigModeAttr,
				Optional:   true,
				ForceNew:   true,
				Set:        resourceArmSiteRecoveryReplicatedVMDiskHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_id": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     validation.StringIsNotEmpty,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"staging_storage_account_id": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     azure.ValidateResourceID,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"target_resource_group_id": {
							Type:             schema.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     azure.ValidateResourceID,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"target_disk_type": {
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
						"target_replica_disk_type": {
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

func resourceArmSiteRecoveryReplicatedItemCreate(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	sourceVmId := d.Get("source_vm_id").(string)
	policyId := d.Get("recovery_replication_policy_id").(string)
	sourceProtectionContainerName := d.Get("source_recovery_protection_container_name").(string)
	targetProtectionContainerId := d.Get("target_recovery_protection_container_id").(string)
	targetResourceGroupId := d.Get("target_resource_group_id").(string)

	var targetAvailabilitySetID *string
	if id, isSet := d.GetOk("target_availability_set_id"); isSet {
		tmp := id.(string)
		targetAvailabilitySetID = &tmp
	} else {
		targetAvailabilitySetID = nil
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, sourceProtectionContainerName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replicated_vm", azure.HandleAzureSdkForGoBug2824(*existing.ID))
		}
	}

	managedDisks := []siterecovery.A2AVMManagedDiskInputDetails{}

	for _, raw := range d.Get("managed_disk").(*schema.Set).List() {
		diskInput := raw.(map[string]interface{})
		diskId := diskInput["disk_id"].(string)
		primaryStagingAzureStorageAccountID := diskInput["staging_storage_account_id"].(string)
		recoveryResourceGroupId := diskInput["target_resource_group_id"].(string)
		targetReplicaDiskType := diskInput["target_replica_disk_type"].(string)
		targetDiskType := diskInput["target_disk_type"].(string)

		managedDisks = append(managedDisks, siterecovery.A2AVMManagedDiskInputDetails{
			DiskID:                              &diskId,
			PrimaryStagingAzureStorageAccountID: &primaryStagingAzureStorageAccountID,
			RecoveryResourceGroupID:             &recoveryResourceGroupId,
			RecoveryReplicaDiskAccountType:      &targetReplicaDiskType,
			RecoveryTargetDiskAccountType:       &targetDiskType,
		})
	}

	var parameters = siterecovery.EnableProtectionInput{
		Properties: &siterecovery.EnableProtectionInputProperties{
			PolicyID: &policyId,
			ProviderSpecificDetails: siterecovery.A2AEnableProtectionInput{
				FabricObjectID:            &sourceVmId,
				RecoveryContainerID:       &targetProtectionContainerId,
				RecoveryResourceGroupID:   &targetResourceGroupId,
				RecoveryAvailabilitySetID: targetAvailabilitySetID,
				VMManagedDisks:            &managedDisks,
			},
		},
	}
	future, err := client.Create(ctx, fabricName, sourceProtectionContainerName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error creating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	resp, err := client.Get(ctx, fabricName, sourceProtectionContainerName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(azure.HandleAzureSdkForGoBug2824(*resp.ID))

	return resourceArmSiteRecoveryReplicatedItemRead(d, meta)
}

func resourceArmSiteRecoveryReplicatedItemRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectedItems"]

	client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, fabricName, protectionContainerName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("recovery_vault_name", vaultName)
	d.Set("source_recovery_fabric_name", fabricName)
	d.Set("target_recovery_fabric_id", resp.Properties.RecoveryFabricID)
	d.Set("recovery_replication_policy_id", resp.Properties.PolicyID)
	d.Set("source_recovery_protection_container_name", protectionContainerName)
	d.Set("target_recovery_protection_container_id", resp.Properties.RecoveryContainerID)

	if a2aDetails, isA2a := resp.Properties.ProviderSpecificDetails.AsA2AReplicationDetails(); isA2a {
		d.Set("source_vm_id", a2aDetails.FabricObjectID)
		d.Set("target_resource_group_id", a2aDetails.RecoveryAzureResourceGroupID)
		d.Set("target_availability_set_id", a2aDetails.RecoveryAvailabilitySet)
		if a2aDetails.ProtectedManagedDisks != nil {
			disksOutput := make([]interface{}, 0)
			for _, disk := range *a2aDetails.ProtectedManagedDisks {
				diskOutput := make(map[string]interface{})
				diskOutput["disk_id"] = *disk.DiskID
				diskOutput["staging_storage_account_id"] = *disk.PrimaryStagingAzureStorageAccountID
				diskOutput["target_resource_group_id"] = *disk.RecoveryResourceGroupID
				diskOutput["target_replica_disk_type"] = *disk.RecoveryReplicaDiskAccountType
				diskOutput["target_disk_type"] = *disk.RecoveryTargetDiskAccountType

				disksOutput = append(disksOutput, diskOutput)
			}
			d.Set("managed_disk", schema.NewSet(resourceArmSiteRecoveryReplicatedVMDiskHash, disksOutput))
		}
	}

	return nil
}

func resourceArmSiteRecoveryReplicatedItemDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
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

	client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	future, err := client.Delete(ctx, fabricName, protectionContainerName, name, disableProtectionInput)
	if err != nil {
		return fmt.Errorf("Error deleting site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
	}
	return nil
}

func resourceArmSiteRecoveryReplicatedVMDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["disk_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
	}

	return hashcode.String(buf.String())
}

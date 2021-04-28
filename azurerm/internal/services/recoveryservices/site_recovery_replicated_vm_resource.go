package recoveryservices

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/recoveryservices/validate"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2020-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSiteRecoveryReplicatedVM() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRecoveryReplicatedItemCreate,
		Read:   resourceSiteRecoveryReplicatedItemRead,
		Update: resourceSiteRecoveryReplicatedItemUpdate,
		Delete: resourceSiteRecoveryReplicatedItemDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
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
				ValidateFunc: validate.RecoveryServicesVaultName,
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
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},
			"target_network_id": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"managed_disk": {
				Type:       schema.TypeSet,
				ConfigMode: schema.SchemaConfigModeAttr,
				Optional:   true,
				ForceNew:   true,
				Set:        resourceSiteRecoveryReplicatedVMDiskHash,
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
			"network_interface": {
				Type:       schema.TypeSet,
				ConfigMode: schema.SchemaConfigModeAttr,
				Computed:   true,
				Optional:   true,
				Elem:       networkInterfaceResource(),
			},
		},
	}
}

func networkInterfaceResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_network_interface_id": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"target_static_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"target_subnet_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"recovery_public_ip_address_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}
}

func resourceSiteRecoveryReplicatedItemCreate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	name := d.Get("name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	sourceVmId := d.Get("source_vm_id").(string)
	policyId := d.Get("recovery_replication_policy_id").(string)
	sourceProtectionContainerName := d.Get("source_recovery_protection_container_name").(string)
	targetProtectionContainerId := d.Get("target_recovery_protection_container_id").(string)
	targetResourceGroupId := d.Get("target_resource_group_id").(string)

	var targetAvailabilitySetID *string
	if id, isSet := d.GetOk("target_availability_set_id"); isSet {
		targetAvailabilitySetID = utils.String(id.(string))
	} else {
		targetAvailabilitySetID = nil
	}

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.IsNewResource() {
		existing, err := client.Get(ctx, fabricName, sourceProtectionContainerName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replicated_vm", handleAzureSdkForGoBug2824(*existing.ID))
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

	parameters := siterecovery.EnableProtectionInput{
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

	d.SetId(handleAzureSdkForGoBug2824(*resp.ID))

	// We are not allowed to configure the NIC on the initial setup, and the VM has to be replicated before
	// we can reconfigure. Hence this call to update when we create.
	return resourceSiteRecoveryReplicatedItemUpdate(d, meta)
}

func resourceSiteRecoveryReplicatedItemUpdate(d *schema.ResourceData, meta interface{}) error {
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)

	// We are only allowed to update the configuration once the VM is fully protected
	state, err := waitForReplicationToBeHealthy(d, meta)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	sourceProtectionContainerName := d.Get("source_recovery_protection_container_name").(string)
	targetNetworkId := d.Get("target_network_id").(string)

	var targetAvailabilitySetID *string
	if id, isSet := d.GetOk("target_availability_set_id"); isSet {
		tmp := id.(string)
		targetAvailabilitySetID = &tmp
	} else {
		targetAvailabilitySetID = nil
	}

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vmNics := []siterecovery.VMNicInputDetails{}
	for _, raw := range d.Get("network_interface").(*schema.Set).List() {
		vmNicInput := raw.(map[string]interface{})
		sourceNicId := vmNicInput["source_network_interface_id"].(string)
		targetStaticIp := vmNicInput["target_static_ip"].(string)
		targetSubnetName := vmNicInput["target_subnet_name"].(string)
		recoveryPublicIPAddressID := vmNicInput["recovery_public_ip_address_id"].(string)

		nicId := findNicId(state, sourceNicId)
		if nicId == nil {
			return fmt.Errorf("Error updating replicated vm %s (vault %s): Trying to update NIC that is not known by Azure %s", name, vaultName, sourceNicId)
		}
		vmNics = append(vmNics, siterecovery.VMNicInputDetails{
			NicID:                     nicId,
			RecoveryVMSubnetName:      &targetSubnetName,
			ReplicaNicStaticIPAddress: &targetStaticIp,
			RecoveryPublicIPAddressID: &recoveryPublicIPAddressID,
		})
	}

	managedDisks := []siterecovery.A2AVMManagedDiskUpdateDetails{}
	for _, raw := range d.Get("managed_disk").(*schema.Set).List() {
		diskInput := raw.(map[string]interface{})
		diskId := diskInput["disk_id"].(string)
		targetReplicaDiskType := diskInput["target_replica_disk_type"].(string)
		targetDiskType := diskInput["target_disk_type"].(string)

		managedDisks = append(managedDisks, siterecovery.A2AVMManagedDiskUpdateDetails{
			DiskID:                         &diskId,
			RecoveryReplicaDiskAccountType: &targetReplicaDiskType,
			RecoveryTargetDiskAccountType:  &targetDiskType,
		})
	}

	if targetNetworkId == "" {
		// No target network id was specified, so we want to preserve what was selected
		if a2aDetails, isA2a := state.Properties.ProviderSpecificDetails.AsA2AReplicationDetails(); isA2a {
			if a2aDetails.SelectedRecoveryAzureNetworkID != nil {
				targetNetworkId = *a2aDetails.SelectedRecoveryAzureNetworkID
			} else {
				return fmt.Errorf("target_network_id must be set when a network_interface is configured")
			}
		} else {
			return fmt.Errorf("target_network_id must be set when a network_interface is configured")
		}
	}

	parameters := siterecovery.UpdateReplicationProtectedItemInput{
		Properties: &siterecovery.UpdateReplicationProtectedItemInputProperties{
			RecoveryAzureVMName:            &name,
			SelectedRecoveryAzureNetworkID: &targetNetworkId,
			VMNics:                         &vmNics,
			RecoveryAvailabilitySetID:      targetAvailabilitySetID,
			ProviderSpecificDetails: siterecovery.A2AUpdateReplicationProtectedItemInput{
				ManagedDiskUpdateDetails: &managedDisks,
			},
		},
	}

	future, err := client.Update(ctx, fabricName, sourceProtectionContainerName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error updating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	return resourceSiteRecoveryReplicatedItemRead(d, meta)
}

func findNicId(state *siterecovery.ReplicationProtectedItem, sourceNicId string) *string {
	if a2aDetails, isA2a := state.Properties.ProviderSpecificDetails.AsA2AReplicationDetails(); isA2a {
		if a2aDetails.VMNics != nil {
			for _, nic := range *a2aDetails.VMNics {
				if nic.SourceNicArmID != nil && *nic.SourceNicArmID == sourceNicId {
					return nic.NicID
				}
			}
		}
	}
	return nil
}

func resourceSiteRecoveryReplicatedItemRead(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectedItems"]

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
		d.Set("target_network_id", a2aDetails.SelectedRecoveryAzureNetworkID)
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
			d.Set("managed_disk", schema.NewSet(resourceSiteRecoveryReplicatedVMDiskHash, disksOutput))
		}

		if a2aDetails.VMNics != nil {
			nicsOutput := make([]interface{}, 0)
			for _, nic := range *a2aDetails.VMNics {
				nicOutput := make(map[string]interface{})
				if nic.SourceNicArmID != nil {
					nicOutput["source_network_interface_id"] = *nic.SourceNicArmID
				}
				if nic.ReplicaNicStaticIPAddress != nil {
					nicOutput["target_static_ip"] = *nic.ReplicaNicStaticIPAddress
				}
				if nic.RecoveryVMSubnetName != nil {
					nicOutput["target_subnet_name"] = *nic.RecoveryVMSubnetName
				}
				if nic.RecoveryPublicIPAddressID != nil {
					nicOutput["recovery_public_ip_address_id"] = *nic.RecoveryPublicIPAddressID
				}
				nicsOutput = append(nicsOutput, nicOutput)
			}
			d.Set("network_interface", schema.NewSet(schema.HashResource(networkInterfaceResource()), nicsOutput))
		}
	}

	return nil
}

func resourceSiteRecoveryReplicatedItemDelete(d *schema.ResourceData, meta interface{}) error {
	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vaultName := id.Path["vaults"]
	client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
	fabricName := id.Path["replicationFabrics"]
	protectionContainerName := id.Path["replicationProtectionContainers"]
	name := id.Path["replicationProtectedItems"]

	disableProtectionInput := siterecovery.DisableProtectionInput{
		Properties: &siterecovery.DisableProtectionInputProperties{
			DisableProtectionReason:  siterecovery.NotSpecified,
			ReplicationProviderInput: siterecovery.DisableProtectionProviderSpecificInput{},
		},
	}

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

func resourceSiteRecoveryReplicatedVMDiskHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["disk_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
	}

	return schema.HashString(buf.String())
}

func waitForReplicationToBeHealthy(d *schema.ResourceData, meta interface{}) (*siterecovery.ReplicationProtectedItem, error) {
	log.Printf("Waiting for Site Recover to replicate VM.")
	stateConf := &resource.StateChangeConf{
		Target:       []string{"Protected"},
		Refresh:      waitForReplicationToBeHealthyRefreshFunc(d, meta),
		PollInterval: time.Minute,
	}

	stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)

	result, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf("Error waiting for site recovery to replicate vm: %+v", err)
	}

	protectedItem, ok := result.(siterecovery.ReplicationProtectedItem)
	if ok {
		return &protectedItem, nil
	} else {
		return nil, fmt.Errorf("Error waiting for site recovery return incompatible tyupe")
	}
}

func waitForReplicationToBeHealthyRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := azure.ParseAzureResourceID(d.Id())
		if err != nil {
			return nil, "", err
		}

		resGroup := id.ResourceGroup
		vaultName := id.Path["vaults"]
		client := meta.(*clients.Client).RecoveryServices.ReplicationMigrationItemsClient(resGroup, vaultName)
		fabricName := id.Path["replicationFabrics"]
		protectionContainerName := id.Path["replicationProtectionContainers"]
		name := id.Path["replicationProtectedItems"]

		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, fabricName, protectionContainerName, name)
		if err != nil {
			return nil, "", fmt.Errorf("Error making Read request on site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
		}

		if resp.Properties == nil {
			return nil, "", fmt.Errorf("Missing Properties in response when making Read request on site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
		}

		if resp.Properties.ProviderSpecificDetails == nil {
			return nil, "", fmt.Errorf("Missing Properties.ProviderSpecificDetails in response when making Read request on site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
		}

		// Find first disk that is not fully replicated yet
		if a2aDetails, isA2a := resp.Properties.ProviderSpecificDetails.AsA2AReplicationDetails(); isA2a {
			if a2aDetails.MonitoringPercentageCompletion != nil {
				log.Printf("Waiting for Site Recover to replicate VM, %d%% complete.", *a2aDetails.MonitoringPercentageCompletion)
			}
			if a2aDetails.VMProtectionState != nil {
				return resp, *a2aDetails.VMProtectionState, nil
			}
		}

		if resp.Properties.ReplicationHealth == nil {
			return nil, "", fmt.Errorf("Missing ReplicationHealth in response when making Read request on site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
		}
		return resp, *resp.Properties.ReplicationHealth, nil
	}
}

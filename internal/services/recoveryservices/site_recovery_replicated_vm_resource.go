package recoveryservices

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2018-07-10/siterecovery" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/edgezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSiteRecoveryReplicatedVM() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSiteRecoveryReplicatedItemCreate,
		Read:   resourceSiteRecoveryReplicatedItemRead,
		Update: resourceSiteRecoveryReplicatedItemUpdate,
		Delete: resourceSiteRecoveryReplicatedItemDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ReplicationProtectedItemID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(180 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(80 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(80 * time.Minute),
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
			"source_recovery_fabric_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"source_vm_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"target_recovery_fabric_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"recovery_replication_policy_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"source_recovery_protection_container_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"target_recovery_protection_container_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"target_resource_group_id": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"target_availability_set_id": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ValidateFunc:     azure.ValidateResourceID,
				DiffSuppressFunc: suppress.CaseDifference,
				ConflictsWith: []string{
					"target_zone",
				},
			},

			"target_zone": commonschema.ZoneSingleOptionalForceNew(),

			"target_network_id": {
				Type:         pluginsdk.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"test_network_id": {
				Type:         pluginsdk.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"target_edge_zone": commonschema.EdgeZoneOptionalForceNew(),

			"unmanaged_disk": {
				Type:       pluginsdk.TypeSet,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Optional:   true,
				ForceNew:   true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disk_uri": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"staging_storage_account_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"target_storage_account_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"multi_vm_group_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"managed_disk": {
				Type:       pluginsdk.TypeSet,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Optional:   true,
				ForceNew:   true,
				Set:        resourceSiteRecoveryReplicatedVMDiskHash,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"disk_id": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     validation.StringIsNotEmpty,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"staging_storage_account_id": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     azure.ValidateResourceID,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"target_resource_group_id": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							ForceNew:         true,
							ValidateFunc:     azure.ValidateResourceID,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"target_disk_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(disks.DiskStorageAccountTypesStandardLRS),
								string(disks.DiskStorageAccountTypesPremiumLRS),
								string(disks.DiskStorageAccountTypesStandardSSDLRS),
								string(disks.DiskStorageAccountTypesUltraSSDLRS),
							}, false),
						},
						"target_replica_disk_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(disks.DiskStorageAccountTypesStandardLRS),
								string(disks.DiskStorageAccountTypesPremiumLRS),
								string(disks.DiskStorageAccountTypesStandardSSDLRS),
								string(disks.DiskStorageAccountTypesUltraSSDLRS),
							}, false),
						},
						"target_disk_encryption_set_id": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ForceNew:         true,
							ValidateFunc:     azure.ValidateResourceID,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"target_disk_encryption": {
							Type:       pluginsdk.TypeList,
							ConfigMode: pluginsdk.SchemaConfigModeAttr,
							Optional:   true,
							MaxItems:   1,
							Elem:       diskEncryptionResource(),
						},
					},
				},
			},
			"target_proximity_placement_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"target_boot_diagnostic_storage_account_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"target_capacity_reservation_group_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"target_virtual_machine_scale_set_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"network_interface": {
				Type:       pluginsdk.TypeSet, // use set to avoid diff caused by different orders.
				Set:        resourceSiteRecoveryReplicatedVMNicHash,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Computed:   true,
				Optional:   true,
				Elem:       networkInterfaceResource(),
			},
		},
	}
}

func networkInterfaceResource() *pluginsdk.Resource {
	out := &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"source_network_interface_id": {
				Type:         pluginsdk.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"failover_test_static_ip": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"target_static_ip": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"failover_test_subnet_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"target_subnet_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"failover_test_public_ip_address_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     false,
				ValidateFunc: azure.ValidateResourceID,
			},
			"recovery_public_ip_address_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
	}

	if !features.FourPointOhBeta() {
		out.Schema["is_primary"] = &pluginsdk.Schema{
			Deprecated: "this property is not used and will be removed in version 4.0 of the provider",
			Type:       pluginsdk.TypeBool,
			Optional:   true,
			Default:    false,
		}
	}
	return out
}

func diskEncryptionResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"disk_encryption_key": {
				Type:       pluginsdk.TypeList,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Required:   true,
				MaxItems:   1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"secret_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: keyVaultValidate.NestedItemId,
						},

						"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(commonids.KeyVaultId{}),
					},
				},
			},
			"key_encryption_key": {
				Type:       pluginsdk.TypeList,
				ConfigMode: pluginsdk.SchemaConfigModeAttr,
				Optional:   true,
				MaxItems:   1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_url": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: keyVaultValidate.NestedItemId,
						},

						"vault_id": commonschema.ResourceIDReferenceRequiredForceNew(commonids.KeyVaultId{}),
					},
				},
			},
		},
	}
}

func resourceSiteRecoveryReplicatedItemCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

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

	var targetAvailabilityZone *string
	if zone, isSet := d.GetOk("target_zone"); isSet {
		targetAvailabilityZone = utils.String(zone.(string))
	} else {
		targetAvailabilityZone = nil
	}

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := replicationprotecteditems.NewReplicationProtectedItemID(subscriptionId, resGroup, vaultName, fabricName, sourceProtectionContainerName, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing site recovery replicated vm %s (vault %s): %+v", name, vaultName, err)
			}
		}

		if existing.Model != nil && existing.Model.Id != nil && *existing.Model.Id != "" {
			return tf.ImportAsExistsError("azurerm_site_recovery_replicated_vm", *existing.Model.Id)
		}
	}

	var managedDisks []replicationprotecteditems.A2AVMManagedDiskInputDetails

	for _, raw := range d.Get("managed_disk").(*pluginsdk.Set).List() {
		diskInput := raw.(map[string]interface{})
		diskId := diskInput["disk_id"].(string)
		primaryStagingAzureStorageAccountID := diskInput["staging_storage_account_id"].(string)
		recoveryResourceGroupId := diskInput["target_resource_group_id"].(string)
		targetReplicaDiskType := diskInput["target_replica_disk_type"].(string)
		targetDiskType := diskInput["target_disk_type"].(string)
		targetEncryptionDiskSetID := diskInput["target_disk_encryption_set_id"].(string)

		managedDisks = append(managedDisks, replicationprotecteditems.A2AVMManagedDiskInputDetails{
			DiskId:                              diskId,
			PrimaryStagingAzureStorageAccountId: primaryStagingAzureStorageAccountID,
			RecoveryResourceGroupId:             recoveryResourceGroupId,
			RecoveryReplicaDiskAccountType:      &targetReplicaDiskType,
			RecoveryTargetDiskAccountType:       &targetDiskType,
			RecoveryDiskEncryptionSetId:         &targetEncryptionDiskSetID,
			DiskEncryptionInfo:                  expandDiskEncryption(diskInput["target_disk_encryption"].([]interface{})),
		})
	}

	var vmDisks []replicationprotecteditems.A2AVMDiskInputDetails
	for _, raw := range d.Get("unmanaged_disk").(*pluginsdk.Set).List() {
		diskInput := raw.(map[string]interface{})
		diskUri := diskInput["disk_uri"].(string)
		primaryStagingAzureStorageAccountID := diskInput["staging_storage_account_id"].(string)
		recoveryAzureStorageAccountId := diskInput["target_storage_account_id"].(string)

		vmDisks = append(vmDisks, replicationprotecteditems.A2AVMDiskInputDetails{
			DiskUri:                             diskUri,
			PrimaryStagingAzureStorageAccountId: primaryStagingAzureStorageAccountID,
			RecoveryAzureStorageAccountId:       recoveryAzureStorageAccountId,
		})

	}

	var protectionInput replicationprotecteditems.EnableProtectionProviderSpecificInput = replicationprotecteditems.A2AEnableProtectionInput{
		FabricObjectId:                     sourceVmId,
		RecoveryContainerId:                &targetProtectionContainerId,
		RecoveryResourceGroupId:            &targetResourceGroupId,
		RecoveryAvailabilitySetId:          targetAvailabilitySetID,
		RecoveryAvailabilityZone:           targetAvailabilityZone,
		MultiVMGroupName:                   utils.String(d.Get("multi_vm_group_name").(string)),
		RecoveryProximityPlacementGroupId:  utils.String(d.Get("target_proximity_placement_group_id").(string)),
		RecoveryBootDiagStorageAccountId:   utils.String(d.Get("target_boot_diagnostic_storage_account_id").(string)),
		RecoveryCapacityReservationGroupId: utils.String(d.Get("target_capacity_reservation_group_id").(string)),
		RecoveryVirtualMachineScaleSetId:   utils.String(d.Get("target_virtual_machine_scale_set_id").(string)),
		VMManagedDisks:                     &managedDisks,
		VMDisks:                            &vmDisks,
		RecoveryExtendedLocation:           expandEdgeZone(d.Get("target_edge_zone").(string)),
	}
	parameters := replicationprotecteditems.EnableProtectionInput{
		Properties: &replicationprotecteditems.EnableProtectionInputProperties{
			PolicyId:                &policyId,
			ProviderSpecificDetails: &protectionInput,
		},
	}

	err := client.CreateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	d.SetId(id.ID())

	// We are not allowed to configure the NIC on the initial setup, and the VM has to be replicated before
	// we can reconfigure. Hence this call to update when we create.
	return resourceSiteRecoveryReplicatedItemUpdateInternal(ctx, d, meta)
}

func resourceSiteRecoveryReplicatedItemUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	return resourceSiteRecoveryReplicatedItemUpdateInternal(ctx, d, meta)
}

func resourceSiteRecoveryReplicatedItemUpdateInternal(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

	// We are only allowed to update the configuration once the VM is fully protected
	state, err := waitForReplicationToBeHealthy(ctx, d, meta)
	if err != nil {
		return err
	}

	// TODO: this should be parsing the Resource ID and not the config
	resGroup := d.Get("resource_group_name").(string)
	vaultName := d.Get("recovery_vault_name").(string)
	name := d.Get("name").(string)
	fabricName := d.Get("source_recovery_fabric_name").(string)
	sourceProtectionContainerName := d.Get("source_recovery_protection_container_name").(string)
	targetNetworkId := d.Get("target_network_id").(string)
	testNetworkId := d.Get("test_network_id").(string)

	id := replicationprotecteditems.NewReplicationProtectedItemID(subscriptionId, resGroup, vaultName, fabricName, sourceProtectionContainerName, name)

	var targetAvailabilitySetID *string
	if id, isSet := d.GetOk("target_availability_set_id"); isSet {
		tmp := id.(string)
		targetAvailabilitySetID = &tmp
	} else {
		targetAvailabilitySetID = nil
	}

	vmNics := make([]replicationprotecteditems.VMNicInputDetails, 0)
	nicList := d.Get("network_interface").(*pluginsdk.Set).List()
	for _, raw := range nicList {
		vmNicInput := raw.(map[string]interface{})
		sourceNicId := vmNicInput["source_network_interface_id"].(string)
		targetStaticIp := vmNicInput["target_static_ip"].(string)
		targetSubnetName := vmNicInput["target_subnet_name"].(string)
		recoveryPublicIPAddressID := vmNicInput["recovery_public_ip_address_id"].(string)
		testStaticIp := vmNicInput["target_static_ip"].(string)
		testSubNetName := vmNicInput["target_subnet_name"].(string)
		testPublicIpAddressID := vmNicInput["recovery_public_ip_address_id"].(string)

		nicId := findNicId(state, sourceNicId)
		if nicId == nil {
			return fmt.Errorf("updating replicated vm %s (vault %s): Trying to update NIC that is not known by Azure %s", name, vaultName, sourceNicId)
		}
		ipConfig := []replicationprotecteditems.IPConfigInputDetails{
			{
				RecoverySubnetName:        &targetSubnetName,
				RecoveryStaticIPAddress:   &targetStaticIp,
				RecoveryPublicIPAddressId: &recoveryPublicIPAddressID,
				TfoStaticIPAddress:        &testStaticIp,
				TfoPublicIPAddressId:      &testPublicIpAddressID,
				TfoSubnetName:             &testSubNetName,
				IsPrimary:                 utils.Bool(true), // currently we can only set one IPconfig for a nic, so we dont need to expose this to users.
			},
		}
		vmNics = append(vmNics, replicationprotecteditems.VMNicInputDetails{
			NicId:     nicId,
			IPConfigs: &ipConfig,
		})
	}

	managedDisks := make([]replicationprotecteditems.A2AVMManagedDiskUpdateDetails, 0)
	for _, raw := range d.Get("managed_disk").(*pluginsdk.Set).List() {
		diskInput := raw.(map[string]interface{})
		diskId := diskInput["disk_id"].(string)
		targetReplicaDiskType := diskInput["target_replica_disk_type"].(string)
		targetDiskType := diskInput["target_disk_type"].(string)

		managedDisks = append(managedDisks, replicationprotecteditems.A2AVMManagedDiskUpdateDetails{
			DiskId:                         &diskId,
			RecoveryReplicaDiskAccountType: &targetReplicaDiskType,
			RecoveryTargetDiskAccountType:  &targetDiskType,
			DiskEncryptionInfo:             expandDiskEncryption(diskInput["target_disk_encryption"].([]interface{})),
		})
	}

	if targetNetworkId == "" {
		// No target network id was specified, so we want to preserve what was selected
		if a2aDetails, isA2a := (*state.Properties.ProviderSpecificDetails).(replicationprotecteditems.A2AReplicationDetails); isA2a {
			if a2aDetails.SelectedRecoveryAzureNetworkId != nil {
				targetNetworkId = *a2aDetails.SelectedRecoveryAzureNetworkId
			} else {
				return fmt.Errorf("target_network_id must be set when a network_interface is configured")
			}
		} else {
			return fmt.Errorf("target_network_id must be set when a network_interface is configured")
		}
	}

	if testNetworkId == "" {
		// No test network id was specified, so we want to preserve what was selected
		if a2aDetails, isA2a := (*state.Properties.ProviderSpecificDetails).(replicationprotecteditems.A2AReplicationDetails); isA2a {
			if a2aDetails.SelectedTfoAzureNetworkId != nil {
				testNetworkId = *a2aDetails.SelectedTfoAzureNetworkId
			}
		}
	}

	var providerInput replicationprotecteditems.UpdateReplicationProtectedItemProviderInput = replicationprotecteditems.A2AUpdateReplicationProtectedItemInput{
		ManagedDiskUpdateDetails:           &managedDisks,
		RecoveryProximityPlacementGroupId:  utils.String(d.Get("target_proximity_placement_group_id").(string)),
		RecoveryBootDiagStorageAccountId:   utils.String(d.Get("target_boot_diagnostic_storage_account_id").(string)),
		RecoveryCapacityReservationGroupId: utils.String(d.Get("target_capacity_reservation_group_id").(string)),
		RecoveryVirtualMachineScaleSetId:   utils.String(d.Get("target_virtual_machine_scale_set_id").(string)),
	}
	parameters := replicationprotecteditems.UpdateReplicationProtectedItemInput{
		Properties: &replicationprotecteditems.UpdateReplicationProtectedItemInputProperties{
			RecoveryAzureVMName:            &name,
			SelectedRecoveryAzureNetworkId: &targetNetworkId,
			SelectedTfoAzureNetworkId:      &testNetworkId,
			VMNics:                         &vmNics,
			RecoveryAvailabilitySetId:      targetAvailabilitySetID,
			ProviderSpecificDetails:        &providerInput,
		},
	}

	err = client.UpdateThenPoll(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("updating replicated vm %s (vault %s): %+v", name, vaultName, err)
	}

	return resourceSiteRecoveryReplicatedItemRead(d, meta)
}

func findNicId(state *replicationprotecteditems.ReplicationProtectedItem, sourceNicId string) *string {
	if a2aDetails, isA2a := (*state.Properties.ProviderSpecificDetails).(replicationprotecteditems.A2AReplicationDetails); isA2a {
		if a2aDetails.VMNics != nil {
			for _, nic := range *a2aDetails.VMNics {
				if nic.SourceNicArmId != nil && *nic.SourceNicArmId == sourceNicId {
					return nic.NicId
				}
			}
		}
	}
	return nil
}

func resourceSiteRecoveryReplicatedItemRead(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on site recovery replicated vm %s: %+v", id.String(), err)
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("making Read request on site recovery replicated vm %s: model is nil", id.String())
	}

	d.Set("name", id.ReplicationProtectedItemName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("recovery_vault_name", id.VaultName)
	d.Set("source_recovery_fabric_name", id.ReplicationFabricName)
	d.Set("source_recovery_protection_container_name", id.ReplicationProtectionContainerName)

	if prop := model.Properties; prop != nil {
		d.Set("target_recovery_fabric_id", prop.RecoveryFabricId)
		d.Set("recovery_replication_policy_id", prop.PolicyId)
		d.Set("target_recovery_protection_container_id", prop.RecoveryContainerId)

		if details := prop.ProviderSpecificDetails; details != nil {
			if a2aDetails, isA2a := (*details).(replicationprotecteditems.A2AReplicationDetails); isA2a {
				d.Set("source_vm_id", a2aDetails.FabricObjectId)
				d.Set("target_resource_group_id", a2aDetails.RecoveryAzureResourceGroupId)
				d.Set("target_availability_set_id", a2aDetails.RecoveryAvailabilitySet)
				d.Set("target_zone", a2aDetails.RecoveryAvailabilityZone)
				d.Set("target_network_id", a2aDetails.SelectedRecoveryAzureNetworkId)
				d.Set("test_network_id", a2aDetails.SelectedTfoAzureNetworkId)
				d.Set("target_proximity_placement_group_id", a2aDetails.RecoveryProximityPlacementGroupId)
				d.Set("target_boot_diagnostic_storage_account_id", a2aDetails.RecoveryBootDiagStorageAccountId)
				d.Set("target_capacity_reservation_group_id", a2aDetails.RecoveryCapacityReservationGroupId)
				d.Set("target_virtual_machine_scale_set_id", a2aDetails.RecoveryVirtualMachineScaleSetId)
				d.Set("target_edge_zone", flattenEdgeZone(a2aDetails.RecoveryExtendedLocation))
				d.Set("multi_vm_group_name", a2aDetails.MultiVMGroupName)

				if a2aDetails.ProtectedDisks != nil {
					disksOutput := make([]interface{}, 0)
					for _, disk := range *a2aDetails.ProtectedDisks {
						disksOutput = append(disksOutput, map[string]interface{}{
							"disk_uri":                   disk.DiskUri,
							"staging_storage_account_id": disk.PrimaryStagingAzureStorageAccountId,
							"target_storage_account_id":  disk.RecoveryAzureStorageAccountId,
						})
					}
					d.Set("unmanaged_disk", disksOutput)
				}

				if a2aDetails.ProtectedManagedDisks != nil {
					disksOutput := make([]interface{}, 0)
					for _, disk := range *a2aDetails.ProtectedManagedDisks {
						diskOutput := make(map[string]interface{})
						diskId := ""
						if disk.DiskId != nil {
							diskId = *disk.DiskId
						}
						diskOutput["disk_id"] = diskId

						primaryStagingAzureStorageAccountID := ""
						if disk.PrimaryStagingAzureStorageAccountId != nil {
							primaryStagingAzureStorageAccountID = *disk.PrimaryStagingAzureStorageAccountId
						}
						diskOutput["staging_storage_account_id"] = primaryStagingAzureStorageAccountID

						recoveryResourceGroupID := ""
						if disk.RecoveryResourceGroupId != nil {
							recoveryResourceGroupID = *disk.RecoveryResourceGroupId
						}
						diskOutput["target_resource_group_id"] = recoveryResourceGroupID

						recoveryReplicaDiskAccountType := ""
						if disk.RecoveryReplicaDiskAccountType != nil {
							recoveryReplicaDiskAccountType = *disk.RecoveryReplicaDiskAccountType
						}
						diskOutput["target_replica_disk_type"] = recoveryReplicaDiskAccountType

						recoveryTargetDiskAccountType := ""
						if disk.RecoveryTargetDiskAccountType != nil {
							recoveryTargetDiskAccountType = *disk.RecoveryTargetDiskAccountType
						}
						diskOutput["target_disk_type"] = recoveryTargetDiskAccountType

						recoveryEncryptionSetId := ""
						if disk.RecoveryDiskEncryptionSetId != nil {
							recoveryEncryptionSetId = *disk.RecoveryDiskEncryptionSetId
						}
						diskOutput["target_disk_encryption_set_id"] = recoveryEncryptionSetId

						diskOutput["target_disk_encryption"] = flattenTargetDiskEncryption(disk)

						disksOutput = append(disksOutput, diskOutput)
					}
					d.Set("managed_disk", pluginsdk.NewSet(resourceSiteRecoveryReplicatedVMDiskHash, disksOutput))
				}

				if a2aDetails.VMNics != nil {
					nicsOutput := make([]interface{}, 0)
					for _, nic := range *a2aDetails.VMNics {
						nicOutput := make(map[string]interface{})
						if nic.SourceNicArmId != nil {
							nicOutput["source_network_interface_id"] = *nic.SourceNicArmId
						}
						if nic.IPConfigs != nil && len(*(nic.IPConfigs)) > 0 {
							ipConfig := (*(nic.IPConfigs))[0]
							if ipConfig.RecoveryStaticIPAddress != nil {
								nicOutput["target_static_ip"] = *ipConfig.RecoveryStaticIPAddress
							}
							if ipConfig.RecoverySubnetName != nil {
								nicOutput["target_subnet_name"] = *ipConfig.RecoverySubnetName
							}
							if ipConfig.RecoveryPublicIPAddressId != nil {
								nicOutput["recovery_public_ip_address_id"] = *ipConfig.RecoveryPublicIPAddressId
							}
							if ipConfig.TfoStaticIPAddress != nil {
								nicOutput["failover_test_static_ip"] = *ipConfig.TfoStaticIPAddress
							}
							if ipConfig.TfoSubnetName != nil {
								nicOutput["failover_test_subnet_name"] = *ipConfig.TfoSubnetName
							}
							if ipConfig.TfoPublicIPAddressId != nil {
								nicOutput["failover_test_public_ip_address_id"] = *ipConfig.TfoPublicIPAddressId
							}
						}
						nicsOutput = append(nicsOutput, nicOutput)
					}
					d.Set("network_interface", pluginsdk.NewSet(pluginsdk.HashResource(networkInterfaceResource()), nicsOutput))
				}
			}
		}
	}

	return nil
}

func resourceSiteRecoveryReplicatedItemDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(d.Id())
	if err != nil {
		return err
	}

	disableProtectionReason := replicationprotecteditems.DisableProtectionReasonNotSpecified
	var providerInput replicationprotecteditems.DisableProtectionProviderSpecificInput = siterecovery.DisableProtectionProviderSpecificInput{
		InstanceType: siterecovery.InstanceTypeDisableProtectionProviderSpecificInput,
	}
	disableProtectionInput := replicationprotecteditems.DisableProtectionInput{
		Properties: replicationprotecteditems.DisableProtectionInputProperties{
			DisableProtectionReason:  &disableProtectionReason,
			ReplicationProviderInput: &providerInput,
		},
	}

	err = client.DeleteThenPoll(ctx, *id, disableProtectionInput)
	if err != nil {
		return fmt.Errorf("deleting site recovery replicated vm %s : %+v", id.String(), err)
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

	return pluginsdk.HashString(buf.String())
}

// the default hash function will not ignore Option + Computed properties, which will cause diff.
func resourceSiteRecoveryReplicatedVMNicHash(v interface{}) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]interface{}); ok {
		if v, ok := m["source_network_interface_id"]; ok {
			buf.WriteString(strings.ToLower(v.(string)))
		}
	}

	return pluginsdk.HashString(buf.String())
}

func waitForReplicationToBeHealthy(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (*replicationprotecteditems.ReplicationProtectedItem, error) {
	log.Printf("Waiting for Site Recover to replicate VM.")
	stateConf := &pluginsdk.StateChangeConf{
		Target:       []string{"Protected"},
		Refresh:      waitForReplicationToBeHealthyRefreshFunc(d, meta),
		PollInterval: 1 * time.Minute,
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return nil, fmt.Errorf("context had no deadline")
	}
	stateConf.Timeout = time.Until(deadline)

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("waiting for site recovery to replicate vm: %+v", err)
	}

	protectedItem, ok := result.(replicationprotecteditems.ReplicationProtectedItem)
	if ok {
		return &protectedItem, nil
	} else {
		return nil, fmt.Errorf("waiting for site recovery return incompatible type")
	}
}

func waitForReplicationToBeHealthyRefreshFunc(d *pluginsdk.ResourceData, meta interface{}) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := replicationprotecteditems.ParseReplicationProtectedItemID(d.Id())
		if err != nil {
			return nil, "", err
		}

		client := meta.(*clients.Client).RecoveryServices.ReplicationProtectedItemsClient

		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, *id)
		if err != nil {
			return nil, "", fmt.Errorf("making Read request on site recovery replicated vm %s : %+v", id.String(), err)
		}

		if resp.Model == nil {
			return nil, "", fmt.Errorf("Missing Model in response when making Read request on site recovery replicated vm %s  %+v", id.String(), err)
		}

		if resp.Model.Properties == nil {
			return nil, "", fmt.Errorf("Missing Properties in response when making Read request on site recovery replicated vm %s  %+v", id.String(), err)
		}

		if resp.Model.Properties.ProviderSpecificDetails == nil {
			return nil, "", fmt.Errorf("Missing Properties.ProviderSpecificDetails in response when making Read request on site recovery replicated vm %s : %+v", id.String(), err)
		}

		// Find first disk that is not fully replicated yet
		if a2aDetails, isA2a := (*resp.Model.Properties.ProviderSpecificDetails).(replicationprotecteditems.A2AReplicationDetails); isA2a {
			if a2aDetails.MonitoringPercentageCompletion != nil {
				log.Printf("Waiting for Site Recover to replicate VM, %d%% complete.", *a2aDetails.MonitoringPercentageCompletion)
			}
			if a2aDetails.VMProtectionState != nil {
				return *resp.Model, *a2aDetails.VMProtectionState, nil
			}
		}

		if resp.Model.Properties.ReplicationHealth == nil {
			return nil, "", fmt.Errorf("missing ReplicationHealth in response when making Read request on site recovery replicated vm %s : %+v", id.String(), err)
		}
		return *resp.Model, *resp.Model.Properties.ReplicationHealth, nil
	}
}

func expandDiskEncryption(diskEncryptionInfoList []interface{}) *replicationprotecteditems.DiskEncryptionInfo {
	if len(diskEncryptionInfoList) == 0 {
		return &replicationprotecteditems.DiskEncryptionInfo{}
	}
	diskEncryptionInfoMap := diskEncryptionInfoList[0].(map[string]interface{})

	dek := diskEncryptionInfoMap["disk_encryption_key"].([]interface{})[0].(map[string]interface{})
	diskEncryptionInfo := &replicationprotecteditems.DiskEncryptionInfo{
		DiskEncryptionKeyInfo: &replicationprotecteditems.DiskEncryptionKeyInfo{
			SecretIdentifier:      utils.String(dek["secret_url"].(string)),
			KeyVaultResourceArmId: utils.String(dek["vault_id"].(string)),
		},
	}

	if keyEncryptionKey := diskEncryptionInfoMap["key_encryption_key"].([]interface{}); len(keyEncryptionKey) > 0 {
		kek := keyEncryptionKey[0].(map[string]interface{})
		diskEncryptionInfo.KeyEncryptionKeyInfo = &replicationprotecteditems.KeyEncryptionKeyInfo{
			KeyIdentifier:         utils.String(kek["key_url"].(string)),
			KeyVaultResourceArmId: utils.String(kek["vault_id"].(string)),
		}
	}

	return diskEncryptionInfo
}

func flattenTargetDiskEncryption(disk replicationprotecteditems.A2AProtectedManagedDiskDetails) []interface{} {
	secretUrl := ""
	dekVaultId := ""
	keyUrl := ""
	kekVaultId := ""

	if disk.SecretIdentifier != nil {
		secretUrl = *disk.SecretIdentifier
	}
	if disk.DekKeyVaultArmId != nil {
		dekVaultId = *disk.DekKeyVaultArmId
	}
	if disk.KeyIdentifier != nil {
		keyUrl = *disk.KeyIdentifier
	}
	if disk.KekKeyVaultArmId != nil {
		kekVaultId = *disk.KekKeyVaultArmId
	}

	if secretUrl == "" && dekVaultId == "" && keyUrl == "" && kekVaultId == "" {
		return []interface{}{}
	}

	diskEncryptionKeys := make([]interface{}, 0)
	if secretUrl != "" || dekVaultId != "" {
		diskEncryptionKeys = append(diskEncryptionKeys, map[string]interface{}{
			"secret_url": secretUrl,
			"vault_id":   dekVaultId,
		})
	}

	keyEncryptionKeys := make([]interface{}, 0)
	if keyUrl != "" || kekVaultId != "" {
		keyEncryptionKeys = append(keyEncryptionKeys, map[string]interface{}{
			"key_url":  keyUrl,
			"vault_id": kekVaultId,
		})
	}

	return []interface{}{
		map[string]interface{}{
			"disk_encryption_key": diskEncryptionKeys,
			"key_encryption_key":  keyEncryptionKeys,
		},
	}
}

func expandEdgeZone(input string) *edgezones.Model {
	normalized := edgezones.Normalize(input)
	if normalized == "" {
		return nil
	}

	return &edgezones.Model{
		Name: normalized,
	}
}

func flattenEdgeZone(input *edgezones.Model) string {
	if input == nil || input.Name == "" {
		return ""
	}
	return edgezones.NormalizeNilable(&input.Name)
}

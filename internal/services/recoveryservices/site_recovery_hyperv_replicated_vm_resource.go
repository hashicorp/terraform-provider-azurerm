package recoveryservices

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotectableitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationprotecteditems"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	// tracdked on https://github.com/Azure/azure-rest-api-specs/issues/22798
	enableRDPNever              = "Never"
	enableRdpOnlyOnTestFailOver = "OnlyOnTestFailover"
	enableRdpAlways             = "Always"
)

type SiteRecoveryHyperVReplicatedVMModel struct {
	Name                            string                                                `tfschema:"name"`
	HyperVSiteId                    string                                                `tfschema:"hyperv_site_id"`
	SourceVMName                    string                                                `tfschema:"source_vm_name"`
	TargetResourceGroupId           string                                                `tfschema:"target_resource_group_id"`
	TargetVMName                    string                                                `tfschema:"target_vm_name"`
	PolicyId                        string                                                `tfschema:"replication_policy_id"`
	OsType                          string                                                `tfschema:"os_type"`
	OSDiskName                      string                                                `tfschema:"os_disk_name"`
	DiskNamesToInclude              []string                                              `tfschema:"disk_to_include"`
	TargetStorageAccountId          string                                                `tfschema:"target_storage_account_id"`
	TargetNetworkId                 string                                                `tfschema:"target_network_id"`
	TargetAvailabilityZone          string                                                `tfschema:"target_availability_zone"`
	NetworkInterface                []SiteRecoveryHyperVReplicatedVMNetworkInterfaceModel `tfschema:"network_interface"`
	UseManagedDiskEnabled           bool                                                  `tfschema:"use_managed_disk_enabled"`
	ManagedDisks                    []SiteRecoveryHyperVReplicatedVMManagedDiskModel      `tfschema:"managed_disk"`
	EnableRdpOnTargetOption         string                                                `tfschema:"enable_rdp_or_ssh_on_target_option"`
	LicenseType                     string                                                `tfschema:"license_type"`
	SqlServerLicenseType            string                                                `tfschema:"sql_server_license_type"`
	TargetAvailabilitySetId         string                                                `tfschema:"target_availability_set_id"`
	TargetManagedDiskTags           map[string]string                                     `tfschema:"target_disk_tags"`
	TargetNicTags                   map[string]string                                     `tfschema:"target_nic_tags"`
	TargetProximityPlacementGroupId string                                                `tfschema:"target_proximity_placement_group_id"`
	TargetVMSize                    string                                                `tfschema:"target_vm_size"`
	TargetVMTags                    map[string]string                                     `tfschema:"target_vm_tags"`
	LogStorageAccountId             string                                                `tfschema:"log_storage_account_id"`
}

type SiteRecoveryHyperVReplicatedVMNetworkInterfaceModel struct {
	NetworkName      string `tfschema:"network_name"`
	TargetSubnetName string `tfschema:"target_subnet_name"`
	TargetStaticIp   string `tfschema:"target_static_ip"`
	IsPrimary        bool   `tfschema:"is_primary"`
	FailoverEnabled  bool   `tfschema:"failover_enabled"`
}

type SiteRecoveryHyperVReplicatedVMManagedDiskModel struct {
	DiskName            string `tfschema:"disk_name"`
	DiskEncryptionSetId string `tfschema:"target_disk_encryption_set_id"`
	DiskType            string `tfschema:"target_disk_type"`
}

type SiteRecoveryHyperVReplicatedVMResource struct{}

var _ sdk.Resource = SiteRecoveryHyperVReplicatedVMResource{}

func (s SiteRecoveryHyperVReplicatedVMResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hyperv_site_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		"source_vm_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_resource_group_id": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     azure.ValidateResourceID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"target_vm_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"target_storage_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"replication_policy_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationpolicies.ValidateReplicationPolicyID,
		},

		"os_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Linux",
				"Windows",
			}, false),
		},

		"os_disk_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"disks_to_include": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"use_managed_disk_enabled", "managed_disk"},
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"use_managed_disk_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
			ForceNew: true,
		},

		"managed_disk": {
			Type:         pluginsdk.TypeSet,
			ConfigMode:   pluginsdk.SchemaConfigModeAttr,
			Optional:     true,
			ForceNew:     true,
			RequiredWith: []string{"use_managed_disk_enabled"},
			Set:          resourceSiteRecoveryReplicatedVMDiskHash,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"disk_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"target_disk_encryption_set_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validate.DiskEncryptionSetID,
					},
					"target_disk_type": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ForceNew: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(replicationprotecteditems.DiskAccountTypePremiumLRS),
							string(replicationprotecteditems.DiskAccountTypeStandardLRS),
							string(replicationprotecteditems.DiskAccountTypeStandardSSDLRS),
						}, false),
					},
				},
			},
		},

		"log_storage_account_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			ValidateFunc:     storageValidate.StorageAccountID,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"enable_rdp_or_ssh_on_target_option": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  enableRDPNever,
			ValidateFunc: validation.StringInSlice([]string{
				enableRDPNever,
				enableRdpOnlyOnTestFailOver,
				enableRdpAlways},
				false),
		},

		"target_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true, // if no target_network_id interface set, the update request will fail.
			ValidateFunc: azure.ValidateResourceID,
		},

		"network_interface": {
			Type:       pluginsdk.TypeSet, // use set to avoid diff caused by different orders.
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Required:   true, // if no network interface set, the update request will fail.
			Elem:       hyperVNetworkInterfaceResource(),
		},

		"target_vm_size": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"target_availability_zone": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ConflictsWith: []string{
				"target_availability_set_id",
			},
		},

		"license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  replicationprotecteditems.LicenseTypeNotSpecified,
			ValidateFunc: validation.StringInSlice([]string{
				string(replicationprotecteditems.LicenseTypeNotSpecified),
				string(replicationprotecteditems.LicenseTypeNoLicenseType),
				string(replicationprotecteditems.LicenseTypeWindowsServer),
			}, false),
		},

		"sql_server_license_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  replicationprotecteditems.SqlServerLicenseTypeNotSpecified,
			ValidateFunc: validation.StringInSlice([]string{
				string(replicationprotecteditems.SqlServerLicenseTypeNotSpecified),
				string(replicationprotecteditems.SqlServerLicenseTypeNoLicenseType),
				string(replicationprotecteditems.SqlServerLicenseTypePAYG),
				string(replicationprotecteditems.SqlServerLicenseTypeAHUB),
			}, false),
		},

		"target_availability_set_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     availabilitysets.ValidateAvailabilitySetID,
			DiffSuppressFunc: suppress.CaseDifference,
			ConflictsWith:    []string{"target_availability_zone"},
		},

		"target_proximity_placement_group_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: proximityplacementgroups.ValidateProximityPlacementGroupID,
		},

		"target_vm_tags": tags.Schema(),

		"target_disk_tags": {
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ValidateFunc: tags.Validate,
			RequiredWith: []string{"use_managed_disk_enabled"},
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"target_network_interface_tags": tags.Schema(),
	}
}

func hyperVNetworkInterfaceResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"network_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"target_static_ip": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"target_subnet_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"is_primary": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"failover_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func (s SiteRecoveryHyperVReplicatedVMResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}

func (s SiteRecoveryHyperVReplicatedVMResource) ModelObject() interface{} {
	return &SiteRecoveryHyperVReplicatedVMModel{}
}

func (s SiteRecoveryHyperVReplicatedVMResource) ResourceType() string {
	return "azurerm_site_recovery_hyperv_replicated_vm"
}

func (s SiteRecoveryHyperVReplicatedVMResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return replicationprotecteditems.ValidateReplicationProtectedItemID
}

func (s SiteRecoveryHyperVReplicatedVMResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var plan SiteRecoveryHyperVReplicatedVMModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

			parsedFabricID, err := replicationfabrics.ParseReplicationFabricID(plan.HyperVSiteId)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", plan.HyperVSiteId, err)
			}
			containerId, err := fetchHyperVContainerIdByFabricId(ctx, metadata.Client.RecoveryServices.ProtectionContainerClient, *parsedFabricID)
			if err != nil {
				return fmt.Errorf("fetching HyperV Container Name by Fabric Id %s: %+v", plan.HyperVSiteId, err)
			}

			parsedContainerId, err := replicationprotecteditems.ParseReplicationProtectionContainerID(containerId)
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", containerId, err)
			}

			id := replicationprotecteditems.NewReplicationProtectedItemID(parsedContainerId.SubscriptionId, parsedContainerId.ResourceGroupName, parsedContainerId.VaultName, parsedContainerId.ReplicationFabricName, parsedContainerId.ReplicationProtectionContainerName, plan.Name)

			protectableItem, err := fetchProtectableItemByVMName(ctx, metadata.Client.RecoveryServices.ReplicationProtectableItemsClient, containerId, plan.SourceVMName)
			if err != nil {
				return fmt.Errorf("fetching Protectable Item by VM Name %s: %+v", plan.SourceVMName, err)
			}

			if protectableItem.Properties == nil {
				return fmt.Errorf("retrieving properties for Protectable Item %q: Properties were nil", plan.SourceVMName)
			}

			customDetail, ok := protectableItem.Properties.CustomDetails.(replicationprotectableitems.HyperVVirtualMachineDetails)
			if !ok {
				return fmt.Errorf("retrieving properties for Protectable Item %q: type mismatch", plan.SourceVMName)
			}

			osVHDId := ""
			diskIdsToInclude := make([]string, 0)
			var diskToIncludeForManagedDisks []replicationprotecteditems.HyperVReplicaAzureDiskInputDetails
			if plan.UseManagedDiskEnabled {
				for _, disk := range plan.ManagedDisks {
					diskId := ""
					for _, d := range *customDetail.DiskDetails {
						if *d.VhdName == disk.DiskName {
							diskId = *d.VhdId
							break
						}
					}
					if diskId == "" {
						return fmt.Errorf("disk %s not found in protectable item", disk.DiskName)
					}
					diskType := replicationprotecteditems.DiskAccountType(disk.DiskType)
					diskToIncludeForManagedDisks = append(diskToIncludeForManagedDisks, replicationprotecteditems.HyperVReplicaAzureDiskInputDetails{
						DiskId:              &diskId,
						DiskEncryptionSetId: &disk.DiskEncryptionSetId,
						DiskType:            &diskType,
					})
				}
			} else {
				for _, disk := range *customDetail.DiskDetails {
					if *disk.VhdName == plan.OSDiskName {
						osVHDId = *disk.VhdId
					}
					if utils.SliceContainsValue(plan.DiskNamesToInclude, *disk.VhdName) {
						diskIdsToInclude = append(diskIdsToInclude, *disk.VhdId)
					}
				}
			}

			licenseType := replicationprotecteditems.LicenseType(plan.LicenseType)
			sqlLicenseType := replicationprotecteditems.SqlServerLicenseType(plan.SqlServerLicenseType)
			input := replicationprotecteditems.EnableProtectionInput{
				Properties: &replicationprotecteditems.EnableProtectionInputProperties{
					PolicyId:          &plan.PolicyId,
					ProtectableItemId: protectableItem.Id,
					ProviderSpecificDetails: &replicationprotecteditems.HyperVReplicaAzureEnableProtectionInput{
						OsType:                        &plan.OsType,
						TargetAzureVMName:             &plan.TargetVMName,
						VhdId:                         &osVHDId,
						DisksToInclude:                &diskIdsToInclude,
						TargetAzureV2ResourceGroupId:  &plan.TargetResourceGroupId,
						TargetAvailabilityZone:        &plan.TargetAvailabilityZone,
						TargetStorageAccountId:        &plan.TargetStorageAccountId,
						TargetAvailabilitySetId:       &plan.TargetAvailabilitySetId,
						UseManagedDisks:               utils.String(strconv.FormatBool(plan.UseManagedDiskEnabled)),
						TargetManagedDiskTags:         &plan.TargetManagedDiskTags,
						TargetVMTags:                  &plan.TargetVMTags,
						TargetVMSize:                  &plan.TargetVMSize,
						DisksToIncludeForManagedDisks: &diskToIncludeForManagedDisks,
						EnableRdpOnTargetOption:       &plan.EnableRdpOnTargetOption,
						LicenseType:                   &licenseType,
						SqlServerLicenseType:          &sqlLicenseType,
						LogStorageAccountId:           &plan.LogStorageAccountId,
					},
				},
			}

			err = client.CreateThenPoll(ctx, id, input)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			err = hyperVReplicatedVMWaitForFullyProtected(ctx, metadata)
			if err != nil {
				return fmt.Errorf("waiting for %s to be fully protected: %+v", id, err)
			}

			return HyperVReplicatedVMUpdateInternal(ctx, metadata)
		},
	}
}

func (s SiteRecoveryHyperVReplicatedVMResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

			id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			hyperVSiteId := replicationfabrics.NewReplicationFabricID(id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.ReplicationFabricName)

			state := SiteRecoveryHyperVReplicatedVMModel{
				Name:         id.ReplicationProtectedItemName,
				HyperVSiteId: hyperVSiteId.ID(),
			}

			if model := resp.Model; model != nil {
				if prop := model.Properties; prop != nil {
					if prop.ProtectableItemId != nil {
						state.SourceVMName, err = fetchVMNameByProtectableItemId(ctx, metadata.Client.RecoveryServices.ReplicationProtectableItemsClient, *prop.ProtectableItemId)
						if err != nil {
							return fmt.Errorf("fetching VM Name by Protectable Item Id %s: %+v", *prop.ProtectableItemId, err)
						}
					}

					if prop.PolicyId != nil {
						state.PolicyId = *prop.PolicyId
					}

					if prop.ProviderSpecificDetails != nil {
						if detail, ok := prop.ProviderSpecificDetails.(replicationprotecteditems.HyperVReplicaAzureReplicationDetails); ok {
							if detail.RecoveryAzureResourceGroupId != nil {
								state.TargetResourceGroupId = *detail.RecoveryAzureResourceGroupId
							}
							if detail.RecoveryAzureVMName != nil {
								state.TargetVMName = *detail.RecoveryAzureVMName
							}
							if detail.SelectedRecoveryAzureNetworkId != nil {
								state.TargetNetworkId = *detail.SelectedRecoveryAzureNetworkId
							}
							if detail.TargetAvailabilityZone != nil {
								state.TargetAvailabilityZone = *detail.TargetAvailabilityZone
							}
							if detail.RecoveryAzureStorageAccount != nil {
								state.TargetStorageAccountId = *detail.RecoveryAzureStorageAccount
							}
							if detail.OSDetails != nil && detail.OSDetails.OsType != nil {
								state.OsType = *detail.OSDetails.OsType
							}
							if detail.AzureVMDiskDetails != nil {
								var diskNames []string
								for _, disk := range *detail.AzureVMDiskDetails {
									if disk.VhdName != nil {
										diskNames = append(diskNames, *disk.VhdName)
									}
									if disk.VhdType != nil && strings.EqualFold(*disk.VhdType, "OperatingSystem") {
										state.OSDiskName = *disk.VhdName
									}
								}

								state.DiskNamesToInclude = diskNames
							}
							if detail.RecoveryAzureStorageAccount != nil {
								state.TargetStorageAccountId = *detail.RecoveryAzureStorageAccount
							}
							if detail.SelectedRecoveryAzureNetworkId != nil {
								state.TargetNetworkId = *detail.SelectedRecoveryAzureNetworkId
							}
							if detail.TargetAvailabilityZone != nil {
								state.TargetAvailabilityZone = *detail.TargetAvailabilityZone
							}
							if detail.VMNics != nil {
								var outputs []SiteRecoveryHyperVReplicatedVMNetworkInterfaceModel
								primaryNicId := ""
								if detail.SelectedSourceNicId != nil {
									primaryNicId = *detail.SelectedSourceNicId
								}
								for _, nic := range *detail.VMNics {
									o := SiteRecoveryHyperVReplicatedVMNetworkInterfaceModel{}
									if nic.VMNetworkName != nil {
										o.NetworkName = *nic.VMNetworkName
									}
									if nic.IPConfigs != nil && len(*nic.IPConfigs) == 1 { // it's only support to set one Ipconfig for now.
										ip := (*nic.IPConfigs)[0]
										if ip.RecoverySubnetName != nil {
											o.TargetSubnetName = *ip.RecoverySubnetName
										}
										if ip.RecoveryStaticIPAddress != nil {
											o.TargetStaticIp = *ip.RecoveryStaticIPAddress
										}
									}
									if nic.SelectionType != nil {
										o.FailoverEnabled = strings.EqualFold(*nic.SelectionType, "NotSelected")
									}
									if nic.NicId != nil && *nic.NicId == primaryNicId {
										o.IsPrimary = true
									}
									outputs = append(outputs, o)
								}
								state.NetworkInterface = outputs
							}
							if detail.UseManagedDisks != nil {
								state.UseManagedDiskEnabled, err = strconv.ParseBool(*detail.UseManagedDisks)
								if err != nil {
									return fmt.Errorf("parsing `use_managed_disk_enabled` %s: %+v", *detail.UseManagedDisks, err)
								}
							}

							if detail.ProtectedManagedDisks != nil {
								diskIdToNameMap, _, err := fetchProtectableDiskNameIdMap(ctx, metadata.Client.RecoveryServices.ReplicationProtectableItemsClient, *prop.ProtectableItemId)
								if err != nil {
									return fmt.Errorf("fetching Disk Name to Id Map by Protectable Item Id %s: %+v", *prop.ProtectableItemId, err)
								}
								var outputs []SiteRecoveryHyperVReplicatedVMManagedDiskModel
								for _, disk := range *detail.ProtectedManagedDisks {
									o := SiteRecoveryHyperVReplicatedVMManagedDiskModel{}
									if disk.DiskEncryptionSetId != nil {
										o.DiskEncryptionSetId = *disk.DiskEncryptionSetId
									}
									if disk.ReplicaDiskType != nil {
										o.DiskType = *disk.ReplicaDiskType
									}
									if v, existing := diskIdToNameMap[*disk.DiskId]; existing {
										o.DiskName = v
									}
								}
								state.ManagedDisks = outputs
							}
							if detail.EnableRdpOnTargetOption != nil {
								state.EnableRdpOnTargetOption = *detail.EnableRdpOnTargetOption
							}
							if detail.LicenseType != nil {
								state.LicenseType = *detail.LicenseType
							}
							if detail.SqlServerLicenseType != nil {
								state.SqlServerLicenseType = *detail.SqlServerLicenseType
							}
							if detail.RecoveryAvailabilitySetId != nil {
								state.TargetAvailabilitySetId = *detail.RecoveryAvailabilitySetId
							}
							if detail.TargetManagedDiskTags != nil {
								state.TargetManagedDiskTags = *detail.TargetManagedDiskTags
							}
							if detail.TargetProximityPlacementGroupId != nil {
								state.TargetProximityPlacementGroupId = *detail.TargetProximityPlacementGroupId
							}
							if detail.TargetVMTags != nil {
								state.TargetVMTags = *detail.TargetVMTags
							}
							if detail.RecoveryAzureLogStorageAccountId != nil {
								state.LogStorageAccountId = *detail.RecoveryAzureLogStorageAccountId
							}
							if detail.RecoveryAzureVMSize != nil {
								state.TargetVMSize = *detail.RecoveryAzureVMSize
							}
						}
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SiteRecoveryHyperVReplicatedVMResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 80 * time.Minute,
		Func:    HyperVReplicatedVMUpdateInternal,
	}
}

func (s SiteRecoveryHyperVReplicatedVMResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 80 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

			id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
			}

			err = client.DeleteThenPoll(ctx, *id, replicationprotecteditems.DisableProtectionInput{})
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func fetchProtectableItemByVMName(ctx context.Context, client *replicationprotectableitems.ReplicationProtectableItemsClient, containerId string, vmName string) (*replicationprotectableitems.ProtectableItem, error) {
	parsedContainerId, err := replicationprotectableitems.ParseReplicationProtectionContainerID(containerId)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %+v", containerId, err)
	}

	protectableItems, err := client.ListByReplicationProtectionContainers(ctx, *parsedContainerId, replicationprotectableitems.DefaultListByReplicationProtectionContainersOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("listing protectable items by container %s: %+v", containerId, err)
	}

	if protectableItems.Model == nil {
		return nil, fmt.Errorf("listing protectable items by container: %s is nil", containerId)
	}

	for _, v := range *protectableItems.Model {
		if v.Properties != nil && v.Properties.FriendlyName != nil && *v.Properties.FriendlyName == vmName {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("protectable item with vm name %s not found", vmName)
}

func fetchProtectableDiskNameIdMap(ctx context.Context, client *replicationprotectableitems.ReplicationProtectableItemsClient, protectableItemId string) (idToName map[string]string, nameToId map[string]string, err error) {
	idToName = make(map[string]string)
	nameToId = make(map[string]string)
	protectableItem, err := fetchProtectableItemById(ctx, client, protectableItemId)
	if err != nil {
		return idToName, nameToId, fmt.Errorf("retrieving %s: %+v", protectableItemId, err)
	}

	if protectableItem.Properties != nil && protectableItem.Properties.CustomDetails != nil {
		if customDetail, ok := protectableItem.Properties.CustomDetails.(replicationprotectableitems.HyperVVirtualMachineDetails); ok {
			if customDetail.DiskDetails != nil {
				for _, disk := range *customDetail.DiskDetails {
					if disk.VhdName != nil && disk.VhdId != nil {
						nameToId[*disk.VhdName] = *disk.VhdId
						idToName[*disk.VhdId] = *disk.VhdName
					}
				}
			}
		}
	}

	return idToName, nameToId, nil
}

func fetchVMNameByProtectableItemId(ctx context.Context, client *replicationprotectableitems.ReplicationProtectableItemsClient, protectableItemId string) (string, error) {
	protectableItem, err := fetchProtectableItemById(ctx, client, protectableItemId)
	if err != nil {
		return "", fmt.Errorf("retrieving %s: %+v", protectableItemId, err)
	}

	if protectableItem != nil && protectableItem.Properties != nil && protectableItem.Properties.FriendlyName != nil {
		return *protectableItem.Properties.FriendlyName, nil
	}

	return "", fmt.Errorf("retrieving %s: properties were nil", protectableItemId)
}

func fetchProtectableItemById(ctx context.Context, client *replicationprotectableitems.ReplicationProtectableItemsClient, protectableItemId string) (*replicationprotectableitems.ProtectableItem, error) {
	parsedProtectableItemId, err := replicationprotectableitems.ParseReplicationProtectableItemID(handleAzureSdkForGoBug2824(protectableItemId))
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %+v", protectableItemId, err)
	}

	protectableItem, err := client.Get(ctx, *parsedProtectableItemId)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", protectableItemId, err)
	}

	if protectableItem.Model != nil {
		return protectableItem.Model, nil
	}

	return nil, fmt.Errorf("retrieving %s: properties were nil", protectableItemId)
}

func HyperVReplicatedVMUpdateInternal(ctx context.Context, metadata sdk.ResourceMetaData) error {
	var plan SiteRecoveryHyperVReplicatedVMModel
	if err := metadata.Decode(&plan); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

	id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", metadata.ResourceData.Id(), err)
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return metadata.MarkAsGone(id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: Model was nil", id)
	}

	if existing.Model.Properties == nil || existing.Model.Properties.ProviderSpecificDetails == nil {
		return fmt.Errorf("retrieving %s: Properties or ProviderSpecificDetails was nil", id)
	}

	detail, ok := existing.Model.Properties.ProviderSpecificDetails.(replicationprotecteditems.HyperVReplicaAzureReplicationDetails)
	if !ok {
		return fmt.Errorf("retrieving %s: ProviderSpecificDetails was not HyperVReplicaAzureReplicationDetails", id)
	}

	if detail.SelectedSourceNicId == nil {
		return fmt.Errorf("retrieving %s: SelectedSourceNicId was nil", id)
	}

	primaryNicId := *detail.SelectedSourceNicId
	vmNics := make([]replicationprotecteditems.VMNicInputDetails, 0)
	for _, nic := range plan.NetworkInterface {
		nicId, err := HyperVReplicatedVMFindNicId(*existing.Model, nic.NetworkName)
		if err != nil {
			return fmt.Errorf("finding nic id: %+v", err)
		}

		if nic.IsPrimary {
			if primaryNicId != "" {
				return fmt.Errorf("only one nic can be primary")
			}
			primaryNicId = nicId
		}

		selectType := "NotSelected"
		if nic.FailoverEnabled {
			selectType = "SelectedByUser"
		}

		vmNics = append(vmNics, replicationprotecteditems.VMNicInputDetails{
			NicId:         &nicId,
			SelectionType: &selectType,
			IPConfigs: &[]replicationprotecteditems.IPConfigInputDetails{
				{
					RecoverySubnetName:      &nic.TargetSubnetName,
					RecoveryStaticIPAddress: &nic.TargetStaticIp,
					// per Portal behaviour, these two property are always set to true.
					// Primary Nic is selected by `selectedSourceNicId` parameter.
					// Whether to create nic for failover is controled by `selectionType`
					IsPrimary:            utils.Bool(true),
					IsSeletedForFailover: utils.Bool(true),
				},
			},
		})
	}

	if existing.Model.Properties == nil || existing.Model.Properties.ProtectableItemId == nil {
		return fmt.Errorf("retrieving %s: properties or protectableItemId was nil", id)
	}
	_, diskNameToIdMap, err := fetchProtectableDiskNameIdMap(ctx, metadata.Client.RecoveryServices.ReplicationProtectableItemsClient, *existing.Model.Properties.ProtectableItemId)
	diskIdToDiskEncryptionMap := make(map[string]string, 0)
	for _, disk := range plan.ManagedDisks {
		vhdId := diskNameToIdMap[disk.DiskName]
		diskIdToDiskEncryptionMap[vhdId] = disk.DiskEncryptionSetId
	}
	licenseType := replicationprotecteditems.LicenseType(plan.LicenseType)
	sqlServerLicenseType := replicationprotecteditems.SqlServerLicenseType(plan.SqlServerLicenseType)
	input := replicationprotecteditems.UpdateReplicationProtectedItemInput{
		Properties: &replicationprotecteditems.UpdateReplicationProtectedItemInputProperties{
			EnableRdpOnTargetOption:        &plan.EnableRdpOnTargetOption,
			LicenseType:                    &licenseType,
			RecoveryAvailabilitySetId:      &plan.TargetAvailabilitySetId,
			RecoveryAzureVMName:            &plan.TargetVMName,
			SelectedRecoveryAzureNetworkId: &plan.TargetNetworkId,
			VMNics:                         &vmNics,
			SelectedSourceNicId:            &primaryNicId,
			ProviderSpecificDetails: replicationprotecteditems.HyperVReplicaAzureUpdateReplicationProtectedItemInput{
				RecoveryAzureV2ResourceGroupId:  &plan.TargetResourceGroupId,
				SqlServerLicenseType:            &sqlServerLicenseType,
				TargetAvailabilityZone:          &plan.TargetAvailabilityZone,
				TargetManagedDiskTags:           &plan.TargetManagedDiskTags,
				TargetNicTags:                   &plan.TargetNicTags,
				TargetProximityPlacementGroupId: &plan.TargetProximityPlacementGroupId,
				TargetVMTags:                    &plan.TargetVMTags,
				DiskIdToDiskEncryptionMap:       &diskIdToDiskEncryptionMap,
			},
		},
	}

	if err := client.UpdateThenPoll(ctx, *id, input); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return nil
}

func HyperVReplicatedVMFindNicId(protectedItem replicationprotecteditems.ReplicationProtectedItem, networkName string) (string, error) {
	if protectedItem.Properties == nil {
		return "", fmt.Errorf("properties were nil")
	}
	if protectedItem.Properties.ProviderSpecificDetails == nil {
		return "", fmt.Errorf("provider specific details were nil")
	}
	detail, ok := protectedItem.Properties.ProviderSpecificDetails.(replicationprotecteditems.HyperVReplicaAzureReplicationDetails)
	if !ok {
		return "", fmt.Errorf("not a HyperVReplicaAzureReplicationDetails")
	}
	if detail.VMNics == nil {
		return "", fmt.Errorf("vm nics were nil")
	}
	for _, nic := range *detail.VMNics {
		if *nic.VMNetworkName == networkName {
			return *nic.NicId, nil
		}
	}
	return "", fmt.Errorf("nic with network name %s not found", networkName)
}

func hyperVReplicatedVMWaitForFullyProtected(ctx context.Context, metadata sdk.ResourceMetaData) error {
	stateConf := &pluginsdk.StateChangeConf{
		Target:       []string{"Protected"},
		Refresh:      hyperVWaitForReplicationToBeHealthyRefreshFunc(ctx, metadata),
		PollInterval: time.Minute,
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}
	stateConf.Timeout = time.Until(deadline)

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func hyperVWaitForReplicationToBeHealthyRefreshFunc(ctx context.Context, metadata sdk.ResourceMetaData) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		id, err := replicationprotecteditems.ParseReplicationProtectedItemID(metadata.ResourceData.Id())
		if err != nil {
			return nil, "", err
		}

		client := metadata.Client.RecoveryServices.ReplicationProtectedItemsClient

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

		if detail, isH2A := resp.Model.Properties.ProviderSpecificDetails.(replicationprotecteditems.HyperVReplicaAzureReplicationDetails); isH2A {
			if detail.VMProtectionState != nil {
				return *resp.Model, *detail.VMProtectionState, nil
			}
		}

		if resp.Model.Properties.ReplicationHealth == nil {
			return nil, "", fmt.Errorf("missing ReplicationHealth in response when making Read request on site recovery replicated vm %s : %+v", id.String(), err)
		}
		return *resp.Model, *resp.Model.Properties.ReplicationHealth, nil
	}
}

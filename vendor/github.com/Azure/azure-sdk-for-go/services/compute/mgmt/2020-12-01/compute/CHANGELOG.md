Generated from https://github.com/Azure/azure-rest-api-specs/tree/e0f8b9ab0f5fe5e71b7429ebfea8a33c19ec9d8d/specification/compute/resource-manager/readme.md tag: `package-2020-12-01`

Code generator @microsoft.azure/autorest.go@2.1.178


## Breaking Changes

## Signature Changes

### Const Types

1. TrustedLaunch changed type from SecurityTypes to DiskSecurityTypes

### New Constants

1. DiskStorageAccountTypes.PremiumZRS
1. DiskStorageAccountTypes.StandardSSDZRS
1. SecurityTypes.SecurityTypesTrustedLaunch
1. StorageAccountTypes.StorageAccountTypesPremiumZRS
1. StorageAccountTypes.StorageAccountTypesStandardSSDZRS

### New Funcs

1. *DedicatedHostsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *DedicatedHostsDeleteFuture.UnmarshalJSON([]byte) error
1. *DedicatedHostsUpdateFuture.UnmarshalJSON([]byte) error
1. *DiskAccessesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *DiskAccessesDeleteAPrivateEndpointConnectionFuture.UnmarshalJSON([]byte) error
1. *DiskAccessesDeleteFuture.UnmarshalJSON([]byte) error
1. *DiskAccessesUpdateAPrivateEndpointConnectionFuture.UnmarshalJSON([]byte) error
1. *DiskAccessesUpdateFuture.UnmarshalJSON([]byte) error
1. *DiskEncryptionSetsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *DiskEncryptionSetsDeleteFuture.UnmarshalJSON([]byte) error
1. *DiskEncryptionSetsUpdateFuture.UnmarshalJSON([]byte) error
1. *DisksCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *DisksDeleteFuture.UnmarshalJSON([]byte) error
1. *DisksGrantAccessFuture.UnmarshalJSON([]byte) error
1. *DisksRevokeAccessFuture.UnmarshalJSON([]byte) error
1. *DisksUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleriesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleriesDeleteFuture.UnmarshalJSON([]byte) error
1. *GalleriesUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryApplicationVersionsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryApplicationVersionsDeleteFuture.UnmarshalJSON([]byte) error
1. *GalleryApplicationVersionsUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryApplicationsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryApplicationsDeleteFuture.UnmarshalJSON([]byte) error
1. *GalleryApplicationsUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryImageVersionsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryImageVersionsDeleteFuture.UnmarshalJSON([]byte) error
1. *GalleryImageVersionsUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryImagesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *GalleryImagesDeleteFuture.UnmarshalJSON([]byte) error
1. *GalleryImagesUpdateFuture.UnmarshalJSON([]byte) error
1. *ImagesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *ImagesDeleteFuture.UnmarshalJSON([]byte) error
1. *ImagesUpdateFuture.UnmarshalJSON([]byte) error
1. *LogAnalyticsExportRequestRateByIntervalFuture.UnmarshalJSON([]byte) error
1. *LogAnalyticsExportThrottledRequestsFuture.UnmarshalJSON([]byte) error
1. *SnapshotsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *SnapshotsDeleteFuture.UnmarshalJSON([]byte) error
1. *SnapshotsGrantAccessFuture.UnmarshalJSON([]byte) error
1. *SnapshotsRevokeAccessFuture.UnmarshalJSON([]byte) error
1. *SnapshotsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineExtensionsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineExtensionsDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineExtensionsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineRunCommandsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineRunCommandsDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineRunCommandsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetExtensionsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetExtensionsDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetExtensionsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetRollingUpgradesCancelFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetRollingUpgradesStartExtensionUpgradeFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetRollingUpgradesStartOSUpgradeFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMExtensionsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMExtensionsDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMExtensionsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMRunCommandsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMRunCommandsDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMRunCommandsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsDeallocateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsPerformMaintenanceFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsPowerOffFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsRedeployFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsReimageAllFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsReimageFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsRestartFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsRunCommandFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsStartFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetVMsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsDeallocateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsDeleteInstancesFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsPerformMaintenanceFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsPowerOffFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsRedeployFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsReimageAllFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsReimageFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsRestartFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsSetOrchestrationServiceStateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsStartFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachineScaleSetsUpdateInstancesFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesAssessPatchesFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesCaptureFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesConvertToManagedDisksFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesCreateOrUpdateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesDeallocateFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesDeleteFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesInstallPatchesFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesPerformMaintenanceFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesPowerOffFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesReapplyFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesRedeployFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesReimageFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesRestartFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesRunCommandFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesStartFuture.UnmarshalJSON([]byte) error
1. *VirtualMachinesUpdateFuture.UnmarshalJSON([]byte) error
1. DiskUpdateProperties.MarshalJSON() ([]byte, error)
1. PossibleDiskSecurityTypesValues() []DiskSecurityTypes
1. PrivateEndpointConnectionProperties.MarshalJSON() ([]byte, error)

## Struct Changes

### New Structs

1. DiskSecurityProfile
1. PropertyUpdatesInProgress

### New Struct Fields

1. DiskEncryptionSetUpdate.Identity
1. DiskEncryptionSetUpdateProperties.RotationToLatestKeyVersionEnabled
1. DiskProperties.PropertyUpdatesInProgress
1. DiskProperties.SecurityProfile
1. DiskProperties.SupportsHibernation
1. DiskRestorePointProperties.SupportsHibernation
1. DiskUpdateProperties.PropertyUpdatesInProgress
1. DiskUpdateProperties.SupportsHibernation
1. EncryptionSetProperties.LastKeyRotationTimestamp
1. EncryptionSetProperties.RotationToLatestKeyVersionEnabled
1. SnapshotProperties.SupportsHibernation
1. SnapshotUpdateProperties.SupportsHibernation

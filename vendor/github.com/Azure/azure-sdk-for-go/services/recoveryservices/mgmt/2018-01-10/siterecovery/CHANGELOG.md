Generated from https://github.com/Azure/azure-rest-api-specs/tree/3c764635e7d442b3e74caf593029fcd440b3ef82

Code generator @microsoft.azure/autorest.go@~2.1.161

## Breaking Changes

- Function `NewMigrationItemCollectionPage` parameter(s) have been changed from `(func(context.Context, MigrationItemCollection) (MigrationItemCollection, error))` to `(MigrationItemCollection, func(context.Context, MigrationItemCollection) (MigrationItemCollection, error))`
- Function `NewTargetComputeSizeCollectionPage` parameter(s) have been changed from `(func(context.Context, TargetComputeSizeCollection) (TargetComputeSizeCollection, error))` to `(TargetComputeSizeCollection, func(context.Context, TargetComputeSizeCollection) (TargetComputeSizeCollection, error))`
- Function `NewStorageClassificationCollectionPage` parameter(s) have been changed from `(func(context.Context, StorageClassificationCollection) (StorageClassificationCollection, error))` to `(StorageClassificationCollection, func(context.Context, StorageClassificationCollection) (StorageClassificationCollection, error))`
- Function `NewLogicalNetworkCollectionPage` parameter(s) have been changed from `(func(context.Context, LogicalNetworkCollection) (LogicalNetworkCollection, error))` to `(LogicalNetworkCollection, func(context.Context, LogicalNetworkCollection) (LogicalNetworkCollection, error))`
- Function `NewNetworkCollectionPage` parameter(s) have been changed from `(func(context.Context, NetworkCollection) (NetworkCollection, error))` to `(NetworkCollection, func(context.Context, NetworkCollection) (NetworkCollection, error))`
- Function `NewNetworkMappingCollectionPage` parameter(s) have been changed from `(func(context.Context, NetworkMappingCollection) (NetworkMappingCollection, error))` to `(NetworkMappingCollection, func(context.Context, NetworkMappingCollection) (NetworkMappingCollection, error))`
- Function `NewRecoveryPointCollectionPage` parameter(s) have been changed from `(func(context.Context, RecoveryPointCollection) (RecoveryPointCollection, error))` to `(RecoveryPointCollection, func(context.Context, RecoveryPointCollection) (RecoveryPointCollection, error))`
- Function `NewProtectableItemCollectionPage` parameter(s) have been changed from `(func(context.Context, ProtectableItemCollection) (ProtectableItemCollection, error))` to `(ProtectableItemCollection, func(context.Context, ProtectableItemCollection) (ProtectableItemCollection, error))`
- Function `NewProtectionContainerCollectionPage` parameter(s) have been changed from `(func(context.Context, ProtectionContainerCollection) (ProtectionContainerCollection, error))` to `(ProtectionContainerCollection, func(context.Context, ProtectionContainerCollection) (ProtectionContainerCollection, error))`
- Function `NewJobCollectionPage` parameter(s) have been changed from `(func(context.Context, JobCollection) (JobCollection, error))` to `(JobCollection, func(context.Context, JobCollection) (JobCollection, error))`
- Function `NewFabricCollectionPage` parameter(s) have been changed from `(func(context.Context, FabricCollection) (FabricCollection, error))` to `(FabricCollection, func(context.Context, FabricCollection) (FabricCollection, error))`
- Function `NewAlertCollectionPage` parameter(s) have been changed from `(func(context.Context, AlertCollection) (AlertCollection, error))` to `(AlertCollection, func(context.Context, AlertCollection) (AlertCollection, error))`
- Function `NewRecoveryPlanCollectionPage` parameter(s) have been changed from `(func(context.Context, RecoveryPlanCollection) (RecoveryPlanCollection, error))` to `(RecoveryPlanCollection, func(context.Context, RecoveryPlanCollection) (RecoveryPlanCollection, error))`
- Function `NewOperationsDiscoveryCollectionPage` parameter(s) have been changed from `(func(context.Context, OperationsDiscoveryCollection) (OperationsDiscoveryCollection, error))` to `(OperationsDiscoveryCollection, func(context.Context, OperationsDiscoveryCollection) (OperationsDiscoveryCollection, error))`
- Function `NewStorageClassificationMappingCollectionPage` parameter(s) have been changed from `(func(context.Context, StorageClassificationMappingCollection) (StorageClassificationMappingCollection, error))` to `(StorageClassificationMappingCollection, func(context.Context, StorageClassificationMappingCollection) (StorageClassificationMappingCollection, error))`
- Function `NewReplicationProtectedItemCollectionPage` parameter(s) have been changed from `(func(context.Context, ReplicationProtectedItemCollection) (ReplicationProtectedItemCollection, error))` to `(ReplicationProtectedItemCollection, func(context.Context, ReplicationProtectedItemCollection) (ReplicationProtectedItemCollection, error))`
- Function `NewVCenterCollectionPage` parameter(s) have been changed from `(func(context.Context, VCenterCollection) (VCenterCollection, error))` to `(VCenterCollection, func(context.Context, VCenterCollection) (VCenterCollection, error))`
- Function `NewEventCollectionPage` parameter(s) have been changed from `(func(context.Context, EventCollection) (EventCollection, error))` to `(EventCollection, func(context.Context, EventCollection) (EventCollection, error))`
- Function `NewMigrationRecoveryPointCollectionPage` parameter(s) have been changed from `(func(context.Context, MigrationRecoveryPointCollection) (MigrationRecoveryPointCollection, error))` to `(MigrationRecoveryPointCollection, func(context.Context, MigrationRecoveryPointCollection) (MigrationRecoveryPointCollection, error))`
- Function `NewProtectionContainerMappingCollectionPage` parameter(s) have been changed from `(func(context.Context, ProtectionContainerMappingCollection) (ProtectionContainerMappingCollection, error))` to `(ProtectionContainerMappingCollection, func(context.Context, ProtectionContainerMappingCollection) (ProtectionContainerMappingCollection, error))`
- Function `NewRecoveryServicesProviderCollectionPage` parameter(s) have been changed from `(func(context.Context, RecoveryServicesProviderCollection) (RecoveryServicesProviderCollection, error))` to `(RecoveryServicesProviderCollection, func(context.Context, RecoveryServicesProviderCollection) (RecoveryServicesProviderCollection, error))`
- Function `NewPolicyCollectionPage` parameter(s) have been changed from `(func(context.Context, PolicyCollection) (PolicyCollection, error))` to `(PolicyCollection, func(context.Context, PolicyCollection) (PolicyCollection, error))`

## New Content

- New field `TargetGeneration` in struct `VMwareCbtMigrationDetails`
- New field `FirmwareType` in struct `VMwareCbtMigrationDetails`
- New field `TargetAvailabilityZone` in struct `VMwareCbtMigrationDetails`
- New field `TargetAvailabilityZone` in struct `VMwareCbtEnableMigrationInput`
- New field `TargetAvailabilityZone` in struct `VMwareCbtUpdateMigrationItemInput`

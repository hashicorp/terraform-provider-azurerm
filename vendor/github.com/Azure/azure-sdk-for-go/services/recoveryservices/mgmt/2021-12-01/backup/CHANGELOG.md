# Change History

## Breaking Changes

### Removed Constants

1. ContainerTypeBasicProtectionContainer.ContainerTypeBasicProtectionContainerContainerTypeIaaSVMContainer

## Additive Changes

### New Constants

1. ContainerTypeBasicProtectionContainer.ContainerTypeBasicProtectionContainerContainerTypeIaasVMContainer
1. TieringMode.TieringModeDoNotTier
1. TieringMode.TieringModeInvalid
1. TieringMode.TieringModeTierAfter
1. TieringMode.TieringModeTierRecommended

### New Funcs

1. PossibleTieringModeValues() []TieringMode
1. SubProtectionPolicy.MarshalJSON() ([]byte, error)

### Struct Changes

#### New Structs

1. TieringPolicy

#### New Struct Fields

1. AzureIaaSVMProtectedItemExtendedInfo.NewestRecoveryPointInArchive
1. AzureIaaSVMProtectedItemExtendedInfo.OldestRecoveryPointInArchive
1. AzureIaaSVMProtectedItemExtendedInfo.OldestRecoveryPointInVault
1. AzureIaaSVMProtectionPolicy.TieringPolicy
1. AzureVMWorkloadProtectedItemExtendedInfo.NewestRecoveryPointInArchive
1. AzureVMWorkloadProtectedItemExtendedInfo.OldestRecoveryPointInArchive
1. AzureVMWorkloadProtectedItemExtendedInfo.OldestRecoveryPointInVault
1. SubProtectionPolicy.TieringPolicy

# Change History

## Breaking Changes

### Signature Changes

#### Funcs

1. GalleriesClient.Get
	- Params
		- From: context.Context, string, string, SelectPermissions
		- To: context.Context, string, string, SelectPermissions, GalleryExpandParams
1. GalleriesClient.GetPreparer
	- Params
		- From: context.Context, string, string, SelectPermissions
		- To: context.Context, string, string, SelectPermissions, GalleryExpandParams

## Additive Changes

### New Constants

1. Architecture.ArchitectureArm64
1. Architecture.ArchitectureX64
1. ArchitectureTypes.ArchitectureTypesArm64
1. ArchitectureTypes.ArchitectureTypesX64
1. ConfidentialVMEncryptionType.ConfidentialVMEncryptionTypeEncryptedVMGuestStateOnlyWithPmk
1. ConfidentialVMEncryptionType.ConfidentialVMEncryptionTypeEncryptedWithCmk
1. ConfidentialVMEncryptionType.ConfidentialVMEncryptionTypeEncryptedWithPmk
1. GalleryExpandParams.GalleryExpandParamsSharingProfileGroups
1. GalleryExtendedLocationType.GalleryExtendedLocationTypeEdgeZone
1. GalleryExtendedLocationType.GalleryExtendedLocationTypeUnknown
1. SharingProfileGroupTypes.SharingProfileGroupTypesCommunity
1. SharingState.SharingStateFailed
1. SharingState.SharingStateInProgress
1. SharingState.SharingStateSucceeded
1. SharingState.SharingStateUnknown
1. SharingUpdateOperationTypes.SharingUpdateOperationTypesEnableCommunity

### New Funcs

1. CommunityGalleryInfo.MarshalJSON() ([]byte, error)
1. PossibleArchitectureTypesValues() []ArchitectureTypes
1. PossibleArchitectureValues() []Architecture
1. PossibleConfidentialVMEncryptionTypeValues() []ConfidentialVMEncryptionType
1. PossibleGalleryExpandParamsValues() []GalleryExpandParams
1. PossibleGalleryExtendedLocationTypeValues() []GalleryExtendedLocationType
1. PossibleSharingStateValues() []SharingState

### Struct Changes

#### New Structs

1. CommunityGalleryInfo
1. GalleryExtendedLocation
1. GalleryTargetExtendedLocation
1. OSDiskImageSecurityProfile
1. RegionalSharingStatus
1. SharingStatus

#### New Struct Fields

1. GalleryApplicationVersionPublishingProfile.TargetExtendedLocations
1. GalleryArtifactPublishingProfileBase.TargetExtendedLocations
1. GalleryImageProperties.Architecture
1. GalleryImageVersionPublishingProfile.TargetExtendedLocations
1. GalleryProperties.SharingStatus
1. OSDiskImageEncryption.SecurityProfile
1. SharingProfile.CommunityGalleryInfo
1. VirtualMachineImageProperties.Architecture

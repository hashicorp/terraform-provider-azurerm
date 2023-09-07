package disks

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessLevel string

const (
	AccessLevelNone  AccessLevel = "None"
	AccessLevelRead  AccessLevel = "Read"
	AccessLevelWrite AccessLevel = "Write"
)

func PossibleValuesForAccessLevel() []string {
	return []string{
		string(AccessLevelNone),
		string(AccessLevelRead),
		string(AccessLevelWrite),
	}
}

func (s *AccessLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAccessLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAccessLevel(input string) (*AccessLevel, error) {
	vals := map[string]AccessLevel{
		"none":  AccessLevelNone,
		"read":  AccessLevelRead,
		"write": AccessLevelWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AccessLevel(input)
	return &out, nil
}

type Architecture string

const (
	ArchitectureArmSixFour Architecture = "Arm64"
	ArchitectureXSixFour   Architecture = "x64"
)

func PossibleValuesForArchitecture() []string {
	return []string{
		string(ArchitectureArmSixFour),
		string(ArchitectureXSixFour),
	}
}

func (s *Architecture) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseArchitecture(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseArchitecture(input string) (*Architecture, error) {
	vals := map[string]Architecture{
		"arm64": ArchitectureArmSixFour,
		"x64":   ArchitectureXSixFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Architecture(input)
	return &out, nil
}

type DataAccessAuthMode string

const (
	DataAccessAuthModeAzureActiveDirectory DataAccessAuthMode = "AzureActiveDirectory"
	DataAccessAuthModeNone                 DataAccessAuthMode = "None"
)

func PossibleValuesForDataAccessAuthMode() []string {
	return []string{
		string(DataAccessAuthModeAzureActiveDirectory),
		string(DataAccessAuthModeNone),
	}
}

func (s *DataAccessAuthMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataAccessAuthMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataAccessAuthMode(input string) (*DataAccessAuthMode, error) {
	vals := map[string]DataAccessAuthMode{
		"azureactivedirectory": DataAccessAuthModeAzureActiveDirectory,
		"none":                 DataAccessAuthModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataAccessAuthMode(input)
	return &out, nil
}

type DiskCreateOption string

const (
	DiskCreateOptionAttach               DiskCreateOption = "Attach"
	DiskCreateOptionCopy                 DiskCreateOption = "Copy"
	DiskCreateOptionCopyStart            DiskCreateOption = "CopyStart"
	DiskCreateOptionEmpty                DiskCreateOption = "Empty"
	DiskCreateOptionFromImage            DiskCreateOption = "FromImage"
	DiskCreateOptionImport               DiskCreateOption = "Import"
	DiskCreateOptionImportSecure         DiskCreateOption = "ImportSecure"
	DiskCreateOptionRestore              DiskCreateOption = "Restore"
	DiskCreateOptionUpload               DiskCreateOption = "Upload"
	DiskCreateOptionUploadPreparedSecure DiskCreateOption = "UploadPreparedSecure"
)

func PossibleValuesForDiskCreateOption() []string {
	return []string{
		string(DiskCreateOptionAttach),
		string(DiskCreateOptionCopy),
		string(DiskCreateOptionCopyStart),
		string(DiskCreateOptionEmpty),
		string(DiskCreateOptionFromImage),
		string(DiskCreateOptionImport),
		string(DiskCreateOptionImportSecure),
		string(DiskCreateOptionRestore),
		string(DiskCreateOptionUpload),
		string(DiskCreateOptionUploadPreparedSecure),
	}
}

func (s *DiskCreateOption) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskCreateOption(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskCreateOption(input string) (*DiskCreateOption, error) {
	vals := map[string]DiskCreateOption{
		"attach":               DiskCreateOptionAttach,
		"copy":                 DiskCreateOptionCopy,
		"copystart":            DiskCreateOptionCopyStart,
		"empty":                DiskCreateOptionEmpty,
		"fromimage":            DiskCreateOptionFromImage,
		"import":               DiskCreateOptionImport,
		"importsecure":         DiskCreateOptionImportSecure,
		"restore":              DiskCreateOptionRestore,
		"upload":               DiskCreateOptionUpload,
		"uploadpreparedsecure": DiskCreateOptionUploadPreparedSecure,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskCreateOption(input)
	return &out, nil
}

type DiskSecurityTypes string

const (
	DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey             DiskSecurityTypes = "ConfidentialVM_DiskEncryptedWithCustomerKey"
	DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey             DiskSecurityTypes = "ConfidentialVM_DiskEncryptedWithPlatformKey"
	DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey DiskSecurityTypes = "ConfidentialVM_VMGuestStateOnlyEncryptedWithPlatformKey"
	DiskSecurityTypesTrustedLaunch                                          DiskSecurityTypes = "TrustedLaunch"
)

func PossibleValuesForDiskSecurityTypes() []string {
	return []string{
		string(DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey),
		string(DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey),
		string(DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey),
		string(DiskSecurityTypesTrustedLaunch),
	}
}

func (s *DiskSecurityTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskSecurityTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskSecurityTypes(input string) (*DiskSecurityTypes, error) {
	vals := map[string]DiskSecurityTypes{
		"confidentialvm_diskencryptedwithcustomerkey":             DiskSecurityTypesConfidentialVMDiskEncryptedWithCustomerKey,
		"confidentialvm_diskencryptedwithplatformkey":             DiskSecurityTypesConfidentialVMDiskEncryptedWithPlatformKey,
		"confidentialvm_vmgueststateonlyencryptedwithplatformkey": DiskSecurityTypesConfidentialVMVMGuestStateOnlyEncryptedWithPlatformKey,
		"trustedlaunch": DiskSecurityTypesTrustedLaunch,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskSecurityTypes(input)
	return &out, nil
}

type DiskState string

const (
	DiskStateActiveSAS       DiskState = "ActiveSAS"
	DiskStateActiveSASFrozen DiskState = "ActiveSASFrozen"
	DiskStateActiveUpload    DiskState = "ActiveUpload"
	DiskStateAttached        DiskState = "Attached"
	DiskStateFrozen          DiskState = "Frozen"
	DiskStateReadyToUpload   DiskState = "ReadyToUpload"
	DiskStateReserved        DiskState = "Reserved"
	DiskStateUnattached      DiskState = "Unattached"
)

func PossibleValuesForDiskState() []string {
	return []string{
		string(DiskStateActiveSAS),
		string(DiskStateActiveSASFrozen),
		string(DiskStateActiveUpload),
		string(DiskStateAttached),
		string(DiskStateFrozen),
		string(DiskStateReadyToUpload),
		string(DiskStateReserved),
		string(DiskStateUnattached),
	}
}

func (s *DiskState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskState(input string) (*DiskState, error) {
	vals := map[string]DiskState{
		"activesas":       DiskStateActiveSAS,
		"activesasfrozen": DiskStateActiveSASFrozen,
		"activeupload":    DiskStateActiveUpload,
		"attached":        DiskStateAttached,
		"frozen":          DiskStateFrozen,
		"readytoupload":   DiskStateReadyToUpload,
		"reserved":        DiskStateReserved,
		"unattached":      DiskStateUnattached,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskState(input)
	return &out, nil
}

type DiskStorageAccountTypes string

const (
	DiskStorageAccountTypesPremiumLRS     DiskStorageAccountTypes = "Premium_LRS"
	DiskStorageAccountTypesPremiumVTwoLRS DiskStorageAccountTypes = "PremiumV2_LRS"
	DiskStorageAccountTypesPremiumZRS     DiskStorageAccountTypes = "Premium_ZRS"
	DiskStorageAccountTypesStandardLRS    DiskStorageAccountTypes = "Standard_LRS"
	DiskStorageAccountTypesStandardSSDLRS DiskStorageAccountTypes = "StandardSSD_LRS"
	DiskStorageAccountTypesStandardSSDZRS DiskStorageAccountTypes = "StandardSSD_ZRS"
	DiskStorageAccountTypesUltraSSDLRS    DiskStorageAccountTypes = "UltraSSD_LRS"
)

func PossibleValuesForDiskStorageAccountTypes() []string {
	return []string{
		string(DiskStorageAccountTypesPremiumLRS),
		string(DiskStorageAccountTypesPremiumVTwoLRS),
		string(DiskStorageAccountTypesPremiumZRS),
		string(DiskStorageAccountTypesStandardLRS),
		string(DiskStorageAccountTypesStandardSSDLRS),
		string(DiskStorageAccountTypesStandardSSDZRS),
		string(DiskStorageAccountTypesUltraSSDLRS),
	}
}

func (s *DiskStorageAccountTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskStorageAccountTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskStorageAccountTypes(input string) (*DiskStorageAccountTypes, error) {
	vals := map[string]DiskStorageAccountTypes{
		"premium_lrs":     DiskStorageAccountTypesPremiumLRS,
		"premiumv2_lrs":   DiskStorageAccountTypesPremiumVTwoLRS,
		"premium_zrs":     DiskStorageAccountTypesPremiumZRS,
		"standard_lrs":    DiskStorageAccountTypesStandardLRS,
		"standardssd_lrs": DiskStorageAccountTypesStandardSSDLRS,
		"standardssd_zrs": DiskStorageAccountTypesStandardSSDZRS,
		"ultrassd_lrs":    DiskStorageAccountTypesUltraSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskStorageAccountTypes(input)
	return &out, nil
}

type EncryptionType string

const (
	EncryptionTypeEncryptionAtRestWithCustomerKey             EncryptionType = "EncryptionAtRestWithCustomerKey"
	EncryptionTypeEncryptionAtRestWithPlatformAndCustomerKeys EncryptionType = "EncryptionAtRestWithPlatformAndCustomerKeys"
	EncryptionTypeEncryptionAtRestWithPlatformKey             EncryptionType = "EncryptionAtRestWithPlatformKey"
)

func PossibleValuesForEncryptionType() []string {
	return []string{
		string(EncryptionTypeEncryptionAtRestWithCustomerKey),
		string(EncryptionTypeEncryptionAtRestWithPlatformAndCustomerKeys),
		string(EncryptionTypeEncryptionAtRestWithPlatformKey),
	}
}

func (s *EncryptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEncryptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEncryptionType(input string) (*EncryptionType, error) {
	vals := map[string]EncryptionType{
		"encryptionatrestwithcustomerkey":             EncryptionTypeEncryptionAtRestWithCustomerKey,
		"encryptionatrestwithplatformandcustomerkeys": EncryptionTypeEncryptionAtRestWithPlatformAndCustomerKeys,
		"encryptionatrestwithplatformkey":             EncryptionTypeEncryptionAtRestWithPlatformKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EncryptionType(input)
	return &out, nil
}

type HyperVGeneration string

const (
	HyperVGenerationVOne HyperVGeneration = "V1"
	HyperVGenerationVTwo HyperVGeneration = "V2"
)

func PossibleValuesForHyperVGeneration() []string {
	return []string{
		string(HyperVGenerationVOne),
		string(HyperVGenerationVTwo),
	}
}

func (s *HyperVGeneration) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHyperVGeneration(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHyperVGeneration(input string) (*HyperVGeneration, error) {
	vals := map[string]HyperVGeneration{
		"v1": HyperVGenerationVOne,
		"v2": HyperVGenerationVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HyperVGeneration(input)
	return &out, nil
}

type NetworkAccessPolicy string

const (
	NetworkAccessPolicyAllowAll     NetworkAccessPolicy = "AllowAll"
	NetworkAccessPolicyAllowPrivate NetworkAccessPolicy = "AllowPrivate"
	NetworkAccessPolicyDenyAll      NetworkAccessPolicy = "DenyAll"
)

func PossibleValuesForNetworkAccessPolicy() []string {
	return []string{
		string(NetworkAccessPolicyAllowAll),
		string(NetworkAccessPolicyAllowPrivate),
		string(NetworkAccessPolicyDenyAll),
	}
}

func (s *NetworkAccessPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkAccessPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkAccessPolicy(input string) (*NetworkAccessPolicy, error) {
	vals := map[string]NetworkAccessPolicy{
		"allowall":     NetworkAccessPolicyAllowAll,
		"allowprivate": NetworkAccessPolicyAllowPrivate,
		"denyall":      NetworkAccessPolicyDenyAll,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkAccessPolicy(input)
	return &out, nil
}

type OperatingSystemTypes string

const (
	OperatingSystemTypesLinux   OperatingSystemTypes = "Linux"
	OperatingSystemTypesWindows OperatingSystemTypes = "Windows"
)

func PossibleValuesForOperatingSystemTypes() []string {
	return []string{
		string(OperatingSystemTypesLinux),
		string(OperatingSystemTypesWindows),
	}
}

func (s *OperatingSystemTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatingSystemTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatingSystemTypes(input string) (*OperatingSystemTypes, error) {
	vals := map[string]OperatingSystemTypes{
		"linux":   OperatingSystemTypesLinux,
		"windows": OperatingSystemTypesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemTypes(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled  PublicNetworkAccess = "Enabled"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
	}
}

func (s *PublicNetworkAccess) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicNetworkAccess(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicNetworkAccess(input string) (*PublicNetworkAccess, error) {
	vals := map[string]PublicNetworkAccess{
		"disabled": PublicNetworkAccessDisabled,
		"enabled":  PublicNetworkAccessEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

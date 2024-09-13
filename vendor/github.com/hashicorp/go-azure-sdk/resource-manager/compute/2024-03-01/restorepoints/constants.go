package restorepoints

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CachingTypes string

const (
	CachingTypesNone      CachingTypes = "None"
	CachingTypesReadOnly  CachingTypes = "ReadOnly"
	CachingTypesReadWrite CachingTypes = "ReadWrite"
)

func PossibleValuesForCachingTypes() []string {
	return []string{
		string(CachingTypesNone),
		string(CachingTypesReadOnly),
		string(CachingTypesReadWrite),
	}
}

func (s *CachingTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCachingTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCachingTypes(input string) (*CachingTypes, error) {
	vals := map[string]CachingTypes{
		"none":      CachingTypesNone,
		"readonly":  CachingTypesReadOnly,
		"readwrite": CachingTypesReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CachingTypes(input)
	return &out, nil
}

type ComponentNames string

const (
	ComponentNamesMicrosoftNegativeWindowsNegativeShellNegativeSetup ComponentNames = "Microsoft-Windows-Shell-Setup"
)

func PossibleValuesForComponentNames() []string {
	return []string{
		string(ComponentNamesMicrosoftNegativeWindowsNegativeShellNegativeSetup),
	}
}

func (s *ComponentNames) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComponentNames(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComponentNames(input string) (*ComponentNames, error) {
	vals := map[string]ComponentNames{
		"microsoft-windows-shell-setup": ComponentNamesMicrosoftNegativeWindowsNegativeShellNegativeSetup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComponentNames(input)
	return &out, nil
}

type ConsistencyModeTypes string

const (
	ConsistencyModeTypesApplicationConsistent ConsistencyModeTypes = "ApplicationConsistent"
	ConsistencyModeTypesCrashConsistent       ConsistencyModeTypes = "CrashConsistent"
	ConsistencyModeTypesFileSystemConsistent  ConsistencyModeTypes = "FileSystemConsistent"
)

func PossibleValuesForConsistencyModeTypes() []string {
	return []string{
		string(ConsistencyModeTypesApplicationConsistent),
		string(ConsistencyModeTypesCrashConsistent),
		string(ConsistencyModeTypesFileSystemConsistent),
	}
}

func (s *ConsistencyModeTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConsistencyModeTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConsistencyModeTypes(input string) (*ConsistencyModeTypes, error) {
	vals := map[string]ConsistencyModeTypes{
		"applicationconsistent": ConsistencyModeTypesApplicationConsistent,
		"crashconsistent":       ConsistencyModeTypesCrashConsistent,
		"filesystemconsistent":  ConsistencyModeTypesFileSystemConsistent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConsistencyModeTypes(input)
	return &out, nil
}

type DiskControllerTypes string

const (
	DiskControllerTypesNVMe DiskControllerTypes = "NVMe"
	DiskControllerTypesSCSI DiskControllerTypes = "SCSI"
)

func PossibleValuesForDiskControllerTypes() []string {
	return []string{
		string(DiskControllerTypesNVMe),
		string(DiskControllerTypesSCSI),
	}
}

func (s *DiskControllerTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskControllerTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskControllerTypes(input string) (*DiskControllerTypes, error) {
	vals := map[string]DiskControllerTypes{
		"nvme": DiskControllerTypesNVMe,
		"scsi": DiskControllerTypesSCSI,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskControllerTypes(input)
	return &out, nil
}

type HyperVGenerationTypes string

const (
	HyperVGenerationTypesVOne HyperVGenerationTypes = "V1"
	HyperVGenerationTypesVTwo HyperVGenerationTypes = "V2"
)

func PossibleValuesForHyperVGenerationTypes() []string {
	return []string{
		string(HyperVGenerationTypesVOne),
		string(HyperVGenerationTypesVTwo),
	}
}

func (s *HyperVGenerationTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHyperVGenerationTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHyperVGenerationTypes(input string) (*HyperVGenerationTypes, error) {
	vals := map[string]HyperVGenerationTypes{
		"v1": HyperVGenerationTypesVOne,
		"v2": HyperVGenerationTypesVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HyperVGenerationTypes(input)
	return &out, nil
}

type LinuxPatchAssessmentMode string

const (
	LinuxPatchAssessmentModeAutomaticByPlatform LinuxPatchAssessmentMode = "AutomaticByPlatform"
	LinuxPatchAssessmentModeImageDefault        LinuxPatchAssessmentMode = "ImageDefault"
)

func PossibleValuesForLinuxPatchAssessmentMode() []string {
	return []string{
		string(LinuxPatchAssessmentModeAutomaticByPlatform),
		string(LinuxPatchAssessmentModeImageDefault),
	}
}

func (s *LinuxPatchAssessmentMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLinuxPatchAssessmentMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLinuxPatchAssessmentMode(input string) (*LinuxPatchAssessmentMode, error) {
	vals := map[string]LinuxPatchAssessmentMode{
		"automaticbyplatform": LinuxPatchAssessmentModeAutomaticByPlatform,
		"imagedefault":        LinuxPatchAssessmentModeImageDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinuxPatchAssessmentMode(input)
	return &out, nil
}

type LinuxVMGuestPatchAutomaticByPlatformRebootSetting string

const (
	LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways     LinuxVMGuestPatchAutomaticByPlatformRebootSetting = "Always"
	LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired LinuxVMGuestPatchAutomaticByPlatformRebootSetting = "IfRequired"
	LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever      LinuxVMGuestPatchAutomaticByPlatformRebootSetting = "Never"
	LinuxVMGuestPatchAutomaticByPlatformRebootSettingUnknown    LinuxVMGuestPatchAutomaticByPlatformRebootSetting = "Unknown"
)

func PossibleValuesForLinuxVMGuestPatchAutomaticByPlatformRebootSetting() []string {
	return []string{
		string(LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways),
		string(LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
		string(LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever),
		string(LinuxVMGuestPatchAutomaticByPlatformRebootSettingUnknown),
	}
}

func (s *LinuxVMGuestPatchAutomaticByPlatformRebootSetting) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLinuxVMGuestPatchAutomaticByPlatformRebootSetting(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLinuxVMGuestPatchAutomaticByPlatformRebootSetting(input string) (*LinuxVMGuestPatchAutomaticByPlatformRebootSetting, error) {
	vals := map[string]LinuxVMGuestPatchAutomaticByPlatformRebootSetting{
		"always":     LinuxVMGuestPatchAutomaticByPlatformRebootSettingAlways,
		"ifrequired": LinuxVMGuestPatchAutomaticByPlatformRebootSettingIfRequired,
		"never":      LinuxVMGuestPatchAutomaticByPlatformRebootSettingNever,
		"unknown":    LinuxVMGuestPatchAutomaticByPlatformRebootSettingUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinuxVMGuestPatchAutomaticByPlatformRebootSetting(input)
	return &out, nil
}

type LinuxVMGuestPatchMode string

const (
	LinuxVMGuestPatchModeAutomaticByPlatform LinuxVMGuestPatchMode = "AutomaticByPlatform"
	LinuxVMGuestPatchModeImageDefault        LinuxVMGuestPatchMode = "ImageDefault"
)

func PossibleValuesForLinuxVMGuestPatchMode() []string {
	return []string{
		string(LinuxVMGuestPatchModeAutomaticByPlatform),
		string(LinuxVMGuestPatchModeImageDefault),
	}
}

func (s *LinuxVMGuestPatchMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLinuxVMGuestPatchMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLinuxVMGuestPatchMode(input string) (*LinuxVMGuestPatchMode, error) {
	vals := map[string]LinuxVMGuestPatchMode{
		"automaticbyplatform": LinuxVMGuestPatchModeAutomaticByPlatform,
		"imagedefault":        LinuxVMGuestPatchModeImageDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LinuxVMGuestPatchMode(input)
	return &out, nil
}

type Mode string

const (
	ModeAudit   Mode = "Audit"
	ModeEnforce Mode = "Enforce"
)

func PossibleValuesForMode() []string {
	return []string{
		string(ModeAudit),
		string(ModeEnforce),
	}
}

func (s *Mode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMode(input string) (*Mode, error) {
	vals := map[string]Mode{
		"audit":   ModeAudit,
		"enforce": ModeEnforce,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Mode(input)
	return &out, nil
}

type OperatingSystemType string

const (
	OperatingSystemTypeLinux   OperatingSystemType = "Linux"
	OperatingSystemTypeWindows OperatingSystemType = "Windows"
)

func PossibleValuesForOperatingSystemType() []string {
	return []string{
		string(OperatingSystemTypeLinux),
		string(OperatingSystemTypeWindows),
	}
}

func (s *OperatingSystemType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperatingSystemType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperatingSystemType(input string) (*OperatingSystemType, error) {
	vals := map[string]OperatingSystemType{
		"linux":   OperatingSystemTypeLinux,
		"windows": OperatingSystemTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OperatingSystemType(input)
	return &out, nil
}

type PassNames string

const (
	PassNamesOobeSystem PassNames = "OobeSystem"
)

func PossibleValuesForPassNames() []string {
	return []string{
		string(PassNamesOobeSystem),
	}
}

func (s *PassNames) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePassNames(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePassNames(input string) (*PassNames, error) {
	vals := map[string]PassNames{
		"oobesystem": PassNamesOobeSystem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PassNames(input)
	return &out, nil
}

type ProtocolTypes string

const (
	ProtocolTypesHTTP  ProtocolTypes = "Http"
	ProtocolTypesHTTPS ProtocolTypes = "Https"
)

func PossibleValuesForProtocolTypes() []string {
	return []string{
		string(ProtocolTypesHTTP),
		string(ProtocolTypesHTTPS),
	}
}

func (s *ProtocolTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProtocolTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProtocolTypes(input string) (*ProtocolTypes, error) {
	vals := map[string]ProtocolTypes{
		"http":  ProtocolTypesHTTP,
		"https": ProtocolTypesHTTPS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProtocolTypes(input)
	return &out, nil
}

type RestorePointEncryptionType string

const (
	RestorePointEncryptionTypeEncryptionAtRestWithCustomerKey             RestorePointEncryptionType = "EncryptionAtRestWithCustomerKey"
	RestorePointEncryptionTypeEncryptionAtRestWithPlatformAndCustomerKeys RestorePointEncryptionType = "EncryptionAtRestWithPlatformAndCustomerKeys"
	RestorePointEncryptionTypeEncryptionAtRestWithPlatformKey             RestorePointEncryptionType = "EncryptionAtRestWithPlatformKey"
)

func PossibleValuesForRestorePointEncryptionType() []string {
	return []string{
		string(RestorePointEncryptionTypeEncryptionAtRestWithCustomerKey),
		string(RestorePointEncryptionTypeEncryptionAtRestWithPlatformAndCustomerKeys),
		string(RestorePointEncryptionTypeEncryptionAtRestWithPlatformKey),
	}
}

func (s *RestorePointEncryptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRestorePointEncryptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRestorePointEncryptionType(input string) (*RestorePointEncryptionType, error) {
	vals := map[string]RestorePointEncryptionType{
		"encryptionatrestwithcustomerkey":             RestorePointEncryptionTypeEncryptionAtRestWithCustomerKey,
		"encryptionatrestwithplatformandcustomerkeys": RestorePointEncryptionTypeEncryptionAtRestWithPlatformAndCustomerKeys,
		"encryptionatrestwithplatformkey":             RestorePointEncryptionTypeEncryptionAtRestWithPlatformKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RestorePointEncryptionType(input)
	return &out, nil
}

type RestorePointExpandOptions string

const (
	RestorePointExpandOptionsInstanceView RestorePointExpandOptions = "instanceView"
)

func PossibleValuesForRestorePointExpandOptions() []string {
	return []string{
		string(RestorePointExpandOptionsInstanceView),
	}
}

func (s *RestorePointExpandOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRestorePointExpandOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRestorePointExpandOptions(input string) (*RestorePointExpandOptions, error) {
	vals := map[string]RestorePointExpandOptions{
		"instanceview": RestorePointExpandOptionsInstanceView,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RestorePointExpandOptions(input)
	return &out, nil
}

type SecurityEncryptionTypes string

const (
	SecurityEncryptionTypesDiskWithVMGuestState SecurityEncryptionTypes = "DiskWithVMGuestState"
	SecurityEncryptionTypesNonPersistedTPM      SecurityEncryptionTypes = "NonPersistedTPM"
	SecurityEncryptionTypesVMGuestStateOnly     SecurityEncryptionTypes = "VMGuestStateOnly"
)

func PossibleValuesForSecurityEncryptionTypes() []string {
	return []string{
		string(SecurityEncryptionTypesDiskWithVMGuestState),
		string(SecurityEncryptionTypesNonPersistedTPM),
		string(SecurityEncryptionTypesVMGuestStateOnly),
	}
}

func (s *SecurityEncryptionTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityEncryptionTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityEncryptionTypes(input string) (*SecurityEncryptionTypes, error) {
	vals := map[string]SecurityEncryptionTypes{
		"diskwithvmgueststate": SecurityEncryptionTypesDiskWithVMGuestState,
		"nonpersistedtpm":      SecurityEncryptionTypesNonPersistedTPM,
		"vmgueststateonly":     SecurityEncryptionTypesVMGuestStateOnly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityEncryptionTypes(input)
	return &out, nil
}

type SecurityTypes string

const (
	SecurityTypesConfidentialVM SecurityTypes = "ConfidentialVM"
	SecurityTypesTrustedLaunch  SecurityTypes = "TrustedLaunch"
)

func PossibleValuesForSecurityTypes() []string {
	return []string{
		string(SecurityTypesConfidentialVM),
		string(SecurityTypesTrustedLaunch),
	}
}

func (s *SecurityTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecurityTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecurityTypes(input string) (*SecurityTypes, error) {
	vals := map[string]SecurityTypes{
		"confidentialvm": SecurityTypesConfidentialVM,
		"trustedlaunch":  SecurityTypesTrustedLaunch,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecurityTypes(input)
	return &out, nil
}

type SettingNames string

const (
	SettingNamesAutoLogon          SettingNames = "AutoLogon"
	SettingNamesFirstLogonCommands SettingNames = "FirstLogonCommands"
)

func PossibleValuesForSettingNames() []string {
	return []string{
		string(SettingNamesAutoLogon),
		string(SettingNamesFirstLogonCommands),
	}
}

func (s *SettingNames) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSettingNames(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSettingNames(input string) (*SettingNames, error) {
	vals := map[string]SettingNames{
		"autologon":          SettingNamesAutoLogon,
		"firstlogoncommands": SettingNamesFirstLogonCommands,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SettingNames(input)
	return &out, nil
}

type StatusLevelTypes string

const (
	StatusLevelTypesError   StatusLevelTypes = "Error"
	StatusLevelTypesInfo    StatusLevelTypes = "Info"
	StatusLevelTypesWarning StatusLevelTypes = "Warning"
)

func PossibleValuesForStatusLevelTypes() []string {
	return []string{
		string(StatusLevelTypesError),
		string(StatusLevelTypesInfo),
		string(StatusLevelTypesWarning),
	}
}

func (s *StatusLevelTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatusLevelTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatusLevelTypes(input string) (*StatusLevelTypes, error) {
	vals := map[string]StatusLevelTypes{
		"error":   StatusLevelTypesError,
		"info":    StatusLevelTypesInfo,
		"warning": StatusLevelTypesWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StatusLevelTypes(input)
	return &out, nil
}

type StorageAccountTypes string

const (
	StorageAccountTypesPremiumLRS     StorageAccountTypes = "Premium_LRS"
	StorageAccountTypesPremiumVTwoLRS StorageAccountTypes = "PremiumV2_LRS"
	StorageAccountTypesPremiumZRS     StorageAccountTypes = "Premium_ZRS"
	StorageAccountTypesStandardLRS    StorageAccountTypes = "Standard_LRS"
	StorageAccountTypesStandardSSDLRS StorageAccountTypes = "StandardSSD_LRS"
	StorageAccountTypesStandardSSDZRS StorageAccountTypes = "StandardSSD_ZRS"
	StorageAccountTypesUltraSSDLRS    StorageAccountTypes = "UltraSSD_LRS"
)

func PossibleValuesForStorageAccountTypes() []string {
	return []string{
		string(StorageAccountTypesPremiumLRS),
		string(StorageAccountTypesPremiumVTwoLRS),
		string(StorageAccountTypesPremiumZRS),
		string(StorageAccountTypesStandardLRS),
		string(StorageAccountTypesStandardSSDLRS),
		string(StorageAccountTypesStandardSSDZRS),
		string(StorageAccountTypesUltraSSDLRS),
	}
}

func (s *StorageAccountTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAccountTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAccountTypes(input string) (*StorageAccountTypes, error) {
	vals := map[string]StorageAccountTypes{
		"premium_lrs":     StorageAccountTypesPremiumLRS,
		"premiumv2_lrs":   StorageAccountTypesPremiumVTwoLRS,
		"premium_zrs":     StorageAccountTypesPremiumZRS,
		"standard_lrs":    StorageAccountTypesStandardLRS,
		"standardssd_lrs": StorageAccountTypesStandardSSDLRS,
		"standardssd_zrs": StorageAccountTypesStandardSSDZRS,
		"ultrassd_lrs":    StorageAccountTypesUltraSSDLRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountTypes(input)
	return &out, nil
}

type VirtualMachineSizeTypes string

const (
	VirtualMachineSizeTypesBasicAFour                              VirtualMachineSizeTypes = "Basic_A4"
	VirtualMachineSizeTypesBasicAOne                               VirtualMachineSizeTypes = "Basic_A1"
	VirtualMachineSizeTypesBasicAThree                             VirtualMachineSizeTypes = "Basic_A3"
	VirtualMachineSizeTypesBasicATwo                               VirtualMachineSizeTypes = "Basic_A2"
	VirtualMachineSizeTypesBasicAZero                              VirtualMachineSizeTypes = "Basic_A0"
	VirtualMachineSizeTypesStandardAEight                          VirtualMachineSizeTypes = "Standard_A8"
	VirtualMachineSizeTypesStandardAEightVTwo                      VirtualMachineSizeTypes = "Standard_A8_v2"
	VirtualMachineSizeTypesStandardAEightmVTwo                     VirtualMachineSizeTypes = "Standard_A8m_v2"
	VirtualMachineSizeTypesStandardAFive                           VirtualMachineSizeTypes = "Standard_A5"
	VirtualMachineSizeTypesStandardAFour                           VirtualMachineSizeTypes = "Standard_A4"
	VirtualMachineSizeTypesStandardAFourVTwo                       VirtualMachineSizeTypes = "Standard_A4_v2"
	VirtualMachineSizeTypesStandardAFourmVTwo                      VirtualMachineSizeTypes = "Standard_A4m_v2"
	VirtualMachineSizeTypesStandardANine                           VirtualMachineSizeTypes = "Standard_A9"
	VirtualMachineSizeTypesStandardAOne                            VirtualMachineSizeTypes = "Standard_A1"
	VirtualMachineSizeTypesStandardAOneOne                         VirtualMachineSizeTypes = "Standard_A11"
	VirtualMachineSizeTypesStandardAOneVTwo                        VirtualMachineSizeTypes = "Standard_A1_v2"
	VirtualMachineSizeTypesStandardAOneZero                        VirtualMachineSizeTypes = "Standard_A10"
	VirtualMachineSizeTypesStandardASeven                          VirtualMachineSizeTypes = "Standard_A7"
	VirtualMachineSizeTypesStandardASix                            VirtualMachineSizeTypes = "Standard_A6"
	VirtualMachineSizeTypesStandardAThree                          VirtualMachineSizeTypes = "Standard_A3"
	VirtualMachineSizeTypesStandardATwo                            VirtualMachineSizeTypes = "Standard_A2"
	VirtualMachineSizeTypesStandardATwoVTwo                        VirtualMachineSizeTypes = "Standard_A2_v2"
	VirtualMachineSizeTypesStandardATwomVTwo                       VirtualMachineSizeTypes = "Standard_A2m_v2"
	VirtualMachineSizeTypesStandardAZero                           VirtualMachineSizeTypes = "Standard_A0"
	VirtualMachineSizeTypesStandardBEightms                        VirtualMachineSizeTypes = "Standard_B8ms"
	VirtualMachineSizeTypesStandardBFourms                         VirtualMachineSizeTypes = "Standard_B4ms"
	VirtualMachineSizeTypesStandardBOnems                          VirtualMachineSizeTypes = "Standard_B1ms"
	VirtualMachineSizeTypesStandardBOnes                           VirtualMachineSizeTypes = "Standard_B1s"
	VirtualMachineSizeTypesStandardBTwoms                          VirtualMachineSizeTypes = "Standard_B2ms"
	VirtualMachineSizeTypesStandardBTwos                           VirtualMachineSizeTypes = "Standard_B2s"
	VirtualMachineSizeTypesStandardDEightVThree                    VirtualMachineSizeTypes = "Standard_D8_v3"
	VirtualMachineSizeTypesStandardDEightsVThree                   VirtualMachineSizeTypes = "Standard_D8s_v3"
	VirtualMachineSizeTypesStandardDFiveVTwo                       VirtualMachineSizeTypes = "Standard_D5_v2"
	VirtualMachineSizeTypesStandardDFour                           VirtualMachineSizeTypes = "Standard_D4"
	VirtualMachineSizeTypesStandardDFourVThree                     VirtualMachineSizeTypes = "Standard_D4_v3"
	VirtualMachineSizeTypesStandardDFourVTwo                       VirtualMachineSizeTypes = "Standard_D4_v2"
	VirtualMachineSizeTypesStandardDFoursVThree                    VirtualMachineSizeTypes = "Standard_D4s_v3"
	VirtualMachineSizeTypesStandardDOne                            VirtualMachineSizeTypes = "Standard_D1"
	VirtualMachineSizeTypesStandardDOneFiveVTwo                    VirtualMachineSizeTypes = "Standard_D15_v2"
	VirtualMachineSizeTypesStandardDOneFour                        VirtualMachineSizeTypes = "Standard_D14"
	VirtualMachineSizeTypesStandardDOneFourVTwo                    VirtualMachineSizeTypes = "Standard_D14_v2"
	VirtualMachineSizeTypesStandardDOneOne                         VirtualMachineSizeTypes = "Standard_D11"
	VirtualMachineSizeTypesStandardDOneOneVTwo                     VirtualMachineSizeTypes = "Standard_D11_v2"
	VirtualMachineSizeTypesStandardDOneSixVThree                   VirtualMachineSizeTypes = "Standard_D16_v3"
	VirtualMachineSizeTypesStandardDOneSixsVThree                  VirtualMachineSizeTypes = "Standard_D16s_v3"
	VirtualMachineSizeTypesStandardDOneThree                       VirtualMachineSizeTypes = "Standard_D13"
	VirtualMachineSizeTypesStandardDOneThreeVTwo                   VirtualMachineSizeTypes = "Standard_D13_v2"
	VirtualMachineSizeTypesStandardDOneTwo                         VirtualMachineSizeTypes = "Standard_D12"
	VirtualMachineSizeTypesStandardDOneTwoVTwo                     VirtualMachineSizeTypes = "Standard_D12_v2"
	VirtualMachineSizeTypesStandardDOneVTwo                        VirtualMachineSizeTypes = "Standard_D1_v2"
	VirtualMachineSizeTypesStandardDSFiveVTwo                      VirtualMachineSizeTypes = "Standard_DS5_v2"
	VirtualMachineSizeTypesStandardDSFour                          VirtualMachineSizeTypes = "Standard_DS4"
	VirtualMachineSizeTypesStandardDSFourVTwo                      VirtualMachineSizeTypes = "Standard_DS4_v2"
	VirtualMachineSizeTypesStandardDSOne                           VirtualMachineSizeTypes = "Standard_DS1"
	VirtualMachineSizeTypesStandardDSOneFiveVTwo                   VirtualMachineSizeTypes = "Standard_DS15_v2"
	VirtualMachineSizeTypesStandardDSOneFour                       VirtualMachineSizeTypes = "Standard_DS14"
	VirtualMachineSizeTypesStandardDSOneFourNegativeEightVTwo      VirtualMachineSizeTypes = "Standard_DS14-8_v2"
	VirtualMachineSizeTypesStandardDSOneFourNegativeFourVTwo       VirtualMachineSizeTypes = "Standard_DS14-4_v2"
	VirtualMachineSizeTypesStandardDSOneFourVTwo                   VirtualMachineSizeTypes = "Standard_DS14_v2"
	VirtualMachineSizeTypesStandardDSOneOne                        VirtualMachineSizeTypes = "Standard_DS11"
	VirtualMachineSizeTypesStandardDSOneOneVTwo                    VirtualMachineSizeTypes = "Standard_DS11_v2"
	VirtualMachineSizeTypesStandardDSOneThree                      VirtualMachineSizeTypes = "Standard_DS13"
	VirtualMachineSizeTypesStandardDSOneThreeNegativeFourVTwo      VirtualMachineSizeTypes = "Standard_DS13-4_v2"
	VirtualMachineSizeTypesStandardDSOneThreeNegativeTwoVTwo       VirtualMachineSizeTypes = "Standard_DS13-2_v2"
	VirtualMachineSizeTypesStandardDSOneThreeVTwo                  VirtualMachineSizeTypes = "Standard_DS13_v2"
	VirtualMachineSizeTypesStandardDSOneTwo                        VirtualMachineSizeTypes = "Standard_DS12"
	VirtualMachineSizeTypesStandardDSOneTwoVTwo                    VirtualMachineSizeTypes = "Standard_DS12_v2"
	VirtualMachineSizeTypesStandardDSOneVTwo                       VirtualMachineSizeTypes = "Standard_DS1_v2"
	VirtualMachineSizeTypesStandardDSThree                         VirtualMachineSizeTypes = "Standard_DS3"
	VirtualMachineSizeTypesStandardDSThreeVTwo                     VirtualMachineSizeTypes = "Standard_DS3_v2"
	VirtualMachineSizeTypesStandardDSTwo                           VirtualMachineSizeTypes = "Standard_DS2"
	VirtualMachineSizeTypesStandardDSTwoVTwo                       VirtualMachineSizeTypes = "Standard_DS2_v2"
	VirtualMachineSizeTypesStandardDSixFourVThree                  VirtualMachineSizeTypes = "Standard_D64_v3"
	VirtualMachineSizeTypesStandardDSixFoursVThree                 VirtualMachineSizeTypes = "Standard_D64s_v3"
	VirtualMachineSizeTypesStandardDThree                          VirtualMachineSizeTypes = "Standard_D3"
	VirtualMachineSizeTypesStandardDThreeTwoVThree                 VirtualMachineSizeTypes = "Standard_D32_v3"
	VirtualMachineSizeTypesStandardDThreeTwosVThree                VirtualMachineSizeTypes = "Standard_D32s_v3"
	VirtualMachineSizeTypesStandardDThreeVTwo                      VirtualMachineSizeTypes = "Standard_D3_v2"
	VirtualMachineSizeTypesStandardDTwo                            VirtualMachineSizeTypes = "Standard_D2"
	VirtualMachineSizeTypesStandardDTwoVThree                      VirtualMachineSizeTypes = "Standard_D2_v3"
	VirtualMachineSizeTypesStandardDTwoVTwo                        VirtualMachineSizeTypes = "Standard_D2_v2"
	VirtualMachineSizeTypesStandardDTwosVThree                     VirtualMachineSizeTypes = "Standard_D2s_v3"
	VirtualMachineSizeTypesStandardEEightVThree                    VirtualMachineSizeTypes = "Standard_E8_v3"
	VirtualMachineSizeTypesStandardEEightsVThree                   VirtualMachineSizeTypes = "Standard_E8s_v3"
	VirtualMachineSizeTypesStandardEFourVThree                     VirtualMachineSizeTypes = "Standard_E4_v3"
	VirtualMachineSizeTypesStandardEFoursVThree                    VirtualMachineSizeTypes = "Standard_E4s_v3"
	VirtualMachineSizeTypesStandardEOneSixVThree                   VirtualMachineSizeTypes = "Standard_E16_v3"
	VirtualMachineSizeTypesStandardEOneSixsVThree                  VirtualMachineSizeTypes = "Standard_E16s_v3"
	VirtualMachineSizeTypesStandardESixFourNegativeOneSixsVThree   VirtualMachineSizeTypes = "Standard_E64-16s_v3"
	VirtualMachineSizeTypesStandardESixFourNegativeThreeTwosVThree VirtualMachineSizeTypes = "Standard_E64-32s_v3"
	VirtualMachineSizeTypesStandardESixFourVThree                  VirtualMachineSizeTypes = "Standard_E64_v3"
	VirtualMachineSizeTypesStandardESixFoursVThree                 VirtualMachineSizeTypes = "Standard_E64s_v3"
	VirtualMachineSizeTypesStandardEThreeTwoNegativeEightsVThree   VirtualMachineSizeTypes = "Standard_E32-8s_v3"
	VirtualMachineSizeTypesStandardEThreeTwoNegativeOneSixVThree   VirtualMachineSizeTypes = "Standard_E32-16_v3"
	VirtualMachineSizeTypesStandardEThreeTwoVThree                 VirtualMachineSizeTypes = "Standard_E32_v3"
	VirtualMachineSizeTypesStandardEThreeTwosVThree                VirtualMachineSizeTypes = "Standard_E32s_v3"
	VirtualMachineSizeTypesStandardETwoVThree                      VirtualMachineSizeTypes = "Standard_E2_v3"
	VirtualMachineSizeTypesStandardETwosVThree                     VirtualMachineSizeTypes = "Standard_E2s_v3"
	VirtualMachineSizeTypesStandardFEight                          VirtualMachineSizeTypes = "Standard_F8"
	VirtualMachineSizeTypesStandardFEights                         VirtualMachineSizeTypes = "Standard_F8s"
	VirtualMachineSizeTypesStandardFEightsVTwo                     VirtualMachineSizeTypes = "Standard_F8s_v2"
	VirtualMachineSizeTypesStandardFFour                           VirtualMachineSizeTypes = "Standard_F4"
	VirtualMachineSizeTypesStandardFFours                          VirtualMachineSizeTypes = "Standard_F4s"
	VirtualMachineSizeTypesStandardFFoursVTwo                      VirtualMachineSizeTypes = "Standard_F4s_v2"
	VirtualMachineSizeTypesStandardFOne                            VirtualMachineSizeTypes = "Standard_F1"
	VirtualMachineSizeTypesStandardFOneSix                         VirtualMachineSizeTypes = "Standard_F16"
	VirtualMachineSizeTypesStandardFOneSixs                        VirtualMachineSizeTypes = "Standard_F16s"
	VirtualMachineSizeTypesStandardFOneSixsVTwo                    VirtualMachineSizeTypes = "Standard_F16s_v2"
	VirtualMachineSizeTypesStandardFOnes                           VirtualMachineSizeTypes = "Standard_F1s"
	VirtualMachineSizeTypesStandardFSevenTwosVTwo                  VirtualMachineSizeTypes = "Standard_F72s_v2"
	VirtualMachineSizeTypesStandardFSixFoursVTwo                   VirtualMachineSizeTypes = "Standard_F64s_v2"
	VirtualMachineSizeTypesStandardFThreeTwosVTwo                  VirtualMachineSizeTypes = "Standard_F32s_v2"
	VirtualMachineSizeTypesStandardFTwo                            VirtualMachineSizeTypes = "Standard_F2"
	VirtualMachineSizeTypesStandardFTwos                           VirtualMachineSizeTypes = "Standard_F2s"
	VirtualMachineSizeTypesStandardFTwosVTwo                       VirtualMachineSizeTypes = "Standard_F2s_v2"
	VirtualMachineSizeTypesStandardGFive                           VirtualMachineSizeTypes = "Standard_G5"
	VirtualMachineSizeTypesStandardGFour                           VirtualMachineSizeTypes = "Standard_G4"
	VirtualMachineSizeTypesStandardGOne                            VirtualMachineSizeTypes = "Standard_G1"
	VirtualMachineSizeTypesStandardGSFive                          VirtualMachineSizeTypes = "Standard_GS5"
	VirtualMachineSizeTypesStandardGSFiveNegativeEight             VirtualMachineSizeTypes = "Standard_GS5-8"
	VirtualMachineSizeTypesStandardGSFiveNegativeOneSix            VirtualMachineSizeTypes = "Standard_GS5-16"
	VirtualMachineSizeTypesStandardGSFour                          VirtualMachineSizeTypes = "Standard_GS4"
	VirtualMachineSizeTypesStandardGSFourNegativeEight             VirtualMachineSizeTypes = "Standard_GS4-8"
	VirtualMachineSizeTypesStandardGSFourNegativeFour              VirtualMachineSizeTypes = "Standard_GS4-4"
	VirtualMachineSizeTypesStandardGSOne                           VirtualMachineSizeTypes = "Standard_GS1"
	VirtualMachineSizeTypesStandardGSThree                         VirtualMachineSizeTypes = "Standard_GS3"
	VirtualMachineSizeTypesStandardGSTwo                           VirtualMachineSizeTypes = "Standard_GS2"
	VirtualMachineSizeTypesStandardGThree                          VirtualMachineSizeTypes = "Standard_G3"
	VirtualMachineSizeTypesStandardGTwo                            VirtualMachineSizeTypes = "Standard_G2"
	VirtualMachineSizeTypesStandardHEight                          VirtualMachineSizeTypes = "Standard_H8"
	VirtualMachineSizeTypesStandardHEightm                         VirtualMachineSizeTypes = "Standard_H8m"
	VirtualMachineSizeTypesStandardHOneSix                         VirtualMachineSizeTypes = "Standard_H16"
	VirtualMachineSizeTypesStandardHOneSixm                        VirtualMachineSizeTypes = "Standard_H16m"
	VirtualMachineSizeTypesStandardHOneSixmr                       VirtualMachineSizeTypes = "Standard_H16mr"
	VirtualMachineSizeTypesStandardHOneSixr                        VirtualMachineSizeTypes = "Standard_H16r"
	VirtualMachineSizeTypesStandardLEights                         VirtualMachineSizeTypes = "Standard_L8s"
	VirtualMachineSizeTypesStandardLFours                          VirtualMachineSizeTypes = "Standard_L4s"
	VirtualMachineSizeTypesStandardLOneSixs                        VirtualMachineSizeTypes = "Standard_L16s"
	VirtualMachineSizeTypesStandardLThreeTwos                      VirtualMachineSizeTypes = "Standard_L32s"
	VirtualMachineSizeTypesStandardMOneTwoEightNegativeSixFourms   VirtualMachineSizeTypes = "Standard_M128-64ms"
	VirtualMachineSizeTypesStandardMOneTwoEightNegativeThreeTwoms  VirtualMachineSizeTypes = "Standard_M128-32ms"
	VirtualMachineSizeTypesStandardMOneTwoEightms                  VirtualMachineSizeTypes = "Standard_M128ms"
	VirtualMachineSizeTypesStandardMOneTwoEights                   VirtualMachineSizeTypes = "Standard_M128s"
	VirtualMachineSizeTypesStandardMSixFourNegativeOneSixms        VirtualMachineSizeTypes = "Standard_M64-16ms"
	VirtualMachineSizeTypesStandardMSixFourNegativeThreeTwoms      VirtualMachineSizeTypes = "Standard_M64-32ms"
	VirtualMachineSizeTypesStandardMSixFourms                      VirtualMachineSizeTypes = "Standard_M64ms"
	VirtualMachineSizeTypesStandardMSixFours                       VirtualMachineSizeTypes = "Standard_M64s"
	VirtualMachineSizeTypesStandardNCOneTwo                        VirtualMachineSizeTypes = "Standard_NC12"
	VirtualMachineSizeTypesStandardNCOneTwosVThree                 VirtualMachineSizeTypes = "Standard_NC12s_v3"
	VirtualMachineSizeTypesStandardNCOneTwosVTwo                   VirtualMachineSizeTypes = "Standard_NC12s_v2"
	VirtualMachineSizeTypesStandardNCSix                           VirtualMachineSizeTypes = "Standard_NC6"
	VirtualMachineSizeTypesStandardNCSixsVThree                    VirtualMachineSizeTypes = "Standard_NC6s_v3"
	VirtualMachineSizeTypesStandardNCSixsVTwo                      VirtualMachineSizeTypes = "Standard_NC6s_v2"
	VirtualMachineSizeTypesStandardNCTwoFour                       VirtualMachineSizeTypes = "Standard_NC24"
	VirtualMachineSizeTypesStandardNCTwoFourr                      VirtualMachineSizeTypes = "Standard_NC24r"
	VirtualMachineSizeTypesStandardNCTwoFourrsVThree               VirtualMachineSizeTypes = "Standard_NC24rs_v3"
	VirtualMachineSizeTypesStandardNCTwoFourrsVTwo                 VirtualMachineSizeTypes = "Standard_NC24rs_v2"
	VirtualMachineSizeTypesStandardNCTwoFoursVThree                VirtualMachineSizeTypes = "Standard_NC24s_v3"
	VirtualMachineSizeTypesStandardNCTwoFoursVTwo                  VirtualMachineSizeTypes = "Standard_NC24s_v2"
	VirtualMachineSizeTypesStandardNDOneTwos                       VirtualMachineSizeTypes = "Standard_ND12s"
	VirtualMachineSizeTypesStandardNDSixs                          VirtualMachineSizeTypes = "Standard_ND6s"
	VirtualMachineSizeTypesStandardNDTwoFourrs                     VirtualMachineSizeTypes = "Standard_ND24rs"
	VirtualMachineSizeTypesStandardNDTwoFours                      VirtualMachineSizeTypes = "Standard_ND24s"
	VirtualMachineSizeTypesStandardNVOneTwo                        VirtualMachineSizeTypes = "Standard_NV12"
	VirtualMachineSizeTypesStandardNVSix                           VirtualMachineSizeTypes = "Standard_NV6"
	VirtualMachineSizeTypesStandardNVTwoFour                       VirtualMachineSizeTypes = "Standard_NV24"
)

func PossibleValuesForVirtualMachineSizeTypes() []string {
	return []string{
		string(VirtualMachineSizeTypesBasicAFour),
		string(VirtualMachineSizeTypesBasicAOne),
		string(VirtualMachineSizeTypesBasicAThree),
		string(VirtualMachineSizeTypesBasicATwo),
		string(VirtualMachineSizeTypesBasicAZero),
		string(VirtualMachineSizeTypesStandardAEight),
		string(VirtualMachineSizeTypesStandardAEightVTwo),
		string(VirtualMachineSizeTypesStandardAEightmVTwo),
		string(VirtualMachineSizeTypesStandardAFive),
		string(VirtualMachineSizeTypesStandardAFour),
		string(VirtualMachineSizeTypesStandardAFourVTwo),
		string(VirtualMachineSizeTypesStandardAFourmVTwo),
		string(VirtualMachineSizeTypesStandardANine),
		string(VirtualMachineSizeTypesStandardAOne),
		string(VirtualMachineSizeTypesStandardAOneOne),
		string(VirtualMachineSizeTypesStandardAOneVTwo),
		string(VirtualMachineSizeTypesStandardAOneZero),
		string(VirtualMachineSizeTypesStandardASeven),
		string(VirtualMachineSizeTypesStandardASix),
		string(VirtualMachineSizeTypesStandardAThree),
		string(VirtualMachineSizeTypesStandardATwo),
		string(VirtualMachineSizeTypesStandardATwoVTwo),
		string(VirtualMachineSizeTypesStandardATwomVTwo),
		string(VirtualMachineSizeTypesStandardAZero),
		string(VirtualMachineSizeTypesStandardBEightms),
		string(VirtualMachineSizeTypesStandardBFourms),
		string(VirtualMachineSizeTypesStandardBOnems),
		string(VirtualMachineSizeTypesStandardBOnes),
		string(VirtualMachineSizeTypesStandardBTwoms),
		string(VirtualMachineSizeTypesStandardBTwos),
		string(VirtualMachineSizeTypesStandardDEightVThree),
		string(VirtualMachineSizeTypesStandardDEightsVThree),
		string(VirtualMachineSizeTypesStandardDFiveVTwo),
		string(VirtualMachineSizeTypesStandardDFour),
		string(VirtualMachineSizeTypesStandardDFourVThree),
		string(VirtualMachineSizeTypesStandardDFourVTwo),
		string(VirtualMachineSizeTypesStandardDFoursVThree),
		string(VirtualMachineSizeTypesStandardDOne),
		string(VirtualMachineSizeTypesStandardDOneFiveVTwo),
		string(VirtualMachineSizeTypesStandardDOneFour),
		string(VirtualMachineSizeTypesStandardDOneFourVTwo),
		string(VirtualMachineSizeTypesStandardDOneOne),
		string(VirtualMachineSizeTypesStandardDOneOneVTwo),
		string(VirtualMachineSizeTypesStandardDOneSixVThree),
		string(VirtualMachineSizeTypesStandardDOneSixsVThree),
		string(VirtualMachineSizeTypesStandardDOneThree),
		string(VirtualMachineSizeTypesStandardDOneThreeVTwo),
		string(VirtualMachineSizeTypesStandardDOneTwo),
		string(VirtualMachineSizeTypesStandardDOneTwoVTwo),
		string(VirtualMachineSizeTypesStandardDOneVTwo),
		string(VirtualMachineSizeTypesStandardDSFiveVTwo),
		string(VirtualMachineSizeTypesStandardDSFour),
		string(VirtualMachineSizeTypesStandardDSFourVTwo),
		string(VirtualMachineSizeTypesStandardDSOne),
		string(VirtualMachineSizeTypesStandardDSOneFiveVTwo),
		string(VirtualMachineSizeTypesStandardDSOneFour),
		string(VirtualMachineSizeTypesStandardDSOneFourNegativeEightVTwo),
		string(VirtualMachineSizeTypesStandardDSOneFourNegativeFourVTwo),
		string(VirtualMachineSizeTypesStandardDSOneFourVTwo),
		string(VirtualMachineSizeTypesStandardDSOneOne),
		string(VirtualMachineSizeTypesStandardDSOneOneVTwo),
		string(VirtualMachineSizeTypesStandardDSOneThree),
		string(VirtualMachineSizeTypesStandardDSOneThreeNegativeFourVTwo),
		string(VirtualMachineSizeTypesStandardDSOneThreeNegativeTwoVTwo),
		string(VirtualMachineSizeTypesStandardDSOneThreeVTwo),
		string(VirtualMachineSizeTypesStandardDSOneTwo),
		string(VirtualMachineSizeTypesStandardDSOneTwoVTwo),
		string(VirtualMachineSizeTypesStandardDSOneVTwo),
		string(VirtualMachineSizeTypesStandardDSThree),
		string(VirtualMachineSizeTypesStandardDSThreeVTwo),
		string(VirtualMachineSizeTypesStandardDSTwo),
		string(VirtualMachineSizeTypesStandardDSTwoVTwo),
		string(VirtualMachineSizeTypesStandardDSixFourVThree),
		string(VirtualMachineSizeTypesStandardDSixFoursVThree),
		string(VirtualMachineSizeTypesStandardDThree),
		string(VirtualMachineSizeTypesStandardDThreeTwoVThree),
		string(VirtualMachineSizeTypesStandardDThreeTwosVThree),
		string(VirtualMachineSizeTypesStandardDThreeVTwo),
		string(VirtualMachineSizeTypesStandardDTwo),
		string(VirtualMachineSizeTypesStandardDTwoVThree),
		string(VirtualMachineSizeTypesStandardDTwoVTwo),
		string(VirtualMachineSizeTypesStandardDTwosVThree),
		string(VirtualMachineSizeTypesStandardEEightVThree),
		string(VirtualMachineSizeTypesStandardEEightsVThree),
		string(VirtualMachineSizeTypesStandardEFourVThree),
		string(VirtualMachineSizeTypesStandardEFoursVThree),
		string(VirtualMachineSizeTypesStandardEOneSixVThree),
		string(VirtualMachineSizeTypesStandardEOneSixsVThree),
		string(VirtualMachineSizeTypesStandardESixFourNegativeOneSixsVThree),
		string(VirtualMachineSizeTypesStandardESixFourNegativeThreeTwosVThree),
		string(VirtualMachineSizeTypesStandardESixFourVThree),
		string(VirtualMachineSizeTypesStandardESixFoursVThree),
		string(VirtualMachineSizeTypesStandardEThreeTwoNegativeEightsVThree),
		string(VirtualMachineSizeTypesStandardEThreeTwoNegativeOneSixVThree),
		string(VirtualMachineSizeTypesStandardEThreeTwoVThree),
		string(VirtualMachineSizeTypesStandardEThreeTwosVThree),
		string(VirtualMachineSizeTypesStandardETwoVThree),
		string(VirtualMachineSizeTypesStandardETwosVThree),
		string(VirtualMachineSizeTypesStandardFEight),
		string(VirtualMachineSizeTypesStandardFEights),
		string(VirtualMachineSizeTypesStandardFEightsVTwo),
		string(VirtualMachineSizeTypesStandardFFour),
		string(VirtualMachineSizeTypesStandardFFours),
		string(VirtualMachineSizeTypesStandardFFoursVTwo),
		string(VirtualMachineSizeTypesStandardFOne),
		string(VirtualMachineSizeTypesStandardFOneSix),
		string(VirtualMachineSizeTypesStandardFOneSixs),
		string(VirtualMachineSizeTypesStandardFOneSixsVTwo),
		string(VirtualMachineSizeTypesStandardFOnes),
		string(VirtualMachineSizeTypesStandardFSevenTwosVTwo),
		string(VirtualMachineSizeTypesStandardFSixFoursVTwo),
		string(VirtualMachineSizeTypesStandardFThreeTwosVTwo),
		string(VirtualMachineSizeTypesStandardFTwo),
		string(VirtualMachineSizeTypesStandardFTwos),
		string(VirtualMachineSizeTypesStandardFTwosVTwo),
		string(VirtualMachineSizeTypesStandardGFive),
		string(VirtualMachineSizeTypesStandardGFour),
		string(VirtualMachineSizeTypesStandardGOne),
		string(VirtualMachineSizeTypesStandardGSFive),
		string(VirtualMachineSizeTypesStandardGSFiveNegativeEight),
		string(VirtualMachineSizeTypesStandardGSFiveNegativeOneSix),
		string(VirtualMachineSizeTypesStandardGSFour),
		string(VirtualMachineSizeTypesStandardGSFourNegativeEight),
		string(VirtualMachineSizeTypesStandardGSFourNegativeFour),
		string(VirtualMachineSizeTypesStandardGSOne),
		string(VirtualMachineSizeTypesStandardGSThree),
		string(VirtualMachineSizeTypesStandardGSTwo),
		string(VirtualMachineSizeTypesStandardGThree),
		string(VirtualMachineSizeTypesStandardGTwo),
		string(VirtualMachineSizeTypesStandardHEight),
		string(VirtualMachineSizeTypesStandardHEightm),
		string(VirtualMachineSizeTypesStandardHOneSix),
		string(VirtualMachineSizeTypesStandardHOneSixm),
		string(VirtualMachineSizeTypesStandardHOneSixmr),
		string(VirtualMachineSizeTypesStandardHOneSixr),
		string(VirtualMachineSizeTypesStandardLEights),
		string(VirtualMachineSizeTypesStandardLFours),
		string(VirtualMachineSizeTypesStandardLOneSixs),
		string(VirtualMachineSizeTypesStandardLThreeTwos),
		string(VirtualMachineSizeTypesStandardMOneTwoEightNegativeSixFourms),
		string(VirtualMachineSizeTypesStandardMOneTwoEightNegativeThreeTwoms),
		string(VirtualMachineSizeTypesStandardMOneTwoEightms),
		string(VirtualMachineSizeTypesStandardMOneTwoEights),
		string(VirtualMachineSizeTypesStandardMSixFourNegativeOneSixms),
		string(VirtualMachineSizeTypesStandardMSixFourNegativeThreeTwoms),
		string(VirtualMachineSizeTypesStandardMSixFourms),
		string(VirtualMachineSizeTypesStandardMSixFours),
		string(VirtualMachineSizeTypesStandardNCOneTwo),
		string(VirtualMachineSizeTypesStandardNCOneTwosVThree),
		string(VirtualMachineSizeTypesStandardNCOneTwosVTwo),
		string(VirtualMachineSizeTypesStandardNCSix),
		string(VirtualMachineSizeTypesStandardNCSixsVThree),
		string(VirtualMachineSizeTypesStandardNCSixsVTwo),
		string(VirtualMachineSizeTypesStandardNCTwoFour),
		string(VirtualMachineSizeTypesStandardNCTwoFourr),
		string(VirtualMachineSizeTypesStandardNCTwoFourrsVThree),
		string(VirtualMachineSizeTypesStandardNCTwoFourrsVTwo),
		string(VirtualMachineSizeTypesStandardNCTwoFoursVThree),
		string(VirtualMachineSizeTypesStandardNCTwoFoursVTwo),
		string(VirtualMachineSizeTypesStandardNDOneTwos),
		string(VirtualMachineSizeTypesStandardNDSixs),
		string(VirtualMachineSizeTypesStandardNDTwoFourrs),
		string(VirtualMachineSizeTypesStandardNDTwoFours),
		string(VirtualMachineSizeTypesStandardNVOneTwo),
		string(VirtualMachineSizeTypesStandardNVSix),
		string(VirtualMachineSizeTypesStandardNVTwoFour),
	}
}

func (s *VirtualMachineSizeTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVirtualMachineSizeTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVirtualMachineSizeTypes(input string) (*VirtualMachineSizeTypes, error) {
	vals := map[string]VirtualMachineSizeTypes{
		"basic_a4":            VirtualMachineSizeTypesBasicAFour,
		"basic_a1":            VirtualMachineSizeTypesBasicAOne,
		"basic_a3":            VirtualMachineSizeTypesBasicAThree,
		"basic_a2":            VirtualMachineSizeTypesBasicATwo,
		"basic_a0":            VirtualMachineSizeTypesBasicAZero,
		"standard_a8":         VirtualMachineSizeTypesStandardAEight,
		"standard_a8_v2":      VirtualMachineSizeTypesStandardAEightVTwo,
		"standard_a8m_v2":     VirtualMachineSizeTypesStandardAEightmVTwo,
		"standard_a5":         VirtualMachineSizeTypesStandardAFive,
		"standard_a4":         VirtualMachineSizeTypesStandardAFour,
		"standard_a4_v2":      VirtualMachineSizeTypesStandardAFourVTwo,
		"standard_a4m_v2":     VirtualMachineSizeTypesStandardAFourmVTwo,
		"standard_a9":         VirtualMachineSizeTypesStandardANine,
		"standard_a1":         VirtualMachineSizeTypesStandardAOne,
		"standard_a11":        VirtualMachineSizeTypesStandardAOneOne,
		"standard_a1_v2":      VirtualMachineSizeTypesStandardAOneVTwo,
		"standard_a10":        VirtualMachineSizeTypesStandardAOneZero,
		"standard_a7":         VirtualMachineSizeTypesStandardASeven,
		"standard_a6":         VirtualMachineSizeTypesStandardASix,
		"standard_a3":         VirtualMachineSizeTypesStandardAThree,
		"standard_a2":         VirtualMachineSizeTypesStandardATwo,
		"standard_a2_v2":      VirtualMachineSizeTypesStandardATwoVTwo,
		"standard_a2m_v2":     VirtualMachineSizeTypesStandardATwomVTwo,
		"standard_a0":         VirtualMachineSizeTypesStandardAZero,
		"standard_b8ms":       VirtualMachineSizeTypesStandardBEightms,
		"standard_b4ms":       VirtualMachineSizeTypesStandardBFourms,
		"standard_b1ms":       VirtualMachineSizeTypesStandardBOnems,
		"standard_b1s":        VirtualMachineSizeTypesStandardBOnes,
		"standard_b2ms":       VirtualMachineSizeTypesStandardBTwoms,
		"standard_b2s":        VirtualMachineSizeTypesStandardBTwos,
		"standard_d8_v3":      VirtualMachineSizeTypesStandardDEightVThree,
		"standard_d8s_v3":     VirtualMachineSizeTypesStandardDEightsVThree,
		"standard_d5_v2":      VirtualMachineSizeTypesStandardDFiveVTwo,
		"standard_d4":         VirtualMachineSizeTypesStandardDFour,
		"standard_d4_v3":      VirtualMachineSizeTypesStandardDFourVThree,
		"standard_d4_v2":      VirtualMachineSizeTypesStandardDFourVTwo,
		"standard_d4s_v3":     VirtualMachineSizeTypesStandardDFoursVThree,
		"standard_d1":         VirtualMachineSizeTypesStandardDOne,
		"standard_d15_v2":     VirtualMachineSizeTypesStandardDOneFiveVTwo,
		"standard_d14":        VirtualMachineSizeTypesStandardDOneFour,
		"standard_d14_v2":     VirtualMachineSizeTypesStandardDOneFourVTwo,
		"standard_d11":        VirtualMachineSizeTypesStandardDOneOne,
		"standard_d11_v2":     VirtualMachineSizeTypesStandardDOneOneVTwo,
		"standard_d16_v3":     VirtualMachineSizeTypesStandardDOneSixVThree,
		"standard_d16s_v3":    VirtualMachineSizeTypesStandardDOneSixsVThree,
		"standard_d13":        VirtualMachineSizeTypesStandardDOneThree,
		"standard_d13_v2":     VirtualMachineSizeTypesStandardDOneThreeVTwo,
		"standard_d12":        VirtualMachineSizeTypesStandardDOneTwo,
		"standard_d12_v2":     VirtualMachineSizeTypesStandardDOneTwoVTwo,
		"standard_d1_v2":      VirtualMachineSizeTypesStandardDOneVTwo,
		"standard_ds5_v2":     VirtualMachineSizeTypesStandardDSFiveVTwo,
		"standard_ds4":        VirtualMachineSizeTypesStandardDSFour,
		"standard_ds4_v2":     VirtualMachineSizeTypesStandardDSFourVTwo,
		"standard_ds1":        VirtualMachineSizeTypesStandardDSOne,
		"standard_ds15_v2":    VirtualMachineSizeTypesStandardDSOneFiveVTwo,
		"standard_ds14":       VirtualMachineSizeTypesStandardDSOneFour,
		"standard_ds14-8_v2":  VirtualMachineSizeTypesStandardDSOneFourNegativeEightVTwo,
		"standard_ds14-4_v2":  VirtualMachineSizeTypesStandardDSOneFourNegativeFourVTwo,
		"standard_ds14_v2":    VirtualMachineSizeTypesStandardDSOneFourVTwo,
		"standard_ds11":       VirtualMachineSizeTypesStandardDSOneOne,
		"standard_ds11_v2":    VirtualMachineSizeTypesStandardDSOneOneVTwo,
		"standard_ds13":       VirtualMachineSizeTypesStandardDSOneThree,
		"standard_ds13-4_v2":  VirtualMachineSizeTypesStandardDSOneThreeNegativeFourVTwo,
		"standard_ds13-2_v2":  VirtualMachineSizeTypesStandardDSOneThreeNegativeTwoVTwo,
		"standard_ds13_v2":    VirtualMachineSizeTypesStandardDSOneThreeVTwo,
		"standard_ds12":       VirtualMachineSizeTypesStandardDSOneTwo,
		"standard_ds12_v2":    VirtualMachineSizeTypesStandardDSOneTwoVTwo,
		"standard_ds1_v2":     VirtualMachineSizeTypesStandardDSOneVTwo,
		"standard_ds3":        VirtualMachineSizeTypesStandardDSThree,
		"standard_ds3_v2":     VirtualMachineSizeTypesStandardDSThreeVTwo,
		"standard_ds2":        VirtualMachineSizeTypesStandardDSTwo,
		"standard_ds2_v2":     VirtualMachineSizeTypesStandardDSTwoVTwo,
		"standard_d64_v3":     VirtualMachineSizeTypesStandardDSixFourVThree,
		"standard_d64s_v3":    VirtualMachineSizeTypesStandardDSixFoursVThree,
		"standard_d3":         VirtualMachineSizeTypesStandardDThree,
		"standard_d32_v3":     VirtualMachineSizeTypesStandardDThreeTwoVThree,
		"standard_d32s_v3":    VirtualMachineSizeTypesStandardDThreeTwosVThree,
		"standard_d3_v2":      VirtualMachineSizeTypesStandardDThreeVTwo,
		"standard_d2":         VirtualMachineSizeTypesStandardDTwo,
		"standard_d2_v3":      VirtualMachineSizeTypesStandardDTwoVThree,
		"standard_d2_v2":      VirtualMachineSizeTypesStandardDTwoVTwo,
		"standard_d2s_v3":     VirtualMachineSizeTypesStandardDTwosVThree,
		"standard_e8_v3":      VirtualMachineSizeTypesStandardEEightVThree,
		"standard_e8s_v3":     VirtualMachineSizeTypesStandardEEightsVThree,
		"standard_e4_v3":      VirtualMachineSizeTypesStandardEFourVThree,
		"standard_e4s_v3":     VirtualMachineSizeTypesStandardEFoursVThree,
		"standard_e16_v3":     VirtualMachineSizeTypesStandardEOneSixVThree,
		"standard_e16s_v3":    VirtualMachineSizeTypesStandardEOneSixsVThree,
		"standard_e64-16s_v3": VirtualMachineSizeTypesStandardESixFourNegativeOneSixsVThree,
		"standard_e64-32s_v3": VirtualMachineSizeTypesStandardESixFourNegativeThreeTwosVThree,
		"standard_e64_v3":     VirtualMachineSizeTypesStandardESixFourVThree,
		"standard_e64s_v3":    VirtualMachineSizeTypesStandardESixFoursVThree,
		"standard_e32-8s_v3":  VirtualMachineSizeTypesStandardEThreeTwoNegativeEightsVThree,
		"standard_e32-16_v3":  VirtualMachineSizeTypesStandardEThreeTwoNegativeOneSixVThree,
		"standard_e32_v3":     VirtualMachineSizeTypesStandardEThreeTwoVThree,
		"standard_e32s_v3":    VirtualMachineSizeTypesStandardEThreeTwosVThree,
		"standard_e2_v3":      VirtualMachineSizeTypesStandardETwoVThree,
		"standard_e2s_v3":     VirtualMachineSizeTypesStandardETwosVThree,
		"standard_f8":         VirtualMachineSizeTypesStandardFEight,
		"standard_f8s":        VirtualMachineSizeTypesStandardFEights,
		"standard_f8s_v2":     VirtualMachineSizeTypesStandardFEightsVTwo,
		"standard_f4":         VirtualMachineSizeTypesStandardFFour,
		"standard_f4s":        VirtualMachineSizeTypesStandardFFours,
		"standard_f4s_v2":     VirtualMachineSizeTypesStandardFFoursVTwo,
		"standard_f1":         VirtualMachineSizeTypesStandardFOne,
		"standard_f16":        VirtualMachineSizeTypesStandardFOneSix,
		"standard_f16s":       VirtualMachineSizeTypesStandardFOneSixs,
		"standard_f16s_v2":    VirtualMachineSizeTypesStandardFOneSixsVTwo,
		"standard_f1s":        VirtualMachineSizeTypesStandardFOnes,
		"standard_f72s_v2":    VirtualMachineSizeTypesStandardFSevenTwosVTwo,
		"standard_f64s_v2":    VirtualMachineSizeTypesStandardFSixFoursVTwo,
		"standard_f32s_v2":    VirtualMachineSizeTypesStandardFThreeTwosVTwo,
		"standard_f2":         VirtualMachineSizeTypesStandardFTwo,
		"standard_f2s":        VirtualMachineSizeTypesStandardFTwos,
		"standard_f2s_v2":     VirtualMachineSizeTypesStandardFTwosVTwo,
		"standard_g5":         VirtualMachineSizeTypesStandardGFive,
		"standard_g4":         VirtualMachineSizeTypesStandardGFour,
		"standard_g1":         VirtualMachineSizeTypesStandardGOne,
		"standard_gs5":        VirtualMachineSizeTypesStandardGSFive,
		"standard_gs5-8":      VirtualMachineSizeTypesStandardGSFiveNegativeEight,
		"standard_gs5-16":     VirtualMachineSizeTypesStandardGSFiveNegativeOneSix,
		"standard_gs4":        VirtualMachineSizeTypesStandardGSFour,
		"standard_gs4-8":      VirtualMachineSizeTypesStandardGSFourNegativeEight,
		"standard_gs4-4":      VirtualMachineSizeTypesStandardGSFourNegativeFour,
		"standard_gs1":        VirtualMachineSizeTypesStandardGSOne,
		"standard_gs3":        VirtualMachineSizeTypesStandardGSThree,
		"standard_gs2":        VirtualMachineSizeTypesStandardGSTwo,
		"standard_g3":         VirtualMachineSizeTypesStandardGThree,
		"standard_g2":         VirtualMachineSizeTypesStandardGTwo,
		"standard_h8":         VirtualMachineSizeTypesStandardHEight,
		"standard_h8m":        VirtualMachineSizeTypesStandardHEightm,
		"standard_h16":        VirtualMachineSizeTypesStandardHOneSix,
		"standard_h16m":       VirtualMachineSizeTypesStandardHOneSixm,
		"standard_h16mr":      VirtualMachineSizeTypesStandardHOneSixmr,
		"standard_h16r":       VirtualMachineSizeTypesStandardHOneSixr,
		"standard_l8s":        VirtualMachineSizeTypesStandardLEights,
		"standard_l4s":        VirtualMachineSizeTypesStandardLFours,
		"standard_l16s":       VirtualMachineSizeTypesStandardLOneSixs,
		"standard_l32s":       VirtualMachineSizeTypesStandardLThreeTwos,
		"standard_m128-64ms":  VirtualMachineSizeTypesStandardMOneTwoEightNegativeSixFourms,
		"standard_m128-32ms":  VirtualMachineSizeTypesStandardMOneTwoEightNegativeThreeTwoms,
		"standard_m128ms":     VirtualMachineSizeTypesStandardMOneTwoEightms,
		"standard_m128s":      VirtualMachineSizeTypesStandardMOneTwoEights,
		"standard_m64-16ms":   VirtualMachineSizeTypesStandardMSixFourNegativeOneSixms,
		"standard_m64-32ms":   VirtualMachineSizeTypesStandardMSixFourNegativeThreeTwoms,
		"standard_m64ms":      VirtualMachineSizeTypesStandardMSixFourms,
		"standard_m64s":       VirtualMachineSizeTypesStandardMSixFours,
		"standard_nc12":       VirtualMachineSizeTypesStandardNCOneTwo,
		"standard_nc12s_v3":   VirtualMachineSizeTypesStandardNCOneTwosVThree,
		"standard_nc12s_v2":   VirtualMachineSizeTypesStandardNCOneTwosVTwo,
		"standard_nc6":        VirtualMachineSizeTypesStandardNCSix,
		"standard_nc6s_v3":    VirtualMachineSizeTypesStandardNCSixsVThree,
		"standard_nc6s_v2":    VirtualMachineSizeTypesStandardNCSixsVTwo,
		"standard_nc24":       VirtualMachineSizeTypesStandardNCTwoFour,
		"standard_nc24r":      VirtualMachineSizeTypesStandardNCTwoFourr,
		"standard_nc24rs_v3":  VirtualMachineSizeTypesStandardNCTwoFourrsVThree,
		"standard_nc24rs_v2":  VirtualMachineSizeTypesStandardNCTwoFourrsVTwo,
		"standard_nc24s_v3":   VirtualMachineSizeTypesStandardNCTwoFoursVThree,
		"standard_nc24s_v2":   VirtualMachineSizeTypesStandardNCTwoFoursVTwo,
		"standard_nd12s":      VirtualMachineSizeTypesStandardNDOneTwos,
		"standard_nd6s":       VirtualMachineSizeTypesStandardNDSixs,
		"standard_nd24rs":     VirtualMachineSizeTypesStandardNDTwoFourrs,
		"standard_nd24s":      VirtualMachineSizeTypesStandardNDTwoFours,
		"standard_nv12":       VirtualMachineSizeTypesStandardNVOneTwo,
		"standard_nv6":        VirtualMachineSizeTypesStandardNVSix,
		"standard_nv24":       VirtualMachineSizeTypesStandardNVTwoFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VirtualMachineSizeTypes(input)
	return &out, nil
}

type WindowsPatchAssessmentMode string

const (
	WindowsPatchAssessmentModeAutomaticByPlatform WindowsPatchAssessmentMode = "AutomaticByPlatform"
	WindowsPatchAssessmentModeImageDefault        WindowsPatchAssessmentMode = "ImageDefault"
)

func PossibleValuesForWindowsPatchAssessmentMode() []string {
	return []string{
		string(WindowsPatchAssessmentModeAutomaticByPlatform),
		string(WindowsPatchAssessmentModeImageDefault),
	}
}

func (s *WindowsPatchAssessmentMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWindowsPatchAssessmentMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWindowsPatchAssessmentMode(input string) (*WindowsPatchAssessmentMode, error) {
	vals := map[string]WindowsPatchAssessmentMode{
		"automaticbyplatform": WindowsPatchAssessmentModeAutomaticByPlatform,
		"imagedefault":        WindowsPatchAssessmentModeImageDefault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WindowsPatchAssessmentMode(input)
	return &out, nil
}

type WindowsVMGuestPatchAutomaticByPlatformRebootSetting string

const (
	WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways     WindowsVMGuestPatchAutomaticByPlatformRebootSetting = "Always"
	WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired WindowsVMGuestPatchAutomaticByPlatformRebootSetting = "IfRequired"
	WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever      WindowsVMGuestPatchAutomaticByPlatformRebootSetting = "Never"
	WindowsVMGuestPatchAutomaticByPlatformRebootSettingUnknown    WindowsVMGuestPatchAutomaticByPlatformRebootSetting = "Unknown"
)

func PossibleValuesForWindowsVMGuestPatchAutomaticByPlatformRebootSetting() []string {
	return []string{
		string(WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways),
		string(WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired),
		string(WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever),
		string(WindowsVMGuestPatchAutomaticByPlatformRebootSettingUnknown),
	}
}

func (s *WindowsVMGuestPatchAutomaticByPlatformRebootSetting) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWindowsVMGuestPatchAutomaticByPlatformRebootSetting(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWindowsVMGuestPatchAutomaticByPlatformRebootSetting(input string) (*WindowsVMGuestPatchAutomaticByPlatformRebootSetting, error) {
	vals := map[string]WindowsVMGuestPatchAutomaticByPlatformRebootSetting{
		"always":     WindowsVMGuestPatchAutomaticByPlatformRebootSettingAlways,
		"ifrequired": WindowsVMGuestPatchAutomaticByPlatformRebootSettingIfRequired,
		"never":      WindowsVMGuestPatchAutomaticByPlatformRebootSettingNever,
		"unknown":    WindowsVMGuestPatchAutomaticByPlatformRebootSettingUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WindowsVMGuestPatchAutomaticByPlatformRebootSetting(input)
	return &out, nil
}

type WindowsVMGuestPatchMode string

const (
	WindowsVMGuestPatchModeAutomaticByOS       WindowsVMGuestPatchMode = "AutomaticByOS"
	WindowsVMGuestPatchModeAutomaticByPlatform WindowsVMGuestPatchMode = "AutomaticByPlatform"
	WindowsVMGuestPatchModeManual              WindowsVMGuestPatchMode = "Manual"
)

func PossibleValuesForWindowsVMGuestPatchMode() []string {
	return []string{
		string(WindowsVMGuestPatchModeAutomaticByOS),
		string(WindowsVMGuestPatchModeAutomaticByPlatform),
		string(WindowsVMGuestPatchModeManual),
	}
}

func (s *WindowsVMGuestPatchMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWindowsVMGuestPatchMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWindowsVMGuestPatchMode(input string) (*WindowsVMGuestPatchMode, error) {
	vals := map[string]WindowsVMGuestPatchMode{
		"automaticbyos":       WindowsVMGuestPatchModeAutomaticByOS,
		"automaticbyplatform": WindowsVMGuestPatchModeAutomaticByPlatform,
		"manual":              WindowsVMGuestPatchModeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WindowsVMGuestPatchMode(input)
	return &out, nil
}

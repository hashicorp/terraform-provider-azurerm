package fleets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AcceleratorManufacturer string

const (
	AcceleratorManufacturerAMD    AcceleratorManufacturer = "AMD"
	AcceleratorManufacturerNvidia AcceleratorManufacturer = "Nvidia"
	AcceleratorManufacturerXilinx AcceleratorManufacturer = "Xilinx"
)

func PossibleValuesForAcceleratorManufacturer() []string {
	return []string{
		string(AcceleratorManufacturerAMD),
		string(AcceleratorManufacturerNvidia),
		string(AcceleratorManufacturerXilinx),
	}
}

func (s *AcceleratorManufacturer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAcceleratorManufacturer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAcceleratorManufacturer(input string) (*AcceleratorManufacturer, error) {
	vals := map[string]AcceleratorManufacturer{
		"amd":    AcceleratorManufacturerAMD,
		"nvidia": AcceleratorManufacturerNvidia,
		"xilinx": AcceleratorManufacturerXilinx,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AcceleratorManufacturer(input)
	return &out, nil
}

type AcceleratorType string

const (
	AcceleratorTypeFPGA AcceleratorType = "FPGA"
	AcceleratorTypeGPU  AcceleratorType = "GPU"
)

func PossibleValuesForAcceleratorType() []string {
	return []string{
		string(AcceleratorTypeFPGA),
		string(AcceleratorTypeGPU),
	}
}

func (s *AcceleratorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAcceleratorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAcceleratorType(input string) (*AcceleratorType, error) {
	vals := map[string]AcceleratorType{
		"fpga": AcceleratorTypeFPGA,
		"gpu":  AcceleratorTypeGPU,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AcceleratorType(input)
	return &out, nil
}

type ArchitectureType string

const (
	ArchitectureTypeARMSixFour ArchitectureType = "ARM64"
	ArchitectureTypeXSixFour   ArchitectureType = "X64"
)

func PossibleValuesForArchitectureType() []string {
	return []string{
		string(ArchitectureTypeARMSixFour),
		string(ArchitectureTypeXSixFour),
	}
}

func (s *ArchitectureType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseArchitectureType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseArchitectureType(input string) (*ArchitectureType, error) {
	vals := map[string]ArchitectureType{
		"arm64": ArchitectureTypeARMSixFour,
		"x64":   ArchitectureTypeXSixFour,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ArchitectureType(input)
	return &out, nil
}

type CPUManufacturer string

const (
	CPUManufacturerAMD       CPUManufacturer = "AMD"
	CPUManufacturerAmpere    CPUManufacturer = "Ampere"
	CPUManufacturerIntel     CPUManufacturer = "Intel"
	CPUManufacturerMicrosoft CPUManufacturer = "Microsoft"
)

func PossibleValuesForCPUManufacturer() []string {
	return []string{
		string(CPUManufacturerAMD),
		string(CPUManufacturerAmpere),
		string(CPUManufacturerIntel),
		string(CPUManufacturerMicrosoft),
	}
}

func (s *CPUManufacturer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCPUManufacturer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCPUManufacturer(input string) (*CPUManufacturer, error) {
	vals := map[string]CPUManufacturer{
		"amd":       CPUManufacturerAMD,
		"ampere":    CPUManufacturerAmpere,
		"intel":     CPUManufacturerIntel,
		"microsoft": CPUManufacturerMicrosoft,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CPUManufacturer(input)
	return &out, nil
}

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

type ComponentName string

const (
	ComponentNameMicrosoftNegativeWindowsNegativeShellNegativeSetup ComponentName = "Microsoft-Windows-Shell-Setup"
)

func PossibleValuesForComponentName() []string {
	return []string{
		string(ComponentNameMicrosoftNegativeWindowsNegativeShellNegativeSetup),
	}
}

func (s *ComponentName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComponentName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComponentName(input string) (*ComponentName, error) {
	vals := map[string]ComponentName{
		"microsoft-windows-shell-setup": ComponentNameMicrosoftNegativeWindowsNegativeShellNegativeSetup,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComponentName(input)
	return &out, nil
}

type DeleteOptions string

const (
	DeleteOptionsDelete DeleteOptions = "Delete"
	DeleteOptionsDetach DeleteOptions = "Detach"
)

func PossibleValuesForDeleteOptions() []string {
	return []string{
		string(DeleteOptionsDelete),
		string(DeleteOptionsDetach),
	}
}

func (s *DeleteOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDeleteOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDeleteOptions(input string) (*DeleteOptions, error) {
	vals := map[string]DeleteOptions{
		"delete": DeleteOptionsDelete,
		"detach": DeleteOptionsDetach,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DeleteOptions(input)
	return &out, nil
}

type DiffDiskOptions string

const (
	DiffDiskOptionsLocal DiffDiskOptions = "Local"
)

func PossibleValuesForDiffDiskOptions() []string {
	return []string{
		string(DiffDiskOptionsLocal),
	}
}

func (s *DiffDiskOptions) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiffDiskOptions(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiffDiskOptions(input string) (*DiffDiskOptions, error) {
	vals := map[string]DiffDiskOptions{
		"local": DiffDiskOptionsLocal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiffDiskOptions(input)
	return &out, nil
}

type DiffDiskPlacement string

const (
	DiffDiskPlacementCacheDisk    DiffDiskPlacement = "CacheDisk"
	DiffDiskPlacementNVMeDisk     DiffDiskPlacement = "NvmeDisk"
	DiffDiskPlacementResourceDisk DiffDiskPlacement = "ResourceDisk"
)

func PossibleValuesForDiffDiskPlacement() []string {
	return []string{
		string(DiffDiskPlacementCacheDisk),
		string(DiffDiskPlacementNVMeDisk),
		string(DiffDiskPlacementResourceDisk),
	}
}

func (s *DiffDiskPlacement) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiffDiskPlacement(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiffDiskPlacement(input string) (*DiffDiskPlacement, error) {
	vals := map[string]DiffDiskPlacement{
		"cachedisk":    DiffDiskPlacementCacheDisk,
		"nvmedisk":     DiffDiskPlacementNVMeDisk,
		"resourcedisk": DiffDiskPlacementResourceDisk,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiffDiskPlacement(input)
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

type DiskCreateOptionTypes string

const (
	DiskCreateOptionTypesAttach    DiskCreateOptionTypes = "Attach"
	DiskCreateOptionTypesCopy      DiskCreateOptionTypes = "Copy"
	DiskCreateOptionTypesEmpty     DiskCreateOptionTypes = "Empty"
	DiskCreateOptionTypesFromImage DiskCreateOptionTypes = "FromImage"
	DiskCreateOptionTypesRestore   DiskCreateOptionTypes = "Restore"
)

func PossibleValuesForDiskCreateOptionTypes() []string {
	return []string{
		string(DiskCreateOptionTypesAttach),
		string(DiskCreateOptionTypesCopy),
		string(DiskCreateOptionTypesEmpty),
		string(DiskCreateOptionTypesFromImage),
		string(DiskCreateOptionTypesRestore),
	}
}

func (s *DiskCreateOptionTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskCreateOptionTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskCreateOptionTypes(input string) (*DiskCreateOptionTypes, error) {
	vals := map[string]DiskCreateOptionTypes{
		"attach":    DiskCreateOptionTypesAttach,
		"copy":      DiskCreateOptionTypesCopy,
		"empty":     DiskCreateOptionTypesEmpty,
		"fromimage": DiskCreateOptionTypesFromImage,
		"restore":   DiskCreateOptionTypesRestore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskCreateOptionTypes(input)
	return &out, nil
}

type DiskDeleteOptionTypes string

const (
	DiskDeleteOptionTypesDelete DiskDeleteOptionTypes = "Delete"
	DiskDeleteOptionTypesDetach DiskDeleteOptionTypes = "Detach"
)

func PossibleValuesForDiskDeleteOptionTypes() []string {
	return []string{
		string(DiskDeleteOptionTypesDelete),
		string(DiskDeleteOptionTypesDetach),
	}
}

func (s *DiskDeleteOptionTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskDeleteOptionTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskDeleteOptionTypes(input string) (*DiskDeleteOptionTypes, error) {
	vals := map[string]DiskDeleteOptionTypes{
		"delete": DiskDeleteOptionTypesDelete,
		"detach": DiskDeleteOptionTypesDetach,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskDeleteOptionTypes(input)
	return &out, nil
}

type DomainNameLabelScopeTypes string

const (
	DomainNameLabelScopeTypesNoReuse            DomainNameLabelScopeTypes = "NoReuse"
	DomainNameLabelScopeTypesResourceGroupReuse DomainNameLabelScopeTypes = "ResourceGroupReuse"
	DomainNameLabelScopeTypesSubscriptionReuse  DomainNameLabelScopeTypes = "SubscriptionReuse"
	DomainNameLabelScopeTypesTenantReuse        DomainNameLabelScopeTypes = "TenantReuse"
)

func PossibleValuesForDomainNameLabelScopeTypes() []string {
	return []string{
		string(DomainNameLabelScopeTypesNoReuse),
		string(DomainNameLabelScopeTypesResourceGroupReuse),
		string(DomainNameLabelScopeTypesSubscriptionReuse),
		string(DomainNameLabelScopeTypesTenantReuse),
	}
}

func (s *DomainNameLabelScopeTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDomainNameLabelScopeTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDomainNameLabelScopeTypes(input string) (*DomainNameLabelScopeTypes, error) {
	vals := map[string]DomainNameLabelScopeTypes{
		"noreuse":            DomainNameLabelScopeTypesNoReuse,
		"resourcegroupreuse": DomainNameLabelScopeTypesResourceGroupReuse,
		"subscriptionreuse":  DomainNameLabelScopeTypesSubscriptionReuse,
		"tenantreuse":        DomainNameLabelScopeTypesTenantReuse,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DomainNameLabelScopeTypes(input)
	return &out, nil
}

type EvictionPolicy string

const (
	EvictionPolicyDeallocate EvictionPolicy = "Deallocate"
	EvictionPolicyDelete     EvictionPolicy = "Delete"
)

func PossibleValuesForEvictionPolicy() []string {
	return []string{
		string(EvictionPolicyDeallocate),
		string(EvictionPolicyDelete),
	}
}

func (s *EvictionPolicy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEvictionPolicy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEvictionPolicy(input string) (*EvictionPolicy, error) {
	vals := map[string]EvictionPolicy{
		"deallocate": EvictionPolicyDeallocate,
		"delete":     EvictionPolicyDelete,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EvictionPolicy(input)
	return &out, nil
}

type IPVersion string

const (
	IPVersionIPvFour IPVersion = "IPv4"
	IPVersionIPvSix  IPVersion = "IPv6"
)

func PossibleValuesForIPVersion() []string {
	return []string{
		string(IPVersionIPvFour),
		string(IPVersionIPvSix),
	}
}

func (s *IPVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIPVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIPVersion(input string) (*IPVersion, error) {
	vals := map[string]IPVersion{
		"ipv4": IPVersionIPvFour,
		"ipv6": IPVersionIPvSix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IPVersion(input)
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

type LocalStorageDiskType string

const (
	LocalStorageDiskTypeHDD LocalStorageDiskType = "HDD"
	LocalStorageDiskTypeSSD LocalStorageDiskType = "SSD"
)

func PossibleValuesForLocalStorageDiskType() []string {
	return []string{
		string(LocalStorageDiskTypeHDD),
		string(LocalStorageDiskTypeSSD),
	}
}

func (s *LocalStorageDiskType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLocalStorageDiskType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLocalStorageDiskType(input string) (*LocalStorageDiskType, error) {
	vals := map[string]LocalStorageDiskType{
		"hdd": LocalStorageDiskTypeHDD,
		"ssd": LocalStorageDiskTypeSSD,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LocalStorageDiskType(input)
	return &out, nil
}

type ManagedServiceIdentityType string

const (
	ManagedServiceIdentityTypeNone                       ManagedServiceIdentityType = "None"
	ManagedServiceIdentityTypeSystemAssigned             ManagedServiceIdentityType = "SystemAssigned"
	ManagedServiceIdentityTypeSystemAssignedUserAssigned ManagedServiceIdentityType = "SystemAssigned,UserAssigned"
	ManagedServiceIdentityTypeUserAssigned               ManagedServiceIdentityType = "UserAssigned"
)

func PossibleValuesForManagedServiceIdentityType() []string {
	return []string{
		string(ManagedServiceIdentityTypeNone),
		string(ManagedServiceIdentityTypeSystemAssigned),
		string(ManagedServiceIdentityTypeSystemAssignedUserAssigned),
		string(ManagedServiceIdentityTypeUserAssigned),
	}
}

func (s *ManagedServiceIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedServiceIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedServiceIdentityType(input string) (*ManagedServiceIdentityType, error) {
	vals := map[string]ManagedServiceIdentityType{
		"none":                        ManagedServiceIdentityTypeNone,
		"systemassigned":              ManagedServiceIdentityTypeSystemAssigned,
		"systemassigned,userassigned": ManagedServiceIdentityTypeSystemAssignedUserAssigned,
		"userassigned":                ManagedServiceIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedServiceIdentityType(input)
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

type NetworkApiVersion string

const (
	NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne NetworkApiVersion = "2020-11-01"
)

func PossibleValuesForNetworkApiVersion() []string {
	return []string{
		string(NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne),
	}
}

func (s *NetworkApiVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkApiVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkApiVersion(input string) (*NetworkApiVersion, error) {
	vals := map[string]NetworkApiVersion{
		"2020-11-01": NetworkApiVersionTwoZeroTwoZeroNegativeOneOneNegativeZeroOne,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkApiVersion(input)
	return &out, nil
}

type NetworkInterfaceAuxiliaryMode string

const (
	NetworkInterfaceAuxiliaryModeAcceleratedConnections NetworkInterfaceAuxiliaryMode = "AcceleratedConnections"
	NetworkInterfaceAuxiliaryModeFloating               NetworkInterfaceAuxiliaryMode = "Floating"
	NetworkInterfaceAuxiliaryModeNone                   NetworkInterfaceAuxiliaryMode = "None"
)

func PossibleValuesForNetworkInterfaceAuxiliaryMode() []string {
	return []string{
		string(NetworkInterfaceAuxiliaryModeAcceleratedConnections),
		string(NetworkInterfaceAuxiliaryModeFloating),
		string(NetworkInterfaceAuxiliaryModeNone),
	}
}

func (s *NetworkInterfaceAuxiliaryMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceAuxiliaryMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceAuxiliaryMode(input string) (*NetworkInterfaceAuxiliaryMode, error) {
	vals := map[string]NetworkInterfaceAuxiliaryMode{
		"acceleratedconnections": NetworkInterfaceAuxiliaryModeAcceleratedConnections,
		"floating":               NetworkInterfaceAuxiliaryModeFloating,
		"none":                   NetworkInterfaceAuxiliaryModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceAuxiliaryMode(input)
	return &out, nil
}

type NetworkInterfaceAuxiliarySku string

const (
	NetworkInterfaceAuxiliarySkuAEight NetworkInterfaceAuxiliarySku = "A8"
	NetworkInterfaceAuxiliarySkuAFour  NetworkInterfaceAuxiliarySku = "A4"
	NetworkInterfaceAuxiliarySkuAOne   NetworkInterfaceAuxiliarySku = "A1"
	NetworkInterfaceAuxiliarySkuATwo   NetworkInterfaceAuxiliarySku = "A2"
	NetworkInterfaceAuxiliarySkuNone   NetworkInterfaceAuxiliarySku = "None"
)

func PossibleValuesForNetworkInterfaceAuxiliarySku() []string {
	return []string{
		string(NetworkInterfaceAuxiliarySkuAEight),
		string(NetworkInterfaceAuxiliarySkuAFour),
		string(NetworkInterfaceAuxiliarySkuAOne),
		string(NetworkInterfaceAuxiliarySkuATwo),
		string(NetworkInterfaceAuxiliarySkuNone),
	}
}

func (s *NetworkInterfaceAuxiliarySku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkInterfaceAuxiliarySku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkInterfaceAuxiliarySku(input string) (*NetworkInterfaceAuxiliarySku, error) {
	vals := map[string]NetworkInterfaceAuxiliarySku{
		"a8":   NetworkInterfaceAuxiliarySkuAEight,
		"a4":   NetworkInterfaceAuxiliarySkuAFour,
		"a1":   NetworkInterfaceAuxiliarySkuAOne,
		"a2":   NetworkInterfaceAuxiliarySkuATwo,
		"none": NetworkInterfaceAuxiliarySkuNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkInterfaceAuxiliarySku(input)
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

type PassName string

const (
	PassNameOobeSystem PassName = "OobeSystem"
)

func PossibleValuesForPassName() []string {
	return []string{
		string(PassNameOobeSystem),
	}
}

func (s *PassName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePassName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePassName(input string) (*PassName, error) {
	vals := map[string]PassName{
		"oobesystem": PassNameOobeSystem,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PassName(input)
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

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateCreating  ProvisioningState = "Creating"
	ProvisioningStateDeleting  ProvisioningState = "Deleting"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateMigrating ProvisioningState = "Migrating"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	ProvisioningStateUpdating  ProvisioningState = "Updating"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreating),
		string(ProvisioningStateDeleting),
		string(ProvisioningStateFailed),
		string(ProvisioningStateMigrating),
		string(ProvisioningStateSucceeded),
		string(ProvisioningStateUpdating),
	}
}

func (s *ProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseProvisioningState(input string) (*ProvisioningState, error) {
	vals := map[string]ProvisioningState{
		"canceled":  ProvisioningStateCanceled,
		"creating":  ProvisioningStateCreating,
		"deleting":  ProvisioningStateDeleting,
		"failed":    ProvisioningStateFailed,
		"migrating": ProvisioningStateMigrating,
		"succeeded": ProvisioningStateSucceeded,
		"updating":  ProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type PublicIPAddressSkuName string

const (
	PublicIPAddressSkuNameBasic    PublicIPAddressSkuName = "Basic"
	PublicIPAddressSkuNameStandard PublicIPAddressSkuName = "Standard"
)

func PossibleValuesForPublicIPAddressSkuName() []string {
	return []string{
		string(PublicIPAddressSkuNameBasic),
		string(PublicIPAddressSkuNameStandard),
	}
}

func (s *PublicIPAddressSkuName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressSkuName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressSkuName(input string) (*PublicIPAddressSkuName, error) {
	vals := map[string]PublicIPAddressSkuName{
		"basic":    PublicIPAddressSkuNameBasic,
		"standard": PublicIPAddressSkuNameStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressSkuName(input)
	return &out, nil
}

type PublicIPAddressSkuTier string

const (
	PublicIPAddressSkuTierGlobal   PublicIPAddressSkuTier = "Global"
	PublicIPAddressSkuTierRegional PublicIPAddressSkuTier = "Regional"
)

func PossibleValuesForPublicIPAddressSkuTier() []string {
	return []string{
		string(PublicIPAddressSkuTierGlobal),
		string(PublicIPAddressSkuTierRegional),
	}
}

func (s *PublicIPAddressSkuTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePublicIPAddressSkuTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePublicIPAddressSkuTier(input string) (*PublicIPAddressSkuTier, error) {
	vals := map[string]PublicIPAddressSkuTier{
		"global":   PublicIPAddressSkuTierGlobal,
		"regional": PublicIPAddressSkuTierRegional,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicIPAddressSkuTier(input)
	return &out, nil
}

type RegularPriorityAllocationStrategy string

const (
	RegularPriorityAllocationStrategyLowestPrice RegularPriorityAllocationStrategy = "LowestPrice"
	RegularPriorityAllocationStrategyPrioritized RegularPriorityAllocationStrategy = "Prioritized"
)

func PossibleValuesForRegularPriorityAllocationStrategy() []string {
	return []string{
		string(RegularPriorityAllocationStrategyLowestPrice),
		string(RegularPriorityAllocationStrategyPrioritized),
	}
}

func (s *RegularPriorityAllocationStrategy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRegularPriorityAllocationStrategy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRegularPriorityAllocationStrategy(input string) (*RegularPriorityAllocationStrategy, error) {
	vals := map[string]RegularPriorityAllocationStrategy{
		"lowestprice": RegularPriorityAllocationStrategyLowestPrice,
		"prioritized": RegularPriorityAllocationStrategyPrioritized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RegularPriorityAllocationStrategy(input)
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

type SpotAllocationStrategy string

const (
	SpotAllocationStrategyCapacityOptimized      SpotAllocationStrategy = "CapacityOptimized"
	SpotAllocationStrategyLowestPrice            SpotAllocationStrategy = "LowestPrice"
	SpotAllocationStrategyPriceCapacityOptimized SpotAllocationStrategy = "PriceCapacityOptimized"
)

func PossibleValuesForSpotAllocationStrategy() []string {
	return []string{
		string(SpotAllocationStrategyCapacityOptimized),
		string(SpotAllocationStrategyLowestPrice),
		string(SpotAllocationStrategyPriceCapacityOptimized),
	}
}

func (s *SpotAllocationStrategy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSpotAllocationStrategy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSpotAllocationStrategy(input string) (*SpotAllocationStrategy, error) {
	vals := map[string]SpotAllocationStrategy{
		"capacityoptimized":      SpotAllocationStrategyCapacityOptimized,
		"lowestprice":            SpotAllocationStrategyLowestPrice,
		"pricecapacityoptimized": SpotAllocationStrategyPriceCapacityOptimized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SpotAllocationStrategy(input)
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

type VMAttributeSupport string

const (
	VMAttributeSupportExcluded VMAttributeSupport = "Excluded"
	VMAttributeSupportIncluded VMAttributeSupport = "Included"
	VMAttributeSupportRequired VMAttributeSupport = "Required"
)

func PossibleValuesForVMAttributeSupport() []string {
	return []string{
		string(VMAttributeSupportExcluded),
		string(VMAttributeSupportIncluded),
		string(VMAttributeSupportRequired),
	}
}

func (s *VMAttributeSupport) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMAttributeSupport(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMAttributeSupport(input string) (*VMAttributeSupport, error) {
	vals := map[string]VMAttributeSupport{
		"excluded": VMAttributeSupportExcluded,
		"included": VMAttributeSupportIncluded,
		"required": VMAttributeSupportRequired,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMAttributeSupport(input)
	return &out, nil
}

type VMCategory string

const (
	VMCategoryComputeOptimized       VMCategory = "ComputeOptimized"
	VMCategoryFpgaAccelerated        VMCategory = "FpgaAccelerated"
	VMCategoryGeneralPurpose         VMCategory = "GeneralPurpose"
	VMCategoryGpuAccelerated         VMCategory = "GpuAccelerated"
	VMCategoryHighPerformanceCompute VMCategory = "HighPerformanceCompute"
	VMCategoryMemoryOptimized        VMCategory = "MemoryOptimized"
	VMCategoryStorageOptimized       VMCategory = "StorageOptimized"
)

func PossibleValuesForVMCategory() []string {
	return []string{
		string(VMCategoryComputeOptimized),
		string(VMCategoryFpgaAccelerated),
		string(VMCategoryGeneralPurpose),
		string(VMCategoryGpuAccelerated),
		string(VMCategoryHighPerformanceCompute),
		string(VMCategoryMemoryOptimized),
		string(VMCategoryStorageOptimized),
	}
}

func (s *VMCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMCategory(input string) (*VMCategory, error) {
	vals := map[string]VMCategory{
		"computeoptimized":       VMCategoryComputeOptimized,
		"fpgaaccelerated":        VMCategoryFpgaAccelerated,
		"generalpurpose":         VMCategoryGeneralPurpose,
		"gpuaccelerated":         VMCategoryGpuAccelerated,
		"highperformancecompute": VMCategoryHighPerformanceCompute,
		"memoryoptimized":        VMCategoryMemoryOptimized,
		"storageoptimized":       VMCategoryStorageOptimized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMCategory(input)
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

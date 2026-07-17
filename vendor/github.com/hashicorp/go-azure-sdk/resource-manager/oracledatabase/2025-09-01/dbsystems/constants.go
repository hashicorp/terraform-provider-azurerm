package dbsystems

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureResourceProvisioningState string

const (
	AzureResourceProvisioningStateCanceled     AzureResourceProvisioningState = "Canceled"
	AzureResourceProvisioningStateFailed       AzureResourceProvisioningState = "Failed"
	AzureResourceProvisioningStateProvisioning AzureResourceProvisioningState = "Provisioning"
	AzureResourceProvisioningStateSucceeded    AzureResourceProvisioningState = "Succeeded"
)

func PossibleValuesForAzureResourceProvisioningState() []string {
	return []string{
		string(AzureResourceProvisioningStateCanceled),
		string(AzureResourceProvisioningStateFailed),
		string(AzureResourceProvisioningStateProvisioning),
		string(AzureResourceProvisioningStateSucceeded),
	}
}

func (s *AzureResourceProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAzureResourceProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAzureResourceProvisioningState(input string) (*AzureResourceProvisioningState, error) {
	vals := map[string]AzureResourceProvisioningState{
		"canceled":     AzureResourceProvisioningStateCanceled,
		"failed":       AzureResourceProvisioningStateFailed,
		"provisioning": AzureResourceProvisioningStateProvisioning,
		"succeeded":    AzureResourceProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AzureResourceProvisioningState(input)
	return &out, nil
}

type ComputeModel string

const (
	ComputeModelECPU ComputeModel = "ECPU"
	ComputeModelOCPU ComputeModel = "OCPU"
)

func PossibleValuesForComputeModel() []string {
	return []string{
		string(ComputeModelECPU),
		string(ComputeModelOCPU),
	}
}

func (s *ComputeModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseComputeModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseComputeModel(input string) (*ComputeModel, error) {
	vals := map[string]ComputeModel{
		"ecpu": ComputeModelECPU,
		"ocpu": ComputeModelOCPU,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ComputeModel(input)
	return &out, nil
}

type DbSystemDatabaseEditionType string

const (
	DbSystemDatabaseEditionTypeEnterpriseEdition                DbSystemDatabaseEditionType = "EnterpriseEdition"
	DbSystemDatabaseEditionTypeEnterpriseEditionDeveloper       DbSystemDatabaseEditionType = "EnterpriseEditionDeveloper"
	DbSystemDatabaseEditionTypeEnterpriseEditionExtreme         DbSystemDatabaseEditionType = "EnterpriseEditionExtreme"
	DbSystemDatabaseEditionTypeEnterpriseEditionHighPerformance DbSystemDatabaseEditionType = "EnterpriseEditionHighPerformance"
	DbSystemDatabaseEditionTypeStandardEdition                  DbSystemDatabaseEditionType = "StandardEdition"
)

func PossibleValuesForDbSystemDatabaseEditionType() []string {
	return []string{
		string(DbSystemDatabaseEditionTypeEnterpriseEdition),
		string(DbSystemDatabaseEditionTypeEnterpriseEditionDeveloper),
		string(DbSystemDatabaseEditionTypeEnterpriseEditionExtreme),
		string(DbSystemDatabaseEditionTypeEnterpriseEditionHighPerformance),
		string(DbSystemDatabaseEditionTypeStandardEdition),
	}
}

func (s *DbSystemDatabaseEditionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbSystemDatabaseEditionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbSystemDatabaseEditionType(input string) (*DbSystemDatabaseEditionType, error) {
	vals := map[string]DbSystemDatabaseEditionType{
		"enterpriseedition":                DbSystemDatabaseEditionTypeEnterpriseEdition,
		"enterpriseeditiondeveloper":       DbSystemDatabaseEditionTypeEnterpriseEditionDeveloper,
		"enterpriseeditionextreme":         DbSystemDatabaseEditionTypeEnterpriseEditionExtreme,
		"enterpriseeditionhighperformance": DbSystemDatabaseEditionTypeEnterpriseEditionHighPerformance,
		"standardedition":                  DbSystemDatabaseEditionTypeStandardEdition,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbSystemDatabaseEditionType(input)
	return &out, nil
}

type DbSystemLifecycleState string

const (
	DbSystemLifecycleStateAvailable             DbSystemLifecycleState = "Available"
	DbSystemLifecycleStateFailed                DbSystemLifecycleState = "Failed"
	DbSystemLifecycleStateMaintenanceInProgress DbSystemLifecycleState = "MaintenanceInProgress"
	DbSystemLifecycleStateMigrated              DbSystemLifecycleState = "Migrated"
	DbSystemLifecycleStateNeedsAttention        DbSystemLifecycleState = "NeedsAttention"
	DbSystemLifecycleStateProvisioning          DbSystemLifecycleState = "Provisioning"
	DbSystemLifecycleStateTerminated            DbSystemLifecycleState = "Terminated"
	DbSystemLifecycleStateTerminating           DbSystemLifecycleState = "Terminating"
	DbSystemLifecycleStateUpdating              DbSystemLifecycleState = "Updating"
	DbSystemLifecycleStateUpgrading             DbSystemLifecycleState = "Upgrading"
)

func PossibleValuesForDbSystemLifecycleState() []string {
	return []string{
		string(DbSystemLifecycleStateAvailable),
		string(DbSystemLifecycleStateFailed),
		string(DbSystemLifecycleStateMaintenanceInProgress),
		string(DbSystemLifecycleStateMigrated),
		string(DbSystemLifecycleStateNeedsAttention),
		string(DbSystemLifecycleStateProvisioning),
		string(DbSystemLifecycleStateTerminated),
		string(DbSystemLifecycleStateTerminating),
		string(DbSystemLifecycleStateUpdating),
		string(DbSystemLifecycleStateUpgrading),
	}
}

func (s *DbSystemLifecycleState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbSystemLifecycleState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbSystemLifecycleState(input string) (*DbSystemLifecycleState, error) {
	vals := map[string]DbSystemLifecycleState{
		"available":             DbSystemLifecycleStateAvailable,
		"failed":                DbSystemLifecycleStateFailed,
		"maintenanceinprogress": DbSystemLifecycleStateMaintenanceInProgress,
		"migrated":              DbSystemLifecycleStateMigrated,
		"needsattention":        DbSystemLifecycleStateNeedsAttention,
		"provisioning":          DbSystemLifecycleStateProvisioning,
		"terminated":            DbSystemLifecycleStateTerminated,
		"terminating":           DbSystemLifecycleStateTerminating,
		"updating":              DbSystemLifecycleStateUpdating,
		"upgrading":             DbSystemLifecycleStateUpgrading,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbSystemLifecycleState(input)
	return &out, nil
}

type DbSystemSourceType string

const (
	DbSystemSourceTypeNone DbSystemSourceType = "None"
)

func PossibleValuesForDbSystemSourceType() []string {
	return []string{
		string(DbSystemSourceTypeNone),
	}
}

func (s *DbSystemSourceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDbSystemSourceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDbSystemSourceType(input string) (*DbSystemSourceType, error) {
	vals := map[string]DbSystemSourceType{
		"none": DbSystemSourceTypeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DbSystemSourceType(input)
	return &out, nil
}

type DiskRedundancyType string

const (
	DiskRedundancyTypeHigh   DiskRedundancyType = "High"
	DiskRedundancyTypeNormal DiskRedundancyType = "Normal"
)

func PossibleValuesForDiskRedundancyType() []string {
	return []string{
		string(DiskRedundancyTypeHigh),
		string(DiskRedundancyTypeNormal),
	}
}

func (s *DiskRedundancyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskRedundancyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskRedundancyType(input string) (*DiskRedundancyType, error) {
	vals := map[string]DiskRedundancyType{
		"high":   DiskRedundancyTypeHigh,
		"normal": DiskRedundancyTypeNormal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskRedundancyType(input)
	return &out, nil
}

type LicenseModel string

const (
	LicenseModelBringYourOwnLicense LicenseModel = "BringYourOwnLicense"
	LicenseModelLicenseIncluded     LicenseModel = "LicenseIncluded"
)

func PossibleValuesForLicenseModel() []string {
	return []string{
		string(LicenseModelBringYourOwnLicense),
		string(LicenseModelLicenseIncluded),
	}
}

func (s *LicenseModel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLicenseModel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLicenseModel(input string) (*LicenseModel, error) {
	vals := map[string]LicenseModel{
		"bringyourownlicense": LicenseModelBringYourOwnLicense,
		"licenseincluded":     LicenseModelLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LicenseModel(input)
	return &out, nil
}

type Source string

const (
	SourceNone Source = "None"
)

func PossibleValuesForSource() []string {
	return []string{
		string(SourceNone),
	}
}

func (s *Source) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSource(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSource(input string) (*Source, error) {
	vals := map[string]Source{
		"none": SourceNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Source(input)
	return &out, nil
}

type StorageManagementType string

const (
	StorageManagementTypeLVM StorageManagementType = "LVM"
)

func PossibleValuesForStorageManagementType() []string {
	return []string{
		string(StorageManagementTypeLVM),
	}
}

func (s *StorageManagementType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageManagementType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageManagementType(input string) (*StorageManagementType, error) {
	vals := map[string]StorageManagementType{
		"lvm": StorageManagementTypeLVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageManagementType(input)
	return &out, nil
}

type StorageVolumePerformanceMode string

const (
	StorageVolumePerformanceModeBalanced        StorageVolumePerformanceMode = "Balanced"
	StorageVolumePerformanceModeHighPerformance StorageVolumePerformanceMode = "HighPerformance"
)

func PossibleValuesForStorageVolumePerformanceMode() []string {
	return []string{
		string(StorageVolumePerformanceModeBalanced),
		string(StorageVolumePerformanceModeHighPerformance),
	}
}

func (s *StorageVolumePerformanceMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageVolumePerformanceMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageVolumePerformanceMode(input string) (*StorageVolumePerformanceMode, error) {
	vals := map[string]StorageVolumePerformanceMode{
		"balanced":        StorageVolumePerformanceModeBalanced,
		"highperformance": StorageVolumePerformanceModeHighPerformance,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageVolumePerformanceMode(input)
	return &out, nil
}

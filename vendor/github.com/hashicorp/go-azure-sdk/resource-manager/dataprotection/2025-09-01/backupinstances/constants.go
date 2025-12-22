package backupinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AKSVolumeTypes string

const (
	AKSVolumeTypesAzureDisk         AKSVolumeTypes = "AzureDisk"
	AKSVolumeTypesAzureFileShareSMB AKSVolumeTypes = "AzureFileShareSMB"
)

func PossibleValuesForAKSVolumeTypes() []string {
	return []string{
		string(AKSVolumeTypesAzureDisk),
		string(AKSVolumeTypesAzureFileShareSMB),
	}
}

func (s *AKSVolumeTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAKSVolumeTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAKSVolumeTypes(input string) (*AKSVolumeTypes, error) {
	vals := map[string]AKSVolumeTypes{
		"azuredisk":         AKSVolumeTypesAzureDisk,
		"azurefilesharesmb": AKSVolumeTypesAzureFileShareSMB,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AKSVolumeTypes(input)
	return &out, nil
}

type CurrentProtectionState string

const (
	CurrentProtectionStateBackupSchedulesSuspended    CurrentProtectionState = "BackupSchedulesSuspended"
	CurrentProtectionStateConfiguringProtection       CurrentProtectionState = "ConfiguringProtection"
	CurrentProtectionStateConfiguringProtectionFailed CurrentProtectionState = "ConfiguringProtectionFailed"
	CurrentProtectionStateInvalid                     CurrentProtectionState = "Invalid"
	CurrentProtectionStateNotProtected                CurrentProtectionState = "NotProtected"
	CurrentProtectionStateProtectionConfigured        CurrentProtectionState = "ProtectionConfigured"
	CurrentProtectionStateProtectionError             CurrentProtectionState = "ProtectionError"
	CurrentProtectionStateProtectionStopped           CurrentProtectionState = "ProtectionStopped"
	CurrentProtectionStateRetentionSchedulesSuspended CurrentProtectionState = "RetentionSchedulesSuspended"
	CurrentProtectionStateSoftDeleted                 CurrentProtectionState = "SoftDeleted"
	CurrentProtectionStateSoftDeleting                CurrentProtectionState = "SoftDeleting"
	CurrentProtectionStateUpdatingProtection          CurrentProtectionState = "UpdatingProtection"
)

func PossibleValuesForCurrentProtectionState() []string {
	return []string{
		string(CurrentProtectionStateBackupSchedulesSuspended),
		string(CurrentProtectionStateConfiguringProtection),
		string(CurrentProtectionStateConfiguringProtectionFailed),
		string(CurrentProtectionStateInvalid),
		string(CurrentProtectionStateNotProtected),
		string(CurrentProtectionStateProtectionConfigured),
		string(CurrentProtectionStateProtectionError),
		string(CurrentProtectionStateProtectionStopped),
		string(CurrentProtectionStateRetentionSchedulesSuspended),
		string(CurrentProtectionStateSoftDeleted),
		string(CurrentProtectionStateSoftDeleting),
		string(CurrentProtectionStateUpdatingProtection),
	}
}

func (s *CurrentProtectionState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCurrentProtectionState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCurrentProtectionState(input string) (*CurrentProtectionState, error) {
	vals := map[string]CurrentProtectionState{
		"backupschedulessuspended":    CurrentProtectionStateBackupSchedulesSuspended,
		"configuringprotection":       CurrentProtectionStateConfiguringProtection,
		"configuringprotectionfailed": CurrentProtectionStateConfiguringProtectionFailed,
		"invalid":                     CurrentProtectionStateInvalid,
		"notprotected":                CurrentProtectionStateNotProtected,
		"protectionconfigured":        CurrentProtectionStateProtectionConfigured,
		"protectionerror":             CurrentProtectionStateProtectionError,
		"protectionstopped":           CurrentProtectionStateProtectionStopped,
		"retentionschedulessuspended": CurrentProtectionStateRetentionSchedulesSuspended,
		"softdeleted":                 CurrentProtectionStateSoftDeleted,
		"softdeleting":                CurrentProtectionStateSoftDeleting,
		"updatingprotection":          CurrentProtectionStateUpdatingProtection,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CurrentProtectionState(input)
	return &out, nil
}

type DataStoreTypes string

const (
	DataStoreTypesArchiveStore     DataStoreTypes = "ArchiveStore"
	DataStoreTypesOperationalStore DataStoreTypes = "OperationalStore"
	DataStoreTypesVaultStore       DataStoreTypes = "VaultStore"
)

func PossibleValuesForDataStoreTypes() []string {
	return []string{
		string(DataStoreTypesArchiveStore),
		string(DataStoreTypesOperationalStore),
		string(DataStoreTypesVaultStore),
	}
}

func (s *DataStoreTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataStoreTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataStoreTypes(input string) (*DataStoreTypes, error) {
	vals := map[string]DataStoreTypes{
		"archivestore":     DataStoreTypesArchiveStore,
		"operationalstore": DataStoreTypesOperationalStore,
		"vaultstore":       DataStoreTypesVaultStore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataStoreTypes(input)
	return &out, nil
}

type ResourcePropertiesObjectType string

const (
	ResourcePropertiesObjectTypeDefaultResourceProperties ResourcePropertiesObjectType = "DefaultResourceProperties"
)

func PossibleValuesForResourcePropertiesObjectType() []string {
	return []string{
		string(ResourcePropertiesObjectTypeDefaultResourceProperties),
	}
}

func (s *ResourcePropertiesObjectType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseResourcePropertiesObjectType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseResourcePropertiesObjectType(input string) (*ResourcePropertiesObjectType, error) {
	vals := map[string]ResourcePropertiesObjectType{
		"defaultresourceproperties": ResourcePropertiesObjectTypeDefaultResourceProperties,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ResourcePropertiesObjectType(input)
	return &out, nil
}

type SecretStoreType string

const (
	SecretStoreTypeAzureKeyVault SecretStoreType = "AzureKeyVault"
	SecretStoreTypeInvalid       SecretStoreType = "Invalid"
)

func PossibleValuesForSecretStoreType() []string {
	return []string{
		string(SecretStoreTypeAzureKeyVault),
		string(SecretStoreTypeInvalid),
	}
}

func (s *SecretStoreType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecretStoreType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecretStoreType(input string) (*SecretStoreType, error) {
	vals := map[string]SecretStoreType{
		"azurekeyvault": SecretStoreTypeAzureKeyVault,
		"invalid":       SecretStoreTypeInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecretStoreType(input)
	return &out, nil
}

type Status string

const (
	StatusConfiguringProtection       Status = "ConfiguringProtection"
	StatusConfiguringProtectionFailed Status = "ConfiguringProtectionFailed"
	StatusProtectionConfigured        Status = "ProtectionConfigured"
	StatusProtectionStopped           Status = "ProtectionStopped"
	StatusSoftDeleted                 Status = "SoftDeleted"
	StatusSoftDeleting                Status = "SoftDeleting"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusConfiguringProtection),
		string(StatusConfiguringProtectionFailed),
		string(StatusProtectionConfigured),
		string(StatusProtectionStopped),
		string(StatusSoftDeleted),
		string(StatusSoftDeleting),
	}
}

func (s *Status) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStatus(input string) (*Status, error) {
	vals := map[string]Status{
		"configuringprotection":       StatusConfiguringProtection,
		"configuringprotectionfailed": StatusConfiguringProtectionFailed,
		"protectionconfigured":        StatusProtectionConfigured,
		"protectionstopped":           StatusProtectionStopped,
		"softdeleted":                 StatusSoftDeleted,
		"softdeleting":                StatusSoftDeleting,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type ValidationType string

const (
	ValidationTypeDeepValidation    ValidationType = "DeepValidation"
	ValidationTypeShallowValidation ValidationType = "ShallowValidation"
)

func PossibleValuesForValidationType() []string {
	return []string{
		string(ValidationTypeDeepValidation),
		string(ValidationTypeShallowValidation),
	}
}

func (s *ValidationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseValidationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseValidationType(input string) (*ValidationType, error) {
	vals := map[string]ValidationType{
		"deepvalidation":    ValidationTypeDeepValidation,
		"shallowvalidation": ValidationTypeShallowValidation,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ValidationType(input)
	return &out, nil
}

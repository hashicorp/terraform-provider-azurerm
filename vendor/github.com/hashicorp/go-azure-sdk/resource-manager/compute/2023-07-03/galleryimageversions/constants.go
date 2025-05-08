package galleryimageversions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AggregatedReplicationState string

const (
	AggregatedReplicationStateCompleted  AggregatedReplicationState = "Completed"
	AggregatedReplicationStateFailed     AggregatedReplicationState = "Failed"
	AggregatedReplicationStateInProgress AggregatedReplicationState = "InProgress"
	AggregatedReplicationStateUnknown    AggregatedReplicationState = "Unknown"
)

func PossibleValuesForAggregatedReplicationState() []string {
	return []string{
		string(AggregatedReplicationStateCompleted),
		string(AggregatedReplicationStateFailed),
		string(AggregatedReplicationStateInProgress),
		string(AggregatedReplicationStateUnknown),
	}
}

func (s *AggregatedReplicationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAggregatedReplicationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAggregatedReplicationState(input string) (*AggregatedReplicationState, error) {
	vals := map[string]AggregatedReplicationState{
		"completed":  AggregatedReplicationStateCompleted,
		"failed":     AggregatedReplicationStateFailed,
		"inprogress": AggregatedReplicationStateInProgress,
		"unknown":    AggregatedReplicationStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AggregatedReplicationState(input)
	return &out, nil
}

type ConfidentialVMEncryptionType string

const (
	ConfidentialVMEncryptionTypeEncryptedVMGuestStateOnlyWithPmk ConfidentialVMEncryptionType = "EncryptedVMGuestStateOnlyWithPmk"
	ConfidentialVMEncryptionTypeEncryptedWithCmk                 ConfidentialVMEncryptionType = "EncryptedWithCmk"
	ConfidentialVMEncryptionTypeEncryptedWithPmk                 ConfidentialVMEncryptionType = "EncryptedWithPmk"
	ConfidentialVMEncryptionTypeNonPersistedTPM                  ConfidentialVMEncryptionType = "NonPersistedTPM"
)

func PossibleValuesForConfidentialVMEncryptionType() []string {
	return []string{
		string(ConfidentialVMEncryptionTypeEncryptedVMGuestStateOnlyWithPmk),
		string(ConfidentialVMEncryptionTypeEncryptedWithCmk),
		string(ConfidentialVMEncryptionTypeEncryptedWithPmk),
		string(ConfidentialVMEncryptionTypeNonPersistedTPM),
	}
}

func (s *ConfidentialVMEncryptionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConfidentialVMEncryptionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConfidentialVMEncryptionType(input string) (*ConfidentialVMEncryptionType, error) {
	vals := map[string]ConfidentialVMEncryptionType{
		"encryptedvmgueststateonlywithpmk": ConfidentialVMEncryptionTypeEncryptedVMGuestStateOnlyWithPmk,
		"encryptedwithcmk":                 ConfidentialVMEncryptionTypeEncryptedWithCmk,
		"encryptedwithpmk":                 ConfidentialVMEncryptionTypeEncryptedWithPmk,
		"nonpersistedtpm":                  ConfidentialVMEncryptionTypeNonPersistedTPM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConfidentialVMEncryptionType(input)
	return &out, nil
}

type EdgeZoneStorageAccountType string

const (
	EdgeZoneStorageAccountTypePremiumLRS     EdgeZoneStorageAccountType = "Premium_LRS"
	EdgeZoneStorageAccountTypeStandardLRS    EdgeZoneStorageAccountType = "Standard_LRS"
	EdgeZoneStorageAccountTypeStandardSSDLRS EdgeZoneStorageAccountType = "StandardSSD_LRS"
	EdgeZoneStorageAccountTypeStandardZRS    EdgeZoneStorageAccountType = "Standard_ZRS"
)

func PossibleValuesForEdgeZoneStorageAccountType() []string {
	return []string{
		string(EdgeZoneStorageAccountTypePremiumLRS),
		string(EdgeZoneStorageAccountTypeStandardLRS),
		string(EdgeZoneStorageAccountTypeStandardSSDLRS),
		string(EdgeZoneStorageAccountTypeStandardZRS),
	}
}

func (s *EdgeZoneStorageAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEdgeZoneStorageAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEdgeZoneStorageAccountType(input string) (*EdgeZoneStorageAccountType, error) {
	vals := map[string]EdgeZoneStorageAccountType{
		"premium_lrs":     EdgeZoneStorageAccountTypePremiumLRS,
		"standard_lrs":    EdgeZoneStorageAccountTypeStandardLRS,
		"standardssd_lrs": EdgeZoneStorageAccountTypeStandardSSDLRS,
		"standard_zrs":    EdgeZoneStorageAccountTypeStandardZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := EdgeZoneStorageAccountType(input)
	return &out, nil
}

type GalleryExtendedLocationType string

const (
	GalleryExtendedLocationTypeEdgeZone GalleryExtendedLocationType = "EdgeZone"
	GalleryExtendedLocationTypeUnknown  GalleryExtendedLocationType = "Unknown"
)

func PossibleValuesForGalleryExtendedLocationType() []string {
	return []string{
		string(GalleryExtendedLocationTypeEdgeZone),
		string(GalleryExtendedLocationTypeUnknown),
	}
}

func (s *GalleryExtendedLocationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGalleryExtendedLocationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGalleryExtendedLocationType(input string) (*GalleryExtendedLocationType, error) {
	vals := map[string]GalleryExtendedLocationType{
		"edgezone": GalleryExtendedLocationTypeEdgeZone,
		"unknown":  GalleryExtendedLocationTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GalleryExtendedLocationType(input)
	return &out, nil
}

type GalleryProvisioningState string

const (
	GalleryProvisioningStateCreating  GalleryProvisioningState = "Creating"
	GalleryProvisioningStateDeleting  GalleryProvisioningState = "Deleting"
	GalleryProvisioningStateFailed    GalleryProvisioningState = "Failed"
	GalleryProvisioningStateMigrating GalleryProvisioningState = "Migrating"
	GalleryProvisioningStateSucceeded GalleryProvisioningState = "Succeeded"
	GalleryProvisioningStateUpdating  GalleryProvisioningState = "Updating"
)

func PossibleValuesForGalleryProvisioningState() []string {
	return []string{
		string(GalleryProvisioningStateCreating),
		string(GalleryProvisioningStateDeleting),
		string(GalleryProvisioningStateFailed),
		string(GalleryProvisioningStateMigrating),
		string(GalleryProvisioningStateSucceeded),
		string(GalleryProvisioningStateUpdating),
	}
}

func (s *GalleryProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseGalleryProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseGalleryProvisioningState(input string) (*GalleryProvisioningState, error) {
	vals := map[string]GalleryProvisioningState{
		"creating":  GalleryProvisioningStateCreating,
		"deleting":  GalleryProvisioningStateDeleting,
		"failed":    GalleryProvisioningStateFailed,
		"migrating": GalleryProvisioningStateMigrating,
		"succeeded": GalleryProvisioningStateSucceeded,
		"updating":  GalleryProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := GalleryProvisioningState(input)
	return &out, nil
}

type HostCaching string

const (
	HostCachingNone      HostCaching = "None"
	HostCachingReadOnly  HostCaching = "ReadOnly"
	HostCachingReadWrite HostCaching = "ReadWrite"
)

func PossibleValuesForHostCaching() []string {
	return []string{
		string(HostCachingNone),
		string(HostCachingReadOnly),
		string(HostCachingReadWrite),
	}
}

func (s *HostCaching) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHostCaching(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHostCaching(input string) (*HostCaching, error) {
	vals := map[string]HostCaching{
		"none":      HostCachingNone,
		"readonly":  HostCachingReadOnly,
		"readwrite": HostCachingReadWrite,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HostCaching(input)
	return &out, nil
}

type PolicyViolationCategory string

const (
	PolicyViolationCategoryCopyrightValidation PolicyViolationCategory = "CopyrightValidation"
	PolicyViolationCategoryIPTheft             PolicyViolationCategory = "IpTheft"
	PolicyViolationCategoryImageFlaggedUnsafe  PolicyViolationCategory = "ImageFlaggedUnsafe"
	PolicyViolationCategoryOther               PolicyViolationCategory = "Other"
)

func PossibleValuesForPolicyViolationCategory() []string {
	return []string{
		string(PolicyViolationCategoryCopyrightValidation),
		string(PolicyViolationCategoryIPTheft),
		string(PolicyViolationCategoryImageFlaggedUnsafe),
		string(PolicyViolationCategoryOther),
	}
}

func (s *PolicyViolationCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePolicyViolationCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePolicyViolationCategory(input string) (*PolicyViolationCategory, error) {
	vals := map[string]PolicyViolationCategory{
		"copyrightvalidation": PolicyViolationCategoryCopyrightValidation,
		"iptheft":             PolicyViolationCategoryIPTheft,
		"imageflaggedunsafe":  PolicyViolationCategoryImageFlaggedUnsafe,
		"other":               PolicyViolationCategoryOther,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyViolationCategory(input)
	return &out, nil
}

type ReplicationMode string

const (
	ReplicationModeFull    ReplicationMode = "Full"
	ReplicationModeShallow ReplicationMode = "Shallow"
)

func PossibleValuesForReplicationMode() []string {
	return []string{
		string(ReplicationModeFull),
		string(ReplicationModeShallow),
	}
}

func (s *ReplicationMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationMode(input string) (*ReplicationMode, error) {
	vals := map[string]ReplicationMode{
		"full":    ReplicationModeFull,
		"shallow": ReplicationModeShallow,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationMode(input)
	return &out, nil
}

type ReplicationState string

const (
	ReplicationStateCompleted   ReplicationState = "Completed"
	ReplicationStateFailed      ReplicationState = "Failed"
	ReplicationStateReplicating ReplicationState = "Replicating"
	ReplicationStateUnknown     ReplicationState = "Unknown"
)

func PossibleValuesForReplicationState() []string {
	return []string{
		string(ReplicationStateCompleted),
		string(ReplicationStateFailed),
		string(ReplicationStateReplicating),
		string(ReplicationStateUnknown),
	}
}

func (s *ReplicationState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationState(input string) (*ReplicationState, error) {
	vals := map[string]ReplicationState{
		"completed":   ReplicationStateCompleted,
		"failed":      ReplicationStateFailed,
		"replicating": ReplicationStateReplicating,
		"unknown":     ReplicationStateUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationState(input)
	return &out, nil
}

type ReplicationStatusTypes string

const (
	ReplicationStatusTypesReplicationStatus ReplicationStatusTypes = "ReplicationStatus"
	ReplicationStatusTypesUefiSettings      ReplicationStatusTypes = "UefiSettings"
)

func PossibleValuesForReplicationStatusTypes() []string {
	return []string{
		string(ReplicationStatusTypesReplicationStatus),
		string(ReplicationStatusTypesUefiSettings),
	}
}

func (s *ReplicationStatusTypes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicationStatusTypes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicationStatusTypes(input string) (*ReplicationStatusTypes, error) {
	vals := map[string]ReplicationStatusTypes{
		"replicationstatus": ReplicationStatusTypesReplicationStatus,
		"uefisettings":      ReplicationStatusTypesUefiSettings,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicationStatusTypes(input)
	return &out, nil
}

type StorageAccountType string

const (
	StorageAccountTypePremiumLRS  StorageAccountType = "Premium_LRS"
	StorageAccountTypeStandardLRS StorageAccountType = "Standard_LRS"
	StorageAccountTypeStandardZRS StorageAccountType = "Standard_ZRS"
)

func PossibleValuesForStorageAccountType() []string {
	return []string{
		string(StorageAccountTypePremiumLRS),
		string(StorageAccountTypeStandardLRS),
		string(StorageAccountTypeStandardZRS),
	}
}

func (s *StorageAccountType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageAccountType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageAccountType(input string) (*StorageAccountType, error) {
	vals := map[string]StorageAccountType{
		"premium_lrs":  StorageAccountTypePremiumLRS,
		"standard_lrs": StorageAccountTypeStandardLRS,
		"standard_zrs": StorageAccountTypeStandardZRS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageAccountType(input)
	return &out, nil
}

type UefiKeyType string

const (
	UefiKeyTypeShaTwoFiveSix UefiKeyType = "sha256"
	UefiKeyTypeXFiveZeroNine UefiKeyType = "x509"
)

func PossibleValuesForUefiKeyType() []string {
	return []string{
		string(UefiKeyTypeShaTwoFiveSix),
		string(UefiKeyTypeXFiveZeroNine),
	}
}

func (s *UefiKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUefiKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUefiKeyType(input string) (*UefiKeyType, error) {
	vals := map[string]UefiKeyType{
		"sha256": UefiKeyTypeShaTwoFiveSix,
		"x509":   UefiKeyTypeXFiveZeroNine,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UefiKeyType(input)
	return &out, nil
}

type UefiSignatureTemplateName string

const (
	UefiSignatureTemplateNameMicrosoftUefiCertificateAuthorityTemplate UefiSignatureTemplateName = "MicrosoftUefiCertificateAuthorityTemplate"
	UefiSignatureTemplateNameMicrosoftWindowsTemplate                  UefiSignatureTemplateName = "MicrosoftWindowsTemplate"
	UefiSignatureTemplateNameNoSignatureTemplate                       UefiSignatureTemplateName = "NoSignatureTemplate"
)

func PossibleValuesForUefiSignatureTemplateName() []string {
	return []string{
		string(UefiSignatureTemplateNameMicrosoftUefiCertificateAuthorityTemplate),
		string(UefiSignatureTemplateNameMicrosoftWindowsTemplate),
		string(UefiSignatureTemplateNameNoSignatureTemplate),
	}
}

func (s *UefiSignatureTemplateName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUefiSignatureTemplateName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUefiSignatureTemplateName(input string) (*UefiSignatureTemplateName, error) {
	vals := map[string]UefiSignatureTemplateName{
		"microsoftueficertificateauthoritytemplate": UefiSignatureTemplateNameMicrosoftUefiCertificateAuthorityTemplate,
		"microsoftwindowstemplate":                  UefiSignatureTemplateNameMicrosoftWindowsTemplate,
		"nosignaturetemplate":                       UefiSignatureTemplateNameNoSignatureTemplate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UefiSignatureTemplateName(input)
	return &out, nil
}

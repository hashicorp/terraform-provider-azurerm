package cosmosdb

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalyticalStorageSchemaType string

const (
	AnalyticalStorageSchemaTypeFullFidelity AnalyticalStorageSchemaType = "FullFidelity"
	AnalyticalStorageSchemaTypeWellDefined  AnalyticalStorageSchemaType = "WellDefined"
)

func PossibleValuesForAnalyticalStorageSchemaType() []string {
	return []string{
		string(AnalyticalStorageSchemaTypeFullFidelity),
		string(AnalyticalStorageSchemaTypeWellDefined),
	}
}

func (s *AnalyticalStorageSchemaType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAnalyticalStorageSchemaType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAnalyticalStorageSchemaType(input string) (*AnalyticalStorageSchemaType, error) {
	vals := map[string]AnalyticalStorageSchemaType{
		"fullfidelity": AnalyticalStorageSchemaTypeFullFidelity,
		"welldefined":  AnalyticalStorageSchemaTypeWellDefined,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AnalyticalStorageSchemaType(input)
	return &out, nil
}

type BackupPolicyMigrationStatus string

const (
	BackupPolicyMigrationStatusCompleted  BackupPolicyMigrationStatus = "Completed"
	BackupPolicyMigrationStatusFailed     BackupPolicyMigrationStatus = "Failed"
	BackupPolicyMigrationStatusInProgress BackupPolicyMigrationStatus = "InProgress"
	BackupPolicyMigrationStatusInvalid    BackupPolicyMigrationStatus = "Invalid"
)

func PossibleValuesForBackupPolicyMigrationStatus() []string {
	return []string{
		string(BackupPolicyMigrationStatusCompleted),
		string(BackupPolicyMigrationStatusFailed),
		string(BackupPolicyMigrationStatusInProgress),
		string(BackupPolicyMigrationStatusInvalid),
	}
}

func (s *BackupPolicyMigrationStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupPolicyMigrationStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupPolicyMigrationStatus(input string) (*BackupPolicyMigrationStatus, error) {
	vals := map[string]BackupPolicyMigrationStatus{
		"completed":  BackupPolicyMigrationStatusCompleted,
		"failed":     BackupPolicyMigrationStatusFailed,
		"inprogress": BackupPolicyMigrationStatusInProgress,
		"invalid":    BackupPolicyMigrationStatusInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupPolicyMigrationStatus(input)
	return &out, nil
}

type BackupPolicyType string

const (
	BackupPolicyTypeContinuous BackupPolicyType = "Continuous"
	BackupPolicyTypePeriodic   BackupPolicyType = "Periodic"
)

func PossibleValuesForBackupPolicyType() []string {
	return []string{
		string(BackupPolicyTypeContinuous),
		string(BackupPolicyTypePeriodic),
	}
}

func (s *BackupPolicyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupPolicyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupPolicyType(input string) (*BackupPolicyType, error) {
	vals := map[string]BackupPolicyType{
		"continuous": BackupPolicyTypeContinuous,
		"periodic":   BackupPolicyTypePeriodic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupPolicyType(input)
	return &out, nil
}

type BackupStorageRedundancy string

const (
	BackupStorageRedundancyGeo   BackupStorageRedundancy = "Geo"
	BackupStorageRedundancyLocal BackupStorageRedundancy = "Local"
	BackupStorageRedundancyZone  BackupStorageRedundancy = "Zone"
)

func PossibleValuesForBackupStorageRedundancy() []string {
	return []string{
		string(BackupStorageRedundancyGeo),
		string(BackupStorageRedundancyLocal),
		string(BackupStorageRedundancyZone),
	}
}

func (s *BackupStorageRedundancy) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupStorageRedundancy(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupStorageRedundancy(input string) (*BackupStorageRedundancy, error) {
	vals := map[string]BackupStorageRedundancy{
		"geo":   BackupStorageRedundancyGeo,
		"local": BackupStorageRedundancyLocal,
		"zone":  BackupStorageRedundancyZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupStorageRedundancy(input)
	return &out, nil
}

type CompositePathSortOrder string

const (
	CompositePathSortOrderAscending  CompositePathSortOrder = "ascending"
	CompositePathSortOrderDescending CompositePathSortOrder = "descending"
)

func PossibleValuesForCompositePathSortOrder() []string {
	return []string{
		string(CompositePathSortOrderAscending),
		string(CompositePathSortOrderDescending),
	}
}

func (s *CompositePathSortOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCompositePathSortOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCompositePathSortOrder(input string) (*CompositePathSortOrder, error) {
	vals := map[string]CompositePathSortOrder{
		"ascending":  CompositePathSortOrderAscending,
		"descending": CompositePathSortOrderDescending,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CompositePathSortOrder(input)
	return &out, nil
}

type ConflictResolutionMode string

const (
	ConflictResolutionModeCustom         ConflictResolutionMode = "Custom"
	ConflictResolutionModeLastWriterWins ConflictResolutionMode = "LastWriterWins"
)

func PossibleValuesForConflictResolutionMode() []string {
	return []string{
		string(ConflictResolutionModeCustom),
		string(ConflictResolutionModeLastWriterWins),
	}
}

func (s *ConflictResolutionMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConflictResolutionMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConflictResolutionMode(input string) (*ConflictResolutionMode, error) {
	vals := map[string]ConflictResolutionMode{
		"custom":         ConflictResolutionModeCustom,
		"lastwriterwins": ConflictResolutionModeLastWriterWins,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConflictResolutionMode(input)
	return &out, nil
}

type ConnectorOffer string

const (
	ConnectorOfferSmall ConnectorOffer = "Small"
)

func PossibleValuesForConnectorOffer() []string {
	return []string{
		string(ConnectorOfferSmall),
	}
}

func (s *ConnectorOffer) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectorOffer(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectorOffer(input string) (*ConnectorOffer, error) {
	vals := map[string]ConnectorOffer{
		"small": ConnectorOfferSmall,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectorOffer(input)
	return &out, nil
}

type ContinuousTier string

const (
	ContinuousTierContinuousSevenDays     ContinuousTier = "Continuous7Days"
	ContinuousTierContinuousThreeZeroDays ContinuousTier = "Continuous30Days"
)

func PossibleValuesForContinuousTier() []string {
	return []string{
		string(ContinuousTierContinuousSevenDays),
		string(ContinuousTierContinuousThreeZeroDays),
	}
}

func (s *ContinuousTier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseContinuousTier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseContinuousTier(input string) (*ContinuousTier, error) {
	vals := map[string]ContinuousTier{
		"continuous7days":  ContinuousTierContinuousSevenDays,
		"continuous30days": ContinuousTierContinuousThreeZeroDays,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ContinuousTier(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeDefault CreateMode = "Default"
	CreateModeRestore CreateMode = "Restore"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeDefault),
		string(CreateModeRestore),
	}
}

func (s *CreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCreateMode(input string) (*CreateMode, error) {
	vals := map[string]CreateMode{
		"default": CreateModeDefault,
		"restore": CreateModeRestore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type DataType string

const (
	DataTypeLineString   DataType = "LineString"
	DataTypeMultiPolygon DataType = "MultiPolygon"
	DataTypeNumber       DataType = "Number"
	DataTypePoint        DataType = "Point"
	DataTypePolygon      DataType = "Polygon"
	DataTypeString       DataType = "String"
)

func PossibleValuesForDataType() []string {
	return []string{
		string(DataTypeLineString),
		string(DataTypeMultiPolygon),
		string(DataTypeNumber),
		string(DataTypePoint),
		string(DataTypePolygon),
		string(DataTypeString),
	}
}

func (s *DataType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataType(input string) (*DataType, error) {
	vals := map[string]DataType{
		"linestring":   DataTypeLineString,
		"multipolygon": DataTypeMultiPolygon,
		"number":       DataTypeNumber,
		"point":        DataTypePoint,
		"polygon":      DataTypePolygon,
		"string":       DataTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataType(input)
	return &out, nil
}

type DatabaseAccountKind string

const (
	DatabaseAccountKindGlobalDocumentDB DatabaseAccountKind = "GlobalDocumentDB"
	DatabaseAccountKindMongoDB          DatabaseAccountKind = "MongoDB"
	DatabaseAccountKindParse            DatabaseAccountKind = "Parse"
)

func PossibleValuesForDatabaseAccountKind() []string {
	return []string{
		string(DatabaseAccountKindGlobalDocumentDB),
		string(DatabaseAccountKindMongoDB),
		string(DatabaseAccountKindParse),
	}
}

func (s *DatabaseAccountKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseAccountKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseAccountKind(input string) (*DatabaseAccountKind, error) {
	vals := map[string]DatabaseAccountKind{
		"globaldocumentdb": DatabaseAccountKindGlobalDocumentDB,
		"mongodb":          DatabaseAccountKindMongoDB,
		"parse":            DatabaseAccountKindParse,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseAccountKind(input)
	return &out, nil
}

type DatabaseAccountOfferType string

const (
	DatabaseAccountOfferTypeStandard DatabaseAccountOfferType = "Standard"
)

func PossibleValuesForDatabaseAccountOfferType() []string {
	return []string{
		string(DatabaseAccountOfferTypeStandard),
	}
}

func (s *DatabaseAccountOfferType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseAccountOfferType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseAccountOfferType(input string) (*DatabaseAccountOfferType, error) {
	vals := map[string]DatabaseAccountOfferType{
		"standard": DatabaseAccountOfferTypeStandard,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseAccountOfferType(input)
	return &out, nil
}

type DefaultConsistencyLevel string

const (
	DefaultConsistencyLevelBoundedStaleness DefaultConsistencyLevel = "BoundedStaleness"
	DefaultConsistencyLevelConsistentPrefix DefaultConsistencyLevel = "ConsistentPrefix"
	DefaultConsistencyLevelEventual         DefaultConsistencyLevel = "Eventual"
	DefaultConsistencyLevelSession          DefaultConsistencyLevel = "Session"
	DefaultConsistencyLevelStrong           DefaultConsistencyLevel = "Strong"
)

func PossibleValuesForDefaultConsistencyLevel() []string {
	return []string{
		string(DefaultConsistencyLevelBoundedStaleness),
		string(DefaultConsistencyLevelConsistentPrefix),
		string(DefaultConsistencyLevelEventual),
		string(DefaultConsistencyLevelSession),
		string(DefaultConsistencyLevelStrong),
	}
}

func (s *DefaultConsistencyLevel) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDefaultConsistencyLevel(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDefaultConsistencyLevel(input string) (*DefaultConsistencyLevel, error) {
	vals := map[string]DefaultConsistencyLevel{
		"boundedstaleness": DefaultConsistencyLevelBoundedStaleness,
		"consistentprefix": DefaultConsistencyLevelConsistentPrefix,
		"eventual":         DefaultConsistencyLevelEventual,
		"session":          DefaultConsistencyLevelSession,
		"strong":           DefaultConsistencyLevelStrong,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DefaultConsistencyLevel(input)
	return &out, nil
}

type IndexKind string

const (
	IndexKindHash    IndexKind = "Hash"
	IndexKindRange   IndexKind = "Range"
	IndexKindSpatial IndexKind = "Spatial"
)

func PossibleValuesForIndexKind() []string {
	return []string{
		string(IndexKindHash),
		string(IndexKindRange),
		string(IndexKindSpatial),
	}
}

func (s *IndexKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexKind(input string) (*IndexKind, error) {
	vals := map[string]IndexKind{
		"hash":    IndexKindHash,
		"range":   IndexKindRange,
		"spatial": IndexKindSpatial,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexKind(input)
	return &out, nil
}

type IndexingMode string

const (
	IndexingModeConsistent IndexingMode = "consistent"
	IndexingModeLazy       IndexingMode = "lazy"
	IndexingModeNone       IndexingMode = "none"
)

func PossibleValuesForIndexingMode() []string {
	return []string{
		string(IndexingModeConsistent),
		string(IndexingModeLazy),
		string(IndexingModeNone),
	}
}

func (s *IndexingMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIndexingMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIndexingMode(input string) (*IndexingMode, error) {
	vals := map[string]IndexingMode{
		"consistent": IndexingModeConsistent,
		"lazy":       IndexingModeLazy,
		"none":       IndexingModeNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IndexingMode(input)
	return &out, nil
}

type KeyKind string

const (
	KeyKindPrimary           KeyKind = "primary"
	KeyKindPrimaryReadonly   KeyKind = "primaryReadonly"
	KeyKindSecondary         KeyKind = "secondary"
	KeyKindSecondaryReadonly KeyKind = "secondaryReadonly"
)

func PossibleValuesForKeyKind() []string {
	return []string{
		string(KeyKindPrimary),
		string(KeyKindPrimaryReadonly),
		string(KeyKindSecondary),
		string(KeyKindSecondaryReadonly),
	}
}

func (s *KeyKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKeyKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKeyKind(input string) (*KeyKind, error) {
	vals := map[string]KeyKind{
		"primary":           KeyKindPrimary,
		"primaryreadonly":   KeyKindPrimaryReadonly,
		"secondary":         KeyKindSecondary,
		"secondaryreadonly": KeyKindSecondaryReadonly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KeyKind(input)
	return &out, nil
}

type Kind string

const (
	KindPrimary           Kind = "Primary"
	KindPrimaryReadonly   Kind = "PrimaryReadonly"
	KindSecondary         Kind = "Secondary"
	KindSecondaryReadonly Kind = "SecondaryReadonly"
)

func PossibleValuesForKind() []string {
	return []string{
		string(KindPrimary),
		string(KindPrimaryReadonly),
		string(KindSecondary),
		string(KindSecondaryReadonly),
	}
}

func (s *Kind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKind(input string) (*Kind, error) {
	vals := map[string]Kind{
		"primary":           KindPrimary,
		"primaryreadonly":   KindPrimaryReadonly,
		"secondary":         KindSecondary,
		"secondaryreadonly": KindSecondaryReadonly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Kind(input)
	return &out, nil
}

type MinimalTlsVersion string

const (
	MinimalTlsVersionTls       MinimalTlsVersion = "Tls"
	MinimalTlsVersionTlsOneOne MinimalTlsVersion = "Tls11"
	MinimalTlsVersionTlsOneTwo MinimalTlsVersion = "Tls12"
)

func PossibleValuesForMinimalTlsVersion() []string {
	return []string{
		string(MinimalTlsVersionTls),
		string(MinimalTlsVersionTlsOneOne),
		string(MinimalTlsVersionTlsOneTwo),
	}
}

func (s *MinimalTlsVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMinimalTlsVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMinimalTlsVersion(input string) (*MinimalTlsVersion, error) {
	vals := map[string]MinimalTlsVersion{
		"tls":   MinimalTlsVersionTls,
		"tls11": MinimalTlsVersionTlsOneOne,
		"tls12": MinimalTlsVersionTlsOneTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MinimalTlsVersion(input)
	return &out, nil
}

type NetworkAclBypass string

const (
	NetworkAclBypassAzureServices NetworkAclBypass = "AzureServices"
	NetworkAclBypassNone          NetworkAclBypass = "None"
)

func PossibleValuesForNetworkAclBypass() []string {
	return []string{
		string(NetworkAclBypassAzureServices),
		string(NetworkAclBypassNone),
	}
}

func (s *NetworkAclBypass) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseNetworkAclBypass(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseNetworkAclBypass(input string) (*NetworkAclBypass, error) {
	vals := map[string]NetworkAclBypass{
		"azureservices": NetworkAclBypassAzureServices,
		"none":          NetworkAclBypassNone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := NetworkAclBypass(input)
	return &out, nil
}

type PartitionKind string

const (
	PartitionKindHash      PartitionKind = "Hash"
	PartitionKindMultiHash PartitionKind = "MultiHash"
	PartitionKindRange     PartitionKind = "Range"
)

func PossibleValuesForPartitionKind() []string {
	return []string{
		string(PartitionKindHash),
		string(PartitionKindMultiHash),
		string(PartitionKindRange),
	}
}

func (s *PartitionKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePartitionKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePartitionKind(input string) (*PartitionKind, error) {
	vals := map[string]PartitionKind{
		"hash":      PartitionKindHash,
		"multihash": PartitionKindMultiHash,
		"range":     PartitionKindRange,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PartitionKind(input)
	return &out, nil
}

type PrimaryAggregationType string

const (
	PrimaryAggregationTypeAverage PrimaryAggregationType = "Average"
	PrimaryAggregationTypeLast    PrimaryAggregationType = "Last"
	PrimaryAggregationTypeMaximum PrimaryAggregationType = "Maximum"
	PrimaryAggregationTypeMinimum PrimaryAggregationType = "Minimum"
	PrimaryAggregationTypeNone    PrimaryAggregationType = "None"
	PrimaryAggregationTypeTotal   PrimaryAggregationType = "Total"
)

func PossibleValuesForPrimaryAggregationType() []string {
	return []string{
		string(PrimaryAggregationTypeAverage),
		string(PrimaryAggregationTypeLast),
		string(PrimaryAggregationTypeMaximum),
		string(PrimaryAggregationTypeMinimum),
		string(PrimaryAggregationTypeNone),
		string(PrimaryAggregationTypeTotal),
	}
}

func (s *PrimaryAggregationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrimaryAggregationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrimaryAggregationType(input string) (*PrimaryAggregationType, error) {
	vals := map[string]PrimaryAggregationType{
		"average": PrimaryAggregationTypeAverage,
		"last":    PrimaryAggregationTypeLast,
		"maximum": PrimaryAggregationTypeMaximum,
		"minimum": PrimaryAggregationTypeMinimum,
		"none":    PrimaryAggregationTypeNone,
		"total":   PrimaryAggregationTypeTotal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrimaryAggregationType(input)
	return &out, nil
}

type PublicNetworkAccess string

const (
	PublicNetworkAccessDisabled           PublicNetworkAccess = "Disabled"
	PublicNetworkAccessEnabled            PublicNetworkAccess = "Enabled"
	PublicNetworkAccessSecuredByPerimeter PublicNetworkAccess = "SecuredByPerimeter"
)

func PossibleValuesForPublicNetworkAccess() []string {
	return []string{
		string(PublicNetworkAccessDisabled),
		string(PublicNetworkAccessEnabled),
		string(PublicNetworkAccessSecuredByPerimeter),
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
		"disabled":           PublicNetworkAccessDisabled,
		"enabled":            PublicNetworkAccessEnabled,
		"securedbyperimeter": PublicNetworkAccessSecuredByPerimeter,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PublicNetworkAccess(input)
	return &out, nil
}

type RestoreMode string

const (
	RestoreModePointInTime RestoreMode = "PointInTime"
)

func PossibleValuesForRestoreMode() []string {
	return []string{
		string(RestoreModePointInTime),
	}
}

func (s *RestoreMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRestoreMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRestoreMode(input string) (*RestoreMode, error) {
	vals := map[string]RestoreMode{
		"pointintime": RestoreModePointInTime,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RestoreMode(input)
	return &out, nil
}

type ServerVersion string

const (
	ServerVersionFivePointZero  ServerVersion = "5.0"
	ServerVersionFourPointTwo   ServerVersion = "4.2"
	ServerVersionFourPointZero  ServerVersion = "4.0"
	ServerVersionSevenPointZero ServerVersion = "7.0"
	ServerVersionSixPointZero   ServerVersion = "6.0"
	ServerVersionThreePointSix  ServerVersion = "3.6"
	ServerVersionThreePointTwo  ServerVersion = "3.2"
)

func PossibleValuesForServerVersion() []string {
	return []string{
		string(ServerVersionFivePointZero),
		string(ServerVersionFourPointTwo),
		string(ServerVersionFourPointZero),
		string(ServerVersionSevenPointZero),
		string(ServerVersionSixPointZero),
		string(ServerVersionThreePointSix),
		string(ServerVersionThreePointTwo),
	}
}

func (s *ServerVersion) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServerVersion(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServerVersion(input string) (*ServerVersion, error) {
	vals := map[string]ServerVersion{
		"5.0": ServerVersionFivePointZero,
		"4.2": ServerVersionFourPointTwo,
		"4.0": ServerVersionFourPointZero,
		"7.0": ServerVersionSevenPointZero,
		"6.0": ServerVersionSixPointZero,
		"3.6": ServerVersionThreePointSix,
		"3.2": ServerVersionThreePointTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServerVersion(input)
	return &out, nil
}

type SpatialType string

const (
	SpatialTypeLineString   SpatialType = "LineString"
	SpatialTypeMultiPolygon SpatialType = "MultiPolygon"
	SpatialTypePoint        SpatialType = "Point"
	SpatialTypePolygon      SpatialType = "Polygon"
)

func PossibleValuesForSpatialType() []string {
	return []string{
		string(SpatialTypeLineString),
		string(SpatialTypeMultiPolygon),
		string(SpatialTypePoint),
		string(SpatialTypePolygon),
	}
}

func (s *SpatialType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSpatialType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSpatialType(input string) (*SpatialType, error) {
	vals := map[string]SpatialType{
		"linestring":   SpatialTypeLineString,
		"multipolygon": SpatialTypeMultiPolygon,
		"point":        SpatialTypePoint,
		"polygon":      SpatialTypePolygon,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SpatialType(input)
	return &out, nil
}

type Status string

const (
	StatusDeleting        Status = "Deleting"
	StatusInitializing    Status = "Initializing"
	StatusInternallyReady Status = "InternallyReady"
	StatusOnline          Status = "Online"
	StatusUninitialized   Status = "Uninitialized"
)

func PossibleValuesForStatus() []string {
	return []string{
		string(StatusDeleting),
		string(StatusInitializing),
		string(StatusInternallyReady),
		string(StatusOnline),
		string(StatusUninitialized),
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
		"deleting":        StatusDeleting,
		"initializing":    StatusInitializing,
		"internallyready": StatusInternallyReady,
		"online":          StatusOnline,
		"uninitialized":   StatusUninitialized,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Status(input)
	return &out, nil
}

type TriggerOperation string

const (
	TriggerOperationAll     TriggerOperation = "All"
	TriggerOperationCreate  TriggerOperation = "Create"
	TriggerOperationDelete  TriggerOperation = "Delete"
	TriggerOperationReplace TriggerOperation = "Replace"
	TriggerOperationUpdate  TriggerOperation = "Update"
)

func PossibleValuesForTriggerOperation() []string {
	return []string{
		string(TriggerOperationAll),
		string(TriggerOperationCreate),
		string(TriggerOperationDelete),
		string(TriggerOperationReplace),
		string(TriggerOperationUpdate),
	}
}

func (s *TriggerOperation) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggerOperation(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggerOperation(input string) (*TriggerOperation, error) {
	vals := map[string]TriggerOperation{
		"all":     TriggerOperationAll,
		"create":  TriggerOperationCreate,
		"delete":  TriggerOperationDelete,
		"replace": TriggerOperationReplace,
		"update":  TriggerOperationUpdate,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerOperation(input)
	return &out, nil
}

type TriggerType string

const (
	TriggerTypePost TriggerType = "Post"
	TriggerTypePre  TriggerType = "Pre"
)

func PossibleValuesForTriggerType() []string {
	return []string{
		string(TriggerTypePost),
		string(TriggerTypePre),
	}
}

func (s *TriggerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTriggerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTriggerType(input string) (*TriggerType, error) {
	vals := map[string]TriggerType{
		"post": TriggerTypePost,
		"pre":  TriggerTypePre,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TriggerType(input)
	return &out, nil
}

type Type string

const (
	TypeCassandra                  Type = "Cassandra"
	TypeCassandraConnectorMetadata Type = "CassandraConnectorMetadata"
	TypeGremlin                    Type = "Gremlin"
	TypeGremlinVTwo                Type = "GremlinV2"
	TypeMongoDB                    Type = "MongoDB"
	TypeSql                        Type = "Sql"
	TypeSqlDedicatedGateway        Type = "SqlDedicatedGateway"
	TypeTable                      Type = "Table"
	TypeUndefined                  Type = "Undefined"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeCassandra),
		string(TypeCassandraConnectorMetadata),
		string(TypeGremlin),
		string(TypeGremlinVTwo),
		string(TypeMongoDB),
		string(TypeSql),
		string(TypeSqlDedicatedGateway),
		string(TypeTable),
		string(TypeUndefined),
	}
}

func (s *Type) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"cassandra":                  TypeCassandra,
		"cassandraconnectormetadata": TypeCassandraConnectorMetadata,
		"gremlin":                    TypeGremlin,
		"gremlinv2":                  TypeGremlinVTwo,
		"mongodb":                    TypeMongoDB,
		"sql":                        TypeSql,
		"sqldedicatedgateway":        TypeSqlDedicatedGateway,
		"table":                      TypeTable,
		"undefined":                  TypeUndefined,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}

type UnitType string

const (
	UnitTypeBytes          UnitType = "Bytes"
	UnitTypeBytesPerSecond UnitType = "BytesPerSecond"
	UnitTypeCount          UnitType = "Count"
	UnitTypeCountPerSecond UnitType = "CountPerSecond"
	UnitTypeMilliseconds   UnitType = "Milliseconds"
	UnitTypePercent        UnitType = "Percent"
	UnitTypeSeconds        UnitType = "Seconds"
)

func PossibleValuesForUnitType() []string {
	return []string{
		string(UnitTypeBytes),
		string(UnitTypeBytesPerSecond),
		string(UnitTypeCount),
		string(UnitTypeCountPerSecond),
		string(UnitTypeMilliseconds),
		string(UnitTypePercent),
		string(UnitTypeSeconds),
	}
}

func (s *UnitType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseUnitType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseUnitType(input string) (*UnitType, error) {
	vals := map[string]UnitType{
		"bytes":          UnitTypeBytes,
		"bytespersecond": UnitTypeBytesPerSecond,
		"count":          UnitTypeCount,
		"countpersecond": UnitTypeCountPerSecond,
		"milliseconds":   UnitTypeMilliseconds,
		"percent":        UnitTypePercent,
		"seconds":        UnitTypeSeconds,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := UnitType(input)
	return &out, nil
}

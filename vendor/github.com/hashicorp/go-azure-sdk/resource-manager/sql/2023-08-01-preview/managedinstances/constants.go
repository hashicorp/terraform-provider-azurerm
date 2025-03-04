package managedinstances

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdministratorType string

const (
	AdministratorTypeActiveDirectory AdministratorType = "ActiveDirectory"
)

func PossibleValuesForAdministratorType() []string {
	return []string{
		string(AdministratorTypeActiveDirectory),
	}
}

func (s *AdministratorType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdministratorType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdministratorType(input string) (*AdministratorType, error) {
	vals := map[string]AdministratorType{
		"activedirectory": AdministratorTypeActiveDirectory,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdministratorType(input)
	return &out, nil
}

type AggregationFunctionType string

const (
	AggregationFunctionTypeAvg   AggregationFunctionType = "avg"
	AggregationFunctionTypeMax   AggregationFunctionType = "max"
	AggregationFunctionTypeMin   AggregationFunctionType = "min"
	AggregationFunctionTypeStdev AggregationFunctionType = "stdev"
	AggregationFunctionTypeSum   AggregationFunctionType = "sum"
)

func PossibleValuesForAggregationFunctionType() []string {
	return []string{
		string(AggregationFunctionTypeAvg),
		string(AggregationFunctionTypeMax),
		string(AggregationFunctionTypeMin),
		string(AggregationFunctionTypeStdev),
		string(AggregationFunctionTypeSum),
	}
}

func (s *AggregationFunctionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAggregationFunctionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAggregationFunctionType(input string) (*AggregationFunctionType, error) {
	vals := map[string]AggregationFunctionType{
		"avg":   AggregationFunctionTypeAvg,
		"max":   AggregationFunctionTypeMax,
		"min":   AggregationFunctionTypeMin,
		"stdev": AggregationFunctionTypeStdev,
		"sum":   AggregationFunctionTypeSum,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AggregationFunctionType(input)
	return &out, nil
}

type AuthMetadataLookupModes string

const (
	AuthMetadataLookupModesAzureAD AuthMetadataLookupModes = "AzureAD"
	AuthMetadataLookupModesPaired  AuthMetadataLookupModes = "Paired"
	AuthMetadataLookupModesWindows AuthMetadataLookupModes = "Windows"
)

func PossibleValuesForAuthMetadataLookupModes() []string {
	return []string{
		string(AuthMetadataLookupModesAzureAD),
		string(AuthMetadataLookupModesPaired),
		string(AuthMetadataLookupModesWindows),
	}
}

func (s *AuthMetadataLookupModes) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAuthMetadataLookupModes(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAuthMetadataLookupModes(input string) (*AuthMetadataLookupModes, error) {
	vals := map[string]AuthMetadataLookupModes{
		"azuread": AuthMetadataLookupModesAzureAD,
		"paired":  AuthMetadataLookupModesPaired,
		"windows": AuthMetadataLookupModesWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AuthMetadataLookupModes(input)
	return &out, nil
}

type BackupStorageRedundancy string

const (
	BackupStorageRedundancyGeo     BackupStorageRedundancy = "Geo"
	BackupStorageRedundancyGeoZone BackupStorageRedundancy = "GeoZone"
	BackupStorageRedundancyLocal   BackupStorageRedundancy = "Local"
	BackupStorageRedundancyZone    BackupStorageRedundancy = "Zone"
)

func PossibleValuesForBackupStorageRedundancy() []string {
	return []string{
		string(BackupStorageRedundancyGeo),
		string(BackupStorageRedundancyGeoZone),
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
		"geo":     BackupStorageRedundancyGeo,
		"geozone": BackupStorageRedundancyGeoZone,
		"local":   BackupStorageRedundancyLocal,
		"zone":    BackupStorageRedundancyZone,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupStorageRedundancy(input)
	return &out, nil
}

type ExternalGovernanceStatus string

const (
	ExternalGovernanceStatusDisabled ExternalGovernanceStatus = "Disabled"
	ExternalGovernanceStatusEnabled  ExternalGovernanceStatus = "Enabled"
)

func PossibleValuesForExternalGovernanceStatus() []string {
	return []string{
		string(ExternalGovernanceStatusDisabled),
		string(ExternalGovernanceStatusEnabled),
	}
}

func (s *ExternalGovernanceStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseExternalGovernanceStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseExternalGovernanceStatus(input string) (*ExternalGovernanceStatus, error) {
	vals := map[string]ExternalGovernanceStatus{
		"disabled": ExternalGovernanceStatusDisabled,
		"enabled":  ExternalGovernanceStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ExternalGovernanceStatus(input)
	return &out, nil
}

type FreemiumType string

const (
	FreemiumTypeFreemium FreemiumType = "Freemium"
	FreemiumTypeRegular  FreemiumType = "Regular"
)

func PossibleValuesForFreemiumType() []string {
	return []string{
		string(FreemiumTypeFreemium),
		string(FreemiumTypeRegular),
	}
}

func (s *FreemiumType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFreemiumType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFreemiumType(input string) (*FreemiumType, error) {
	vals := map[string]FreemiumType{
		"freemium": FreemiumTypeFreemium,
		"regular":  FreemiumTypeRegular,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FreemiumType(input)
	return &out, nil
}

type HybridSecondaryUsage string

const (
	HybridSecondaryUsageActive  HybridSecondaryUsage = "Active"
	HybridSecondaryUsagePassive HybridSecondaryUsage = "Passive"
)

func PossibleValuesForHybridSecondaryUsage() []string {
	return []string{
		string(HybridSecondaryUsageActive),
		string(HybridSecondaryUsagePassive),
	}
}

func (s *HybridSecondaryUsage) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHybridSecondaryUsage(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHybridSecondaryUsage(input string) (*HybridSecondaryUsage, error) {
	vals := map[string]HybridSecondaryUsage{
		"active":  HybridSecondaryUsageActive,
		"passive": HybridSecondaryUsagePassive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HybridSecondaryUsage(input)
	return &out, nil
}

type HybridSecondaryUsageDetected string

const (
	HybridSecondaryUsageDetectedActive  HybridSecondaryUsageDetected = "Active"
	HybridSecondaryUsageDetectedPassive HybridSecondaryUsageDetected = "Passive"
)

func PossibleValuesForHybridSecondaryUsageDetected() []string {
	return []string{
		string(HybridSecondaryUsageDetectedActive),
		string(HybridSecondaryUsageDetectedPassive),
	}
}

func (s *HybridSecondaryUsageDetected) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseHybridSecondaryUsageDetected(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseHybridSecondaryUsageDetected(input string) (*HybridSecondaryUsageDetected, error) {
	vals := map[string]HybridSecondaryUsageDetected{
		"active":  HybridSecondaryUsageDetectedActive,
		"passive": HybridSecondaryUsageDetectedPassive,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := HybridSecondaryUsageDetected(input)
	return &out, nil
}

type ManagedInstanceDatabaseFormat string

const (
	ManagedInstanceDatabaseFormatAlwaysUpToDate         ManagedInstanceDatabaseFormat = "AlwaysUpToDate"
	ManagedInstanceDatabaseFormatSQLServerTwoZeroTwoTwo ManagedInstanceDatabaseFormat = "SQLServer2022"
)

func PossibleValuesForManagedInstanceDatabaseFormat() []string {
	return []string{
		string(ManagedInstanceDatabaseFormatAlwaysUpToDate),
		string(ManagedInstanceDatabaseFormatSQLServerTwoZeroTwoTwo),
	}
}

func (s *ManagedInstanceDatabaseFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedInstanceDatabaseFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedInstanceDatabaseFormat(input string) (*ManagedInstanceDatabaseFormat, error) {
	vals := map[string]ManagedInstanceDatabaseFormat{
		"alwaysuptodate": ManagedInstanceDatabaseFormatAlwaysUpToDate,
		"sqlserver2022":  ManagedInstanceDatabaseFormatSQLServerTwoZeroTwoTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedInstanceDatabaseFormat(input)
	return &out, nil
}

type ManagedInstanceLicenseType string

const (
	ManagedInstanceLicenseTypeBasePrice       ManagedInstanceLicenseType = "BasePrice"
	ManagedInstanceLicenseTypeLicenseIncluded ManagedInstanceLicenseType = "LicenseIncluded"
)

func PossibleValuesForManagedInstanceLicenseType() []string {
	return []string{
		string(ManagedInstanceLicenseTypeBasePrice),
		string(ManagedInstanceLicenseTypeLicenseIncluded),
	}
}

func (s *ManagedInstanceLicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedInstanceLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedInstanceLicenseType(input string) (*ManagedInstanceLicenseType, error) {
	vals := map[string]ManagedInstanceLicenseType{
		"baseprice":       ManagedInstanceLicenseTypeBasePrice,
		"licenseincluded": ManagedInstanceLicenseTypeLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedInstanceLicenseType(input)
	return &out, nil
}

type ManagedInstanceProxyOverride string

const (
	ManagedInstanceProxyOverrideDefault  ManagedInstanceProxyOverride = "Default"
	ManagedInstanceProxyOverrideProxy    ManagedInstanceProxyOverride = "Proxy"
	ManagedInstanceProxyOverrideRedirect ManagedInstanceProxyOverride = "Redirect"
)

func PossibleValuesForManagedInstanceProxyOverride() []string {
	return []string{
		string(ManagedInstanceProxyOverrideDefault),
		string(ManagedInstanceProxyOverrideProxy),
		string(ManagedInstanceProxyOverrideRedirect),
	}
}

func (s *ManagedInstanceProxyOverride) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedInstanceProxyOverride(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedInstanceProxyOverride(input string) (*ManagedInstanceProxyOverride, error) {
	vals := map[string]ManagedInstanceProxyOverride{
		"default":  ManagedInstanceProxyOverrideDefault,
		"proxy":    ManagedInstanceProxyOverrideProxy,
		"redirect": ManagedInstanceProxyOverrideRedirect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedInstanceProxyOverride(input)
	return &out, nil
}

type ManagedServerCreateMode string

const (
	ManagedServerCreateModeDefault            ManagedServerCreateMode = "Default"
	ManagedServerCreateModePointInTimeRestore ManagedServerCreateMode = "PointInTimeRestore"
)

func PossibleValuesForManagedServerCreateMode() []string {
	return []string{
		string(ManagedServerCreateModeDefault),
		string(ManagedServerCreateModePointInTimeRestore),
	}
}

func (s *ManagedServerCreateMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedServerCreateMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedServerCreateMode(input string) (*ManagedServerCreateMode, error) {
	vals := map[string]ManagedServerCreateMode{
		"default":            ManagedServerCreateModeDefault,
		"pointintimerestore": ManagedServerCreateModePointInTimeRestore,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedServerCreateMode(input)
	return &out, nil
}

type MetricType string

const (
	MetricTypeCpu      MetricType = "cpu"
	MetricTypeDtu      MetricType = "dtu"
	MetricTypeDuration MetricType = "duration"
	MetricTypeIo       MetricType = "io"
	MetricTypeLogIo    MetricType = "logIo"
)

func PossibleValuesForMetricType() []string {
	return []string{
		string(MetricTypeCpu),
		string(MetricTypeDtu),
		string(MetricTypeDuration),
		string(MetricTypeIo),
		string(MetricTypeLogIo),
	}
}

func (s *MetricType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMetricType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMetricType(input string) (*MetricType, error) {
	vals := map[string]MetricType{
		"cpu":      MetricTypeCpu,
		"dtu":      MetricTypeDtu,
		"duration": MetricTypeDuration,
		"io":       MetricTypeIo,
		"logio":    MetricTypeLogIo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetricType(input)
	return &out, nil
}

type PrincipalType string

const (
	PrincipalTypeApplication PrincipalType = "Application"
	PrincipalTypeGroup       PrincipalType = "Group"
	PrincipalTypeUser        PrincipalType = "User"
)

func PossibleValuesForPrincipalType() []string {
	return []string{
		string(PrincipalTypeApplication),
		string(PrincipalTypeGroup),
		string(PrincipalTypeUser),
	}
}

func (s *PrincipalType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parsePrincipalType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parsePrincipalType(input string) (*PrincipalType, error) {
	vals := map[string]PrincipalType{
		"application": PrincipalTypeApplication,
		"group":       PrincipalTypeGroup,
		"user":        PrincipalTypeUser,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PrincipalType(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled   ProvisioningState = "Canceled"
	ProvisioningStateCreated    ProvisioningState = "Created"
	ProvisioningStateFailed     ProvisioningState = "Failed"
	ProvisioningStateInProgress ProvisioningState = "InProgress"
	ProvisioningStateSucceeded  ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateCreated),
		string(ProvisioningStateFailed),
		string(ProvisioningStateInProgress),
		string(ProvisioningStateSucceeded),
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
		"canceled":   ProvisioningStateCanceled,
		"created":    ProvisioningStateCreated,
		"failed":     ProvisioningStateFailed,
		"inprogress": ProvisioningStateInProgress,
		"succeeded":  ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type QueryMetricUnitType string

const (
	QueryMetricUnitTypeCount        QueryMetricUnitType = "count"
	QueryMetricUnitTypeKB           QueryMetricUnitType = "KB"
	QueryMetricUnitTypeMicroseconds QueryMetricUnitType = "microseconds"
	QueryMetricUnitTypePercentage   QueryMetricUnitType = "percentage"
)

func PossibleValuesForQueryMetricUnitType() []string {
	return []string{
		string(QueryMetricUnitTypeCount),
		string(QueryMetricUnitTypeKB),
		string(QueryMetricUnitTypeMicroseconds),
		string(QueryMetricUnitTypePercentage),
	}
}

func (s *QueryMetricUnitType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryMetricUnitType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryMetricUnitType(input string) (*QueryMetricUnitType, error) {
	vals := map[string]QueryMetricUnitType{
		"count":        QueryMetricUnitTypeCount,
		"kb":           QueryMetricUnitTypeKB,
		"microseconds": QueryMetricUnitTypeMicroseconds,
		"percentage":   QueryMetricUnitTypePercentage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryMetricUnitType(input)
	return &out, nil
}

type QueryTimeGrainType string

const (
	QueryTimeGrainTypePOneD  QueryTimeGrainType = "P1D"
	QueryTimeGrainTypePTOneH QueryTimeGrainType = "PT1H"
)

func PossibleValuesForQueryTimeGrainType() []string {
	return []string{
		string(QueryTimeGrainTypePOneD),
		string(QueryTimeGrainTypePTOneH),
	}
}

func (s *QueryTimeGrainType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryTimeGrainType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryTimeGrainType(input string) (*QueryTimeGrainType, error) {
	vals := map[string]QueryTimeGrainType{
		"p1d":  QueryTimeGrainTypePOneD,
		"pt1h": QueryTimeGrainTypePTOneH,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryTimeGrainType(input)
	return &out, nil
}

type ReplicaType string

const (
	ReplicaTypePrimary           ReplicaType = "Primary"
	ReplicaTypeReadableSecondary ReplicaType = "ReadableSecondary"
)

func PossibleValuesForReplicaType() []string {
	return []string{
		string(ReplicaTypePrimary),
		string(ReplicaTypeReadableSecondary),
	}
}

func (s *ReplicaType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseReplicaType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseReplicaType(input string) (*ReplicaType, error) {
	vals := map[string]ReplicaType{
		"primary":           ReplicaTypePrimary,
		"readablesecondary": ReplicaTypeReadableSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ReplicaType(input)
	return &out, nil
}

type ServicePrincipalType string

const (
	ServicePrincipalTypeNone           ServicePrincipalType = "None"
	ServicePrincipalTypeSystemAssigned ServicePrincipalType = "SystemAssigned"
)

func PossibleValuesForServicePrincipalType() []string {
	return []string{
		string(ServicePrincipalTypeNone),
		string(ServicePrincipalTypeSystemAssigned),
	}
}

func (s *ServicePrincipalType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseServicePrincipalType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseServicePrincipalType(input string) (*ServicePrincipalType, error) {
	vals := map[string]ServicePrincipalType{
		"none":           ServicePrincipalTypeNone,
		"systemassigned": ServicePrincipalTypeSystemAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ServicePrincipalType(input)
	return &out, nil
}

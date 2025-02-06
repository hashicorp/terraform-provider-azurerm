package databases

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlwaysEncryptedEnclaveType string

const (
	AlwaysEncryptedEnclaveTypeDefault AlwaysEncryptedEnclaveType = "Default"
	AlwaysEncryptedEnclaveTypeVBS     AlwaysEncryptedEnclaveType = "VBS"
)

func PossibleValuesForAlwaysEncryptedEnclaveType() []string {
	return []string{
		string(AlwaysEncryptedEnclaveTypeDefault),
		string(AlwaysEncryptedEnclaveTypeVBS),
	}
}

func (s *AlwaysEncryptedEnclaveType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlwaysEncryptedEnclaveType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlwaysEncryptedEnclaveType(input string) (*AlwaysEncryptedEnclaveType, error) {
	vals := map[string]AlwaysEncryptedEnclaveType{
		"default": AlwaysEncryptedEnclaveTypeDefault,
		"vbs":     AlwaysEncryptedEnclaveTypeVBS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlwaysEncryptedEnclaveType(input)
	return &out, nil
}

type AvailabilityZoneType string

const (
	AvailabilityZoneTypeNoPreference AvailabilityZoneType = "NoPreference"
	AvailabilityZoneTypeOne          AvailabilityZoneType = "1"
	AvailabilityZoneTypeThree        AvailabilityZoneType = "3"
	AvailabilityZoneTypeTwo          AvailabilityZoneType = "2"
)

func PossibleValuesForAvailabilityZoneType() []string {
	return []string{
		string(AvailabilityZoneTypeNoPreference),
		string(AvailabilityZoneTypeOne),
		string(AvailabilityZoneTypeThree),
		string(AvailabilityZoneTypeTwo),
	}
}

func (s *AvailabilityZoneType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAvailabilityZoneType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAvailabilityZoneType(input string) (*AvailabilityZoneType, error) {
	vals := map[string]AvailabilityZoneType{
		"nopreference": AvailabilityZoneTypeNoPreference,
		"1":            AvailabilityZoneTypeOne,
		"3":            AvailabilityZoneTypeThree,
		"2":            AvailabilityZoneTypeTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AvailabilityZoneType(input)
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

type CatalogCollationType string

const (
	CatalogCollationTypeDATABASEDEFAULT             CatalogCollationType = "DATABASE_DEFAULT"
	CatalogCollationTypeSQLLatinOneGeneralCPOneCIAS CatalogCollationType = "SQL_Latin1_General_CP1_CI_AS"
)

func PossibleValuesForCatalogCollationType() []string {
	return []string{
		string(CatalogCollationTypeDATABASEDEFAULT),
		string(CatalogCollationTypeSQLLatinOneGeneralCPOneCIAS),
	}
}

func (s *CatalogCollationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCatalogCollationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCatalogCollationType(input string) (*CatalogCollationType, error) {
	vals := map[string]CatalogCollationType{
		"database_default":             CatalogCollationTypeDATABASEDEFAULT,
		"sql_latin1_general_cp1_ci_as": CatalogCollationTypeSQLLatinOneGeneralCPOneCIAS,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CatalogCollationType(input)
	return &out, nil
}

type CreateMode string

const (
	CreateModeCopy                           CreateMode = "Copy"
	CreateModeDefault                        CreateMode = "Default"
	CreateModeOnlineSecondary                CreateMode = "OnlineSecondary"
	CreateModePointInTimeRestore             CreateMode = "PointInTimeRestore"
	CreateModeRecovery                       CreateMode = "Recovery"
	CreateModeRestore                        CreateMode = "Restore"
	CreateModeRestoreExternalBackup          CreateMode = "RestoreExternalBackup"
	CreateModeRestoreExternalBackupSecondary CreateMode = "RestoreExternalBackupSecondary"
	CreateModeRestoreLongTermRetentionBackup CreateMode = "RestoreLongTermRetentionBackup"
	CreateModeSecondary                      CreateMode = "Secondary"
)

func PossibleValuesForCreateMode() []string {
	return []string{
		string(CreateModeCopy),
		string(CreateModeDefault),
		string(CreateModeOnlineSecondary),
		string(CreateModePointInTimeRestore),
		string(CreateModeRecovery),
		string(CreateModeRestore),
		string(CreateModeRestoreExternalBackup),
		string(CreateModeRestoreExternalBackupSecondary),
		string(CreateModeRestoreLongTermRetentionBackup),
		string(CreateModeSecondary),
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
		"copy":                           CreateModeCopy,
		"default":                        CreateModeDefault,
		"onlinesecondary":                CreateModeOnlineSecondary,
		"pointintimerestore":             CreateModePointInTimeRestore,
		"recovery":                       CreateModeRecovery,
		"restore":                        CreateModeRestore,
		"restoreexternalbackup":          CreateModeRestoreExternalBackup,
		"restoreexternalbackupsecondary": CreateModeRestoreExternalBackupSecondary,
		"restorelongtermretentionbackup": CreateModeRestoreLongTermRetentionBackup,
		"secondary":                      CreateModeSecondary,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CreateMode(input)
	return &out, nil
}

type DatabaseKeyType string

const (
	DatabaseKeyTypeAzureKeyVault DatabaseKeyType = "AzureKeyVault"
)

func PossibleValuesForDatabaseKeyType() []string {
	return []string{
		string(DatabaseKeyTypeAzureKeyVault),
	}
}

func (s *DatabaseKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseKeyType(input string) (*DatabaseKeyType, error) {
	vals := map[string]DatabaseKeyType{
		"azurekeyvault": DatabaseKeyTypeAzureKeyVault,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseKeyType(input)
	return &out, nil
}

type DatabaseLicenseType string

const (
	DatabaseLicenseTypeBasePrice       DatabaseLicenseType = "BasePrice"
	DatabaseLicenseTypeLicenseIncluded DatabaseLicenseType = "LicenseIncluded"
)

func PossibleValuesForDatabaseLicenseType() []string {
	return []string{
		string(DatabaseLicenseTypeBasePrice),
		string(DatabaseLicenseTypeLicenseIncluded),
	}
}

func (s *DatabaseLicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseLicenseType(input string) (*DatabaseLicenseType, error) {
	vals := map[string]DatabaseLicenseType{
		"baseprice":       DatabaseLicenseTypeBasePrice,
		"licenseincluded": DatabaseLicenseTypeLicenseIncluded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseLicenseType(input)
	return &out, nil
}

type DatabaseReadScale string

const (
	DatabaseReadScaleDisabled DatabaseReadScale = "Disabled"
	DatabaseReadScaleEnabled  DatabaseReadScale = "Enabled"
)

func PossibleValuesForDatabaseReadScale() []string {
	return []string{
		string(DatabaseReadScaleDisabled),
		string(DatabaseReadScaleEnabled),
	}
}

func (s *DatabaseReadScale) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseReadScale(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseReadScale(input string) (*DatabaseReadScale, error) {
	vals := map[string]DatabaseReadScale{
		"disabled": DatabaseReadScaleDisabled,
		"enabled":  DatabaseReadScaleEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseReadScale(input)
	return &out, nil
}

type DatabaseStatus string

const (
	DatabaseStatusAutoClosed                        DatabaseStatus = "AutoClosed"
	DatabaseStatusCopying                           DatabaseStatus = "Copying"
	DatabaseStatusCreating                          DatabaseStatus = "Creating"
	DatabaseStatusDisabled                          DatabaseStatus = "Disabled"
	DatabaseStatusEmergencyMode                     DatabaseStatus = "EmergencyMode"
	DatabaseStatusInaccessible                      DatabaseStatus = "Inaccessible"
	DatabaseStatusOffline                           DatabaseStatus = "Offline"
	DatabaseStatusOfflineChangingDwPerformanceTiers DatabaseStatus = "OfflineChangingDwPerformanceTiers"
	DatabaseStatusOfflineSecondary                  DatabaseStatus = "OfflineSecondary"
	DatabaseStatusOnline                            DatabaseStatus = "Online"
	DatabaseStatusOnlineChangingDwPerformanceTiers  DatabaseStatus = "OnlineChangingDwPerformanceTiers"
	DatabaseStatusPaused                            DatabaseStatus = "Paused"
	DatabaseStatusPausing                           DatabaseStatus = "Pausing"
	DatabaseStatusRecovering                        DatabaseStatus = "Recovering"
	DatabaseStatusRecoveryPending                   DatabaseStatus = "RecoveryPending"
	DatabaseStatusRestoring                         DatabaseStatus = "Restoring"
	DatabaseStatusResuming                          DatabaseStatus = "Resuming"
	DatabaseStatusScaling                           DatabaseStatus = "Scaling"
	DatabaseStatusShutdown                          DatabaseStatus = "Shutdown"
	DatabaseStatusStandby                           DatabaseStatus = "Standby"
	DatabaseStatusStarting                          DatabaseStatus = "Starting"
	DatabaseStatusStopped                           DatabaseStatus = "Stopped"
	DatabaseStatusStopping                          DatabaseStatus = "Stopping"
	DatabaseStatusSuspect                           DatabaseStatus = "Suspect"
)

func PossibleValuesForDatabaseStatus() []string {
	return []string{
		string(DatabaseStatusAutoClosed),
		string(DatabaseStatusCopying),
		string(DatabaseStatusCreating),
		string(DatabaseStatusDisabled),
		string(DatabaseStatusEmergencyMode),
		string(DatabaseStatusInaccessible),
		string(DatabaseStatusOffline),
		string(DatabaseStatusOfflineChangingDwPerformanceTiers),
		string(DatabaseStatusOfflineSecondary),
		string(DatabaseStatusOnline),
		string(DatabaseStatusOnlineChangingDwPerformanceTiers),
		string(DatabaseStatusPaused),
		string(DatabaseStatusPausing),
		string(DatabaseStatusRecovering),
		string(DatabaseStatusRecoveryPending),
		string(DatabaseStatusRestoring),
		string(DatabaseStatusResuming),
		string(DatabaseStatusScaling),
		string(DatabaseStatusShutdown),
		string(DatabaseStatusStandby),
		string(DatabaseStatusStarting),
		string(DatabaseStatusStopped),
		string(DatabaseStatusStopping),
		string(DatabaseStatusSuspect),
	}
}

func (s *DatabaseStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDatabaseStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDatabaseStatus(input string) (*DatabaseStatus, error) {
	vals := map[string]DatabaseStatus{
		"autoclosed":                        DatabaseStatusAutoClosed,
		"copying":                           DatabaseStatusCopying,
		"creating":                          DatabaseStatusCreating,
		"disabled":                          DatabaseStatusDisabled,
		"emergencymode":                     DatabaseStatusEmergencyMode,
		"inaccessible":                      DatabaseStatusInaccessible,
		"offline":                           DatabaseStatusOffline,
		"offlinechangingdwperformancetiers": DatabaseStatusOfflineChangingDwPerformanceTiers,
		"offlinesecondary":                  DatabaseStatusOfflineSecondary,
		"online":                            DatabaseStatusOnline,
		"onlinechangingdwperformancetiers":  DatabaseStatusOnlineChangingDwPerformanceTiers,
		"paused":                            DatabaseStatusPaused,
		"pausing":                           DatabaseStatusPausing,
		"recovering":                        DatabaseStatusRecovering,
		"recoverypending":                   DatabaseStatusRecoveryPending,
		"restoring":                         DatabaseStatusRestoring,
		"resuming":                          DatabaseStatusResuming,
		"scaling":                           DatabaseStatusScaling,
		"shutdown":                          DatabaseStatusShutdown,
		"standby":                           DatabaseStatusStandby,
		"starting":                          DatabaseStatusStarting,
		"stopped":                           DatabaseStatusStopped,
		"stopping":                          DatabaseStatusStopping,
		"suspect":                           DatabaseStatusSuspect,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DatabaseStatus(input)
	return &out, nil
}

type FreeLimitExhaustionBehavior string

const (
	FreeLimitExhaustionBehaviorAutoPause     FreeLimitExhaustionBehavior = "AutoPause"
	FreeLimitExhaustionBehaviorBillOverUsage FreeLimitExhaustionBehavior = "BillOverUsage"
)

func PossibleValuesForFreeLimitExhaustionBehavior() []string {
	return []string{
		string(FreeLimitExhaustionBehaviorAutoPause),
		string(FreeLimitExhaustionBehaviorBillOverUsage),
	}
}

func (s *FreeLimitExhaustionBehavior) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFreeLimitExhaustionBehavior(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFreeLimitExhaustionBehavior(input string) (*FreeLimitExhaustionBehavior, error) {
	vals := map[string]FreeLimitExhaustionBehavior{
		"autopause":     FreeLimitExhaustionBehaviorAutoPause,
		"billoverusage": FreeLimitExhaustionBehaviorBillOverUsage,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FreeLimitExhaustionBehavior(input)
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

type SampleName string

const (
	SampleNameAdventureWorksLT       SampleName = "AdventureWorksLT"
	SampleNameWideWorldImportersFull SampleName = "WideWorldImportersFull"
	SampleNameWideWorldImportersStd  SampleName = "WideWorldImportersStd"
)

func PossibleValuesForSampleName() []string {
	return []string{
		string(SampleNameAdventureWorksLT),
		string(SampleNameWideWorldImportersFull),
		string(SampleNameWideWorldImportersStd),
	}
}

func (s *SampleName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSampleName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSampleName(input string) (*SampleName, error) {
	vals := map[string]SampleName{
		"adventureworkslt":       SampleNameAdventureWorksLT,
		"wideworldimportersfull": SampleNameWideWorldImportersFull,
		"wideworldimportersstd":  SampleNameWideWorldImportersStd,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SampleName(input)
	return &out, nil
}

type SecondaryType string

const (
	SecondaryTypeGeo     SecondaryType = "Geo"
	SecondaryTypeNamed   SecondaryType = "Named"
	SecondaryTypeStandby SecondaryType = "Standby"
)

func PossibleValuesForSecondaryType() []string {
	return []string{
		string(SecondaryTypeGeo),
		string(SecondaryTypeNamed),
		string(SecondaryTypeStandby),
	}
}

func (s *SecondaryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSecondaryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSecondaryType(input string) (*SecondaryType, error) {
	vals := map[string]SecondaryType{
		"geo":     SecondaryTypeGeo,
		"named":   SecondaryTypeNamed,
		"standby": SecondaryTypeStandby,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SecondaryType(input)
	return &out, nil
}

type StorageKeyType string

const (
	StorageKeyTypeSharedAccessKey  StorageKeyType = "SharedAccessKey"
	StorageKeyTypeStorageAccessKey StorageKeyType = "StorageAccessKey"
)

func PossibleValuesForStorageKeyType() []string {
	return []string{
		string(StorageKeyTypeSharedAccessKey),
		string(StorageKeyTypeStorageAccessKey),
	}
}

func (s *StorageKeyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageKeyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageKeyType(input string) (*StorageKeyType, error) {
	vals := map[string]StorageKeyType{
		"sharedaccesskey":  StorageKeyTypeSharedAccessKey,
		"storageaccesskey": StorageKeyTypeStorageAccessKey,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageKeyType(input)
	return &out, nil
}

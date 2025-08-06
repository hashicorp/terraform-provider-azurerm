package sqlvirtualmachines

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdditionalOsPatch string

const (
	AdditionalOsPatchWSUS AdditionalOsPatch = "WSUS"
	AdditionalOsPatchWU   AdditionalOsPatch = "WU"
	AdditionalOsPatchWUMU AdditionalOsPatch = "WUMU"
)

func PossibleValuesForAdditionalOsPatch() []string {
	return []string{
		string(AdditionalOsPatchWSUS),
		string(AdditionalOsPatchWU),
		string(AdditionalOsPatchWUMU),
	}
}

func (s *AdditionalOsPatch) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdditionalOsPatch(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdditionalOsPatch(input string) (*AdditionalOsPatch, error) {
	vals := map[string]AdditionalOsPatch{
		"wsus": AdditionalOsPatchWSUS,
		"wu":   AdditionalOsPatchWU,
		"wumu": AdditionalOsPatchWUMU,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdditionalOsPatch(input)
	return &out, nil
}

type AdditionalVMPatch string

const (
	AdditionalVMPatchMicrosoftUpdate AdditionalVMPatch = "MicrosoftUpdate"
	AdditionalVMPatchNotSet          AdditionalVMPatch = "NotSet"
)

func PossibleValuesForAdditionalVMPatch() []string {
	return []string{
		string(AdditionalVMPatchMicrosoftUpdate),
		string(AdditionalVMPatchNotSet),
	}
}

func (s *AdditionalVMPatch) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAdditionalVMPatch(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAdditionalVMPatch(input string) (*AdditionalVMPatch, error) {
	vals := map[string]AdditionalVMPatch{
		"microsoftupdate": AdditionalVMPatchMicrosoftUpdate,
		"notset":          AdditionalVMPatchNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AdditionalVMPatch(input)
	return &out, nil
}

type AssessmentDayOfWeek string

const (
	AssessmentDayOfWeekFriday    AssessmentDayOfWeek = "Friday"
	AssessmentDayOfWeekMonday    AssessmentDayOfWeek = "Monday"
	AssessmentDayOfWeekSaturday  AssessmentDayOfWeek = "Saturday"
	AssessmentDayOfWeekSunday    AssessmentDayOfWeek = "Sunday"
	AssessmentDayOfWeekThursday  AssessmentDayOfWeek = "Thursday"
	AssessmentDayOfWeekTuesday   AssessmentDayOfWeek = "Tuesday"
	AssessmentDayOfWeekWednesday AssessmentDayOfWeek = "Wednesday"
)

func PossibleValuesForAssessmentDayOfWeek() []string {
	return []string{
		string(AssessmentDayOfWeekFriday),
		string(AssessmentDayOfWeekMonday),
		string(AssessmentDayOfWeekSaturday),
		string(AssessmentDayOfWeekSunday),
		string(AssessmentDayOfWeekThursday),
		string(AssessmentDayOfWeekTuesday),
		string(AssessmentDayOfWeekWednesday),
	}
}

func (s *AssessmentDayOfWeek) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAssessmentDayOfWeek(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAssessmentDayOfWeek(input string) (*AssessmentDayOfWeek, error) {
	vals := map[string]AssessmentDayOfWeek{
		"friday":    AssessmentDayOfWeekFriday,
		"monday":    AssessmentDayOfWeekMonday,
		"saturday":  AssessmentDayOfWeekSaturday,
		"sunday":    AssessmentDayOfWeekSunday,
		"thursday":  AssessmentDayOfWeekThursday,
		"tuesday":   AssessmentDayOfWeekTuesday,
		"wednesday": AssessmentDayOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AssessmentDayOfWeek(input)
	return &out, nil
}

type AutoBackupDaysOfWeek string

const (
	AutoBackupDaysOfWeekFriday    AutoBackupDaysOfWeek = "Friday"
	AutoBackupDaysOfWeekMonday    AutoBackupDaysOfWeek = "Monday"
	AutoBackupDaysOfWeekSaturday  AutoBackupDaysOfWeek = "Saturday"
	AutoBackupDaysOfWeekSunday    AutoBackupDaysOfWeek = "Sunday"
	AutoBackupDaysOfWeekThursday  AutoBackupDaysOfWeek = "Thursday"
	AutoBackupDaysOfWeekTuesday   AutoBackupDaysOfWeek = "Tuesday"
	AutoBackupDaysOfWeekWednesday AutoBackupDaysOfWeek = "Wednesday"
)

func PossibleValuesForAutoBackupDaysOfWeek() []string {
	return []string{
		string(AutoBackupDaysOfWeekFriday),
		string(AutoBackupDaysOfWeekMonday),
		string(AutoBackupDaysOfWeekSaturday),
		string(AutoBackupDaysOfWeekSunday),
		string(AutoBackupDaysOfWeekThursday),
		string(AutoBackupDaysOfWeekTuesday),
		string(AutoBackupDaysOfWeekWednesday),
	}
}

func (s *AutoBackupDaysOfWeek) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAutoBackupDaysOfWeek(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAutoBackupDaysOfWeek(input string) (*AutoBackupDaysOfWeek, error) {
	vals := map[string]AutoBackupDaysOfWeek{
		"friday":    AutoBackupDaysOfWeekFriday,
		"monday":    AutoBackupDaysOfWeekMonday,
		"saturday":  AutoBackupDaysOfWeekSaturday,
		"sunday":    AutoBackupDaysOfWeekSunday,
		"thursday":  AutoBackupDaysOfWeekThursday,
		"tuesday":   AutoBackupDaysOfWeekTuesday,
		"wednesday": AutoBackupDaysOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AutoBackupDaysOfWeek(input)
	return &out, nil
}

type BackupScheduleType string

const (
	BackupScheduleTypeAutomated BackupScheduleType = "Automated"
	BackupScheduleTypeManual    BackupScheduleType = "Manual"
)

func PossibleValuesForBackupScheduleType() []string {
	return []string{
		string(BackupScheduleTypeAutomated),
		string(BackupScheduleTypeManual),
	}
}

func (s *BackupScheduleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseBackupScheduleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseBackupScheduleType(input string) (*BackupScheduleType, error) {
	vals := map[string]BackupScheduleType{
		"automated": BackupScheduleTypeAutomated,
		"manual":    BackupScheduleTypeManual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := BackupScheduleType(input)
	return &out, nil
}

type ConnectivityType string

const (
	ConnectivityTypeLOCAL   ConnectivityType = "LOCAL"
	ConnectivityTypePRIVATE ConnectivityType = "PRIVATE"
	ConnectivityTypePUBLIC  ConnectivityType = "PUBLIC"
)

func PossibleValuesForConnectivityType() []string {
	return []string{
		string(ConnectivityTypeLOCAL),
		string(ConnectivityTypePRIVATE),
		string(ConnectivityTypePUBLIC),
	}
}

func (s *ConnectivityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConnectivityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConnectivityType(input string) (*ConnectivityType, error) {
	vals := map[string]ConnectivityType{
		"local":   ConnectivityTypeLOCAL,
		"private": ConnectivityTypePRIVATE,
		"public":  ConnectivityTypePUBLIC,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConnectivityType(input)
	return &out, nil
}

type DayOfWeek string

const (
	DayOfWeekEveryday  DayOfWeek = "Everyday"
	DayOfWeekFriday    DayOfWeek = "Friday"
	DayOfWeekMonday    DayOfWeek = "Monday"
	DayOfWeekSaturday  DayOfWeek = "Saturday"
	DayOfWeekSunday    DayOfWeek = "Sunday"
	DayOfWeekThursday  DayOfWeek = "Thursday"
	DayOfWeekTuesday   DayOfWeek = "Tuesday"
	DayOfWeekWednesday DayOfWeek = "Wednesday"
)

func PossibleValuesForDayOfWeek() []string {
	return []string{
		string(DayOfWeekEveryday),
		string(DayOfWeekFriday),
		string(DayOfWeekMonday),
		string(DayOfWeekSaturday),
		string(DayOfWeekSunday),
		string(DayOfWeekThursday),
		string(DayOfWeekTuesday),
		string(DayOfWeekWednesday),
	}
}

func (s *DayOfWeek) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDayOfWeek(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDayOfWeek(input string) (*DayOfWeek, error) {
	vals := map[string]DayOfWeek{
		"everyday":  DayOfWeekEveryday,
		"friday":    DayOfWeekFriday,
		"monday":    DayOfWeekMonday,
		"saturday":  DayOfWeekSaturday,
		"sunday":    DayOfWeekSunday,
		"thursday":  DayOfWeekThursday,
		"tuesday":   DayOfWeekTuesday,
		"wednesday": DayOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DayOfWeek(input)
	return &out, nil
}

type DiskConfigurationType string

const (
	DiskConfigurationTypeADD    DiskConfigurationType = "ADD"
	DiskConfigurationTypeEXTEND DiskConfigurationType = "EXTEND"
	DiskConfigurationTypeNEW    DiskConfigurationType = "NEW"
)

func PossibleValuesForDiskConfigurationType() []string {
	return []string{
		string(DiskConfigurationTypeADD),
		string(DiskConfigurationTypeEXTEND),
		string(DiskConfigurationTypeNEW),
	}
}

func (s *DiskConfigurationType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDiskConfigurationType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDiskConfigurationType(input string) (*DiskConfigurationType, error) {
	vals := map[string]DiskConfigurationType{
		"add":    DiskConfigurationTypeADD,
		"extend": DiskConfigurationTypeEXTEND,
		"new":    DiskConfigurationTypeNEW,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DiskConfigurationType(input)
	return &out, nil
}

type FullBackupFrequencyType string

const (
	FullBackupFrequencyTypeDaily  FullBackupFrequencyType = "Daily"
	FullBackupFrequencyTypeWeekly FullBackupFrequencyType = "Weekly"
)

func PossibleValuesForFullBackupFrequencyType() []string {
	return []string{
		string(FullBackupFrequencyTypeDaily),
		string(FullBackupFrequencyTypeWeekly),
	}
}

func (s *FullBackupFrequencyType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFullBackupFrequencyType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFullBackupFrequencyType(input string) (*FullBackupFrequencyType, error) {
	vals := map[string]FullBackupFrequencyType{
		"daily":  FullBackupFrequencyTypeDaily,
		"weekly": FullBackupFrequencyTypeWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FullBackupFrequencyType(input)
	return &out, nil
}

type IdentityType string

const (
	IdentityTypeNone                       IdentityType = "None"
	IdentityTypeSystemAssigned             IdentityType = "SystemAssigned"
	IdentityTypeSystemAssignedUserAssigned IdentityType = "SystemAssigned,UserAssigned"
	IdentityTypeUserAssigned               IdentityType = "UserAssigned"
)

func PossibleValuesForIdentityType() []string {
	return []string{
		string(IdentityTypeNone),
		string(IdentityTypeSystemAssigned),
		string(IdentityTypeSystemAssignedUserAssigned),
		string(IdentityTypeUserAssigned),
	}
}

func (s *IdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentityType(input string) (*IdentityType, error) {
	vals := map[string]IdentityType{
		"none":                        IdentityTypeNone,
		"systemassigned":              IdentityTypeSystemAssigned,
		"systemassigned,userassigned": IdentityTypeSystemAssignedUserAssigned,
		"userassigned":                IdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IdentityType(input)
	return &out, nil
}

type LeastPrivilegeMode string

const (
	LeastPrivilegeModeEnabled LeastPrivilegeMode = "Enabled"
	LeastPrivilegeModeNotSet  LeastPrivilegeMode = "NotSet"
)

func PossibleValuesForLeastPrivilegeMode() []string {
	return []string{
		string(LeastPrivilegeModeEnabled),
		string(LeastPrivilegeModeNotSet),
	}
}

func (s *LeastPrivilegeMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseLeastPrivilegeMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseLeastPrivilegeMode(input string) (*LeastPrivilegeMode, error) {
	vals := map[string]LeastPrivilegeMode{
		"enabled": LeastPrivilegeModeEnabled,
		"notset":  LeastPrivilegeModeNotSet,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := LeastPrivilegeMode(input)
	return &out, nil
}

type OsType string

const (
	OsTypeLinux   OsType = "Linux"
	OsTypeWindows OsType = "Windows"
)

func PossibleValuesForOsType() []string {
	return []string{
		string(OsTypeLinux),
		string(OsTypeWindows),
	}
}

func (s *OsType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOsType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOsType(input string) (*OsType, error) {
	vals := map[string]OsType{
		"linux":   OsTypeLinux,
		"windows": OsTypeWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := OsType(input)
	return &out, nil
}

type SqlImageSku string

const (
	SqlImageSkuDeveloper  SqlImageSku = "Developer"
	SqlImageSkuEnterprise SqlImageSku = "Enterprise"
	SqlImageSkuExpress    SqlImageSku = "Express"
	SqlImageSkuStandard   SqlImageSku = "Standard"
	SqlImageSkuWeb        SqlImageSku = "Web"
)

func PossibleValuesForSqlImageSku() []string {
	return []string{
		string(SqlImageSkuDeveloper),
		string(SqlImageSkuEnterprise),
		string(SqlImageSkuExpress),
		string(SqlImageSkuStandard),
		string(SqlImageSkuWeb),
	}
}

func (s *SqlImageSku) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSqlImageSku(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSqlImageSku(input string) (*SqlImageSku, error) {
	vals := map[string]SqlImageSku{
		"developer":  SqlImageSkuDeveloper,
		"enterprise": SqlImageSkuEnterprise,
		"express":    SqlImageSkuExpress,
		"standard":   SqlImageSkuStandard,
		"web":        SqlImageSkuWeb,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SqlImageSku(input)
	return &out, nil
}

type SqlManagementMode string

const (
	SqlManagementModeFull        SqlManagementMode = "Full"
	SqlManagementModeLightWeight SqlManagementMode = "LightWeight"
	SqlManagementModeNoAgent     SqlManagementMode = "NoAgent"
)

func PossibleValuesForSqlManagementMode() []string {
	return []string{
		string(SqlManagementModeFull),
		string(SqlManagementModeLightWeight),
		string(SqlManagementModeNoAgent),
	}
}

func (s *SqlManagementMode) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSqlManagementMode(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSqlManagementMode(input string) (*SqlManagementMode, error) {
	vals := map[string]SqlManagementMode{
		"full":        SqlManagementModeFull,
		"lightweight": SqlManagementModeLightWeight,
		"noagent":     SqlManagementModeNoAgent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SqlManagementMode(input)
	return &out, nil
}

type SqlServerLicenseType string

const (
	SqlServerLicenseTypeAHUB SqlServerLicenseType = "AHUB"
	SqlServerLicenseTypeDR   SqlServerLicenseType = "DR"
	SqlServerLicenseTypePAYG SqlServerLicenseType = "PAYG"
)

func PossibleValuesForSqlServerLicenseType() []string {
	return []string{
		string(SqlServerLicenseTypeAHUB),
		string(SqlServerLicenseTypeDR),
		string(SqlServerLicenseTypePAYG),
	}
}

func (s *SqlServerLicenseType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSqlServerLicenseType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSqlServerLicenseType(input string) (*SqlServerLicenseType, error) {
	vals := map[string]SqlServerLicenseType{
		"ahub": SqlServerLicenseTypeAHUB,
		"dr":   SqlServerLicenseTypeDR,
		"payg": SqlServerLicenseTypePAYG,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SqlServerLicenseType(input)
	return &out, nil
}

type SqlWorkloadType string

const (
	SqlWorkloadTypeDW      SqlWorkloadType = "DW"
	SqlWorkloadTypeGENERAL SqlWorkloadType = "GENERAL"
	SqlWorkloadTypeOLTP    SqlWorkloadType = "OLTP"
)

func PossibleValuesForSqlWorkloadType() []string {
	return []string{
		string(SqlWorkloadTypeDW),
		string(SqlWorkloadTypeGENERAL),
		string(SqlWorkloadTypeOLTP),
	}
}

func (s *SqlWorkloadType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSqlWorkloadType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSqlWorkloadType(input string) (*SqlWorkloadType, error) {
	vals := map[string]SqlWorkloadType{
		"dw":      SqlWorkloadTypeDW,
		"general": SqlWorkloadTypeGENERAL,
		"oltp":    SqlWorkloadTypeOLTP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SqlWorkloadType(input)
	return &out, nil
}

type StorageWorkloadType string

const (
	StorageWorkloadTypeDW      StorageWorkloadType = "DW"
	StorageWorkloadTypeGENERAL StorageWorkloadType = "GENERAL"
	StorageWorkloadTypeOLTP    StorageWorkloadType = "OLTP"
)

func PossibleValuesForStorageWorkloadType() []string {
	return []string{
		string(StorageWorkloadTypeDW),
		string(StorageWorkloadTypeGENERAL),
		string(StorageWorkloadTypeOLTP),
	}
}

func (s *StorageWorkloadType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseStorageWorkloadType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseStorageWorkloadType(input string) (*StorageWorkloadType, error) {
	vals := map[string]StorageWorkloadType{
		"dw":      StorageWorkloadTypeDW,
		"general": StorageWorkloadTypeGENERAL,
		"oltp":    StorageWorkloadTypeOLTP,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := StorageWorkloadType(input)
	return &out, nil
}

type TroubleshootingScenario string

const (
	TroubleshootingScenarioUnhealthyReplica TroubleshootingScenario = "UnhealthyReplica"
)

func PossibleValuesForTroubleshootingScenario() []string {
	return []string{
		string(TroubleshootingScenarioUnhealthyReplica),
	}
}

func (s *TroubleshootingScenario) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTroubleshootingScenario(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTroubleshootingScenario(input string) (*TroubleshootingScenario, error) {
	vals := map[string]TroubleshootingScenario{
		"unhealthyreplica": TroubleshootingScenarioUnhealthyReplica,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TroubleshootingScenario(input)
	return &out, nil
}

type VMIdentityType string

const (
	VMIdentityTypeNone           VMIdentityType = "None"
	VMIdentityTypeSystemAssigned VMIdentityType = "SystemAssigned"
	VMIdentityTypeUserAssigned   VMIdentityType = "UserAssigned"
)

func PossibleValuesForVMIdentityType() []string {
	return []string{
		string(VMIdentityTypeNone),
		string(VMIdentityTypeSystemAssigned),
		string(VMIdentityTypeUserAssigned),
	}
}

func (s *VMIdentityType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseVMIdentityType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseVMIdentityType(input string) (*VMIdentityType, error) {
	vals := map[string]VMIdentityType{
		"none":           VMIdentityTypeNone,
		"systemassigned": VMIdentityTypeSystemAssigned,
		"userassigned":   VMIdentityTypeUserAssigned,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := VMIdentityType(input)
	return &out, nil
}

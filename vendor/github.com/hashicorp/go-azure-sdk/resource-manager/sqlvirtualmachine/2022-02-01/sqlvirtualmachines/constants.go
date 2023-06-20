package sqlvirtualmachines

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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

package protectionpolicies

import (
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DayOfWeek string

const (
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

type IAASVMPolicyType string

const (
	IAASVMPolicyTypeInvalid IAASVMPolicyType = "Invalid"
	IAASVMPolicyTypeVOne    IAASVMPolicyType = "V1"
	IAASVMPolicyTypeVTwo    IAASVMPolicyType = "V2"
)

func PossibleValuesForIAASVMPolicyType() []string {
	return []string{
		string(IAASVMPolicyTypeInvalid),
		string(IAASVMPolicyTypeVOne),
		string(IAASVMPolicyTypeVTwo),
	}
}

func parseIAASVMPolicyType(input string) (*IAASVMPolicyType, error) {
	vals := map[string]IAASVMPolicyType{
		"invalid": IAASVMPolicyTypeInvalid,
		"v1":      IAASVMPolicyTypeVOne,
		"v2":      IAASVMPolicyTypeVTwo,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := IAASVMPolicyType(input)
	return &out, nil
}

type MonthOfYear string

const (
	MonthOfYearApril     MonthOfYear = "April"
	MonthOfYearAugust    MonthOfYear = "August"
	MonthOfYearDecember  MonthOfYear = "December"
	MonthOfYearFebruary  MonthOfYear = "February"
	MonthOfYearInvalid   MonthOfYear = "Invalid"
	MonthOfYearJanuary   MonthOfYear = "January"
	MonthOfYearJuly      MonthOfYear = "July"
	MonthOfYearJune      MonthOfYear = "June"
	MonthOfYearMarch     MonthOfYear = "March"
	MonthOfYearMay       MonthOfYear = "May"
	MonthOfYearNovember  MonthOfYear = "November"
	MonthOfYearOctober   MonthOfYear = "October"
	MonthOfYearSeptember MonthOfYear = "September"
)

func PossibleValuesForMonthOfYear() []string {
	return []string{
		string(MonthOfYearApril),
		string(MonthOfYearAugust),
		string(MonthOfYearDecember),
		string(MonthOfYearFebruary),
		string(MonthOfYearInvalid),
		string(MonthOfYearJanuary),
		string(MonthOfYearJuly),
		string(MonthOfYearJune),
		string(MonthOfYearMarch),
		string(MonthOfYearMay),
		string(MonthOfYearNovember),
		string(MonthOfYearOctober),
		string(MonthOfYearSeptember),
	}
}

func parseMonthOfYear(input string) (*MonthOfYear, error) {
	vals := map[string]MonthOfYear{
		"april":     MonthOfYearApril,
		"august":    MonthOfYearAugust,
		"december":  MonthOfYearDecember,
		"february":  MonthOfYearFebruary,
		"invalid":   MonthOfYearInvalid,
		"january":   MonthOfYearJanuary,
		"july":      MonthOfYearJuly,
		"june":      MonthOfYearJune,
		"march":     MonthOfYearMarch,
		"may":       MonthOfYearMay,
		"november":  MonthOfYearNovember,
		"october":   MonthOfYearOctober,
		"september": MonthOfYearSeptember,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonthOfYear(input)
	return &out, nil
}

type PolicyType string

const (
	PolicyTypeCopyOnlyFull         PolicyType = "CopyOnlyFull"
	PolicyTypeDifferential         PolicyType = "Differential"
	PolicyTypeFull                 PolicyType = "Full"
	PolicyTypeIncremental          PolicyType = "Incremental"
	PolicyTypeInvalid              PolicyType = "Invalid"
	PolicyTypeLog                  PolicyType = "Log"
	PolicyTypeSnapshotCopyOnlyFull PolicyType = "SnapshotCopyOnlyFull"
	PolicyTypeSnapshotFull         PolicyType = "SnapshotFull"
)

func PossibleValuesForPolicyType() []string {
	return []string{
		string(PolicyTypeCopyOnlyFull),
		string(PolicyTypeDifferential),
		string(PolicyTypeFull),
		string(PolicyTypeIncremental),
		string(PolicyTypeInvalid),
		string(PolicyTypeLog),
		string(PolicyTypeSnapshotCopyOnlyFull),
		string(PolicyTypeSnapshotFull),
	}
}

func parsePolicyType(input string) (*PolicyType, error) {
	vals := map[string]PolicyType{
		"copyonlyfull":         PolicyTypeCopyOnlyFull,
		"differential":         PolicyTypeDifferential,
		"full":                 PolicyTypeFull,
		"incremental":          PolicyTypeIncremental,
		"invalid":              PolicyTypeInvalid,
		"log":                  PolicyTypeLog,
		"snapshotcopyonlyfull": PolicyTypeSnapshotCopyOnlyFull,
		"snapshotfull":         PolicyTypeSnapshotFull,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := PolicyType(input)
	return &out, nil
}

type RetentionDurationType string

const (
	RetentionDurationTypeDays    RetentionDurationType = "Days"
	RetentionDurationTypeInvalid RetentionDurationType = "Invalid"
	RetentionDurationTypeMonths  RetentionDurationType = "Months"
	RetentionDurationTypeWeeks   RetentionDurationType = "Weeks"
	RetentionDurationTypeYears   RetentionDurationType = "Years"
)

func PossibleValuesForRetentionDurationType() []string {
	return []string{
		string(RetentionDurationTypeDays),
		string(RetentionDurationTypeInvalid),
		string(RetentionDurationTypeMonths),
		string(RetentionDurationTypeWeeks),
		string(RetentionDurationTypeYears),
	}
}

func parseRetentionDurationType(input string) (*RetentionDurationType, error) {
	vals := map[string]RetentionDurationType{
		"days":    RetentionDurationTypeDays,
		"invalid": RetentionDurationTypeInvalid,
		"months":  RetentionDurationTypeMonths,
		"weeks":   RetentionDurationTypeWeeks,
		"years":   RetentionDurationTypeYears,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RetentionDurationType(input)
	return &out, nil
}

type RetentionScheduleFormat string

const (
	RetentionScheduleFormatDaily   RetentionScheduleFormat = "Daily"
	RetentionScheduleFormatInvalid RetentionScheduleFormat = "Invalid"
	RetentionScheduleFormatWeekly  RetentionScheduleFormat = "Weekly"
)

func PossibleValuesForRetentionScheduleFormat() []string {
	return []string{
		string(RetentionScheduleFormatDaily),
		string(RetentionScheduleFormatInvalid),
		string(RetentionScheduleFormatWeekly),
	}
}

func parseRetentionScheduleFormat(input string) (*RetentionScheduleFormat, error) {
	vals := map[string]RetentionScheduleFormat{
		"daily":   RetentionScheduleFormatDaily,
		"invalid": RetentionScheduleFormatInvalid,
		"weekly":  RetentionScheduleFormatWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RetentionScheduleFormat(input)
	return &out, nil
}

type ScheduleRunType string

const (
	ScheduleRunTypeDaily   ScheduleRunType = "Daily"
	ScheduleRunTypeHourly  ScheduleRunType = "Hourly"
	ScheduleRunTypeInvalid ScheduleRunType = "Invalid"
	ScheduleRunTypeWeekly  ScheduleRunType = "Weekly"
)

func PossibleValuesForScheduleRunType() []string {
	return []string{
		string(ScheduleRunTypeDaily),
		string(ScheduleRunTypeHourly),
		string(ScheduleRunTypeInvalid),
		string(ScheduleRunTypeWeekly),
	}
}

func parseScheduleRunType(input string) (*ScheduleRunType, error) {
	vals := map[string]ScheduleRunType{
		"daily":   ScheduleRunTypeDaily,
		"hourly":  ScheduleRunTypeHourly,
		"invalid": ScheduleRunTypeInvalid,
		"weekly":  ScheduleRunTypeWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleRunType(input)
	return &out, nil
}

type TieringMode string

const (
	TieringModeDoNotTier       TieringMode = "DoNotTier"
	TieringModeInvalid         TieringMode = "Invalid"
	TieringModeTierAfter       TieringMode = "TierAfter"
	TieringModeTierRecommended TieringMode = "TierRecommended"
)

func PossibleValuesForTieringMode() []string {
	return []string{
		string(TieringModeDoNotTier),
		string(TieringModeInvalid),
		string(TieringModeTierAfter),
		string(TieringModeTierRecommended),
	}
}

func parseTieringMode(input string) (*TieringMode, error) {
	vals := map[string]TieringMode{
		"donottier":       TieringModeDoNotTier,
		"invalid":         TieringModeInvalid,
		"tierafter":       TieringModeTierAfter,
		"tierrecommended": TieringModeTierRecommended,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TieringMode(input)
	return &out, nil
}

type WeekOfMonth string

const (
	WeekOfMonthFirst   WeekOfMonth = "First"
	WeekOfMonthFourth  WeekOfMonth = "Fourth"
	WeekOfMonthInvalid WeekOfMonth = "Invalid"
	WeekOfMonthLast    WeekOfMonth = "Last"
	WeekOfMonthSecond  WeekOfMonth = "Second"
	WeekOfMonthThird   WeekOfMonth = "Third"
)

func PossibleValuesForWeekOfMonth() []string {
	return []string{
		string(WeekOfMonthFirst),
		string(WeekOfMonthFourth),
		string(WeekOfMonthInvalid),
		string(WeekOfMonthLast),
		string(WeekOfMonthSecond),
		string(WeekOfMonthThird),
	}
}

func parseWeekOfMonth(input string) (*WeekOfMonth, error) {
	vals := map[string]WeekOfMonth{
		"first":   WeekOfMonthFirst,
		"fourth":  WeekOfMonthFourth,
		"invalid": WeekOfMonthInvalid,
		"last":    WeekOfMonthLast,
		"second":  WeekOfMonthSecond,
		"third":   WeekOfMonthThird,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WeekOfMonth(input)
	return &out, nil
}

type WorkloadType string

const (
	WorkloadTypeAzureFileShare    WorkloadType = "AzureFileShare"
	WorkloadTypeAzureSqlDb        WorkloadType = "AzureSqlDb"
	WorkloadTypeClient            WorkloadType = "Client"
	WorkloadTypeExchange          WorkloadType = "Exchange"
	WorkloadTypeFileFolder        WorkloadType = "FileFolder"
	WorkloadTypeGenericDataSource WorkloadType = "GenericDataSource"
	WorkloadTypeInvalid           WorkloadType = "Invalid"
	WorkloadTypeSAPAseDatabase    WorkloadType = "SAPAseDatabase"
	WorkloadTypeSAPHanaDBInstance WorkloadType = "SAPHanaDBInstance"
	WorkloadTypeSAPHanaDatabase   WorkloadType = "SAPHanaDatabase"
	WorkloadTypeSQLDB             WorkloadType = "SQLDB"
	WorkloadTypeSQLDataBase       WorkloadType = "SQLDataBase"
	WorkloadTypeSharepoint        WorkloadType = "Sharepoint"
	WorkloadTypeSystemState       WorkloadType = "SystemState"
	WorkloadTypeVM                WorkloadType = "VM"
	WorkloadTypeVMwareVM          WorkloadType = "VMwareVM"
)

func PossibleValuesForWorkloadType() []string {
	return []string{
		string(WorkloadTypeAzureFileShare),
		string(WorkloadTypeAzureSqlDb),
		string(WorkloadTypeClient),
		string(WorkloadTypeExchange),
		string(WorkloadTypeFileFolder),
		string(WorkloadTypeGenericDataSource),
		string(WorkloadTypeInvalid),
		string(WorkloadTypeSAPAseDatabase),
		string(WorkloadTypeSAPHanaDBInstance),
		string(WorkloadTypeSAPHanaDatabase),
		string(WorkloadTypeSQLDB),
		string(WorkloadTypeSQLDataBase),
		string(WorkloadTypeSharepoint),
		string(WorkloadTypeSystemState),
		string(WorkloadTypeVM),
		string(WorkloadTypeVMwareVM),
	}
}

func parseWorkloadType(input string) (*WorkloadType, error) {
	vals := map[string]WorkloadType{
		"azurefileshare":    WorkloadTypeAzureFileShare,
		"azuresqldb":        WorkloadTypeAzureSqlDb,
		"client":            WorkloadTypeClient,
		"exchange":          WorkloadTypeExchange,
		"filefolder":        WorkloadTypeFileFolder,
		"genericdatasource": WorkloadTypeGenericDataSource,
		"invalid":           WorkloadTypeInvalid,
		"sapasedatabase":    WorkloadTypeSAPAseDatabase,
		"saphanadbinstance": WorkloadTypeSAPHanaDBInstance,
		"saphanadatabase":   WorkloadTypeSAPHanaDatabase,
		"sqldb":             WorkloadTypeSQLDB,
		"sqldatabase":       WorkloadTypeSQLDataBase,
		"sharepoint":        WorkloadTypeSharepoint,
		"systemstate":       WorkloadTypeSystemState,
		"vm":                WorkloadTypeVM,
		"vmwarevm":          WorkloadTypeVMwareVM,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WorkloadType(input)
	return &out, nil
}

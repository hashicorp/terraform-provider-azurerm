package backuppolicies

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AbsoluteMarker string

const (
	AbsoluteMarkerAllBackup    AbsoluteMarker = "AllBackup"
	AbsoluteMarkerFirstOfDay   AbsoluteMarker = "FirstOfDay"
	AbsoluteMarkerFirstOfMonth AbsoluteMarker = "FirstOfMonth"
	AbsoluteMarkerFirstOfWeek  AbsoluteMarker = "FirstOfWeek"
	AbsoluteMarkerFirstOfYear  AbsoluteMarker = "FirstOfYear"
)

func PossibleValuesForAbsoluteMarker() []string {
	return []string{
		string(AbsoluteMarkerAllBackup),
		string(AbsoluteMarkerFirstOfDay),
		string(AbsoluteMarkerFirstOfMonth),
		string(AbsoluteMarkerFirstOfWeek),
		string(AbsoluteMarkerFirstOfYear),
	}
}

func (s *AbsoluteMarker) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAbsoluteMarker(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAbsoluteMarker(input string) (*AbsoluteMarker, error) {
	vals := map[string]AbsoluteMarker{
		"allbackup":    AbsoluteMarkerAllBackup,
		"firstofday":   AbsoluteMarkerFirstOfDay,
		"firstofmonth": AbsoluteMarkerFirstOfMonth,
		"firstofweek":  AbsoluteMarkerFirstOfWeek,
		"firstofyear":  AbsoluteMarkerFirstOfYear,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AbsoluteMarker(input)
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

type Month string

const (
	MonthApril     Month = "April"
	MonthAugust    Month = "August"
	MonthDecember  Month = "December"
	MonthFebruary  Month = "February"
	MonthJanuary   Month = "January"
	MonthJuly      Month = "July"
	MonthJune      Month = "June"
	MonthMarch     Month = "March"
	MonthMay       Month = "May"
	MonthNovember  Month = "November"
	MonthOctober   Month = "October"
	MonthSeptember Month = "September"
)

func PossibleValuesForMonth() []string {
	return []string{
		string(MonthApril),
		string(MonthAugust),
		string(MonthDecember),
		string(MonthFebruary),
		string(MonthJanuary),
		string(MonthJuly),
		string(MonthJune),
		string(MonthMarch),
		string(MonthMay),
		string(MonthNovember),
		string(MonthOctober),
		string(MonthSeptember),
	}
}

func (s *Month) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonth(input string) (*Month, error) {
	vals := map[string]Month{
		"april":     MonthApril,
		"august":    MonthAugust,
		"december":  MonthDecember,
		"february":  MonthFebruary,
		"january":   MonthJanuary,
		"july":      MonthJuly,
		"june":      MonthJune,
		"march":     MonthMarch,
		"may":       MonthMay,
		"november":  MonthNovember,
		"october":   MonthOctober,
		"september": MonthSeptember,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Month(input)
	return &out, nil
}

type WeekNumber string

const (
	WeekNumberFirst  WeekNumber = "First"
	WeekNumberFourth WeekNumber = "Fourth"
	WeekNumberLast   WeekNumber = "Last"
	WeekNumberSecond WeekNumber = "Second"
	WeekNumberThird  WeekNumber = "Third"
)

func PossibleValuesForWeekNumber() []string {
	return []string{
		string(WeekNumberFirst),
		string(WeekNumberFourth),
		string(WeekNumberLast),
		string(WeekNumberSecond),
		string(WeekNumberThird),
	}
}

func (s *WeekNumber) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWeekNumber(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWeekNumber(input string) (*WeekNumber, error) {
	vals := map[string]WeekNumber{
		"first":  WeekNumberFirst,
		"fourth": WeekNumberFourth,
		"last":   WeekNumberLast,
		"second": WeekNumberSecond,
		"third":  WeekNumberThird,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WeekNumber(input)
	return &out, nil
}

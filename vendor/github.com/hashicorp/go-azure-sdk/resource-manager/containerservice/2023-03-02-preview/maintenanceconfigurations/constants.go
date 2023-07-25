package maintenanceconfigurations

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Type string

const (
	TypeFirst  Type = "First"
	TypeFourth Type = "Fourth"
	TypeLast   Type = "Last"
	TypeSecond Type = "Second"
	TypeThird  Type = "Third"
)

func PossibleValuesForType() []string {
	return []string{
		string(TypeFirst),
		string(TypeFourth),
		string(TypeLast),
		string(TypeSecond),
		string(TypeThird),
	}
}

func parseType(input string) (*Type, error) {
	vals := map[string]Type{
		"first":  TypeFirst,
		"fourth": TypeFourth,
		"last":   TypeLast,
		"second": TypeSecond,
		"third":  TypeThird,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Type(input)
	return &out, nil
}

type WeekDay string

const (
	WeekDayFriday    WeekDay = "Friday"
	WeekDayMonday    WeekDay = "Monday"
	WeekDaySaturday  WeekDay = "Saturday"
	WeekDaySunday    WeekDay = "Sunday"
	WeekDayThursday  WeekDay = "Thursday"
	WeekDayTuesday   WeekDay = "Tuesday"
	WeekDayWednesday WeekDay = "Wednesday"
)

func PossibleValuesForWeekDay() []string {
	return []string{
		string(WeekDayFriday),
		string(WeekDayMonday),
		string(WeekDaySaturday),
		string(WeekDaySunday),
		string(WeekDayThursday),
		string(WeekDayTuesday),
		string(WeekDayWednesday),
	}
}

func parseWeekDay(input string) (*WeekDay, error) {
	vals := map[string]WeekDay{
		"friday":    WeekDayFriday,
		"monday":    WeekDayMonday,
		"saturday":  WeekDaySaturday,
		"sunday":    WeekDaySunday,
		"thursday":  WeekDayThursday,
		"tuesday":   WeekDayTuesday,
		"wednesday": WeekDayWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WeekDay(input)
	return &out, nil
}

package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduleDay string

const (
	ScheduleDayFriday    ScheduleDay = "Friday"
	ScheduleDayMonday    ScheduleDay = "Monday"
	ScheduleDaySaturday  ScheduleDay = "Saturday"
	ScheduleDaySunday    ScheduleDay = "Sunday"
	ScheduleDayThursday  ScheduleDay = "Thursday"
	ScheduleDayTuesday   ScheduleDay = "Tuesday"
	ScheduleDayWednesday ScheduleDay = "Wednesday"
)

func PossibleValuesForScheduleDay() []string {
	return []string{
		string(ScheduleDayFriday),
		string(ScheduleDayMonday),
		string(ScheduleDaySaturday),
		string(ScheduleDaySunday),
		string(ScheduleDayThursday),
		string(ScheduleDayTuesday),
		string(ScheduleDayWednesday),
	}
}

func (s *ScheduleDay) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduleDay(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduleDay(input string) (*ScheduleDay, error) {
	vals := map[string]ScheduleDay{
		"friday":    ScheduleDayFriday,
		"monday":    ScheduleDayMonday,
		"saturday":  ScheduleDaySaturday,
		"sunday":    ScheduleDaySunday,
		"thursday":  ScheduleDayThursday,
		"tuesday":   ScheduleDayTuesday,
		"wednesday": ScheduleDayWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleDay(input)
	return &out, nil
}

type ScheduleFrequency string

const (
	ScheduleFrequencyDay     ScheduleFrequency = "Day"
	ScheduleFrequencyHour    ScheduleFrequency = "Hour"
	ScheduleFrequencyMinute  ScheduleFrequency = "Minute"
	ScheduleFrequencyMonth   ScheduleFrequency = "Month"
	ScheduleFrequencyOneTime ScheduleFrequency = "OneTime"
	ScheduleFrequencyWeek    ScheduleFrequency = "Week"
)

func PossibleValuesForScheduleFrequency() []string {
	return []string{
		string(ScheduleFrequencyDay),
		string(ScheduleFrequencyHour),
		string(ScheduleFrequencyMinute),
		string(ScheduleFrequencyMonth),
		string(ScheduleFrequencyOneTime),
		string(ScheduleFrequencyWeek),
	}
}

func (s *ScheduleFrequency) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduleFrequency(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduleFrequency(input string) (*ScheduleFrequency, error) {
	vals := map[string]ScheduleFrequency{
		"day":     ScheduleFrequencyDay,
		"hour":    ScheduleFrequencyHour,
		"minute":  ScheduleFrequencyMinute,
		"month":   ScheduleFrequencyMonth,
		"onetime": ScheduleFrequencyOneTime,
		"week":    ScheduleFrequencyWeek,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleFrequency(input)
	return &out, nil
}

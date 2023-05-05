package scheduledactions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CheckNameAvailabilityReason string

const (
	CheckNameAvailabilityReasonAlreadyExists CheckNameAvailabilityReason = "AlreadyExists"
	CheckNameAvailabilityReasonInvalid       CheckNameAvailabilityReason = "Invalid"
)

func PossibleValuesForCheckNameAvailabilityReason() []string {
	return []string{
		string(CheckNameAvailabilityReasonAlreadyExists),
		string(CheckNameAvailabilityReasonInvalid),
	}
}

func (s *CheckNameAvailabilityReason) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseCheckNameAvailabilityReason(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseCheckNameAvailabilityReason(input string) (*CheckNameAvailabilityReason, error) {
	vals := map[string]CheckNameAvailabilityReason{
		"alreadyexists": CheckNameAvailabilityReasonAlreadyExists,
		"invalid":       CheckNameAvailabilityReasonInvalid,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := CheckNameAvailabilityReason(input)
	return &out, nil
}

type DaysOfWeek string

const (
	DaysOfWeekFriday    DaysOfWeek = "Friday"
	DaysOfWeekMonday    DaysOfWeek = "Monday"
	DaysOfWeekSaturday  DaysOfWeek = "Saturday"
	DaysOfWeekSunday    DaysOfWeek = "Sunday"
	DaysOfWeekThursday  DaysOfWeek = "Thursday"
	DaysOfWeekTuesday   DaysOfWeek = "Tuesday"
	DaysOfWeekWednesday DaysOfWeek = "Wednesday"
)

func PossibleValuesForDaysOfWeek() []string {
	return []string{
		string(DaysOfWeekFriday),
		string(DaysOfWeekMonday),
		string(DaysOfWeekSaturday),
		string(DaysOfWeekSunday),
		string(DaysOfWeekThursday),
		string(DaysOfWeekTuesday),
		string(DaysOfWeekWednesday),
	}
}

func (s *DaysOfWeek) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDaysOfWeek(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDaysOfWeek(input string) (*DaysOfWeek, error) {
	vals := map[string]DaysOfWeek{
		"friday":    DaysOfWeekFriday,
		"monday":    DaysOfWeekMonday,
		"saturday":  DaysOfWeekSaturday,
		"sunday":    DaysOfWeekSunday,
		"thursday":  DaysOfWeekThursday,
		"tuesday":   DaysOfWeekTuesday,
		"wednesday": DaysOfWeekWednesday,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DaysOfWeek(input)
	return &out, nil
}

type FileFormat string

const (
	FileFormatCsv FileFormat = "Csv"
)

func PossibleValuesForFileFormat() []string {
	return []string{
		string(FileFormatCsv),
	}
}

func (s *FileFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseFileFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseFileFormat(input string) (*FileFormat, error) {
	vals := map[string]FileFormat{
		"csv": FileFormatCsv,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := FileFormat(input)
	return &out, nil
}

type ScheduleFrequency string

const (
	ScheduleFrequencyDaily   ScheduleFrequency = "Daily"
	ScheduleFrequencyMonthly ScheduleFrequency = "Monthly"
	ScheduleFrequencyWeekly  ScheduleFrequency = "Weekly"
)

func PossibleValuesForScheduleFrequency() []string {
	return []string{
		string(ScheduleFrequencyDaily),
		string(ScheduleFrequencyMonthly),
		string(ScheduleFrequencyWeekly),
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
		"daily":   ScheduleFrequencyDaily,
		"monthly": ScheduleFrequencyMonthly,
		"weekly":  ScheduleFrequencyWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduleFrequency(input)
	return &out, nil
}

type ScheduledActionKind string

const (
	ScheduledActionKindEmail        ScheduledActionKind = "Email"
	ScheduledActionKindInsightAlert ScheduledActionKind = "InsightAlert"
)

func PossibleValuesForScheduledActionKind() []string {
	return []string{
		string(ScheduledActionKindEmail),
		string(ScheduledActionKindInsightAlert),
	}
}

func (s *ScheduledActionKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduledActionKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduledActionKind(input string) (*ScheduledActionKind, error) {
	vals := map[string]ScheduledActionKind{
		"email":        ScheduledActionKindEmail,
		"insightalert": ScheduledActionKindInsightAlert,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduledActionKind(input)
	return &out, nil
}

type ScheduledActionStatus string

const (
	ScheduledActionStatusDisabled ScheduledActionStatus = "Disabled"
	ScheduledActionStatusEnabled  ScheduledActionStatus = "Enabled"
	ScheduledActionStatusExpired  ScheduledActionStatus = "Expired"
)

func PossibleValuesForScheduledActionStatus() []string {
	return []string{
		string(ScheduledActionStatusDisabled),
		string(ScheduledActionStatusEnabled),
		string(ScheduledActionStatusExpired),
	}
}

func (s *ScheduledActionStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScheduledActionStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScheduledActionStatus(input string) (*ScheduledActionStatus, error) {
	vals := map[string]ScheduledActionStatus{
		"disabled": ScheduledActionStatusDisabled,
		"enabled":  ScheduledActionStatusEnabled,
		"expired":  ScheduledActionStatusExpired,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScheduledActionStatus(input)
	return &out, nil
}

type WeeksOfMonth string

const (
	WeeksOfMonthFirst  WeeksOfMonth = "First"
	WeeksOfMonthFourth WeeksOfMonth = "Fourth"
	WeeksOfMonthLast   WeeksOfMonth = "Last"
	WeeksOfMonthSecond WeeksOfMonth = "Second"
	WeeksOfMonthThird  WeeksOfMonth = "Third"
)

func PossibleValuesForWeeksOfMonth() []string {
	return []string{
		string(WeeksOfMonthFirst),
		string(WeeksOfMonthFourth),
		string(WeeksOfMonthLast),
		string(WeeksOfMonthSecond),
		string(WeeksOfMonthThird),
	}
}

func (s *WeeksOfMonth) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseWeeksOfMonth(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseWeeksOfMonth(input string) (*WeeksOfMonth, error) {
	vals := map[string]WeeksOfMonth{
		"first":  WeeksOfMonthFirst,
		"fourth": WeeksOfMonthFourth,
		"last":   WeeksOfMonthLast,
		"second": WeeksOfMonthSecond,
		"third":  WeeksOfMonthThird,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := WeeksOfMonth(input)
	return &out, nil
}

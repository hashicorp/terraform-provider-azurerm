package alertprocessingrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionType string

const (
	ActionTypeAddActionGroups       ActionType = "AddActionGroups"
	ActionTypeRemoveAllActionGroups ActionType = "RemoveAllActionGroups"
)

func PossibleValuesForActionType() []string {
	return []string{
		string(ActionTypeAddActionGroups),
		string(ActionTypeRemoveAllActionGroups),
	}
}

func (s *ActionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionType(input string) (*ActionType, error) {
	vals := map[string]ActionType{
		"addactiongroups":       ActionTypeAddActionGroups,
		"removeallactiongroups": ActionTypeRemoveAllActionGroups,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionType(input)
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

type Field string

const (
	FieldAlertContext        Field = "AlertContext"
	FieldAlertRuleId         Field = "AlertRuleId"
	FieldAlertRuleName       Field = "AlertRuleName"
	FieldDescription         Field = "Description"
	FieldMonitorCondition    Field = "MonitorCondition"
	FieldMonitorService      Field = "MonitorService"
	FieldSeverity            Field = "Severity"
	FieldSignalType          Field = "SignalType"
	FieldTargetResource      Field = "TargetResource"
	FieldTargetResourceGroup Field = "TargetResourceGroup"
	FieldTargetResourceType  Field = "TargetResourceType"
)

func PossibleValuesForField() []string {
	return []string{
		string(FieldAlertContext),
		string(FieldAlertRuleId),
		string(FieldAlertRuleName),
		string(FieldDescription),
		string(FieldMonitorCondition),
		string(FieldMonitorService),
		string(FieldSeverity),
		string(FieldSignalType),
		string(FieldTargetResource),
		string(FieldTargetResourceGroup),
		string(FieldTargetResourceType),
	}
}

func (s *Field) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseField(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseField(input string) (*Field, error) {
	vals := map[string]Field{
		"alertcontext":        FieldAlertContext,
		"alertruleid":         FieldAlertRuleId,
		"alertrulename":       FieldAlertRuleName,
		"description":         FieldDescription,
		"monitorcondition":    FieldMonitorCondition,
		"monitorservice":      FieldMonitorService,
		"severity":            FieldSeverity,
		"signaltype":          FieldSignalType,
		"targetresource":      FieldTargetResource,
		"targetresourcegroup": FieldTargetResourceGroup,
		"targetresourcetype":  FieldTargetResourceType,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Field(input)
	return &out, nil
}

type Operator string

const (
	OperatorContains       Operator = "Contains"
	OperatorDoesNotContain Operator = "DoesNotContain"
	OperatorEquals         Operator = "Equals"
	OperatorNotEquals      Operator = "NotEquals"
)

func PossibleValuesForOperator() []string {
	return []string{
		string(OperatorContains),
		string(OperatorDoesNotContain),
		string(OperatorEquals),
		string(OperatorNotEquals),
	}
}

func (s *Operator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseOperator(input string) (*Operator, error) {
	vals := map[string]Operator{
		"contains":       OperatorContains,
		"doesnotcontain": OperatorDoesNotContain,
		"equals":         OperatorEquals,
		"notequals":      OperatorNotEquals,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operator(input)
	return &out, nil
}

type RecurrenceType string

const (
	RecurrenceTypeDaily   RecurrenceType = "Daily"
	RecurrenceTypeMonthly RecurrenceType = "Monthly"
	RecurrenceTypeWeekly  RecurrenceType = "Weekly"
)

func PossibleValuesForRecurrenceType() []string {
	return []string{
		string(RecurrenceTypeDaily),
		string(RecurrenceTypeMonthly),
		string(RecurrenceTypeWeekly),
	}
}

func (s *RecurrenceType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseRecurrenceType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseRecurrenceType(input string) (*RecurrenceType, error) {
	vals := map[string]RecurrenceType{
		"daily":   RecurrenceTypeDaily,
		"monthly": RecurrenceTypeMonthly,
		"weekly":  RecurrenceTypeWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := RecurrenceType(input)
	return &out, nil
}

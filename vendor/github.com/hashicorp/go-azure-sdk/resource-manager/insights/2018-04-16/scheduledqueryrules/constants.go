package scheduledqueryrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertSeverity string

const (
	AlertSeverityFour  AlertSeverity = "4"
	AlertSeverityOne   AlertSeverity = "1"
	AlertSeverityThree AlertSeverity = "3"
	AlertSeverityTwo   AlertSeverity = "2"
	AlertSeverityZero  AlertSeverity = "0"
)

func PossibleValuesForAlertSeverity() []string {
	return []string{
		string(AlertSeverityFour),
		string(AlertSeverityOne),
		string(AlertSeverityThree),
		string(AlertSeverityTwo),
		string(AlertSeverityZero),
	}
}

func (s *AlertSeverity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlertSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlertSeverity(input string) (*AlertSeverity, error) {
	vals := map[string]AlertSeverity{
		"4": AlertSeverityFour,
		"1": AlertSeverityOne,
		"3": AlertSeverityThree,
		"2": AlertSeverityTwo,
		"0": AlertSeverityZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertSeverity(input)
	return &out, nil
}

type ConditionalOperator string

const (
	ConditionalOperatorEqual              ConditionalOperator = "Equal"
	ConditionalOperatorGreaterThan        ConditionalOperator = "GreaterThan"
	ConditionalOperatorGreaterThanOrEqual ConditionalOperator = "GreaterThanOrEqual"
	ConditionalOperatorLessThan           ConditionalOperator = "LessThan"
	ConditionalOperatorLessThanOrEqual    ConditionalOperator = "LessThanOrEqual"
)

func PossibleValuesForConditionalOperator() []string {
	return []string{
		string(ConditionalOperatorEqual),
		string(ConditionalOperatorGreaterThan),
		string(ConditionalOperatorGreaterThanOrEqual),
		string(ConditionalOperatorLessThan),
		string(ConditionalOperatorLessThanOrEqual),
	}
}

func (s *ConditionalOperator) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseConditionalOperator(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseConditionalOperator(input string) (*ConditionalOperator, error) {
	vals := map[string]ConditionalOperator{
		"equal":              ConditionalOperatorEqual,
		"greaterthan":        ConditionalOperatorGreaterThan,
		"greaterthanorequal": ConditionalOperatorGreaterThanOrEqual,
		"lessthan":           ConditionalOperatorLessThan,
		"lessthanorequal":    ConditionalOperatorLessThanOrEqual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ConditionalOperator(input)
	return &out, nil
}

type Enabled string

const (
	EnabledFalse Enabled = "false"
	EnabledTrue  Enabled = "true"
)

func PossibleValuesForEnabled() []string {
	return []string{
		string(EnabledFalse),
		string(EnabledTrue),
	}
}

func (s *Enabled) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseEnabled(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseEnabled(input string) (*Enabled, error) {
	vals := map[string]Enabled{
		"false": EnabledFalse,
		"true":  EnabledTrue,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Enabled(input)
	return &out, nil
}

type MetricTriggerType string

const (
	MetricTriggerTypeConsecutive MetricTriggerType = "Consecutive"
	MetricTriggerTypeTotal       MetricTriggerType = "Total"
)

func PossibleValuesForMetricTriggerType() []string {
	return []string{
		string(MetricTriggerTypeConsecutive),
		string(MetricTriggerTypeTotal),
	}
}

func (s *MetricTriggerType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMetricTriggerType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMetricTriggerType(input string) (*MetricTriggerType, error) {
	vals := map[string]MetricTriggerType{
		"consecutive": MetricTriggerTypeConsecutive,
		"total":       MetricTriggerTypeTotal,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetricTriggerType(input)
	return &out, nil
}

type Operator string

const (
	OperatorInclude Operator = "Include"
)

func PossibleValuesForOperator() []string {
	return []string{
		string(OperatorInclude),
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
		"include": OperatorInclude,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Operator(input)
	return &out, nil
}

type ProvisioningState string

const (
	ProvisioningStateCanceled  ProvisioningState = "Canceled"
	ProvisioningStateDeploying ProvisioningState = "Deploying"
	ProvisioningStateFailed    ProvisioningState = "Failed"
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
)

func PossibleValuesForProvisioningState() []string {
	return []string{
		string(ProvisioningStateCanceled),
		string(ProvisioningStateDeploying),
		string(ProvisioningStateFailed),
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
		"canceled":  ProvisioningStateCanceled,
		"deploying": ProvisioningStateDeploying,
		"failed":    ProvisioningStateFailed,
		"succeeded": ProvisioningStateSucceeded,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ProvisioningState(input)
	return &out, nil
}

type QueryType string

const (
	QueryTypeResultCount QueryType = "ResultCount"
)

func PossibleValuesForQueryType() []string {
	return []string{
		string(QueryTypeResultCount),
	}
}

func (s *QueryType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseQueryType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseQueryType(input string) (*QueryType, error) {
	vals := map[string]QueryType{
		"resultcount": QueryTypeResultCount,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := QueryType(input)
	return &out, nil
}

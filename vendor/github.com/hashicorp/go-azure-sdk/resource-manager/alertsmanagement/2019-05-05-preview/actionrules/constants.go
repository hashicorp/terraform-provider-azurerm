package actionrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ActionRuleStatus string

const (
	ActionRuleStatusDisabled ActionRuleStatus = "Disabled"
	ActionRuleStatusEnabled  ActionRuleStatus = "Enabled"
)

func PossibleValuesForActionRuleStatus() []string {
	return []string{
		string(ActionRuleStatusDisabled),
		string(ActionRuleStatusEnabled),
	}
}

func (s *ActionRuleStatus) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionRuleStatus(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionRuleStatus(input string) (*ActionRuleStatus, error) {
	vals := map[string]ActionRuleStatus{
		"disabled": ActionRuleStatusDisabled,
		"enabled":  ActionRuleStatusEnabled,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionRuleStatus(input)
	return &out, nil
}

type ActionRuleType string

const (
	ActionRuleTypeActionGroup ActionRuleType = "ActionGroup"
	ActionRuleTypeDiagnostics ActionRuleType = "Diagnostics"
	ActionRuleTypeSuppression ActionRuleType = "Suppression"
)

func PossibleValuesForActionRuleType() []string {
	return []string{
		string(ActionRuleTypeActionGroup),
		string(ActionRuleTypeDiagnostics),
		string(ActionRuleTypeSuppression),
	}
}

func (s *ActionRuleType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseActionRuleType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseActionRuleType(input string) (*ActionRuleType, error) {
	vals := map[string]ActionRuleType{
		"actiongroup": ActionRuleTypeActionGroup,
		"diagnostics": ActionRuleTypeDiagnostics,
		"suppression": ActionRuleTypeSuppression,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ActionRuleType(input)
	return &out, nil
}

type MonitorService string

const (
	MonitorServiceActivityLogAdministrative MonitorService = "ActivityLog Administrative"
	MonitorServiceActivityLogAutoscale      MonitorService = "ActivityLog Autoscale"
	MonitorServiceActivityLogPolicy         MonitorService = "ActivityLog Policy"
	MonitorServiceActivityLogRecommendation MonitorService = "ActivityLog Recommendation"
	MonitorServiceActivityLogSecurity       MonitorService = "ActivityLog Security"
	MonitorServiceApplicationInsights       MonitorService = "Application Insights"
	MonitorServiceLogAnalytics              MonitorService = "Log Analytics"
	MonitorServiceNagios                    MonitorService = "Nagios"
	MonitorServicePlatform                  MonitorService = "Platform"
	MonitorServiceSCOM                      MonitorService = "SCOM"
	MonitorServiceServiceHealth             MonitorService = "ServiceHealth"
	MonitorServiceSmartDetector             MonitorService = "SmartDetector"
	MonitorServiceVMInsights                MonitorService = "VM Insights"
	MonitorServiceZabbix                    MonitorService = "Zabbix"
)

func PossibleValuesForMonitorService() []string {
	return []string{
		string(MonitorServiceActivityLogAdministrative),
		string(MonitorServiceActivityLogAutoscale),
		string(MonitorServiceActivityLogPolicy),
		string(MonitorServiceActivityLogRecommendation),
		string(MonitorServiceActivityLogSecurity),
		string(MonitorServiceApplicationInsights),
		string(MonitorServiceLogAnalytics),
		string(MonitorServiceNagios),
		string(MonitorServicePlatform),
		string(MonitorServiceSCOM),
		string(MonitorServiceServiceHealth),
		string(MonitorServiceSmartDetector),
		string(MonitorServiceVMInsights),
		string(MonitorServiceZabbix),
	}
}

func (s *MonitorService) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitorService(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitorService(input string) (*MonitorService, error) {
	vals := map[string]MonitorService{
		"activitylog administrative": MonitorServiceActivityLogAdministrative,
		"activitylog autoscale":      MonitorServiceActivityLogAutoscale,
		"activitylog policy":         MonitorServiceActivityLogPolicy,
		"activitylog recommendation": MonitorServiceActivityLogRecommendation,
		"activitylog security":       MonitorServiceActivityLogSecurity,
		"application insights":       MonitorServiceApplicationInsights,
		"log analytics":              MonitorServiceLogAnalytics,
		"nagios":                     MonitorServiceNagios,
		"platform":                   MonitorServicePlatform,
		"scom":                       MonitorServiceSCOM,
		"servicehealth":              MonitorServiceServiceHealth,
		"smartdetector":              MonitorServiceSmartDetector,
		"vm insights":                MonitorServiceVMInsights,
		"zabbix":                     MonitorServiceZabbix,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitorService(input)
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

type ScopeType string

const (
	ScopeTypeResource      ScopeType = "Resource"
	ScopeTypeResourceGroup ScopeType = "ResourceGroup"
	ScopeTypeSubscription  ScopeType = "Subscription"
)

func PossibleValuesForScopeType() []string {
	return []string{
		string(ScopeTypeResource),
		string(ScopeTypeResourceGroup),
		string(ScopeTypeSubscription),
	}
}

func (s *ScopeType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseScopeType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseScopeType(input string) (*ScopeType, error) {
	vals := map[string]ScopeType{
		"resource":      ScopeTypeResource,
		"resourcegroup": ScopeTypeResourceGroup,
		"subscription":  ScopeTypeSubscription,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ScopeType(input)
	return &out, nil
}

type Severity string

const (
	SeveritySevFour  Severity = "Sev4"
	SeveritySevOne   Severity = "Sev1"
	SeveritySevThree Severity = "Sev3"
	SeveritySevTwo   Severity = "Sev2"
	SeveritySevZero  Severity = "Sev0"
)

func PossibleValuesForSeverity() []string {
	return []string{
		string(SeveritySevFour),
		string(SeveritySevOne),
		string(SeveritySevThree),
		string(SeveritySevTwo),
		string(SeveritySevZero),
	}
}

func (s *Severity) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSeverity(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSeverity(input string) (*Severity, error) {
	vals := map[string]Severity{
		"sev4": SeveritySevFour,
		"sev1": SeveritySevOne,
		"sev3": SeveritySevThree,
		"sev2": SeveritySevTwo,
		"sev0": SeveritySevZero,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Severity(input)
	return &out, nil
}

type SuppressionType string

const (
	SuppressionTypeAlways  SuppressionType = "Always"
	SuppressionTypeDaily   SuppressionType = "Daily"
	SuppressionTypeMonthly SuppressionType = "Monthly"
	SuppressionTypeOnce    SuppressionType = "Once"
	SuppressionTypeWeekly  SuppressionType = "Weekly"
)

func PossibleValuesForSuppressionType() []string {
	return []string{
		string(SuppressionTypeAlways),
		string(SuppressionTypeDaily),
		string(SuppressionTypeMonthly),
		string(SuppressionTypeOnce),
		string(SuppressionTypeWeekly),
	}
}

func (s *SuppressionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSuppressionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSuppressionType(input string) (*SuppressionType, error) {
	vals := map[string]SuppressionType{
		"always":  SuppressionTypeAlways,
		"daily":   SuppressionTypeDaily,
		"monthly": SuppressionTypeMonthly,
		"once":    SuppressionTypeOnce,
		"weekly":  SuppressionTypeWeekly,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SuppressionType(input)
	return &out, nil
}

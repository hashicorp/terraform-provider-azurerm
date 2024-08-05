package alertsmanagements

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertModificationEvent string

const (
	AlertModificationEventActionRuleSuppressed   AlertModificationEvent = "ActionRuleSuppressed"
	AlertModificationEventActionRuleTriggered    AlertModificationEvent = "ActionRuleTriggered"
	AlertModificationEventActionsFailed          AlertModificationEvent = "ActionsFailed"
	AlertModificationEventActionsSuppressed      AlertModificationEvent = "ActionsSuppressed"
	AlertModificationEventActionsTriggered       AlertModificationEvent = "ActionsTriggered"
	AlertModificationEventAlertCreated           AlertModificationEvent = "AlertCreated"
	AlertModificationEventMonitorConditionChange AlertModificationEvent = "MonitorConditionChange"
	AlertModificationEventSeverityChange         AlertModificationEvent = "SeverityChange"
	AlertModificationEventStateChange            AlertModificationEvent = "StateChange"
)

func PossibleValuesForAlertModificationEvent() []string {
	return []string{
		string(AlertModificationEventActionRuleSuppressed),
		string(AlertModificationEventActionRuleTriggered),
		string(AlertModificationEventActionsFailed),
		string(AlertModificationEventActionsSuppressed),
		string(AlertModificationEventActionsTriggered),
		string(AlertModificationEventAlertCreated),
		string(AlertModificationEventMonitorConditionChange),
		string(AlertModificationEventSeverityChange),
		string(AlertModificationEventStateChange),
	}
}

func (s *AlertModificationEvent) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlertModificationEvent(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlertModificationEvent(input string) (*AlertModificationEvent, error) {
	vals := map[string]AlertModificationEvent{
		"actionrulesuppressed":   AlertModificationEventActionRuleSuppressed,
		"actionruletriggered":    AlertModificationEventActionRuleTriggered,
		"actionsfailed":          AlertModificationEventActionsFailed,
		"actionssuppressed":      AlertModificationEventActionsSuppressed,
		"actionstriggered":       AlertModificationEventActionsTriggered,
		"alertcreated":           AlertModificationEventAlertCreated,
		"monitorconditionchange": AlertModificationEventMonitorConditionChange,
		"severitychange":         AlertModificationEventSeverityChange,
		"statechange":            AlertModificationEventStateChange,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertModificationEvent(input)
	return &out, nil
}

type AlertState string

const (
	AlertStateAcknowledged AlertState = "Acknowledged"
	AlertStateClosed       AlertState = "Closed"
	AlertStateNew          AlertState = "New"
)

func PossibleValuesForAlertState() []string {
	return []string{
		string(AlertStateAcknowledged),
		string(AlertStateClosed),
		string(AlertStateNew),
	}
}

func (s *AlertState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlertState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlertState(input string) (*AlertState, error) {
	vals := map[string]AlertState{
		"acknowledged": AlertStateAcknowledged,
		"closed":       AlertStateClosed,
		"new":          AlertStateNew,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertState(input)
	return &out, nil
}

type AlertsSortByFields string

const (
	AlertsSortByFieldsAlertState           AlertsSortByFields = "alertState"
	AlertsSortByFieldsLastModifiedDateTime AlertsSortByFields = "lastModifiedDateTime"
	AlertsSortByFieldsMonitorCondition     AlertsSortByFields = "monitorCondition"
	AlertsSortByFieldsName                 AlertsSortByFields = "name"
	AlertsSortByFieldsSeverity             AlertsSortByFields = "severity"
	AlertsSortByFieldsStartDateTime        AlertsSortByFields = "startDateTime"
	AlertsSortByFieldsTargetResource       AlertsSortByFields = "targetResource"
	AlertsSortByFieldsTargetResourceGroup  AlertsSortByFields = "targetResourceGroup"
	AlertsSortByFieldsTargetResourceName   AlertsSortByFields = "targetResourceName"
	AlertsSortByFieldsTargetResourceType   AlertsSortByFields = "targetResourceType"
)

func PossibleValuesForAlertsSortByFields() []string {
	return []string{
		string(AlertsSortByFieldsAlertState),
		string(AlertsSortByFieldsLastModifiedDateTime),
		string(AlertsSortByFieldsMonitorCondition),
		string(AlertsSortByFieldsName),
		string(AlertsSortByFieldsSeverity),
		string(AlertsSortByFieldsStartDateTime),
		string(AlertsSortByFieldsTargetResource),
		string(AlertsSortByFieldsTargetResourceGroup),
		string(AlertsSortByFieldsTargetResourceName),
		string(AlertsSortByFieldsTargetResourceType),
	}
}

func (s *AlertsSortByFields) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlertsSortByFields(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlertsSortByFields(input string) (*AlertsSortByFields, error) {
	vals := map[string]AlertsSortByFields{
		"alertstate":           AlertsSortByFieldsAlertState,
		"lastmodifieddatetime": AlertsSortByFieldsLastModifiedDateTime,
		"monitorcondition":     AlertsSortByFieldsMonitorCondition,
		"name":                 AlertsSortByFieldsName,
		"severity":             AlertsSortByFieldsSeverity,
		"startdatetime":        AlertsSortByFieldsStartDateTime,
		"targetresource":       AlertsSortByFieldsTargetResource,
		"targetresourcegroup":  AlertsSortByFieldsTargetResourceGroup,
		"targetresourcename":   AlertsSortByFieldsTargetResourceName,
		"targetresourcetype":   AlertsSortByFieldsTargetResourceType,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertsSortByFields(input)
	return &out, nil
}

type AlertsSummaryGroupByFields string

const (
	AlertsSummaryGroupByFieldsAlertRule        AlertsSummaryGroupByFields = "alertRule"
	AlertsSummaryGroupByFieldsAlertState       AlertsSummaryGroupByFields = "alertState"
	AlertsSummaryGroupByFieldsMonitorCondition AlertsSummaryGroupByFields = "monitorCondition"
	AlertsSummaryGroupByFieldsMonitorService   AlertsSummaryGroupByFields = "monitorService"
	AlertsSummaryGroupByFieldsSeverity         AlertsSummaryGroupByFields = "severity"
	AlertsSummaryGroupByFieldsSignalType       AlertsSummaryGroupByFields = "signalType"
)

func PossibleValuesForAlertsSummaryGroupByFields() []string {
	return []string{
		string(AlertsSummaryGroupByFieldsAlertRule),
		string(AlertsSummaryGroupByFieldsAlertState),
		string(AlertsSummaryGroupByFieldsMonitorCondition),
		string(AlertsSummaryGroupByFieldsMonitorService),
		string(AlertsSummaryGroupByFieldsSeverity),
		string(AlertsSummaryGroupByFieldsSignalType),
	}
}

func (s *AlertsSummaryGroupByFields) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseAlertsSummaryGroupByFields(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseAlertsSummaryGroupByFields(input string) (*AlertsSummaryGroupByFields, error) {
	vals := map[string]AlertsSummaryGroupByFields{
		"alertrule":        AlertsSummaryGroupByFieldsAlertRule,
		"alertstate":       AlertsSummaryGroupByFieldsAlertState,
		"monitorcondition": AlertsSummaryGroupByFieldsMonitorCondition,
		"monitorservice":   AlertsSummaryGroupByFieldsMonitorService,
		"severity":         AlertsSummaryGroupByFieldsSeverity,
		"signaltype":       AlertsSummaryGroupByFieldsSignalType,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := AlertsSummaryGroupByFields(input)
	return &out, nil
}

type Identifier string

const (
	IdentifierMonitorServiceList Identifier = "MonitorServiceList"
)

func PossibleValuesForIdentifier() []string {
	return []string{
		string(IdentifierMonitorServiceList),
	}
}

func (s *Identifier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseIdentifier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseIdentifier(input string) (*Identifier, error) {
	vals := map[string]Identifier{
		"monitorservicelist": IdentifierMonitorServiceList,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := Identifier(input)
	return &out, nil
}

type MetadataIdentifier string

const (
	MetadataIdentifierMonitorServiceList MetadataIdentifier = "MonitorServiceList"
)

func PossibleValuesForMetadataIdentifier() []string {
	return []string{
		string(MetadataIdentifierMonitorServiceList),
	}
}

func (s *MetadataIdentifier) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMetadataIdentifier(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMetadataIdentifier(input string) (*MetadataIdentifier, error) {
	vals := map[string]MetadataIdentifier{
		"monitorservicelist": MetadataIdentifierMonitorServiceList,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MetadataIdentifier(input)
	return &out, nil
}

type MonitorCondition string

const (
	MonitorConditionFired    MonitorCondition = "Fired"
	MonitorConditionResolved MonitorCondition = "Resolved"
)

func PossibleValuesForMonitorCondition() []string {
	return []string{
		string(MonitorConditionFired),
		string(MonitorConditionResolved),
	}
}

func (s *MonitorCondition) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseMonitorCondition(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseMonitorCondition(input string) (*MonitorCondition, error) {
	vals := map[string]MonitorCondition{
		"fired":    MonitorConditionFired,
		"resolved": MonitorConditionResolved,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := MonitorCondition(input)
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

type SignalType string

const (
	SignalTypeLog     SignalType = "Log"
	SignalTypeMetric  SignalType = "Metric"
	SignalTypeUnknown SignalType = "Unknown"
)

func PossibleValuesForSignalType() []string {
	return []string{
		string(SignalTypeLog),
		string(SignalTypeMetric),
		string(SignalTypeUnknown),
	}
}

func (s *SignalType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSignalType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSignalType(input string) (*SignalType, error) {
	vals := map[string]SignalType{
		"log":     SignalTypeLog,
		"metric":  SignalTypeMetric,
		"unknown": SignalTypeUnknown,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SignalType(input)
	return &out, nil
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func PossibleValuesForSortOrder() []string {
	return []string{
		string(SortOrderAsc),
		string(SortOrderDesc),
	}
}

func (s *SortOrder) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseSortOrder(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseSortOrder(input string) (*SortOrder, error) {
	vals := map[string]SortOrder{
		"asc":  SortOrderAsc,
		"desc": SortOrderDesc,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := SortOrder(input)
	return &out, nil
}

type TimeRange string

const (
	TimeRangeOned       TimeRange = "1d"
	TimeRangeOneh       TimeRange = "1h"
	TimeRangeSevend     TimeRange = "7d"
	TimeRangeThreeZerod TimeRange = "30d"
)

func PossibleValuesForTimeRange() []string {
	return []string{
		string(TimeRangeOned),
		string(TimeRangeOneh),
		string(TimeRangeSevend),
		string(TimeRangeThreeZerod),
	}
}

func (s *TimeRange) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseTimeRange(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseTimeRange(input string) (*TimeRange, error) {
	vals := map[string]TimeRange{
		"1d":  TimeRangeOned,
		"1h":  TimeRangeOneh,
		"7d":  TimeRangeSevend,
		"30d": TimeRangeThreeZerod,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := TimeRange(input)
	return &out, nil
}

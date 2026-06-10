package datacollectionrules

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KnownAgentSettingName string

const (
	KnownAgentSettingNameMaxDiskQuotaInMB                  KnownAgentSettingName = "MaxDiskQuotaInMB"
	KnownAgentSettingNameUseTimeReceivedForForwardedEvents KnownAgentSettingName = "UseTimeReceivedForForwardedEvents"
)

func PossibleValuesForKnownAgentSettingName() []string {
	return []string{
		string(KnownAgentSettingNameMaxDiskQuotaInMB),
		string(KnownAgentSettingNameUseTimeReceivedForForwardedEvents),
	}
}

func (s *KnownAgentSettingName) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownAgentSettingName(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownAgentSettingName(input string) (*KnownAgentSettingName, error) {
	vals := map[string]KnownAgentSettingName{
		"maxdiskquotainmb":                  KnownAgentSettingNameMaxDiskQuotaInMB,
		"usetimereceivedforforwardedevents": KnownAgentSettingNameUseTimeReceivedForForwardedEvents,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownAgentSettingName(input)
	return &out, nil
}

type KnownColumnDefinitionType string

const (
	KnownColumnDefinitionTypeBoolean  KnownColumnDefinitionType = "boolean"
	KnownColumnDefinitionTypeDatetime KnownColumnDefinitionType = "datetime"
	KnownColumnDefinitionTypeDynamic  KnownColumnDefinitionType = "dynamic"
	KnownColumnDefinitionTypeInt      KnownColumnDefinitionType = "int"
	KnownColumnDefinitionTypeLong     KnownColumnDefinitionType = "long"
	KnownColumnDefinitionTypeReal     KnownColumnDefinitionType = "real"
	KnownColumnDefinitionTypeString   KnownColumnDefinitionType = "string"
)

func PossibleValuesForKnownColumnDefinitionType() []string {
	return []string{
		string(KnownColumnDefinitionTypeBoolean),
		string(KnownColumnDefinitionTypeDatetime),
		string(KnownColumnDefinitionTypeDynamic),
		string(KnownColumnDefinitionTypeInt),
		string(KnownColumnDefinitionTypeLong),
		string(KnownColumnDefinitionTypeReal),
		string(KnownColumnDefinitionTypeString),
	}
}

func (s *KnownColumnDefinitionType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownColumnDefinitionType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownColumnDefinitionType(input string) (*KnownColumnDefinitionType, error) {
	vals := map[string]KnownColumnDefinitionType{
		"boolean":  KnownColumnDefinitionTypeBoolean,
		"datetime": KnownColumnDefinitionTypeDatetime,
		"dynamic":  KnownColumnDefinitionTypeDynamic,
		"int":      KnownColumnDefinitionTypeInt,
		"long":     KnownColumnDefinitionTypeLong,
		"real":     KnownColumnDefinitionTypeReal,
		"string":   KnownColumnDefinitionTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownColumnDefinitionType(input)
	return &out, nil
}

type KnownDataCollectionRuleProvisioningState string

const (
	KnownDataCollectionRuleProvisioningStateCanceled  KnownDataCollectionRuleProvisioningState = "Canceled"
	KnownDataCollectionRuleProvisioningStateCreating  KnownDataCollectionRuleProvisioningState = "Creating"
	KnownDataCollectionRuleProvisioningStateDeleting  KnownDataCollectionRuleProvisioningState = "Deleting"
	KnownDataCollectionRuleProvisioningStateFailed    KnownDataCollectionRuleProvisioningState = "Failed"
	KnownDataCollectionRuleProvisioningStateSucceeded KnownDataCollectionRuleProvisioningState = "Succeeded"
	KnownDataCollectionRuleProvisioningStateUpdating  KnownDataCollectionRuleProvisioningState = "Updating"
)

func PossibleValuesForKnownDataCollectionRuleProvisioningState() []string {
	return []string{
		string(KnownDataCollectionRuleProvisioningStateCanceled),
		string(KnownDataCollectionRuleProvisioningStateCreating),
		string(KnownDataCollectionRuleProvisioningStateDeleting),
		string(KnownDataCollectionRuleProvisioningStateFailed),
		string(KnownDataCollectionRuleProvisioningStateSucceeded),
		string(KnownDataCollectionRuleProvisioningStateUpdating),
	}
}

func (s *KnownDataCollectionRuleProvisioningState) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownDataCollectionRuleProvisioningState(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownDataCollectionRuleProvisioningState(input string) (*KnownDataCollectionRuleProvisioningState, error) {
	vals := map[string]KnownDataCollectionRuleProvisioningState{
		"canceled":  KnownDataCollectionRuleProvisioningStateCanceled,
		"creating":  KnownDataCollectionRuleProvisioningStateCreating,
		"deleting":  KnownDataCollectionRuleProvisioningStateDeleting,
		"failed":    KnownDataCollectionRuleProvisioningStateFailed,
		"succeeded": KnownDataCollectionRuleProvisioningStateSucceeded,
		"updating":  KnownDataCollectionRuleProvisioningStateUpdating,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownDataCollectionRuleProvisioningState(input)
	return &out, nil
}

type KnownDataCollectionRuleResourceKind string

const (
	KnownDataCollectionRuleResourceKindLinux   KnownDataCollectionRuleResourceKind = "Linux"
	KnownDataCollectionRuleResourceKindWindows KnownDataCollectionRuleResourceKind = "Windows"
)

func PossibleValuesForKnownDataCollectionRuleResourceKind() []string {
	return []string{
		string(KnownDataCollectionRuleResourceKindLinux),
		string(KnownDataCollectionRuleResourceKindWindows),
	}
}

func (s *KnownDataCollectionRuleResourceKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownDataCollectionRuleResourceKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownDataCollectionRuleResourceKind(input string) (*KnownDataCollectionRuleResourceKind, error) {
	vals := map[string]KnownDataCollectionRuleResourceKind{
		"linux":   KnownDataCollectionRuleResourceKindLinux,
		"windows": KnownDataCollectionRuleResourceKindWindows,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownDataCollectionRuleResourceKind(input)
	return &out, nil
}

type KnownDataFlowStreams string

const (
	KnownDataFlowStreamsMicrosoftNegativeEvent           KnownDataFlowStreams = "Microsoft-Event"
	KnownDataFlowStreamsMicrosoftNegativeInsightsMetrics KnownDataFlowStreams = "Microsoft-InsightsMetrics"
	KnownDataFlowStreamsMicrosoftNegativePerf            KnownDataFlowStreams = "Microsoft-Perf"
	KnownDataFlowStreamsMicrosoftNegativeSyslog          KnownDataFlowStreams = "Microsoft-Syslog"
	KnownDataFlowStreamsMicrosoftNegativeWindowsEvent    KnownDataFlowStreams = "Microsoft-WindowsEvent"
)

func PossibleValuesForKnownDataFlowStreams() []string {
	return []string{
		string(KnownDataFlowStreamsMicrosoftNegativeEvent),
		string(KnownDataFlowStreamsMicrosoftNegativeInsightsMetrics),
		string(KnownDataFlowStreamsMicrosoftNegativePerf),
		string(KnownDataFlowStreamsMicrosoftNegativeSyslog),
		string(KnownDataFlowStreamsMicrosoftNegativeWindowsEvent),
	}
}

func (s *KnownDataFlowStreams) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownDataFlowStreams(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownDataFlowStreams(input string) (*KnownDataFlowStreams, error) {
	vals := map[string]KnownDataFlowStreams{
		"microsoft-event":           KnownDataFlowStreamsMicrosoftNegativeEvent,
		"microsoft-insightsmetrics": KnownDataFlowStreamsMicrosoftNegativeInsightsMetrics,
		"microsoft-perf":            KnownDataFlowStreamsMicrosoftNegativePerf,
		"microsoft-syslog":          KnownDataFlowStreamsMicrosoftNegativeSyslog,
		"microsoft-windowsevent":    KnownDataFlowStreamsMicrosoftNegativeWindowsEvent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownDataFlowStreams(input)
	return &out, nil
}

type KnownExtensionDataSourceStreams string

const (
	KnownExtensionDataSourceStreamsMicrosoftNegativeEvent           KnownExtensionDataSourceStreams = "Microsoft-Event"
	KnownExtensionDataSourceStreamsMicrosoftNegativeInsightsMetrics KnownExtensionDataSourceStreams = "Microsoft-InsightsMetrics"
	KnownExtensionDataSourceStreamsMicrosoftNegativePerf            KnownExtensionDataSourceStreams = "Microsoft-Perf"
	KnownExtensionDataSourceStreamsMicrosoftNegativeSyslog          KnownExtensionDataSourceStreams = "Microsoft-Syslog"
	KnownExtensionDataSourceStreamsMicrosoftNegativeWindowsEvent    KnownExtensionDataSourceStreams = "Microsoft-WindowsEvent"
)

func PossibleValuesForKnownExtensionDataSourceStreams() []string {
	return []string{
		string(KnownExtensionDataSourceStreamsMicrosoftNegativeEvent),
		string(KnownExtensionDataSourceStreamsMicrosoftNegativeInsightsMetrics),
		string(KnownExtensionDataSourceStreamsMicrosoftNegativePerf),
		string(KnownExtensionDataSourceStreamsMicrosoftNegativeSyslog),
		string(KnownExtensionDataSourceStreamsMicrosoftNegativeWindowsEvent),
	}
}

func (s *KnownExtensionDataSourceStreams) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownExtensionDataSourceStreams(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownExtensionDataSourceStreams(input string) (*KnownExtensionDataSourceStreams, error) {
	vals := map[string]KnownExtensionDataSourceStreams{
		"microsoft-event":           KnownExtensionDataSourceStreamsMicrosoftNegativeEvent,
		"microsoft-insightsmetrics": KnownExtensionDataSourceStreamsMicrosoftNegativeInsightsMetrics,
		"microsoft-perf":            KnownExtensionDataSourceStreamsMicrosoftNegativePerf,
		"microsoft-syslog":          KnownExtensionDataSourceStreamsMicrosoftNegativeSyslog,
		"microsoft-windowsevent":    KnownExtensionDataSourceStreamsMicrosoftNegativeWindowsEvent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownExtensionDataSourceStreams(input)
	return &out, nil
}

type KnownLogFileTextSettingsRecordStartTimestampFormat string

const (
	KnownLogFileTextSettingsRecordStartTimestampFormatDdMMMYyyyHHMmSsZzz               KnownLogFileTextSettingsRecordStartTimestampFormat = "dd/MMM/yyyy:HH:mm:ss zzz"
	KnownLogFileTextSettingsRecordStartTimestampFormatDdMMyyHHMmSs                     KnownLogFileTextSettingsRecordStartTimestampFormat = "ddMMyy HH:mm:ss"
	KnownLogFileTextSettingsRecordStartTimestampFormatISOEightSixZeroOne               KnownLogFileTextSettingsRecordStartTimestampFormat = "ISO 8601"
	KnownLogFileTextSettingsRecordStartTimestampFormatMDYYYYHHMMSSAMPM                 KnownLogFileTextSettingsRecordStartTimestampFormat = "M/D/YYYY HH:MM:SS AM/PM"
	KnownLogFileTextSettingsRecordStartTimestampFormatMMMDHhMmSs                       KnownLogFileTextSettingsRecordStartTimestampFormat = "MMM d hh:mm:ss"
	KnownLogFileTextSettingsRecordStartTimestampFormatMonDDYYYYHHMMSS                  KnownLogFileTextSettingsRecordStartTimestampFormat = "Mon DD, YYYY HH:MM:SS"
	KnownLogFileTextSettingsRecordStartTimestampFormatYYYYNegativeMMNegativeDDHHMMSS   KnownLogFileTextSettingsRecordStartTimestampFormat = "YYYY-MM-DD HH:MM:SS"
	KnownLogFileTextSettingsRecordStartTimestampFormatYyMMddHHMmSs                     KnownLogFileTextSettingsRecordStartTimestampFormat = "yyMMdd HH:mm:ss"
	KnownLogFileTextSettingsRecordStartTimestampFormatYyyyNegativeMMNegativeddTHHMmSsK KnownLogFileTextSettingsRecordStartTimestampFormat = "yyyy-MM-ddTHH:mm:ssK"
)

func PossibleValuesForKnownLogFileTextSettingsRecordStartTimestampFormat() []string {
	return []string{
		string(KnownLogFileTextSettingsRecordStartTimestampFormatDdMMMYyyyHHMmSsZzz),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatDdMMyyHHMmSs),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatISOEightSixZeroOne),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatMDYYYYHHMMSSAMPM),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatMMMDHhMmSs),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatMonDDYYYYHHMMSS),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatYYYYNegativeMMNegativeDDHHMMSS),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatYyMMddHHMmSs),
		string(KnownLogFileTextSettingsRecordStartTimestampFormatYyyyNegativeMMNegativeddTHHMmSsK),
	}
}

func (s *KnownLogFileTextSettingsRecordStartTimestampFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownLogFileTextSettingsRecordStartTimestampFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownLogFileTextSettingsRecordStartTimestampFormat(input string) (*KnownLogFileTextSettingsRecordStartTimestampFormat, error) {
	vals := map[string]KnownLogFileTextSettingsRecordStartTimestampFormat{
		"dd/mmm/yyyy:hh:mm:ss zzz": KnownLogFileTextSettingsRecordStartTimestampFormatDdMMMYyyyHHMmSsZzz,
		"ddmmyy hh:mm:ss":          KnownLogFileTextSettingsRecordStartTimestampFormatDdMMyyHHMmSs,
		"iso 8601":                 KnownLogFileTextSettingsRecordStartTimestampFormatISOEightSixZeroOne,
		"m/d/yyyy hh:mm:ss am/pm":  KnownLogFileTextSettingsRecordStartTimestampFormatMDYYYYHHMMSSAMPM,
		"mmm d hh:mm:ss":           KnownLogFileTextSettingsRecordStartTimestampFormatMMMDHhMmSs,
		"mon dd, yyyy hh:mm:ss":    KnownLogFileTextSettingsRecordStartTimestampFormatMonDDYYYYHHMMSS,
		"yyyy-mm-dd hh:mm:ss":      KnownLogFileTextSettingsRecordStartTimestampFormatYYYYNegativeMMNegativeDDHHMMSS,
		"yymmdd hh:mm:ss":          KnownLogFileTextSettingsRecordStartTimestampFormatYyMMddHHMmSs,
		"yyyy-mm-ddthh:mm:ssk":     KnownLogFileTextSettingsRecordStartTimestampFormatYyyyNegativeMMNegativeddTHHMmSsK,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownLogFileTextSettingsRecordStartTimestampFormat(input)
	return &out, nil
}

type KnownLogFilesDataSourceFormat string

const (
	KnownLogFilesDataSourceFormatJson KnownLogFilesDataSourceFormat = "json"
	KnownLogFilesDataSourceFormatText KnownLogFilesDataSourceFormat = "text"
)

func PossibleValuesForKnownLogFilesDataSourceFormat() []string {
	return []string{
		string(KnownLogFilesDataSourceFormatJson),
		string(KnownLogFilesDataSourceFormatText),
	}
}

func (s *KnownLogFilesDataSourceFormat) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownLogFilesDataSourceFormat(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownLogFilesDataSourceFormat(input string) (*KnownLogFilesDataSourceFormat, error) {
	vals := map[string]KnownLogFilesDataSourceFormat{
		"json": KnownLogFilesDataSourceFormatJson,
		"text": KnownLogFilesDataSourceFormatText,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownLogFilesDataSourceFormat(input)
	return &out, nil
}

type KnownPerfCounterDataSourceStreams string

const (
	KnownPerfCounterDataSourceStreamsMicrosoftNegativeInsightsMetrics KnownPerfCounterDataSourceStreams = "Microsoft-InsightsMetrics"
	KnownPerfCounterDataSourceStreamsMicrosoftNegativePerf            KnownPerfCounterDataSourceStreams = "Microsoft-Perf"
)

func PossibleValuesForKnownPerfCounterDataSourceStreams() []string {
	return []string{
		string(KnownPerfCounterDataSourceStreamsMicrosoftNegativeInsightsMetrics),
		string(KnownPerfCounterDataSourceStreamsMicrosoftNegativePerf),
	}
}

func (s *KnownPerfCounterDataSourceStreams) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownPerfCounterDataSourceStreams(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownPerfCounterDataSourceStreams(input string) (*KnownPerfCounterDataSourceStreams, error) {
	vals := map[string]KnownPerfCounterDataSourceStreams{
		"microsoft-insightsmetrics": KnownPerfCounterDataSourceStreamsMicrosoftNegativeInsightsMetrics,
		"microsoft-perf":            KnownPerfCounterDataSourceStreamsMicrosoftNegativePerf,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownPerfCounterDataSourceStreams(input)
	return &out, nil
}

type KnownPrometheusForwarderDataSourceStreams string

const (
	KnownPrometheusForwarderDataSourceStreamsMicrosoftNegativePrometheusMetrics KnownPrometheusForwarderDataSourceStreams = "Microsoft-PrometheusMetrics"
)

func PossibleValuesForKnownPrometheusForwarderDataSourceStreams() []string {
	return []string{
		string(KnownPrometheusForwarderDataSourceStreamsMicrosoftNegativePrometheusMetrics),
	}
}

func (s *KnownPrometheusForwarderDataSourceStreams) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownPrometheusForwarderDataSourceStreams(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownPrometheusForwarderDataSourceStreams(input string) (*KnownPrometheusForwarderDataSourceStreams, error) {
	vals := map[string]KnownPrometheusForwarderDataSourceStreams{
		"microsoft-prometheusmetrics": KnownPrometheusForwarderDataSourceStreamsMicrosoftNegativePrometheusMetrics,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownPrometheusForwarderDataSourceStreams(input)
	return &out, nil
}

type KnownStorageBlobLookupType string

const (
	KnownStorageBlobLookupTypeCidr   KnownStorageBlobLookupType = "Cidr"
	KnownStorageBlobLookupTypeString KnownStorageBlobLookupType = "String"
)

func PossibleValuesForKnownStorageBlobLookupType() []string {
	return []string{
		string(KnownStorageBlobLookupTypeCidr),
		string(KnownStorageBlobLookupTypeString),
	}
}

func (s *KnownStorageBlobLookupType) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownStorageBlobLookupType(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownStorageBlobLookupType(input string) (*KnownStorageBlobLookupType, error) {
	vals := map[string]KnownStorageBlobLookupType{
		"cidr":   KnownStorageBlobLookupTypeCidr,
		"string": KnownStorageBlobLookupTypeString,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownStorageBlobLookupType(input)
	return &out, nil
}

type KnownSyslogDataSourceFacilityNames string

const (
	KnownSyslogDataSourceFacilityNamesAlert      KnownSyslogDataSourceFacilityNames = "alert"
	KnownSyslogDataSourceFacilityNamesAny        KnownSyslogDataSourceFacilityNames = "*"
	KnownSyslogDataSourceFacilityNamesAudit      KnownSyslogDataSourceFacilityNames = "audit"
	KnownSyslogDataSourceFacilityNamesAuth       KnownSyslogDataSourceFacilityNames = "auth"
	KnownSyslogDataSourceFacilityNamesAuthpriv   KnownSyslogDataSourceFacilityNames = "authpriv"
	KnownSyslogDataSourceFacilityNamesClock      KnownSyslogDataSourceFacilityNames = "clock"
	KnownSyslogDataSourceFacilityNamesCron       KnownSyslogDataSourceFacilityNames = "cron"
	KnownSyslogDataSourceFacilityNamesDaemon     KnownSyslogDataSourceFacilityNames = "daemon"
	KnownSyslogDataSourceFacilityNamesFtp        KnownSyslogDataSourceFacilityNames = "ftp"
	KnownSyslogDataSourceFacilityNamesKern       KnownSyslogDataSourceFacilityNames = "kern"
	KnownSyslogDataSourceFacilityNamesLocalFive  KnownSyslogDataSourceFacilityNames = "local5"
	KnownSyslogDataSourceFacilityNamesLocalFour  KnownSyslogDataSourceFacilityNames = "local4"
	KnownSyslogDataSourceFacilityNamesLocalOne   KnownSyslogDataSourceFacilityNames = "local1"
	KnownSyslogDataSourceFacilityNamesLocalSeven KnownSyslogDataSourceFacilityNames = "local7"
	KnownSyslogDataSourceFacilityNamesLocalSix   KnownSyslogDataSourceFacilityNames = "local6"
	KnownSyslogDataSourceFacilityNamesLocalThree KnownSyslogDataSourceFacilityNames = "local3"
	KnownSyslogDataSourceFacilityNamesLocalTwo   KnownSyslogDataSourceFacilityNames = "local2"
	KnownSyslogDataSourceFacilityNamesLocalZero  KnownSyslogDataSourceFacilityNames = "local0"
	KnownSyslogDataSourceFacilityNamesLpr        KnownSyslogDataSourceFacilityNames = "lpr"
	KnownSyslogDataSourceFacilityNamesMail       KnownSyslogDataSourceFacilityNames = "mail"
	KnownSyslogDataSourceFacilityNamesMark       KnownSyslogDataSourceFacilityNames = "mark"
	KnownSyslogDataSourceFacilityNamesNews       KnownSyslogDataSourceFacilityNames = "news"
	KnownSyslogDataSourceFacilityNamesNopri      KnownSyslogDataSourceFacilityNames = "nopri"
	KnownSyslogDataSourceFacilityNamesNtp        KnownSyslogDataSourceFacilityNames = "ntp"
	KnownSyslogDataSourceFacilityNamesSyslog     KnownSyslogDataSourceFacilityNames = "syslog"
	KnownSyslogDataSourceFacilityNamesUser       KnownSyslogDataSourceFacilityNames = "user"
	KnownSyslogDataSourceFacilityNamesUucp       KnownSyslogDataSourceFacilityNames = "uucp"
)

func PossibleValuesForKnownSyslogDataSourceFacilityNames() []string {
	return []string{
		string(KnownSyslogDataSourceFacilityNamesAlert),
		string(KnownSyslogDataSourceFacilityNamesAny),
		string(KnownSyslogDataSourceFacilityNamesAudit),
		string(KnownSyslogDataSourceFacilityNamesAuth),
		string(KnownSyslogDataSourceFacilityNamesAuthpriv),
		string(KnownSyslogDataSourceFacilityNamesClock),
		string(KnownSyslogDataSourceFacilityNamesCron),
		string(KnownSyslogDataSourceFacilityNamesDaemon),
		string(KnownSyslogDataSourceFacilityNamesFtp),
		string(KnownSyslogDataSourceFacilityNamesKern),
		string(KnownSyslogDataSourceFacilityNamesLocalFive),
		string(KnownSyslogDataSourceFacilityNamesLocalFour),
		string(KnownSyslogDataSourceFacilityNamesLocalOne),
		string(KnownSyslogDataSourceFacilityNamesLocalSeven),
		string(KnownSyslogDataSourceFacilityNamesLocalSix),
		string(KnownSyslogDataSourceFacilityNamesLocalThree),
		string(KnownSyslogDataSourceFacilityNamesLocalTwo),
		string(KnownSyslogDataSourceFacilityNamesLocalZero),
		string(KnownSyslogDataSourceFacilityNamesLpr),
		string(KnownSyslogDataSourceFacilityNamesMail),
		string(KnownSyslogDataSourceFacilityNamesMark),
		string(KnownSyslogDataSourceFacilityNamesNews),
		string(KnownSyslogDataSourceFacilityNamesNopri),
		string(KnownSyslogDataSourceFacilityNamesNtp),
		string(KnownSyslogDataSourceFacilityNamesSyslog),
		string(KnownSyslogDataSourceFacilityNamesUser),
		string(KnownSyslogDataSourceFacilityNamesUucp),
	}
}

func (s *KnownSyslogDataSourceFacilityNames) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownSyslogDataSourceFacilityNames(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownSyslogDataSourceFacilityNames(input string) (*KnownSyslogDataSourceFacilityNames, error) {
	vals := map[string]KnownSyslogDataSourceFacilityNames{
		"alert":    KnownSyslogDataSourceFacilityNamesAlert,
		"*":        KnownSyslogDataSourceFacilityNamesAny,
		"audit":    KnownSyslogDataSourceFacilityNamesAudit,
		"auth":     KnownSyslogDataSourceFacilityNamesAuth,
		"authpriv": KnownSyslogDataSourceFacilityNamesAuthpriv,
		"clock":    KnownSyslogDataSourceFacilityNamesClock,
		"cron":     KnownSyslogDataSourceFacilityNamesCron,
		"daemon":   KnownSyslogDataSourceFacilityNamesDaemon,
		"ftp":      KnownSyslogDataSourceFacilityNamesFtp,
		"kern":     KnownSyslogDataSourceFacilityNamesKern,
		"local5":   KnownSyslogDataSourceFacilityNamesLocalFive,
		"local4":   KnownSyslogDataSourceFacilityNamesLocalFour,
		"local1":   KnownSyslogDataSourceFacilityNamesLocalOne,
		"local7":   KnownSyslogDataSourceFacilityNamesLocalSeven,
		"local6":   KnownSyslogDataSourceFacilityNamesLocalSix,
		"local3":   KnownSyslogDataSourceFacilityNamesLocalThree,
		"local2":   KnownSyslogDataSourceFacilityNamesLocalTwo,
		"local0":   KnownSyslogDataSourceFacilityNamesLocalZero,
		"lpr":      KnownSyslogDataSourceFacilityNamesLpr,
		"mail":     KnownSyslogDataSourceFacilityNamesMail,
		"mark":     KnownSyslogDataSourceFacilityNamesMark,
		"news":     KnownSyslogDataSourceFacilityNamesNews,
		"nopri":    KnownSyslogDataSourceFacilityNamesNopri,
		"ntp":      KnownSyslogDataSourceFacilityNamesNtp,
		"syslog":   KnownSyslogDataSourceFacilityNamesSyslog,
		"user":     KnownSyslogDataSourceFacilityNamesUser,
		"uucp":     KnownSyslogDataSourceFacilityNamesUucp,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownSyslogDataSourceFacilityNames(input)
	return &out, nil
}

type KnownSyslogDataSourceLogLevels string

const (
	KnownSyslogDataSourceLogLevelsAlert     KnownSyslogDataSourceLogLevels = "Alert"
	KnownSyslogDataSourceLogLevelsAny       KnownSyslogDataSourceLogLevels = "*"
	KnownSyslogDataSourceLogLevelsCritical  KnownSyslogDataSourceLogLevels = "Critical"
	KnownSyslogDataSourceLogLevelsDebug     KnownSyslogDataSourceLogLevels = "Debug"
	KnownSyslogDataSourceLogLevelsEmergency KnownSyslogDataSourceLogLevels = "Emergency"
	KnownSyslogDataSourceLogLevelsError     KnownSyslogDataSourceLogLevels = "Error"
	KnownSyslogDataSourceLogLevelsInfo      KnownSyslogDataSourceLogLevels = "Info"
	KnownSyslogDataSourceLogLevelsNotice    KnownSyslogDataSourceLogLevels = "Notice"
	KnownSyslogDataSourceLogLevelsWarning   KnownSyslogDataSourceLogLevels = "Warning"
)

func PossibleValuesForKnownSyslogDataSourceLogLevels() []string {
	return []string{
		string(KnownSyslogDataSourceLogLevelsAlert),
		string(KnownSyslogDataSourceLogLevelsAny),
		string(KnownSyslogDataSourceLogLevelsCritical),
		string(KnownSyslogDataSourceLogLevelsDebug),
		string(KnownSyslogDataSourceLogLevelsEmergency),
		string(KnownSyslogDataSourceLogLevelsError),
		string(KnownSyslogDataSourceLogLevelsInfo),
		string(KnownSyslogDataSourceLogLevelsNotice),
		string(KnownSyslogDataSourceLogLevelsWarning),
	}
}

func (s *KnownSyslogDataSourceLogLevels) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownSyslogDataSourceLogLevels(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownSyslogDataSourceLogLevels(input string) (*KnownSyslogDataSourceLogLevels, error) {
	vals := map[string]KnownSyslogDataSourceLogLevels{
		"alert":     KnownSyslogDataSourceLogLevelsAlert,
		"*":         KnownSyslogDataSourceLogLevelsAny,
		"critical":  KnownSyslogDataSourceLogLevelsCritical,
		"debug":     KnownSyslogDataSourceLogLevelsDebug,
		"emergency": KnownSyslogDataSourceLogLevelsEmergency,
		"error":     KnownSyslogDataSourceLogLevelsError,
		"info":      KnownSyslogDataSourceLogLevelsInfo,
		"notice":    KnownSyslogDataSourceLogLevelsNotice,
		"warning":   KnownSyslogDataSourceLogLevelsWarning,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownSyslogDataSourceLogLevels(input)
	return &out, nil
}

type KnownSyslogDataSourceStreams string

const (
	KnownSyslogDataSourceStreamsMicrosoftNegativeSyslog KnownSyslogDataSourceStreams = "Microsoft-Syslog"
)

func PossibleValuesForKnownSyslogDataSourceStreams() []string {
	return []string{
		string(KnownSyslogDataSourceStreamsMicrosoftNegativeSyslog),
	}
}

func (s *KnownSyslogDataSourceStreams) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownSyslogDataSourceStreams(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownSyslogDataSourceStreams(input string) (*KnownSyslogDataSourceStreams, error) {
	vals := map[string]KnownSyslogDataSourceStreams{
		"microsoft-syslog": KnownSyslogDataSourceStreamsMicrosoftNegativeSyslog,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownSyslogDataSourceStreams(input)
	return &out, nil
}

type KnownWindowsEventLogDataSourceStreams string

const (
	KnownWindowsEventLogDataSourceStreamsMicrosoftNegativeEvent        KnownWindowsEventLogDataSourceStreams = "Microsoft-Event"
	KnownWindowsEventLogDataSourceStreamsMicrosoftNegativeWindowsEvent KnownWindowsEventLogDataSourceStreams = "Microsoft-WindowsEvent"
)

func PossibleValuesForKnownWindowsEventLogDataSourceStreams() []string {
	return []string{
		string(KnownWindowsEventLogDataSourceStreamsMicrosoftNegativeEvent),
		string(KnownWindowsEventLogDataSourceStreamsMicrosoftNegativeWindowsEvent),
	}
}

func (s *KnownWindowsEventLogDataSourceStreams) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownWindowsEventLogDataSourceStreams(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownWindowsEventLogDataSourceStreams(input string) (*KnownWindowsEventLogDataSourceStreams, error) {
	vals := map[string]KnownWindowsEventLogDataSourceStreams{
		"microsoft-event":        KnownWindowsEventLogDataSourceStreamsMicrosoftNegativeEvent,
		"microsoft-windowsevent": KnownWindowsEventLogDataSourceStreamsMicrosoftNegativeWindowsEvent,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownWindowsEventLogDataSourceStreams(input)
	return &out, nil
}

type KnownWindowsFirewallLogsDataSourceProfileFilter string

const (
	KnownWindowsFirewallLogsDataSourceProfileFilterDomain  KnownWindowsFirewallLogsDataSourceProfileFilter = "Domain"
	KnownWindowsFirewallLogsDataSourceProfileFilterPrivate KnownWindowsFirewallLogsDataSourceProfileFilter = "Private"
	KnownWindowsFirewallLogsDataSourceProfileFilterPublic  KnownWindowsFirewallLogsDataSourceProfileFilter = "Public"
)

func PossibleValuesForKnownWindowsFirewallLogsDataSourceProfileFilter() []string {
	return []string{
		string(KnownWindowsFirewallLogsDataSourceProfileFilterDomain),
		string(KnownWindowsFirewallLogsDataSourceProfileFilterPrivate),
		string(KnownWindowsFirewallLogsDataSourceProfileFilterPublic),
	}
}

func (s *KnownWindowsFirewallLogsDataSourceProfileFilter) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseKnownWindowsFirewallLogsDataSourceProfileFilter(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseKnownWindowsFirewallLogsDataSourceProfileFilter(input string) (*KnownWindowsFirewallLogsDataSourceProfileFilter, error) {
	vals := map[string]KnownWindowsFirewallLogsDataSourceProfileFilter{
		"domain":  KnownWindowsFirewallLogsDataSourceProfileFilterDomain,
		"private": KnownWindowsFirewallLogsDataSourceProfileFilterPrivate,
		"public":  KnownWindowsFirewallLogsDataSourceProfileFilterPublic,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := KnownWindowsFirewallLogsDataSourceProfileFilter(input)
	return &out, nil
}

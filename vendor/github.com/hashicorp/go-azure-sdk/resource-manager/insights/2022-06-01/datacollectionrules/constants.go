package datacollectionrules

import "strings"

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

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
	KnownLogFilesDataSourceFormatText KnownLogFilesDataSourceFormat = "text"
)

func PossibleValuesForKnownLogFilesDataSourceFormat() []string {
	return []string{
		string(KnownLogFilesDataSourceFormatText),
	}
}

func parseKnownLogFilesDataSourceFormat(input string) (*KnownLogFilesDataSourceFormat, error) {
	vals := map[string]KnownLogFilesDataSourceFormat{
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

type KnownSyslogDataSourceFacilityNames string

const (
	KnownSyslogDataSourceFacilityNamesAny        KnownSyslogDataSourceFacilityNames = "*"
	KnownSyslogDataSourceFacilityNamesAuth       KnownSyslogDataSourceFacilityNames = "auth"
	KnownSyslogDataSourceFacilityNamesAuthpriv   KnownSyslogDataSourceFacilityNames = "authpriv"
	KnownSyslogDataSourceFacilityNamesCron       KnownSyslogDataSourceFacilityNames = "cron"
	KnownSyslogDataSourceFacilityNamesDaemon     KnownSyslogDataSourceFacilityNames = "daemon"
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
	KnownSyslogDataSourceFacilityNamesSyslog     KnownSyslogDataSourceFacilityNames = "syslog"
	KnownSyslogDataSourceFacilityNamesUser       KnownSyslogDataSourceFacilityNames = "user"
	KnownSyslogDataSourceFacilityNamesUucp       KnownSyslogDataSourceFacilityNames = "uucp"
)

func PossibleValuesForKnownSyslogDataSourceFacilityNames() []string {
	return []string{
		string(KnownSyslogDataSourceFacilityNamesAny),
		string(KnownSyslogDataSourceFacilityNamesAuth),
		string(KnownSyslogDataSourceFacilityNamesAuthpriv),
		string(KnownSyslogDataSourceFacilityNamesCron),
		string(KnownSyslogDataSourceFacilityNamesDaemon),
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
		string(KnownSyslogDataSourceFacilityNamesSyslog),
		string(KnownSyslogDataSourceFacilityNamesUser),
		string(KnownSyslogDataSourceFacilityNamesUucp),
	}
}

func parseKnownSyslogDataSourceFacilityNames(input string) (*KnownSyslogDataSourceFacilityNames, error) {
	vals := map[string]KnownSyslogDataSourceFacilityNames{
		"*":        KnownSyslogDataSourceFacilityNamesAny,
		"auth":     KnownSyslogDataSourceFacilityNamesAuth,
		"authpriv": KnownSyslogDataSourceFacilityNamesAuthpriv,
		"cron":     KnownSyslogDataSourceFacilityNamesCron,
		"daemon":   KnownSyslogDataSourceFacilityNamesDaemon,
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

package datasources

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataSourceKind string

const (
	DataSourceKindApplicationInsights                                  DataSourceKind = "ApplicationInsights"
	DataSourceKindAzureActivityLog                                     DataSourceKind = "AzureActivityLog"
	DataSourceKindAzureAuditLog                                        DataSourceKind = "AzureAuditLog"
	DataSourceKindChangeTrackingContentLocation                        DataSourceKind = "ChangeTrackingContentLocation"
	DataSourceKindChangeTrackingCustomPath                             DataSourceKind = "ChangeTrackingCustomPath"
	DataSourceKindChangeTrackingDataTypeConfiguration                  DataSourceKind = "ChangeTrackingDataTypeConfiguration"
	DataSourceKindChangeTrackingDefaultRegistry                        DataSourceKind = "ChangeTrackingDefaultRegistry"
	DataSourceKindChangeTrackingLinuxPath                              DataSourceKind = "ChangeTrackingLinuxPath"
	DataSourceKindChangeTrackingPath                                   DataSourceKind = "ChangeTrackingPath"
	DataSourceKindChangeTrackingRegistry                               DataSourceKind = "ChangeTrackingRegistry"
	DataSourceKindChangeTrackingServices                               DataSourceKind = "ChangeTrackingServices"
	DataSourceKindCustomLog                                            DataSourceKind = "CustomLog"
	DataSourceKindCustomLogCollection                                  DataSourceKind = "CustomLogCollection"
	DataSourceKindDnsAnalytics                                         DataSourceKind = "DnsAnalytics"
	DataSourceKindGenericDataSource                                    DataSourceKind = "GenericDataSource"
	DataSourceKindIISLogs                                              DataSourceKind = "IISLogs"
	DataSourceKindImportComputerGroup                                  DataSourceKind = "ImportComputerGroup"
	DataSourceKindItsm                                                 DataSourceKind = "Itsm"
	DataSourceKindLinuxChangeTrackingPath                              DataSourceKind = "LinuxChangeTrackingPath"
	DataSourceKindLinuxPerformanceCollection                           DataSourceKind = "LinuxPerformanceCollection"
	DataSourceKindLinuxPerformanceObject                               DataSourceKind = "LinuxPerformanceObject"
	DataSourceKindLinuxSyslog                                          DataSourceKind = "LinuxSyslog"
	DataSourceKindLinuxSyslogCollection                                DataSourceKind = "LinuxSyslogCollection"
	DataSourceKindNetworkMonitoring                                    DataSourceKind = "NetworkMonitoring"
	DataSourceKindOfficeThreeSixFive                                   DataSourceKind = "Office365"
	DataSourceKindSecurityCenterSecurityWindowsBaselineConfiguration   DataSourceKind = "SecurityCenterSecurityWindowsBaselineConfiguration"
	DataSourceKindSecurityEventCollectionConfiguration                 DataSourceKind = "SecurityEventCollectionConfiguration"
	DataSourceKindSecurityInsightsSecurityEventCollectionConfiguration DataSourceKind = "SecurityInsightsSecurityEventCollectionConfiguration"
	DataSourceKindSecurityWindowsBaselineConfiguration                 DataSourceKind = "SecurityWindowsBaselineConfiguration"
	DataSourceKindSqlDataClassification                                DataSourceKind = "SqlDataClassification"
	DataSourceKindWindowsEvent                                         DataSourceKind = "WindowsEvent"
	DataSourceKindWindowsPerformanceCounter                            DataSourceKind = "WindowsPerformanceCounter"
	DataSourceKindWindowsTelemetry                                     DataSourceKind = "WindowsTelemetry"
)

func PossibleValuesForDataSourceKind() []string {
	return []string{
		string(DataSourceKindApplicationInsights),
		string(DataSourceKindAzureActivityLog),
		string(DataSourceKindAzureAuditLog),
		string(DataSourceKindChangeTrackingContentLocation),
		string(DataSourceKindChangeTrackingCustomPath),
		string(DataSourceKindChangeTrackingDataTypeConfiguration),
		string(DataSourceKindChangeTrackingDefaultRegistry),
		string(DataSourceKindChangeTrackingLinuxPath),
		string(DataSourceKindChangeTrackingPath),
		string(DataSourceKindChangeTrackingRegistry),
		string(DataSourceKindChangeTrackingServices),
		string(DataSourceKindCustomLog),
		string(DataSourceKindCustomLogCollection),
		string(DataSourceKindDnsAnalytics),
		string(DataSourceKindGenericDataSource),
		string(DataSourceKindIISLogs),
		string(DataSourceKindImportComputerGroup),
		string(DataSourceKindItsm),
		string(DataSourceKindLinuxChangeTrackingPath),
		string(DataSourceKindLinuxPerformanceCollection),
		string(DataSourceKindLinuxPerformanceObject),
		string(DataSourceKindLinuxSyslog),
		string(DataSourceKindLinuxSyslogCollection),
		string(DataSourceKindNetworkMonitoring),
		string(DataSourceKindOfficeThreeSixFive),
		string(DataSourceKindSecurityCenterSecurityWindowsBaselineConfiguration),
		string(DataSourceKindSecurityEventCollectionConfiguration),
		string(DataSourceKindSecurityInsightsSecurityEventCollectionConfiguration),
		string(DataSourceKindSecurityWindowsBaselineConfiguration),
		string(DataSourceKindSqlDataClassification),
		string(DataSourceKindWindowsEvent),
		string(DataSourceKindWindowsPerformanceCounter),
		string(DataSourceKindWindowsTelemetry),
	}
}

func (s *DataSourceKind) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseDataSourceKind(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseDataSourceKind(input string) (*DataSourceKind, error) {
	vals := map[string]DataSourceKind{
		"applicationinsights":                 DataSourceKindApplicationInsights,
		"azureactivitylog":                    DataSourceKindAzureActivityLog,
		"azureauditlog":                       DataSourceKindAzureAuditLog,
		"changetrackingcontentlocation":       DataSourceKindChangeTrackingContentLocation,
		"changetrackingcustompath":            DataSourceKindChangeTrackingCustomPath,
		"changetrackingdatatypeconfiguration": DataSourceKindChangeTrackingDataTypeConfiguration,
		"changetrackingdefaultregistry":       DataSourceKindChangeTrackingDefaultRegistry,
		"changetrackinglinuxpath":             DataSourceKindChangeTrackingLinuxPath,
		"changetrackingpath":                  DataSourceKindChangeTrackingPath,
		"changetrackingregistry":              DataSourceKindChangeTrackingRegistry,
		"changetrackingservices":              DataSourceKindChangeTrackingServices,
		"customlog":                           DataSourceKindCustomLog,
		"customlogcollection":                 DataSourceKindCustomLogCollection,
		"dnsanalytics":                        DataSourceKindDnsAnalytics,
		"genericdatasource":                   DataSourceKindGenericDataSource,
		"iislogs":                             DataSourceKindIISLogs,
		"importcomputergroup":                 DataSourceKindImportComputerGroup,
		"itsm":                                DataSourceKindItsm,
		"linuxchangetrackingpath":             DataSourceKindLinuxChangeTrackingPath,
		"linuxperformancecollection":          DataSourceKindLinuxPerformanceCollection,
		"linuxperformanceobject":              DataSourceKindLinuxPerformanceObject,
		"linuxsyslog":                         DataSourceKindLinuxSyslog,
		"linuxsyslogcollection":               DataSourceKindLinuxSyslogCollection,
		"networkmonitoring":                   DataSourceKindNetworkMonitoring,
		"office365":                           DataSourceKindOfficeThreeSixFive,
		"securitycentersecuritywindowsbaselineconfiguration":   DataSourceKindSecurityCenterSecurityWindowsBaselineConfiguration,
		"securityeventcollectionconfiguration":                 DataSourceKindSecurityEventCollectionConfiguration,
		"securityinsightssecurityeventcollectionconfiguration": DataSourceKindSecurityInsightsSecurityEventCollectionConfiguration,
		"securitywindowsbaselineconfiguration":                 DataSourceKindSecurityWindowsBaselineConfiguration,
		"sqldataclassification":                                DataSourceKindSqlDataClassification,
		"windowsevent":                                         DataSourceKindWindowsEvent,
		"windowsperformancecounter":                            DataSourceKindWindowsPerformanceCounter,
		"windowstelemetry":                                     DataSourceKindWindowsTelemetry,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := DataSourceKind(input)
	return &out, nil
}

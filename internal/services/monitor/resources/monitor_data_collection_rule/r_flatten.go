// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func flattenDataCollectionRuleDataFlows(input *[]datacollectionrules.DataFlow) []DataFlow {
	if input == nil {
		return make([]DataFlow, 0)
	}

	result := make([]DataFlow, 0)
	for _, v := range *input {
		result = append(result, DataFlow{
			BuiltInTransform: pointer.From(v.BuiltInTransform),
			Destinations:     pointer.From(v.Destinations),
			OutputStream:     pointer.From(v.OutputStream),
			Streams:          pointer.FromEnumSlice(v.Streams),
			TransformKql:     pointer.From(v.TransformKql),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSources(input *datacollectionrules.DataSourcesSpec) []DataSourceModel {
	if input == nil {
		return make([]DataSourceModel, 0)
	}

	return []DataSourceModel{{
		DataImport:          flattenDataCollectionRuleDataSourceDataImports(input.DataImports),
		Extensions:          flattenDataCollectionRuleDataSourceExtensions(input.Extensions),
		IisLog:              flattenDataCollectionRuleDataSourceIisLog(input.IisLogs),
		LogFile:             flattenDataCollectionRuleDataSourceLogFiles(input.LogFiles),
		PerformanceCounters: flattenDataCollectionRuleDataSourcePerfCounters(input.PerformanceCounters),
		PlatformTelemetry:   flattenDataCollectionRuleDataSourcePlatformTelemetry(input.PlatformTelemetry),
		PrometheusForwarder: flattenDataCollectionRuleDataSourcePrometheusForwarder(input.PrometheusForwarder),
		Syslog:              flattenDataCollectionRuleDataSourceSyslog(input.Syslog),
		WindowsEventLogs:    flattenDataCollectionRuleWindowsEventLogs(input.WindowsEventLogs),
		WindowsFirewallLog:  flattenDataCollectionRuleWindowsFirewallLog(input.WindowsFirewallLogs),
	}}
}

func flattenDataCollectionRuleDataSourceDataImports(input *datacollectionrules.DataImportSources) []DataImport {
	if input == nil || input.EventHub == nil {
		return make([]DataImport, 0)
	}

	return []DataImport{
		{
			EventHubDataSource: []EventHubDataSource{
				{
					ConsumerGroup: pointer.From(input.EventHub.ConsumerGroup),
					Name:          pointer.From(input.EventHub.Name),
					Stream:        pointer.From(input.EventHub.Stream),
				},
			},
		},
	}
}

func flattenDataCollectionRuleDataSourceExtensions(input *[]datacollectionrules.ExtensionDataSource) []Extension {
	if input == nil {
		return make([]Extension, 0)
	}

	result := make([]Extension, 0)
	for _, v := range *input {
		extensionSettings := ""
		if v.ExtensionSettings != nil {
			settingString, _ := pluginsdk.FlattenJsonToString((*v.ExtensionSettings).(map[string]any))
			extensionSettings = settingString
		}
		result = append(result, Extension{
			ExtensionName:     v.ExtensionName,
			Name:              pointer.From(v.Name),
			ExtensionSettings: extensionSettings,
			InputDataSources:  pointer.From(v.InputDataSources),
			Streams:           pointer.FromEnumSlice(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceIisLog(input *[]datacollectionrules.IisLogsDataSource) []IisLog {
	if input == nil {
		return make([]IisLog, 0)
	}

	result := make([]IisLog, 0)
	for _, v := range *input {
		result = append(result, IisLog{
			Name:           pointer.From(v.Name),
			LogDirectories: pointer.From(v.LogDirectories),
			Streams:        v.Streams,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceLogFiles(input *[]datacollectionrules.LogFilesDataSource) []LogFile {
	if input == nil {
		return make([]LogFile, 0)
	}

	result := make([]LogFile, 0)
	for _, v := range *input {
		setting := make([]LogFileSetting, 0)
		if v.Settings != nil && v.Settings.Text != nil {
			setting = append(setting, LogFileSetting{
				Text: []TextSetting{
					{
						RecordStartTimestampFormat: string(v.Settings.Text.RecordStartTimestampFormat),
					},
				},
			})
		}

		result = append(result, LogFile{
			Name:         pointer.From(v.Name),
			Format:       string(v.Format),
			FilePatterns: v.FilePatterns,
			Streams:      v.Streams,
			Settings:     setting,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePerfCounters(input *[]datacollectionrules.PerfCounterDataSource) []PerfCounter {
	if input == nil {
		return make([]PerfCounter, 0)
	}

	result := make([]PerfCounter, 0)
	for _, v := range *input {
		result = append(result, PerfCounter{
			Name:                       pointer.From(v.Name),
			CounterSpecifiers:          pointer.From(v.CounterSpecifiers),
			SamplingFrequencyInSeconds: pointer.From(v.SamplingFrequencyInSeconds),
			Streams:                    pointer.FromEnumSlice(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePlatformTelemetry(input *[]datacollectionrules.PlatformTelemetryDataSource) []PlatformTelemetry {
	if input == nil {
		return make([]PlatformTelemetry, 0)
	}

	result := make([]PlatformTelemetry, 0)
	for _, v := range *input {
		result = append(result, PlatformTelemetry{
			Name:    pointer.From(v.Name),
			Streams: v.Streams,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePrometheusForwarder(input *[]datacollectionrules.PrometheusForwarderDataSource) []PrometheusForwarder {
	if input == nil {
		return make([]PrometheusForwarder, 0)
	}

	result := make([]PrometheusForwarder, 0)
	for _, v := range *input {
		labelIncludeFilter := make([]LabelIncludeFilter, 0)
		if v.LabelIncludeFilter != nil {
			for label, value := range *v.LabelIncludeFilter {
				labelIncludeFilter = append(labelIncludeFilter, LabelIncludeFilter{
					Label: label,
					Value: value,
				})
			}
		}

		result = append(result, PrometheusForwarder{
			Name:               pointer.From(v.Name),
			Streams:            pointer.FromEnumSlice(v.Streams),
			LabelIncludeFilter: labelIncludeFilter,
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceSyslog(input *[]datacollectionrules.SyslogDataSource) []Syslog {
	if input == nil {
		return make([]Syslog, 0)
	}

	result := make([]Syslog, 0)
	for _, v := range *input {
		result = append(result, Syslog{
			Name:          pointer.From(v.Name),
			FacilityNames: pointer.FromEnumSlice(v.FacilityNames),
			LogLevels:     pointer.FromEnumSlice(v.LogLevels),
			Streams:       pointer.FromEnumSlice(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleWindowsEventLogs(input *[]datacollectionrules.WindowsEventLogDataSource) []WindowsEventLog {
	if input == nil {
		return make([]WindowsEventLog, 0)
	}

	result := make([]WindowsEventLog, 0)
	for _, v := range *input {
		result = append(result, WindowsEventLog{
			Name:         pointer.From(v.Name),
			XPathQueries: pointer.From(v.XPathQueries),
			Streams:      pointer.FromEnumSlice(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleWindowsFirewallLog(input *[]datacollectionrules.WindowsFirewallLogsDataSource) []WindowsFirewallLog {
	if input == nil {
		return make([]WindowsFirewallLog, 0)
	}

	result := make([]WindowsFirewallLog, 0)
	for _, v := range *input {
		result = append(result, WindowsFirewallLog{
			Name:    pointer.From(v.Name),
			Streams: v.Streams,
		})
	}
	return result
}

func flattenDataCollectionRuleDestinations(input *datacollectionrules.DestinationsSpec) []Destination {
	if input == nil {
		return make([]Destination, 0)
	}

	return []Destination{{
		AzureMonitorMetrics: flattenDataCollectionRuleDestinationMetrics(input.AzureMonitorMetrics),
		EventHub:            flattenDataCollectionRuleDestinationEventHubs(input.EventHubs),
		EventHubDirect:      flattenDataCollectionRuleDestinationEventHubDirect(input.EventHubsDirect),
		LogAnalytics:        flattenDataCollectionRuleDestinationLogAnalytics(input.LogAnalytics),
		MonitorAccount:      flattenDataCollectionRuleDestinationMonitorAccount(input.MonitoringAccounts),
		StorageBlob:         flattenDataCollectionRuleDestinationStorageBlob(input.StorageAccounts),
		StorageBlobDirect:   flattenDataCollectionRuleDestinationStorageBlob(input.StorageBlobsDirect),
		StorageTableDirect:  flattenDataCollectionRuleDestinationStorageTableDirect(input.StorageTablesDirect),
	}}
}

func flattenDataCollectionRuleEndpoints(input *datacollectionrules.EndpointsSpec) (string, string) {
	if input == nil {
		return "", ""
	}

	return pointer.From(input.LogsIngestion), pointer.From(input.MetricsIngestion)
}

func flattenDataCollectionRuleDestinationMetrics(input *datacollectionrules.AzureMonitorMetricsDestination) []AzureMonitorMetric {
	if input == nil {
		return make([]AzureMonitorMetric, 0)
	}

	return []AzureMonitorMetric{{
		Name: pointer.From(input.Name),
	}}
}

func flattenDataCollectionRuleDestinationEventHubs(input *[]datacollectionrules.EventHubDestination) []EventHub {
	if input == nil {
		return make([]EventHub, 0)
	}

	result := make([]EventHub, 0)
	for _, v := range *input {
		result = append(result, EventHub{
			Name:               pointer.From(v.Name),
			EventHubResourceId: pointer.From(v.EventHubResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationEventHubDirect(input *[]datacollectionrules.EventHubDirectDestination) []EventHub {
	if input == nil {
		return make([]EventHub, 0)
	}

	result := make([]EventHub, 0)
	for _, v := range *input {
		result = append(result, EventHub{
			Name:               pointer.From(v.Name),
			EventHubResourceId: pointer.From(v.EventHubResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationLogAnalytics(input *[]datacollectionrules.LogAnalyticsDestination) []LogAnalytic {
	if input == nil {
		return make([]LogAnalytic, 0)
	}

	result := make([]LogAnalytic, 0)
	for _, v := range *input {
		result = append(result, LogAnalytic{
			Name:                pointer.From(v.Name),
			WorkspaceResourceId: pointer.From(v.WorkspaceResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationMonitorAccount(input *[]datacollectionrules.MonitoringAccountDestination) []MonitorAccount {
	if input == nil {
		return make([]MonitorAccount, 0)
	}

	result := make([]MonitorAccount, 0)
	for _, v := range *input {
		result = append(result, MonitorAccount{
			Name:      pointer.From(v.Name),
			AccountId: pointer.From(v.AccountResourceId),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationStorageBlob(input *[]datacollectionrules.StorageBlobDestination) []StorageBlob {
	if input == nil {
		return make([]StorageBlob, 0)
	}

	result := make([]StorageBlob, 0)
	for _, v := range *input {
		result = append(result, StorageBlob{
			Name:             pointer.From(v.Name),
			StorageAccountId: pointer.From(v.StorageAccountResourceId),
			ContainerName:    pointer.From(v.ContainerName),
		})
	}
	return result
}

func flattenDataCollectionRuleDestinationStorageTableDirect(input *[]datacollectionrules.StorageTableDestination) []StorageTableDirect {
	if input == nil {
		return make([]StorageTableDirect, 0)
	}

	result := make([]StorageTableDirect, 0)
	for _, v := range *input {
		result = append(result, StorageTableDirect{
			Name:             pointer.From(v.Name),
			StorageAccountId: pointer.From(v.StorageAccountResourceId),
			TableName:        pointer.From(v.TableName),
		})
	}
	return result
}

func flattenDataCollectionRuleStreamDeclarations(input *map[string]datacollectionrules.StreamDeclaration) []StreamDeclaration {
	if input == nil {
		return make([]StreamDeclaration, 0)
	}

	result := make([]StreamDeclaration, 0)
	for name, stream := range *input {
		if stream.Columns == nil {
			continue
		}

		columns := make([]StreamDeclarationColumn, 0)
		for _, column := range *stream.Columns {
			if column.Name == nil || column.Type == nil {
				continue
			}
			columns = append(columns, StreamDeclarationColumn{
				Name: *column.Name,
				Type: string(*column.Type),
			})
		}

		streamDeclaration := StreamDeclaration{
			StreamName: name,
			Column:     columns,
		}
		result = append(result, streamDeclaration)
	}

	return result
}

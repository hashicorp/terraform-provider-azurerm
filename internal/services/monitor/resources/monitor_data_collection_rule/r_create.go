// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_data_collection_rule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-11/datacollectionrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := datacollectionrules.NewDataCollectionRuleID(subscriptionId, state.ResourceGroupName, state.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			dataSources, err := expandDataCollectionRuleDataSources(state.DataSources)
			if err != nil {
				return err
			}

			identityValue, err := identity.ExpandLegacySystemAndUserAssignedMap(metadata.ResourceData.Get("identity").([]any))
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			input := datacollectionrules.DataCollectionRuleResource{
				Identity: identityValue,
				Kind:     expandDataCollectionRuleKind(state.Kind),
				Location: location.Normalize(state.Location),
				Name:     pointer.To(state.Name),
				Properties: &datacollectionrules.DataCollectionRule{
					DataFlows:          expandDataCollectionRuleDataFlows(state.DataFlows),
					DataSources:        dataSources,
					Description:        pointer.To(state.Description),
					Destinations:       expandDataCollectionRuleDestinations(state.Destinations),
					StreamDeclarations: expandDataCollectionRuleStreamDeclarations(state.StreamDeclaration),
				},
				Tags: tags.Expand(state.Tags),
			}

			if state.DataCollectionEndpointId != "" {
				input.Properties.DataCollectionEndpointId = pointer.To(state.DataCollectionEndpointId)
			}

			if _, err := client.Create(ctx, id, input); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandDataCollectionRuleKind(input string) *datacollectionrules.KnownDataCollectionRuleResourceKind {
	if input == "" {
		return nil
	}

	return pointer.ToEnum[datacollectionrules.KnownDataCollectionRuleResourceKind](input)
}

func expandDataCollectionRuleDataFlows(input []DataFlow) *[]datacollectionrules.DataFlow {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.DataFlow, 0)
	for _, v := range input {
		dataFlow := datacollectionrules.DataFlow{
			Destinations: pointer.To(v.Destinations),
			Streams:      expandDataCollectionRuleDataFlowStreams(v.Streams),
		}

		if v.BuiltInTransform != "" {
			dataFlow.BuiltInTransform = pointer.To(v.BuiltInTransform)
		}

		if v.OutputStream != "" {
			dataFlow.OutputStream = pointer.To(v.OutputStream)
		}

		if v.TransformKql != "" {
			dataFlow.TransformKql = pointer.To(v.TransformKql)
		}

		result = append(result, dataFlow)
	}
	return &result
}

func expandDataCollectionRuleDataFlowStreams(input []string) *[]datacollectionrules.KnownDataFlowStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownDataFlowStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownDataFlowStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSources(input []DataSourceModel) (*datacollectionrules.DataSourcesSpec, error) {
	if len(input) == 0 {
		return nil, nil
	}

	extension, err := expandDataCollectionRuleDataSourceExtensions(input[0].Extensions)
	if err != nil {
		return nil, err
	}

	return &datacollectionrules.DataSourcesSpec{
		DataImports:         expandDataCollectionRuleDataSourceDataImports(input[0].DataImport),
		Extensions:          extension,
		IisLogs:             expandDataCollectionRuleDataSourceIisLogs(input[0].IisLog),
		LogFiles:            expandDataCollectionRuleDataSourceLogFiles(input[0].LogFile),
		PerformanceCounters: expandDataCollectionRuleDataSourcePerfCounters(input[0].PerformanceCounters),
		PlatformTelemetry:   expandDataCollectionRuleDataSourcePlatformTelemetry(input[0].PlatformTelemetry),
		PrometheusForwarder: expandDataCollectionRuleDataSourcePrometheusForwarder(input[0].PrometheusForwarder),
		Syslog:              expandDataCollectionRuleDataSourceSyslog(input[0].Syslog),
		WindowsEventLogs:    expandDataCollectionRuleDataSourceWindowsEventLogs(input[0].WindowsEventLogs),
		WindowsFirewallLogs: expandDataCollectionRuleDataSourceWindowsFirewallLogs(input[0].WindowsFirewallLog),
	}, nil
}

func expandDataCollectionRuleDataSourceDataImports(input []DataImport) *datacollectionrules.DataImportSources {
	if len(input) == 0 || len(input[0].EventHubDataSource) == 0 {
		return nil
	}

	result := &datacollectionrules.DataImportSources{
		EventHub: &datacollectionrules.EventHubDataSource{
			Name:   pointer.To(input[0].EventHubDataSource[0].Name),
			Stream: pointer.To(input[0].EventHubDataSource[0].Stream),
		},
	}

	if consumerGroup := input[0].EventHubDataSource[0].ConsumerGroup; consumerGroup != "" {
		result.EventHub.ConsumerGroup = pointer.To(consumerGroup)
	}

	return result
}

func expandDataCollectionRuleDataSourceExtensions(input []Extension) (*[]datacollectionrules.ExtensionDataSource, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]datacollectionrules.ExtensionDataSource, 0)
	for _, v := range input {
		var extensionSettings any
		if v.ExtensionSettings != "" {
			settings, err := pluginsdk.ExpandJsonFromString(v.ExtensionSettings)
			if err != nil {
				return nil, err
			}
			extensionSettings = settings
		}

		result = append(result, datacollectionrules.ExtensionDataSource{
			ExtensionName:     v.ExtensionName,
			ExtensionSettings: &extensionSettings,
			InputDataSources:  pointer.To(v.InputDataSources),
			Name:              pointer.To(v.Name),
			Streams:           expandDataCollectionRuleDataSourceExtensionStreams(v.Streams),
		})
	}
	return &result, nil
}

func expandDataCollectionRuleDataSourceExtensionStreams(input []string) *[]datacollectionrules.KnownExtensionDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownExtensionDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownExtensionDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceIisLogs(input []IisLog) *[]datacollectionrules.IisLogsDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.IisLogsDataSource, 0)
	for _, v := range input {
		iisLog := datacollectionrules.IisLogsDataSource{
			Name:    pointer.To(v.Name),
			Streams: v.Streams,
		}

		if len(v.LogDirectories) != 0 {
			iisLog.LogDirectories = pointer.To(v.LogDirectories)
		}

		result = append(result, iisLog)
	}

	return &result
}

func expandDataCollectionRuleDataSourceLogFiles(input []LogFile) *[]datacollectionrules.LogFilesDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.LogFilesDataSource, 0)
	for _, v := range input {
		logFile := datacollectionrules.LogFilesDataSource{
			Name:         pointer.To(v.Name),
			Streams:      v.Streams,
			FilePatterns: v.FilePatterns,
			Format:       datacollectionrules.KnownLogFilesDataSourceFormat(v.Format),
		}

		if len(v.Settings) != 0 && len(v.Settings[0].Text) != 0 {
			logFile.Settings = &datacollectionrules.LogFileSettings{
				Text: &datacollectionrules.LogFileTextSettings{
					RecordStartTimestampFormat: datacollectionrules.KnownLogFileTextSettingsRecordStartTimestampFormat(v.Settings[0].Text[0].RecordStartTimestampFormat),
				},
			}
		}

		result = append(result, logFile)
	}

	return &result
}

func expandDataCollectionRuleDataSourcePerfCounters(input []PerfCounter) *[]datacollectionrules.PerfCounterDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.PerfCounterDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.PerfCounterDataSource{
			CounterSpecifiers:          pointer.To(v.CounterSpecifiers),
			Name:                       pointer.To(v.Name),
			SamplingFrequencyInSeconds: pointer.To(v.SamplingFrequencyInSeconds),
			Streams:                    expandDataCollectionRuleDataSourcePerfCounterStreams(v.Streams),
		})
	}

	return &result
}

func expandDataCollectionRuleDataSourcePerfCounterStreams(input []string) *[]datacollectionrules.KnownPerfCounterDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownPerfCounterDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownPerfCounterDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourcePlatformTelemetry(input []PlatformTelemetry) *[]datacollectionrules.PlatformTelemetryDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.PlatformTelemetryDataSource, 0)
	for _, v := range input {
		platformTelemetry := datacollectionrules.PlatformTelemetryDataSource{
			Name:    pointer.To(v.Name),
			Streams: v.Streams,
		}

		result = append(result, platformTelemetry)
	}

	return &result
}

func expandDataCollectionRuleDataSourcePrometheusForwarder(input []PrometheusForwarder) *[]datacollectionrules.PrometheusForwarderDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.PrometheusForwarderDataSource, 0)
	for _, v := range input {
		streams := make([]datacollectionrules.KnownPrometheusForwarderDataSourceStreams, 0)
		for _, stream := range v.Streams {
			streams = append(streams, datacollectionrules.KnownPrometheusForwarderDataSourceStreams(stream))
		}

		prometheusForwarder := datacollectionrules.PrometheusForwarderDataSource{
			Name:    pointer.To(v.Name),
			Streams: &streams,
		}

		if len(v.LabelIncludeFilter) != 0 {
			labelIncludeFilter := make(map[string]string, 0)
			for _, filter := range v.LabelIncludeFilter {
				labelIncludeFilter[filter.Label] = filter.Value
			}

			prometheusForwarder.LabelIncludeFilter = &labelIncludeFilter
		}

		result = append(result, prometheusForwarder)
	}

	return &result
}

func expandDataCollectionRuleDataSourceSyslog(input []Syslog) *[]datacollectionrules.SyslogDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.SyslogDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.SyslogDataSource{
			FacilityNames: expandDataCollectionRuleDataSourceSyslogFacilityNames(v.FacilityNames),
			LogLevels:     expandDataCollectionRuleDataSourceSyslogLogLevels(v.LogLevels),
			Name:          pointer.To(v.Name),
			Streams:       expandDataCollectionRuleDataSourceSyslogStreams(v.Streams),
		})
	}
	return &result
}

func expandDataCollectionRuleDataSourceSyslogStreams(input []string) *[]datacollectionrules.KnownSyslogDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownSyslogDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownSyslogDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceSyslogFacilityNames(input []string) *[]datacollectionrules.KnownSyslogDataSourceFacilityNames {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownSyslogDataSourceFacilityNames, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownSyslogDataSourceFacilityNames(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceSyslogLogLevels(input []string) *[]datacollectionrules.KnownSyslogDataSourceLogLevels {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownSyslogDataSourceLogLevels, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownSyslogDataSourceLogLevels(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceWindowsEventLogs(input []WindowsEventLog) *[]datacollectionrules.WindowsEventLogDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.WindowsEventLogDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.WindowsEventLogDataSource{
			Name:         pointer.To(v.Name),
			Streams:      expandDataCollectionRuleDataSourceWindowsEventLogsStreams(v.Streams),
			XPathQueries: pointer.To(v.XPathQueries),
		})
	}
	return &result
}

func expandDataCollectionRuleDataSourceWindowsEventLogsStreams(input []string) *[]datacollectionrules.KnownWindowsEventLogDataSourceStreams {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.KnownWindowsEventLogDataSourceStreams, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.KnownWindowsEventLogDataSourceStreams(v))
	}
	return &result
}

func expandDataCollectionRuleDataSourceWindowsFirewallLogs(input []WindowsFirewallLog) *[]datacollectionrules.WindowsFirewallLogsDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.WindowsFirewallLogsDataSource, 0)
	for _, v := range input {
		windowsFirewallLog := datacollectionrules.WindowsFirewallLogsDataSource{
			Name:    pointer.To(v.Name),
			Streams: v.Streams,
		}

		result = append(result, windowsFirewallLog)
	}

	return &result
}

func expandDataCollectionRuleDestinations(input []Destination) *datacollectionrules.DestinationsSpec {
	if len(input) == 0 {
		return nil
	}

	return &datacollectionrules.DestinationsSpec{
		AzureMonitorMetrics: expandDataCollectionRuleDestinationMetrics(input[0].AzureMonitorMetrics),
		EventHubsDirect:     expandDataCollectionRuleDestinationEventHubsDirect(input[0].EventHubDirect),
		EventHubs:           expandDataCollectionRuleDestinationEventHubs(input[0].EventHub),
		LogAnalytics:        expandDataCollectionRuleDestinationLogAnalytics(input[0].LogAnalytics),
		MonitoringAccounts:  expandDataCollectionRuleDestinationMonitoringAccounts(input[0].MonitorAccount),
		StorageAccounts:     expandDataCollectionRuleDestinationStorageBlobs(input[0].StorageBlob),
		StorageBlobsDirect:  expandDataCollectionRuleDestinationStorageBlobs(input[0].StorageBlobDirect),
		StorageTablesDirect: expandDataCollectionRuleDestinationStorageTableDirect(input[0].StorageTableDirect),
	}
}

func expandDataCollectionRuleDestinationMetrics(input []AzureMonitorMetric) *datacollectionrules.AzureMonitorMetricsDestination {
	if len(input) == 0 {
		return nil
	}

	return &datacollectionrules.AzureMonitorMetricsDestination{
		Name: pointer.To(input[0].Name),
	}
}

func expandDataCollectionRuleDestinationEventHubs(input []EventHub) *[]datacollectionrules.EventHubDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.EventHubDestination, 0)
	for _, v := range input {
		eventhub := datacollectionrules.EventHubDestination{
			Name:               pointer.To(v.Name),
			EventHubResourceId: pointer.To(v.EventHubResourceId),
		}

		result = append(result, eventhub)
	}

	return &result
}

func expandDataCollectionRuleDestinationEventHubsDirect(input []EventHub) *[]datacollectionrules.EventHubDirectDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.EventHubDirectDestination, 0)
	for _, v := range input {
		eventhub := datacollectionrules.EventHubDirectDestination{
			Name:               pointer.To(v.Name),
			EventHubResourceId: pointer.To(v.EventHubResourceId),
		}

		result = append(result, eventhub)
	}

	return &result
}

func expandDataCollectionRuleDestinationLogAnalytics(input []LogAnalytic) *[]datacollectionrules.LogAnalyticsDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.LogAnalyticsDestination, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.LogAnalyticsDestination{
			Name:                pointer.To(v.Name),
			WorkspaceResourceId: pointer.To(v.WorkspaceResourceId),
		})
	}
	return &result
}

func expandDataCollectionRuleDestinationMonitoringAccounts(input []MonitorAccount) *[]datacollectionrules.MonitoringAccountDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.MonitoringAccountDestination, 0)
	for _, v := range input {
		monitorAccount := datacollectionrules.MonitoringAccountDestination{
			Name:              pointer.To(v.Name),
			AccountResourceId: pointer.To(v.AccountId),
		}

		result = append(result, monitorAccount)
	}

	return &result
}

func expandDataCollectionRuleDestinationStorageBlobs(input []StorageBlob) *[]datacollectionrules.StorageBlobDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.StorageBlobDestination, 0)
	for _, v := range input {
		monitorAccount := datacollectionrules.StorageBlobDestination{
			Name:                     pointer.To(v.Name),
			StorageAccountResourceId: pointer.To(v.StorageAccountId),
			ContainerName:            pointer.To(v.ContainerName),
		}

		result = append(result, monitorAccount)
	}

	return &result
}

func expandDataCollectionRuleDestinationStorageTableDirect(input []StorageTableDirect) *[]datacollectionrules.StorageTableDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.StorageTableDestination, 0)
	for _, v := range input {
		monitorAccount := datacollectionrules.StorageTableDestination{
			Name:                     pointer.To(v.Name),
			StorageAccountResourceId: pointer.To(v.StorageAccountId),
			TableName:                pointer.To(v.TableName),
		}

		result = append(result, monitorAccount)
	}

	return &result
}

func expandDataCollectionRuleStreamDeclarations(input []StreamDeclaration) *map[string]datacollectionrules.StreamDeclaration {
	if len(input) == 0 {
		return nil
	}

	result := make(map[string]datacollectionrules.StreamDeclaration, 0)
	for _, v := range input {
		columns := make([]datacollectionrules.ColumnDefinition, 0)
		for _, column := range v.Column {
			columns = append(columns, datacollectionrules.ColumnDefinition{
				Name: pointer.To(column.Name),
				Type: pointer.ToEnum[datacollectionrules.KnownColumnDefinitionType](column.Type),
			})
		}

		result[v.StreamName] = datacollectionrules.StreamDeclaration{
			Columns: &columns,
		}
	}

	return &result
}

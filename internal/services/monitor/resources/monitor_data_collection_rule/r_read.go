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
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRulesClient
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var dataCollectionEndpointId, description, immutableId, kind, loc string
			var tag map[string]any
			var dataFlows []DataFlow
			var dataSources []DataSourceModel
			var destinations []Destination
			var logsIngestionEndpoint, metricsIngestionEndpoint string
			var streamDeclaration []StreamDeclaration

			if model := resp.Model; model != nil {
				kind = pointer.FromEnum(model.Kind)
				loc = location.Normalize(model.Location)
				tag = tags.Flatten(model.Tags)

				identityValue, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
				if err != nil {
					return fmt.Errorf("flattening `identity`: %+v", err)
				}

				if err := metadata.ResourceData.Set("identity", identityValue); err != nil {
					return fmt.Errorf("setting `identity`: %+v", err)
				}

				if prop := model.Properties; prop != nil {
					dataCollectionEndpointId = pointer.From(prop.DataCollectionEndpointId)
					description = pointer.From(prop.Description)
					dataFlows = flattenDataCollectionRuleDataFlows(prop.DataFlows)
					dataSources = flattenDataCollectionRuleDataSources(prop.DataSources)
					destinations = flattenDataCollectionRuleDestinations(prop.Destinations)
					logsIngestionEndpoint, metricsIngestionEndpoint = flattenDataCollectionRuleEndpoints(prop.Endpoints)
					immutableId = pointer.From(prop.ImmutableId)
					streamDeclaration = flattenDataCollectionRuleStreamDeclarations(prop.StreamDeclarations)
				}
			}

			return metadata.Encode(&DataCollectionRule{
				Name:                     id.DataCollectionRuleName,
				ResourceGroupName:        id.ResourceGroupName,
				DataCollectionEndpointId: dataCollectionEndpointId,
				DataFlows:                dataFlows,
				DataSources:              dataSources,
				Description:              description,
				Destinations:             destinations,
				ImmutableId:              immutableId,
				Kind:                     kind,
				Location:                 loc,
				LogsIngestionEndpoint:    logsIngestionEndpoint,
				MetricsIngestionEndpoint: metricsIngestionEndpoint,
				StreamDeclaration:        streamDeclaration,
				Tags:                     tag,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

type DataCollectionRule struct {
	DataCollectionEndpointId string              `tfschema:"data_collection_endpoint_id"`
	DataFlows                []DataFlow          `tfschema:"data_flow"`
	DataSources              []DataSourceModel   `tfschema:"data_sources"`
	Description              string              `tfschema:"description"`
	Destinations             []Destination       `tfschema:"destinations"`
	ImmutableId              string              `tfschema:"immutable_id"`
	Kind                     string              `tfschema:"kind"`
	LogsIngestionEndpoint    string              `tfschema:"logs_ingestion_endpoint"`
	MetricsIngestionEndpoint string              `tfschema:"metrics_ingestion_endpoint"`
	Name                     string              `tfschema:"name"`
	Location                 string              `tfschema:"location"`
	ResourceGroupName        string              `tfschema:"resource_group_name"`
	StreamDeclaration        []StreamDeclaration `tfschema:"stream_declaration"`
	Tags                     map[string]any      `tfschema:"tags"`
}

type DataFlow struct {
	BuiltInTransform string   `tfschema:"built_in_transform"`
	Destinations     []string `tfschema:"destinations"`
	OutputStream     string   `tfschema:"output_stream"`
	Streams          []string `tfschema:"streams"`
	TransformKql     string   `tfschema:"transform_kql"`
}

type DataSourceModel struct {
	DataImport          []DataImport          `tfschema:"data_import"`
	Extensions          []Extension           `tfschema:"extension"`
	IisLog              []IisLog              `tfschema:"iis_log"`
	LogFile             []LogFile             `tfschema:"log_file"`
	PerformanceCounters []PerfCounter         `tfschema:"performance_counter"`
	PlatformTelemetry   []PlatformTelemetry   `tfschema:"platform_telemetry"`
	PrometheusForwarder []PrometheusForwarder `tfschema:"prometheus_forwarder"`
	Syslog              []Syslog              `tfschema:"syslog"`
	WindowsEventLogs    []WindowsEventLog     `tfschema:"windows_event_log"`
	WindowsFirewallLog  []WindowsFirewallLog  `tfschema:"windows_firewall_log"`
}

type Destination struct {
	AzureMonitorMetrics []AzureMonitorMetric `tfschema:"azure_monitor_metrics"`
	EventHubDirect      []EventHub           `tfschema:"event_hub_direct"`
	EventHub            []EventHub           `tfschema:"event_hub"`
	LogAnalytics        []LogAnalytic        `tfschema:"log_analytics"`
	MonitorAccount      []MonitorAccount     `tfschema:"monitor_account"`
	StorageBlob         []StorageBlob        `tfschema:"storage_blob"`
	StorageBlobDirect   []StorageBlob        `tfschema:"storage_blob_direct"`
	StorageTableDirect  []StorageTableDirect `tfschema:"storage_table_direct"`
}

type DataImport struct {
	EventHubDataSource []EventHubDataSource `tfschema:"event_hub_data_source"`
}

type EventHubDataSource struct {
	ConsumerGroup string `tfschema:"consumer_group"`
	Name          string `tfschema:"name"`
	Stream        string `tfschema:"stream"`
}

type Extension struct {
	ExtensionName     string   `tfschema:"extension_name"`
	ExtensionSettings string   `tfschema:"extension_json"`
	InputDataSources  []string `tfschema:"input_data_sources"`
	Name              string   `tfschema:"name"`
	Streams           []string `tfschema:"streams"`
}

type IisLog struct {
	Name           string   `tfschema:"name"`
	Streams        []string `tfschema:"streams"`
	LogDirectories []string `tfschema:"log_directories"`
}

type LogFile struct {
	Name         string           `tfschema:"name"`
	Streams      []string         `tfschema:"streams"`
	FilePatterns []string         `tfschema:"file_patterns"`
	Format       string           `tfschema:"format"`
	Settings     []LogFileSetting `tfschema:"settings"`
}

type PerfCounter struct {
	CounterSpecifiers          []string `tfschema:"counter_specifiers"`
	Name                       string   `tfschema:"name"`
	SamplingFrequencyInSeconds int64    `tfschema:"sampling_frequency_in_seconds"`
	Streams                    []string `tfschema:"streams"`
}

type PlatformTelemetry struct {
	Name    string   `tfschema:"name"`
	Streams []string `tfschema:"streams"`
}

type PrometheusForwarder struct {
	Name               string               `tfschema:"name"`
	Streams            []string             `tfschema:"streams"`
	LabelIncludeFilter []LabelIncludeFilter `tfschema:"label_include_filter"`
}

type LabelIncludeFilter struct {
	Label string `tfschema:"label"`
	Value string `tfschema:"value"`
}

type Syslog struct {
	FacilityNames []string `tfschema:"facility_names"`
	LogLevels     []string `tfschema:"log_levels"`
	Name          string   `tfschema:"name"`
	Streams       []string `tfschema:"streams"`
}

type WindowsEventLog struct {
	Name         string   `tfschema:"name"`
	Streams      []string `tfschema:"streams"`
	XPathQueries []string `tfschema:"x_path_queries"`
}

type WindowsFirewallLog struct {
	Name    string   `tfschema:"name"`
	Streams []string `tfschema:"streams"`
}

type AzureMonitorMetric struct {
	Name string `tfschema:"name"`
}

type EventHub struct {
	EventHubResourceId string `tfschema:"event_hub_id"`
	Name               string `tfschema:"name"`
}

type LogAnalytic struct {
	Name                string `tfschema:"name"`
	WorkspaceResourceId string `tfschema:"workspace_resource_id"`
}

type MonitorAccount struct {
	AccountId string `tfschema:"monitor_account_id"`
	Name      string `tfschema:"name"`
}

type StorageBlob struct {
	ContainerName    string `tfschema:"container_name"`
	Name             string `tfschema:"name"`
	StorageAccountId string `tfschema:"storage_account_id"`
}

type StorageTableDirect struct {
	Name             string `tfschema:"name"`
	StorageAccountId string `tfschema:"storage_account_id"`
	TableName        string `tfschema:"table_name"`
}

type LogFileSetting struct {
	Text []TextSetting `tfschema:"text"`
}

type TextSetting struct {
	RecordStartTimestampFormat string `tfschema:"record_start_timestamp_format"`
}

type StreamDeclaration struct {
	StreamName string                    `tfschema:"stream_name"`
	Column     []StreamDeclarationColumn `tfschema:"column"`
}

type StreamDeclarationColumn struct {
	Name string `tfschema:"name"`
	Type string `tfschema:"type"`
}

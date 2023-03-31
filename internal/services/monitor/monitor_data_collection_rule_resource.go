package monitor

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-06-01/datacollectionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DataCollectionRule struct {
	DataFlows         []DataFlow             `tfschema:"data_flow"`
	DataSources       []DataSource           `tfschema:"data_sources"`
	Description       string                 `tfschema:"description"`
	Destinations      []Destination          `tfschema:"destinations"`
	Kind              string                 `tfschema:"kind"`
	Name              string                 `tfschema:"name"`
	Location          string                 `tfschema:"location"`
	ResourceGroupName string                 `tfschema:"resource_group_name"`
	Tags              map[string]interface{} `tfschema:"tags"`
}

type DataFlow struct {
	Destinations []string `tfschema:"destinations"`
	Streams      []string `tfschema:"streams"`
}

type DataSource struct {
	Extensions          []Extension       `tfschema:"extension"`
	PerformanceCounters []PerfCounter     `tfschema:"performance_counter"`
	Syslog              []Syslog          `tfschema:"syslog"`
	WindowsEventLogs    []WindowsEventLog `tfschema:"windows_event_log"`
}

type Destination struct {
	AzureMonitorMetrics []AzureMonitorMetric `tfschema:"azure_monitor_metrics"`
	LogAnalytics        []LogAnalytic        `tfschema:"log_analytics"`
}

type Extension struct {
	ExtensionName     string   `tfschema:"extension_name"`
	ExtensionSettings string   `tfschema:"extension_json"`
	InputDataSources  []string `tfschema:"input_data_sources"`
	Name              string   `tfschema:"name"`
	Streams           []string `tfschema:"streams"`
}

type PerfCounter struct {
	CounterSpecifiers          []string `tfschema:"counter_specifiers"`
	Name                       string   `tfschema:"name"`
	SamplingFrequencyInSeconds int64    `tfschema:"sampling_frequency_in_seconds"`
	Streams                    []string `tfschema:"streams"`
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

type AzureMonitorMetric struct {
	Name string `tfschema:"name"`
}

type LogAnalytic struct {
	Name                string `tfschema:"name"`
	WorkspaceResourceId string `tfschema:"workspace_resource_id"`
}

type DataCollectionRuleResource struct{}

func (r DataCollectionRuleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"data_flow": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destinations": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
					"streams": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
		},

		"destinations": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"azure_monitor_metrics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
						AtLeastOneOf: []string{"destinations.0.azure_monitor_metrics", "destinations.0.log_analytics"},
					},
					"log_analytics": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"workspace_resource_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: workspaces.ValidateWorkspaceID,
								},
							},
						},
						AtLeastOneOf: []string{"destinations.0.azure_monitor_metrics", "destinations.0.log_analytics"},
					},
				},
			},
		},

		"data_sources": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"extension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"extension_name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"extension_json": {
									Type:             pluginsdk.TypeString,
									Optional:         true,
									ValidateFunc:     validation.StringIsJSON,
									DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
								},
								"input_data_sources": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"performance_counter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"sampling_frequency_in_seconds": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 300),
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"counter_specifiers": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"syslog": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"facility_names": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceFacilityNames(),
											false),
									},
								},
								"log_levels": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											datacollectionrules.PossibleValuesForKnownSyslogDataSourceLogLevels(), false),
									},
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Optional: !features.FourPointOhBeta(),
									Computed: !features.FourPointOhBeta(),
									Required: features.FourPointOhBeta(),
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"windows_event_log": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"streams": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
								"x_path_queries": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"kind": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice(
				datacollectionrules.PossibleValuesForKnownDataCollectionRuleResourceKind(), false),
		},

		"tags": commonschema.Tags(),
	}
}

func (r DataCollectionRuleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DataCollectionRuleResource) ResourceType() string {
	return "azurerm_monitor_data_collection_rule"
}

func (r DataCollectionRuleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return datacollectionrules.ValidateDataCollectionRuleID
}

func (r DataCollectionRuleResource) ModelObject() interface{} {
	return &DataCollectionRule{}
}

func (r DataCollectionRuleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			metadata.Logger.Info("Decoding state..")
			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			client := metadata.Client.Monitor.DataCollectionRulesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			id := datacollectionrules.NewDataCollectionRuleID(subscriptionId, state.ResourceGroupName, state.Name)
			metadata.Logger.Infof("creating %s", id)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			dataSources, err := expandDataCollectionRuleDataSources(state.DataSources)
			if err != nil {
				return err
			}

			input := datacollectionrules.DataCollectionRuleResource{
				Kind:     expandDataCollectionRuleKind(state.Kind),
				Location: azure.NormalizeLocation(state.Location),
				Name:     utils.String(state.Name),
				Properties: &datacollectionrules.DataCollectionRule{
					DataFlows:    expandDataCollectionRuleDataFlows(state.DataFlows),
					DataSources:  dataSources,
					Description:  utils.String(state.Description),
					Destinations: expandDataCollectionRuleDestinations(state.Destinations),
				},
				Tags: tags.Expand(state.Tags),
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

func (r DataCollectionRuleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRulesClient
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("retrieving %s", *id)
			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					metadata.Logger.Infof("%s was not found - removing from state!", *id)
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			var description, kind, location string
			var tag map[string]interface{}
			var dataFlows []DataFlow
			var dataSources []DataSource
			var destinations []Destination

			if model := resp.Model; model != nil {
				kind = flattenDataCollectionRuleKind(model.Kind)
				location = azure.NormalizeLocation(model.Location)
				tag = tags.Flatten(model.Tags)
				if prop := model.Properties; prop != nil {
					description = flattenStringPtr(prop.Description)
					dataFlows = flattenDataCollectionRuleDataFlows(prop.DataFlows)
					dataSources = flattenDataCollectionRuleDataSources(prop.DataSources)
					destinations = flattenDataCollectionRuleDestinations(prop.Destinations)
				}
			}

			return metadata.Encode(&DataCollectionRule{
				Name:              id.DataCollectionRuleName,
				ResourceGroupName: id.ResourceGroupName,
				DataFlows:         dataFlows,
				DataSources:       dataSources,
				Description:       description,
				Destinations:      destinations,
				Kind:              kind,
				Location:          location,
				Tags:              tag,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r DataCollectionRuleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("updating %s..", *id)
			client := metadata.Client.Monitor.DataCollectionRulesClient
			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("unexpected null model of %s", *id)
			}
			existing := resp.Model
			if existing.Properties == nil {
				return fmt.Errorf("unexpected null properties of %s", *id)
			}

			var state DataCollectionRule
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			if metadata.ResourceData.HasChange("kind") {
				existing.Kind = expandDataCollectionRuleKind(state.Kind)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Tags = tags.Expand(state.Tags)
			}

			if metadata.ResourceData.HasChange("data_flow") {
				existing.Properties.DataFlows = expandDataCollectionRuleDataFlows(state.DataFlows)
			}

			if metadata.ResourceData.HasChange("data_sources") {
				dataSource, err := expandDataCollectionRuleDataSources(state.DataSources)
				if err != nil {
					return err
				}
				existing.Properties.DataSources = dataSource
			}

			if metadata.ResourceData.HasChange("description") {
				existing.Properties.Description = utils.String(state.Description)
			}

			if metadata.ResourceData.HasChange("destinations") {
				existing.Properties.Destinations = expandDataCollectionRuleDestinations(state.Destinations)
			}

			if _, err := client.Create(ctx, *id, *existing); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r DataCollectionRuleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.DataCollectionRulesClient
			id, err := datacollectionrules.ParseDataCollectionRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			metadata.Logger.Infof("deleting %s..", *id)
			resp, err := client.Delete(ctx, *id)
			if err != nil && !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func expandDataCollectionRuleKind(input string) *datacollectionrules.KnownDataCollectionRuleResourceKind {
	if input == "" {
		return nil
	}

	result := datacollectionrules.KnownDataCollectionRuleResourceKind(input)
	return &result
}

func stringSlice(input []string) *[]string {
	return &input
}

func expandDataCollectionRuleDataFlows(input []DataFlow) *[]datacollectionrules.DataFlow {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.DataFlow, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.DataFlow{
			Destinations: stringSlice(v.Destinations),
			Streams:      expandDataCollectionRuleDataFlowStreams(v.Streams),
		})
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

func expandDataCollectionRuleDataSources(input []DataSource) (*datacollectionrules.DataSourcesSpec, error) {
	if len(input) == 0 {
		return nil, nil
	}

	extension, err := expandDataCollectionRuleDataSourceExtensions(input[0].Extensions)
	if err != nil {
		return nil, err
	}
	return &datacollectionrules.DataSourcesSpec{
		Extensions:          extension,
		PerformanceCounters: expandDataCollectionRuleDataSourcePerfCounters(input[0].PerformanceCounters),
		Syslog:              expandDataCollectionRuleDataSourceSyslog(input[0].Syslog),
		WindowsEventLogs:    expandDataCollectionRuleDataSourceWindowsEventLogs(input[0].WindowsEventLogs),
	}, nil
}

func expandDataCollectionRuleDataSourceExtensions(input []Extension) (*[]datacollectionrules.ExtensionDataSource, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]datacollectionrules.ExtensionDataSource, 0)
	for _, v := range input {
		var extensionSettings interface{}
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
			InputDataSources:  stringSlice(v.InputDataSources),
			Name:              utils.String(v.Name),
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

func expandDataCollectionRuleDataSourcePerfCounters(input []PerfCounter) *[]datacollectionrules.PerfCounterDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.PerfCounterDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.PerfCounterDataSource{
			CounterSpecifiers:          stringSlice(v.CounterSpecifiers),
			Name:                       utils.String(v.Name),
			SamplingFrequencyInSeconds: utils.Int64(v.SamplingFrequencyInSeconds),
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

func expandDataCollectionRuleDataSourceSyslog(input []Syslog) *[]datacollectionrules.SyslogDataSource {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.SyslogDataSource, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.SyslogDataSource{
			FacilityNames: expandDataCollectionRuleDataSourceSyslogFacilityNames(v.FacilityNames),
			LogLevels:     expandDataCollectionRuleDataSourceSyslogLogLevels(v.LogLevels),
			Name:          utils.String(v.Name),
			Streams:       expandDataCollectionRuleDataSourceSyslogStreams(v.Streams),
		})
	}
	return &result
}

func expandDataCollectionRuleDataSourceSyslogStreams(input []string) *[]datacollectionrules.KnownSyslogDataSourceStreams {
	if len(input) == 0 {
		if !features.FourPointOhBeta() {
			return &[]datacollectionrules.KnownSyslogDataSourceStreams{datacollectionrules.KnownSyslogDataSourceStreamsMicrosoftNegativeSyslog}
		}
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
			Name:         utils.String(v.Name),
			Streams:      expandDataCollectionRuleDataSourceWindowsEventLogsStreams(v.Streams),
			XPathQueries: stringSlice(v.XPathQueries),
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

func expandDataCollectionRuleDestinations(input []Destination) *datacollectionrules.DestinationsSpec {
	if len(input) == 0 {
		return nil
	}

	return &datacollectionrules.DestinationsSpec{
		AzureMonitorMetrics: expandDataCollectionRuleDestinationMetrics(input[0].AzureMonitorMetrics),
		LogAnalytics:        expandDataCollectionRuleDestinationLogAnalytics(input[0].LogAnalytics),
	}
}

func expandDataCollectionRuleDestinationMetrics(input []AzureMonitorMetric) *datacollectionrules.AzureMonitorMetricsDestination {
	if len(input) == 0 {
		return nil
	}

	return &datacollectionrules.AzureMonitorMetricsDestination{
		Name: utils.String(input[0].Name),
	}
}

func expandDataCollectionRuleDestinationLogAnalytics(input []LogAnalytic) *[]datacollectionrules.LogAnalyticsDestination {
	if len(input) == 0 {
		return nil
	}

	result := make([]datacollectionrules.LogAnalyticsDestination, 0)
	for _, v := range input {
		result = append(result, datacollectionrules.LogAnalyticsDestination{
			Name:                utils.String(v.Name),
			WorkspaceResourceId: utils.String(v.WorkspaceResourceId),
		})
	}
	return &result
}

func flattenDataCollectionRuleKind(input *datacollectionrules.KnownDataCollectionRuleResourceKind) string {
	if input == nil {
		return ""
	}
	return string(*input)
}

func flattenStringPtr(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func flattenStringSlicePtr(input *[]string) []string {
	if input == nil {
		return make([]string, 0)
	}
	return *input
}

func flattenDataCollectionRuleDataFlows(input *[]datacollectionrules.DataFlow) []DataFlow {
	if input == nil {
		return make([]DataFlow, 0)
	}

	result := make([]DataFlow, 0)
	for _, v := range *input {
		result = append(result, DataFlow{
			Destinations: flattenStringSlicePtr(v.Destinations),
			Streams:      flattenDataCollectionRuleDataFlowStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataFlowStreams(input *[]datacollectionrules.KnownDataFlowStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDataSources(input *datacollectionrules.DataSourcesSpec) []DataSource {
	if input == nil {
		return make([]DataSource, 0)
	}

	return []DataSource{{
		Extensions:          flattenDataCollectionRuleDataSourceExtensions(input.Extensions),
		PerformanceCounters: flattenDataCollectionRuleDataSourcePerfCounters(input.PerformanceCounters),
		Syslog:              flattenDataCollectionRuleDataSourceSyslog(input.Syslog),
		WindowsEventLogs:    flattenDataCollectionRuleWindowsEventLogs(input.WindowsEventLogs),
	}}
}

func flattenDataCollectionRuleDataSourceExtensions(input *[]datacollectionrules.ExtensionDataSource) []Extension {
	if input == nil {
		return make([]Extension, 0)
	}

	result := make([]Extension, 0)
	for _, v := range *input {
		extensionSettings := ""
		if v.ExtensionSettings != nil {
			settingString, _ := pluginsdk.FlattenJsonToString((*v.ExtensionSettings).(map[string]interface{}))
			extensionSettings = settingString
		}
		result = append(result, Extension{
			ExtensionName:     v.ExtensionName,
			Name:              flattenStringPtr(v.Name),
			ExtensionSettings: extensionSettings,
			InputDataSources:  flattenStringSlicePtr(v.InputDataSources),
			Streams:           flattenDataCollectionRuleDataSourceExtensionStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceExtensionStreams(input *[]datacollectionrules.KnownExtensionDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
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
			Name:                       flattenStringPtr(v.Name),
			CounterSpecifiers:          flattenStringSlicePtr(v.CounterSpecifiers),
			SamplingFrequencyInSeconds: utils.NormaliseNilableInt64(v.SamplingFrequencyInSeconds),
			Streams:                    flattenDataCollectionRuleDataSourcePerfCounterStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourcePerfCounterStreams(input *[]datacollectionrules.KnownPerfCounterDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
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
			Name:          flattenStringPtr(v.Name),
			FacilityNames: flattenDataCollectionRuleDataSourceSyslogFacilityNames(v.FacilityNames),
			LogLevels:     flattenDataCollectionRuleDataSourceSyslogLogLevels(v.LogLevels),
			Streams:       flattenDataCollectionRuleSyslogStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleDataSourceSyslogFacilityNames(input *[]datacollectionrules.KnownSyslogDataSourceFacilityNames) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDataSourceSyslogLogLevels(input *[]datacollectionrules.KnownSyslogDataSourceLogLevels) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleSyslogStreams(input *[]datacollectionrules.KnownSyslogDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
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
			Name:         flattenStringPtr(v.Name),
			XPathQueries: flattenStringSlicePtr(v.XPathQueries),
			Streams:      flattenDataCollectionRuleWindowsEventLogStreams(v.Streams),
		})
	}
	return result
}

func flattenDataCollectionRuleWindowsEventLogStreams(input *[]datacollectionrules.KnownWindowsEventLogDataSourceStreams) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, v := range *input {
		result = append(result, string(v))
	}
	return result
}

func flattenDataCollectionRuleDestinations(input *datacollectionrules.DestinationsSpec) []Destination {
	if input == nil {
		return make([]Destination, 0)
	}

	return []Destination{{
		AzureMonitorMetrics: flattenDataCollectionRuleDestinationMetrics(input.AzureMonitorMetrics),
		LogAnalytics:        flattenDataCollectionRuleDestinationLogAnalytics(input.LogAnalytics),
	}}
}

func flattenDataCollectionRuleDestinationMetrics(input *datacollectionrules.AzureMonitorMetricsDestination) []AzureMonitorMetric {
	if input == nil {
		return make([]AzureMonitorMetric, 0)
	}

	return []AzureMonitorMetric{{
		Name: flattenStringPtr(input.Name),
	}}
}

func flattenDataCollectionRuleDestinationLogAnalytics(input *[]datacollectionrules.LogAnalyticsDestination) []LogAnalytic {
	if input == nil {
		return make([]LogAnalytic, 0)
	}

	result := make([]LogAnalytic, 0)
	for _, v := range *input {
		result = append(result, LogAnalytic{
			Name:                flattenStringPtr(v.Name),
			WorkspaceResourceId: flattenStringPtr(v.WorkspaceResourceId),
		})
	}
	return result
}

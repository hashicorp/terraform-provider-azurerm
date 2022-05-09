package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-09-01-preview/datacollection"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	logAnalyticsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	azSchema "github.com/hashicorp/terraform-provider-azurerm/internal/tf/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorDataCollectionRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorDataCollectionRuleCreateUpdate,
		Read:   resourceMonitorDataCollectionRuleRead,
		Update: resourceMonitorDataCollectionRuleCreateUpdate,
		Delete: resourceMonitorDataCollectionRuleDelete,
		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DataCollectionRuleID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"data_flows": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"streams": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(datacollection.KnownDataFlowStreamsMicrosoftEvent),
									string(datacollection.KnownDataFlowStreamsMicrosoftInsightsMetrics),
									string(datacollection.KnownDataFlowStreamsMicrosoftPerf),
									string(datacollection.KnownDataFlowStreamsMicrosoftSyslog),
									string(datacollection.KnownDataFlowStreamsMicrosoftWindowsEvent),
								}, false),
							},
						},

						"destinations": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"kind": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(datacollection.KnownDataCollectionRuleResourceKindWindows),
					string(datacollection.KnownDataCollectionRuleResourceKindLinux),
				}, false),
			},

			"azure_monitor_metrics_destination": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				AtLeastOneOf: []string{"log_analytics_destination", "azure_monitor_metrics_destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"extension_data_source": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				AtLeastOneOf: []string{"performance_counter_data_source", "windows_event_log_data_source", "syslog_data_source", "extension_data_source"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"extension_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"streams": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(datacollection.KnownExtensionDataSourceStreamsMicrosoftEvent),
									string(datacollection.KnownExtensionDataSourceStreamsMicrosoftInsightsMetrics),
									string(datacollection.KnownExtensionDataSourceStreamsMicrosoftPerf),
									string(datacollection.KnownExtensionDataSourceStreamsMicrosoftSyslog),
									string(datacollection.KnownExtensionDataSourceStreamsMicrosoftWindowsEvent),
								}, false),
							},
						},

						"input_data_sources": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"extension_setting": {
							Type:             pluginsdk.TypeString,
							Optional:         true,
							ValidateFunc:     validation.StringIsJSON,
							DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
						},
					},
				},
			},

			"log_analytics_destination": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				AtLeastOneOf: []string{"log_analytics_destination", "azure_monitor_metrics_destination"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"workspace_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: logAnalyticsValidate.LogAnalyticsWorkspaceID,
						},
					},
				},
			},

			"performance_counter_data_source": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				AtLeastOneOf: []string{"performance_counter_data_source", "windows_event_log_data_source", "syslog_data_source", "extension_data_source"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"streams": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(datacollection.KnownPerfCounterDataSourceStreamsMicrosoftPerf),
									string(datacollection.KnownPerfCounterDataSourceStreamsMicrosoftInsightsMetrics),
								}, false),
							},
						},

						"specifiers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"sampling_frequency": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"syslog_data_source": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				AtLeastOneOf: []string{"performance_counter_data_source", "windows_event_log_data_source", "syslog_data_source", "extension_data_source"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"streams": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(datacollection.KnownSyslogDataSourceStreamsMicrosoftSyslog),
								}, false),
							},
						},

						"log_levels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(datacollection.KnownSyslogDataSourceLogLevelsAlert),
									string(datacollection.KnownSyslogDataSourceLogLevelsAsterisk),
									string(datacollection.KnownSyslogDataSourceLogLevelsCritical),
									string(datacollection.KnownSyslogDataSourceLogLevelsDebug),
									string(datacollection.KnownSyslogDataSourceLogLevelsEmergency),
									string(datacollection.KnownSyslogDataSourceLogLevelsError),
									string(datacollection.KnownSyslogDataSourceLogLevelsInfo),
									string(datacollection.KnownSyslogDataSourceLogLevelsNotice),
									string(datacollection.KnownSyslogDataSourceLogLevelsWarning),
								}, false),
							},
						},

						"facility_names": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(datacollection.KnownSyslogDataSourceFacilityNamesAsterisk),
									string(datacollection.KnownSyslogDataSourceFacilityNamesAuth),
									string(datacollection.KnownSyslogDataSourceFacilityNamesAuthpriv),
									string(datacollection.KnownSyslogDataSourceFacilityNamesCron),
									string(datacollection.KnownSyslogDataSourceFacilityNamesDaemon),
									string(datacollection.KnownSyslogDataSourceFacilityNamesKern),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal0),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal1),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal2),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal3),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal4),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal5),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal6),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLocal7),
									string(datacollection.KnownSyslogDataSourceFacilityNamesLpr),
									string(datacollection.KnownSyslogDataSourceFacilityNamesMail),
									string(datacollection.KnownSyslogDataSourceFacilityNamesMark),
									string(datacollection.KnownSyslogDataSourceFacilityNamesNews),
									string(datacollection.KnownSyslogDataSourceFacilityNamesSyslog),
									string(datacollection.KnownSyslogDataSourceFacilityNamesUser),
									string(datacollection.KnownSyslogDataSourceFacilityNamesUucp),
								}, false),
							},
						},
					},
				},
			},

			"windows_event_log_data_source": {
				Type:         schema.TypeList,
				Optional:     true,
				MinItems:     1,
				AtLeastOneOf: []string{"performance_counter_data_source", "windows_event_log_data_source", "syslog_data_source", "extension_data_source"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"streams": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(datacollection.KnownWindowsEventLogDataSourceStreamsMicrosoftWindowsEvent),
									string(datacollection.KnownWindowsEventLogDataSourceStreamsMicrosoftEvent),
								}, false),
							},
						},

						"xpath_queries": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorDataCollectionRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Monitor.DataCollectionRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewDataCollectionRuleID(subscriptionId, resourceGroup, name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing Monitor DataCollectionRule %q: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_monitor_data_collection_rule", id.ID())
		}

	}

	extensions, err := expandExtensionDataSources(d.Get("extension_data_source").([]interface{}))
	if err != nil {
		return fmt.Errorf("error expanding extension_data_source: %+v", err)
	}
	body := datacollection.RuleResource{
		RuleResourceProperties: &datacollection.RuleResourceProperties{
			Description: utils.String(d.Get("description").(string)),
			DataSources: &datacollection.RuleDataSources{
				PerformanceCounters: expandPerformanceCounterDataSources(d.Get("performance_counter_data_source").([]interface{})),
				WindowsEventLogs:    expandWindowsEventLogDataSources(d.Get("windows_event_log_data_source").([]interface{})),
				Syslog:              expandSyslogDataSources(d.Get("syslog_data_source").([]interface{})),
				Extensions:          extensions,
			},
			Destinations: &datacollection.RuleDestinations{
				LogAnalytics:        expandLogAnalyticsDestinations(d.Get("log_analytics_destination").([]interface{})),
				AzureMonitorMetrics: expandAzureMonitorMetricsDestinations(d.Get("azure_monitor_metrics_destination").([]interface{})),
			},
			DataFlows: expandDataFlows(d.Get("data_flows").([]interface{})),
		},
		Kind:     datacollection.KnownDataCollectionRuleResourceKind(d.Get("kind").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Name:     utils.String(name),
		Location: utils.String(location.Normalize(d.Get("location").(string))),
	}

	_, err = client.Create(ctx, id.ResourceGroup, id.Name, &body)
	if err != nil {
		return fmt.Errorf("creating Monitor DataCollectionRule %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorDataCollectionRuleRead(d, meta)
}

func resourceMonitorDataCollectionRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.DataCollectionRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataCollectionRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Monitor DataCollectionRule %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Monitor DataCollectionRule %q: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	if resp.Description != nil {
		d.Set("description", *resp.Description)
	}
	d.Set("kind", resp.Kind)
	if props := resp.RuleResourceProperties; props != nil {
		if props.DataSources != nil {
			d.Set("performance_counter_data_source", flattenPerformanceCounterDataSources(props.DataSources.PerformanceCounters))
			d.Set("windows_event_log_data_source", flattenWindowsEventLogDataSources(props.DataSources.WindowsEventLogs))
			d.Set("syslog_data_source", flattenSyslogDataSources(props.DataSources.Syslog))

			extensionDataSource, err := flattenExtensionDataSources(props.DataSources.Extensions)
			if err != nil {
				return fmt.Errorf("error setting extension_data_source: %+v", err)
			}
			d.Set("extension_data_source", extensionDataSource)
		}
		if props.Destinations != nil {
			d.Set("log_analytics_destination", flattenLogAnalyticsDestinations(props.Destinations.LogAnalytics))
			d.Set("azure_monitor_metrics_destination", flattenAzureMonitorMetricsDestinations(props.Destinations.AzureMonitorMetrics))
		}
		d.Set("data_flows", flattenDataFlows(props.DataFlows))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMonitorDataCollectionRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.DataCollectionRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DataCollectionRuleID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Monitor DataCollectionRule %q: %+v", id, err)
	}

	return nil
}

func expandPerformanceCounterDataSources(p []interface{}) *[]datacollection.PerfCounterDataSource {
	if len(p) == 0 {
		return nil
	}
	dataSources := make([]datacollection.PerfCounterDataSource, 0)

	for _, v := range p {
		value := v.(map[string]interface{})
		streams := make([]datacollection.KnownPerfCounterDataSourceStreams, 0)
		if value["streams"] != nil {
			for _, streamRaw := range value["streams"].([]interface{}) {
				streams = append(streams, datacollection.KnownPerfCounterDataSourceStreams(streamRaw.(string)))
			}
		}
		dataSources = append(dataSources, datacollection.PerfCounterDataSource{
			Streams:                    &streams,
			SamplingFrequencyInSeconds: utils.Int32(int32(value["sampling_frequency"].(int))),
			CounterSpecifiers:          utils.ExpandStringSlice(value["specifiers"].([]interface{})),
			Name:                       utils.String(value["name"].(string)),
		})
	}

	return &dataSources
}

func expandWindowsEventLogDataSources(p []interface{}) *[]datacollection.WindowsEventLogDataSource {
	if len(p) == 0 {
		return nil
	}
	dataSources := make([]datacollection.WindowsEventLogDataSource, 0)

	for _, v := range p {
		value := v.(map[string]interface{})
		streams := make([]datacollection.KnownWindowsEventLogDataSourceStreams, 0)
		if value["streams"] != nil {
			for _, streamRaw := range value["streams"].([]interface{}) {
				streams = append(streams, datacollection.KnownWindowsEventLogDataSourceStreams(streamRaw.(string)))
			}
		}
		dataSources = append(dataSources, datacollection.WindowsEventLogDataSource{
			Streams:      &streams,
			XPathQueries: utils.ExpandStringSlice(value["xpath_queries"].([]interface{})),
			Name:         utils.String(value["name"].(string)),
		})
	}

	return &dataSources
}

func expandSyslogDataSources(p []interface{}) *[]datacollection.SyslogDataSource {
	if len(p) == 0 {
		return nil
	}
	dataSources := make([]datacollection.SyslogDataSource, 0)

	for _, v := range p {
		value := v.(map[string]interface{})
		streams := make([]datacollection.KnownSyslogDataSourceStreams, 0)
		if value["streams"] != nil {
			for _, streamRaw := range value["streams"].([]interface{}) {
				streams = append(streams, datacollection.KnownSyslogDataSourceStreams(streamRaw.(string)))
			}
		}
		logLevels := make([]datacollection.KnownSyslogDataSourceLogLevels, 0)
		if value["log_levels"] != nil {
			for _, streamRaw := range value["log_levels"].([]interface{}) {
				logLevels = append(logLevels, datacollection.KnownSyslogDataSourceLogLevels(streamRaw.(string)))
			}
		}
		facilityNames := make([]datacollection.KnownSyslogDataSourceFacilityNames, 0)
		if value["facility_names"] != nil {
			for _, streamRaw := range value["facility_names"].([]interface{}) {
				facilityNames = append(facilityNames, datacollection.KnownSyslogDataSourceFacilityNames(streamRaw.(string)))
			}
		}
		dataSources = append(dataSources, datacollection.SyslogDataSource{
			Streams:       &streams,
			LogLevels:     &logLevels,
			FacilityNames: &facilityNames,
			Name:          utils.String(value["name"].(string)),
		})
	}

	return &dataSources
}

func expandExtensionDataSources(p []interface{}) (*[]datacollection.ExtensionDataSource, error) {
	if len(p) == 0 {
		return nil, nil
	}
	dataSources := make([]datacollection.ExtensionDataSource, 0)

	for _, v := range p {
		value := v.(map[string]interface{})
		streams := make([]datacollection.KnownExtensionDataSourceStreams, 0)
		if value["streams"] != nil {
			for _, streamRaw := range value["streams"].([]interface{}) {
				streams = append(streams, datacollection.KnownExtensionDataSourceStreams(streamRaw.(string)))
			}
		}
		dataSource := datacollection.ExtensionDataSource{
			Streams:          &streams,
			InputDataSources: utils.ExpandStringSlice(value["input_data_sources"].([]interface{})),
			ExtensionName:    utils.String(value["extension_name"].(string)),
			Name:             utils.String(value["name"].(string)),
		}
		if value["extension_setting"] != nil {
			extensionSettings, err := pluginsdk.ExpandJsonFromString(value["extension_setting"].(string))
			if err != nil {
				return nil, err
			}
			dataSource.ExtensionSettings = extensionSettings
		}
		dataSources = append(dataSources, dataSource)
	}

	return &dataSources, nil
}

func expandLogAnalyticsDestinations(p []interface{}) *[]datacollection.LogAnalyticsDestination {
	if len(p) == 0 {
		return nil
	}
	destinations := make([]datacollection.LogAnalyticsDestination, 0)

	for _, v := range p {
		value := v.(map[string]interface{})
		destinations = append(destinations, datacollection.LogAnalyticsDestination{
			Name:                utils.String(value["name"].(string)),
			WorkspaceResourceID: utils.String(value["workspace_resource_id"].(string)),
		})
	}

	return &destinations
}

func expandAzureMonitorMetricsDestinations(p []interface{}) *datacollection.DestinationsSpecAzureMonitorMetrics {
	if len(p) == 0 {
		return nil
	}
	value := p[0].(map[string]interface{})
	return &datacollection.DestinationsSpecAzureMonitorMetrics{
		Name: utils.String(value["name"].(string)),
	}
}

func expandDataFlows(p []interface{}) *[]datacollection.DataFlow {
	if len(p) == 0 {
		return nil
	}
	dataFlows := make([]datacollection.DataFlow, 0)

	for _, v := range p {
		value := v.(map[string]interface{})
		streams := make([]datacollection.KnownDataFlowStreams, 0)
		if value["streams"] != nil {
			for _, streamRaw := range value["streams"].([]interface{}) {
				streams = append(streams, datacollection.KnownDataFlowStreams(streamRaw.(string)))
			}
		}
		dataFlows = append(dataFlows, datacollection.DataFlow{
			Streams:      &streams,
			Destinations: utils.ExpandStringSlice(value["destinations"].([]interface{})),
		})
	}

	return &dataFlows
}

func flattenPerformanceCounterDataSources(input *[]datacollection.PerfCounterDataSource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		dataSource := make(map[string]interface{})
		if v.Streams != nil {
			streams := make([]string, 0)
			for _, stream := range *v.Streams {
				streams = append(streams, string(stream))
			}
			dataSource["streams"] = streams
		}
		if v.SamplingFrequencyInSeconds != nil {
			dataSource["sampling_frequency"] = *v.SamplingFrequencyInSeconds
		}
		if v.CounterSpecifiers != nil {
			dataSource["specifiers"] = utils.FlattenStringSlice(v.CounterSpecifiers)
		}
		if v.Name != nil {
			dataSource["name"] = *v.Name
		}
		results = append(results, dataSource)
	}

	return results
}

func flattenWindowsEventLogDataSources(input *[]datacollection.WindowsEventLogDataSource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		dataSource := make(map[string]interface{})
		if v.Streams != nil {
			streams := make([]string, 0)
			for _, stream := range *v.Streams {
				streams = append(streams, string(stream))
			}
			dataSource["streams"] = streams
		}
		if v.XPathQueries != nil {
			dataSource["xpath_queries"] = utils.FlattenStringSlice(v.XPathQueries)
		}
		if v.Name != nil {
			dataSource["name"] = *v.Name
		}
		results = append(results, dataSource)
	}

	return results
}

func flattenSyslogDataSources(input *[]datacollection.SyslogDataSource) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		dataSource := make(map[string]interface{})
		if v.Streams != nil {
			streams := make([]string, 0)
			for _, stream := range *v.Streams {
				streams = append(streams, string(stream))
			}
			dataSource["streams"] = streams
		}
		if v.LogLevels != nil {
			logLevels := make([]string, 0)
			for _, logLevel := range *v.LogLevels {
				logLevels = append(logLevels, string(logLevel))
			}
			dataSource["log_levels"] = logLevels
		}
		if v.FacilityNames != nil {
			facilityNames := make([]string, 0)
			for _, facilityName := range *v.FacilityNames {
				facilityNames = append(facilityNames, string(facilityName))
			}
			dataSource["facility_names"] = facilityNames
		}
		if v.Name != nil {
			dataSource["name"] = *v.Name
		}
		results = append(results, dataSource)
	}

	return results
}

func flattenExtensionDataSources(input *[]datacollection.ExtensionDataSource) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		dataSource := make(map[string]interface{})
		if v.Streams != nil {
			streams := make([]string, 0)
			for _, stream := range *v.Streams {
				streams = append(streams, string(stream))
			}
			dataSource["streams"] = streams
		}
		if v.InputDataSources != nil {
			dataSource["input_data_sources"] = utils.FlattenStringSlice(v.InputDataSources)
		}
		if v.ExtensionName != nil {
			dataSource["extension_name"] = *v.ExtensionName
		}
		if v.Name != nil {
			dataSource["name"] = *v.Name
		}
		if v.ExtensionSettings != nil {
			extensionSetting, err := pluginsdk.FlattenJsonToString(v.ExtensionSettings.(map[string]interface{}))
			if err != nil {
				return []interface{}{}, err
			}
			dataSource["extension_setting"] = extensionSetting
		}
		results = append(results, dataSource)
	}

	return results, nil
}

func flattenLogAnalyticsDestinations(input *[]datacollection.LogAnalyticsDestination) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		destination := make(map[string]interface{})
		if v.Name != nil {
			destination["name"] = *v.Name
		}
		if v.WorkspaceResourceID != nil {
			destination["workspace_resource_id"] = *v.WorkspaceResourceID
		}
		results = append(results, destination)
	}

	return results
}

func flattenAzureMonitorMetricsDestinations(input *datacollection.DestinationsSpecAzureMonitorMetrics) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	destination := make(map[string]interface{})
	if input.Name != nil {
		destination["name"] = *input.Name
	}

	return []interface{}{destination}
}

func flattenDataFlows(input *[]datacollection.DataFlow) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	results := make([]interface{}, 0)
	for _, v := range *input {
		dataFlow := make(map[string]interface{})
		if v.Streams != nil {
			streams := make([]string, 0)
			for _, stream := range *v.Streams {
				streams = append(streams, string(stream))
			}
			dataFlow["streams"] = streams
		}
		if v.Destinations != nil {
			dataFlow["destinations"] = utils.FlattenStringSlice(v.Destinations)
		}
		results = append(results, dataFlow)
	}

	return results
}

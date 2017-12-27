package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	DelimSpace       = " "
	DelimComma       = ","
	DelimTab         = "\t"
	DelimSemiColon   = ";"
	DelimVerticalBar = "|"
)

// Allbut returns a list all possible datasource other than the specified arg
func Allbut(source string) (ret []string) {
	sources := []string{
		"blob",
		"table",
		"event_hub",
		"sql_database",
		"service_bus_queues",
		"service_bus_topics",
		"documentdb",
	}
	for _, val := range sources {
		if val != source {
			ret = append(ret, val)
		}
	}

	return

}

var streamAnalyticsOutputType = "Microsoft.StreamAnalytics/streamingjobs/outputs"

func streamAnalyticsOutputSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
		Elem: map[string]*schema.Schema{
			"output_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"serialization": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamanalytics.TypeCsv),
								string(streamanalytics.TypeAvro),
								string(streamanalytics.TypeJSON),
							}, false),
						},
						"field_delimiter": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								DelimTab,
								DelimComma,
								DelimSemiColon,
								DelimVerticalBar,
								DelimSpace,
							}, false),
						},
						"encoding": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(streamanalytics.Array),
								string(streamanalytics.LineSeparated),
							}, false),
						},
					},
				},
			},
			"blob": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
				// ConflictsWith: Allbut("blob"),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_account_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"storage_account_key": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"container": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"path_pattern": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"date_format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"time_format": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"table": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
				// ConflictsWith: Allbut("table"),

				Elem: map[string]*schema.Schema{
					"account_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"account_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"table_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"partition_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"row_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"batch_size": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
						Default:  100,
					},
				},
			},
			"event_hub": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101

				// ConflictsWith: Allbut("event_hub"),
				Elem: map[string]*schema.Schema{

					"namespace": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"shared_access_policy_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"shared_access_policy_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"event_hub_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"partition_key": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			"sql_database": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
				// ConflictsWith: Allbut("sql_database"),
				Elem: map[string]*schema.Schema{
					"server": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"database": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"user": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"password": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"table": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
			"service_bus_queues": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
				// ConflictsWith: Allbut("service_bus_queues"),
				Elem: map[string]*schema.Schema{
					"namespace": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"shared_access_policy_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"shared_access_policy_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"queue_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"property_columns": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			"service_bus_topics": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
				// ConflictsWith: Allbut("service_bus_topics"),
				Elem: map[string]*schema.Schema{
					"namespace": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"shared_access_policy_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"shared_access_policy_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"topic_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"property_columns": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			"documentdb": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
				// ConflictsWith: Allbut("documentdb"),
				Elem: map[string]*schema.Schema{
					"account_id": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"account_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"database": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"collection_name_pattern": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"partition_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"document_id": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
	}
}

func streamAnalyticsOutputFromSchema(outputSchema interface{}) (*streamanalytics.Output, error) {

	outputMap := outputSchema.(map[string]interface{})
	outputName := outputMap["output_name"].(string)
	output := &streamanalytics.Output{
		Name: &outputName,
		Type: &streamAnalyticsOutputType,
	}

	serialization, err := extractOutputSerialization(outputMap)

	if err != nil {
		return nil, err
	}

	datasource, err := extractOutputDatasource(outputMap)

	if err != nil {
		return nil, err
	}
	outputProperties := streamanalytics.OutputProperties{
		Datasource:    datasource,
		Serialization: serialization,
	}

	output.OutputProperties = &outputProperties

	return output, nil
}

func extractOutputDatasource(outputMap map[string]interface{}) (streamanalytics.OutputDataSource, error) {

	var outputDatasource streamanalytics.OutputDataSource

	if blobSchema := outputMap["blob"].([]interface{}); len(blobSchema) != 0 {
		blobMap := blobSchema[0].(map[string]interface{})
		storageAccountName := blobMap["storage_account_name"].(string)
		storageAccountKey := blobMap["storage_account_key"].(string)
		container := blobMap["container"].(string)
		pathPattern := blobMap["path_pattern"].(string)

		sAccounts := []streamanalytics.StorageAccount{
			{
				AccountKey:  &storageAccountKey,
				AccountName: &storageAccountName,
			},
		}
		outputDatasourceProps := streamanalytics.BlobOutputDataSourceProperties{
			StorageAccounts: &sAccounts,
			Container:       &container,
			PathPattern:     &pathPattern,
		}

		if dateFormat, ok := blobMap["date_format"].(string); ok && dateFormat != "" {
			outputDatasourceProps.DateFormat = &dateFormat
		}
		if timeFormat, ok := blobMap["time_format"].(string); ok && timeFormat != "" {
			outputDatasourceProps.TimeFormat = &timeFormat
		}

		outputDatasource = streamanalytics.BlobOutputDataSource{
			Type: streamanalytics.TypeMicrosoftStorageBlob,
			BlobOutputDataSourceProperties: &outputDatasourceProps,
		}

	} else if tableSchema := outputMap["table"].([]interface{}); len(tableSchema) != 0 {
		tableMap := tableSchema[0].(map[string]interface{})
		accountName := tableMap["account_name"].(string)
		accountKey := tableMap["account_key"].(string)
		tableName := tableMap["table_name"].(string)
		partitionKey := tableMap["partition_key"].(string)
		rowKey := tableMap["row_key"].(string)
		batchSize := tableMap["batch_size"].(int)
		batchSize32 := int32(batchSize)

		tableProps := streamanalytics.AzureTableOutputDataSourceProperties{
			AccountName:  &accountName,
			AccountKey:   &accountKey,
			Table:        &tableName,
			PartitionKey: &partitionKey,
			RowKey:       &rowKey,
			BatchSize:    &batchSize32,
		}

		outputDatasource = streamanalytics.AzureTableOutputDataSource{
			Type: streamanalytics.TypeMicrosoftStorageTable,
			AzureTableOutputDataSourceProperties: &tableProps,
		}

	} else if eventhubSchema := outputMap["event_hub"].([]interface{}); len(eventhubSchema) != 0 {
		eventhubMap := eventhubSchema[0].(map[string]interface{})
		namespace := eventhubMap["namespace"].(string)
		sharedPolicyName := eventhubMap["shared_access_policy_name"].(string)
		sharedPolicyKey := eventhubMap["shared_access_policy_key"].(string)
		eventHubName := eventhubMap["event_hub_name"].(string)

		eventhubProps := streamanalytics.EventHubOutputDataSourceProperties{
			ServiceBusNamespace:    &namespace,
			SharedAccessPolicyKey:  &sharedPolicyKey,
			SharedAccessPolicyName: &sharedPolicyName,
			EventHubName:           &eventHubName,
		}

		if partitionKey, ok := eventhubMap["partition_key"].(string); ok && partitionKey != "" {
			eventhubProps.PartitionKey = &partitionKey
		}

		outputDatasource = streamanalytics.EventHubOutputDataSource{
			Type: streamanalytics.TypeMicrosoftServiceBusEventHub,
			EventHubOutputDataSourceProperties: &eventhubProps,
		}

	} else if SqlSchema := outputMap["sql_database"].([]interface{}); len(SqlSchema) != 0 {
		SqlMap := SqlSchema[0].(map[string]interface{})

	} else if blobSchema := outputMap["blob"].([]interface{}); len(blobSchema) != 0 {
		blobMap := blobSchema[0].(map[string]interface{})
	} else if blobSchema := outputMap["blob"].([]interface{}); len(blobSchema) != 0 {
		blobMap := blobSchema[0].(map[string]interface{})
	} else if blobSchema := outputMap["blob"].([]interface{}); len(blobSchema) != 0 {
		blobMap := blobSchema[0].(map[string]interface{})
	}

}
func extractOutputSerialization(outputMap map[string]interface{}) (streamanalytics.Serialization, error) {

}

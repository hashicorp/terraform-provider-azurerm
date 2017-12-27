package azurerm

import (
	"errors"
	"log"

	"github.com/Azure/azure-sdk-for-go/arm/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var streamAnalyticsInputType = "Microsoft.StreamAnalytics/streamingjobs/inputs"

func streamAnalyticsInputSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"input_name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(streamanalytics.TypeReference),
						string(streamanalytics.TypeStream),
					}, false),
				},
				"serialization": &schema.Schema{
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 1,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
								ValidateFunc: validation.StringInSlice([]string{
									string(streamanalytics.TypeAvro),
									string(streamanalytics.TypeCsv),
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
								Default:  string(streamanalytics.UTF8),
								ValidateFunc: validation.StringInSlice([]string{
									string(streamanalytics.UTF8),
								}, false),
							},
						},
					},
				},
				"datasource": &schema.Schema{
					Type:     schema.TypeList,
					MaxItems: 1,
					Required: true,
					MinItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"blob": &schema.Schema{
								Type:     schema.TypeList,
								MaxItems: 1,
								Optional: true,
								// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
								// ConflictsWith: []string{
								// 	"event_hub",
								// 	"iot_hub",
								// },
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"storage_account_name": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										// This is being sent as null so for now keeping this as optional
										"storage_account_key": &schema.Schema{
											Type:     schema.TypeString,
											Optional: true,
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
										"source_partition_count": &schema.Schema{
											Type:     schema.TypeInt,
											Optional: true,
										},
									},
								},
							},
							"event_hub": &schema.Schema{
								Type:     schema.TypeList,
								MaxItems: 1,
								Optional: true,
								// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
								// ConflictsWith: []string{
								// 	"blob",
								// 	"iot_hub",
								// },
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
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
										"consumer_group_name": &schema.Schema{
											Type:     schema.TypeString,
											Optional: true,
										},
									},
								},
							},
							"iot_hub": &schema.Schema{
								Type:     schema.TypeList,
								MaxItems: 1,
								Optional: true,
								// ConfictWith doesnt work on list https://github.com/hashicorp/terraform/issues/11101
								// ConflictsWith: []string{
								// 	"blob",
								// 	"event_hub",
								// },
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
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
										"consumer_group_name": &schema.Schema{
											Type:     schema.TypeString,
											Optional: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func extractReferenceDataSource(dataMap map[string]interface{}) (streamanalytics.ReferenceInputDataSource, error) {

	datasourceList := dataMap["datasource"].([]interface{})
	datasourceMap := datasourceList[0].(map[string]interface{})

	blobsourceList, _ := datasourceMap["blob"].([]interface{})
	if blobsourceList == nil {
		return nil, errors.New("Blob Datasource returned empty")
	}
	blobsourceMap := blobsourceList[0].(map[string]interface{})

	container := blobsourceMap["container"].(string)
	pathPattern := blobsourceMap["path_pattern"].(string)
	storageAccountName := blobsourceMap["storage_account_name"].(string)
	storageAccountKey := blobsourceMap["storage_account_key"].(string)

	var storageAccounts []streamanalytics.StorageAccount
	storageAccounts = append(storageAccounts, streamanalytics.StorageAccount{
		AccountName: &storageAccountName,
		AccountKey:  &storageAccountKey,
	})

	datasourceProperties := &streamanalytics.BlobReferenceInputDataSourceProperties{
		Container:       &container,
		PathPattern:     &pathPattern,
		StorageAccounts: &storageAccounts,
	}

	if dateFormat, ok := blobsourceMap["date_format"].(string); ok && dateFormat != "" {
		datasourceProperties.DateFormat = &dateFormat
	}

	if timeFormat, ok := blobsourceMap["time_format"].(string); ok && timeFormat != "" {
		datasourceProperties.TimeFormat = &timeFormat
	}

	datasource := streamanalytics.BlobReferenceInputDataSource{
		Type: streamanalytics.TypeReferenceInputDataSourceTypeMicrosoftStorageBlob,
		BlobReferenceInputDataSourceProperties: datasourceProperties,
	}

	return datasource, nil
}

func extractStreamDataSource(dataMap map[string]interface{}) (streamanalytics.StreamInputDataSource, error) {
	datasourceList := dataMap["datasource"].([]interface{})
	datasourceMap := datasourceList[0].(map[string]interface{})

	var streamInputSource streamanalytics.StreamInputDataSource

	if blobsourceList, ok := datasourceMap["blob"].([]interface{}); ok && len(blobsourceList) != 0 {

		blobsourceMap := blobsourceList[0].(map[string]interface{})
		container := blobsourceMap["container"].(string)
		pathPattern := blobsourceMap["path_pattern"].(string)
		storageAccountName := blobsourceMap["storage_account_name"].(string)
		storageAccountKey := blobsourceMap["storage_account_key"].(string)

		var storageAccounts []streamanalytics.StorageAccount
		storageAccounts = append(storageAccounts, streamanalytics.StorageAccount{
			AccountName: &storageAccountName,
			AccountKey:  &storageAccountKey,
		})

		datasourceProperties := &streamanalytics.BlobStreamInputDataSourceProperties{
			Container:       &container,
			PathPattern:     &pathPattern,
			StorageAccounts: &storageAccounts,
		}

		if dateFormat, ok := blobsourceMap["date_format"].(string); ok && dateFormat != "" {
			datasourceProperties.DateFormat = &dateFormat
		}

		if timeFormat, ok := blobsourceMap["time_format"].(string); ok && timeFormat != "" {
			datasourceProperties.TimeFormat = &timeFormat
		}

		if sourcePartionCount, ok := blobsourceMap["source_partition_count"].(int); ok {
			sourceCount32 := int32(sourcePartionCount)
			datasourceProperties.SourcePartitionCount = &sourceCount32
		}
		streamInputSource = streamanalytics.BlobStreamInputDataSource{
			Type: streamanalytics.TypeStreamInputDataSourceTypeMicrosoftStorageBlob,
			BlobStreamInputDataSourceProperties: datasourceProperties,
		}

	} else if eventhubList, ok := datasourceMap["event_hub"].([]interface{}); ok && len(eventhubList) != 0 {
		eventhubMap := eventhubList[0].(map[string]interface{})

		namespace := eventhubMap["namespace"].(string)
		sharedPolicyName := eventhubMap["shared_access_policy_name"].(string)
		sharedPolicyKey := eventhubMap["shared_access_policy_key"].(string)
		eventHubName := eventhubMap["event_hub_name"].(string)
		log.Printf("[INFO] Creating Eventhub input source using alias %s", eventHubName)

		eventhubStreamProps := streamanalytics.EventHubStreamInputDataSourceProperties{
			ServiceBusNamespace:    &namespace,
			SharedAccessPolicyName: &sharedPolicyName,
			SharedAccessPolicyKey:  &sharedPolicyKey,
			EventHubName:           &eventHubName,
		}

		if consumerGroup, ok := eventhubMap["consumer_group_name"].(string); ok && consumerGroup != "" {
			eventhubStreamProps.ConsumerGroupName = &consumerGroup
		}

		streamInputSource = streamanalytics.EventHubStreamInputDataSource{
			Type: streamanalytics.TypeStreamInputDataSourceTypeMicrosoftServiceBusEventHub,
			EventHubStreamInputDataSourceProperties: &eventhubStreamProps,
		}
	} else if iothubList, ok := datasourceMap["iot_hub"].([]interface{}); ok && len(iothubList) != 0 {
		iothubMap := iothubList[0].(map[string]interface{})
		namespace := iothubMap["namespace"].(string)
		sharedPolicyName := iothubMap["shared_access_policy_name"].(string)
		sharedPolicyKey := iothubMap["shared_access_policy_key"].(string)

		iothubStreamProps := streamanalytics.IoTHubStreamInputDataSourceProperties{
			IotHubNamespace:        &namespace,
			SharedAccessPolicyName: &sharedPolicyName,
			SharedAccessPolicyKey:  &sharedPolicyKey,
		}

		if consumerGroup, ok := iothubMap["consumer_group_name"].(string); ok && consumerGroup != "" {
			iothubStreamProps.ConsumerGroupName = &consumerGroup
		}

		streamInputSource = streamanalytics.IoTHubStreamInputDataSource{
			Type: streamanalytics.TypeStreamInputDataSourceTypeMicrosoftDevicesIotHubs,
			IoTHubStreamInputDataSourceProperties: &iothubStreamProps,
		}

	} else {
		// just to keep the logical structure conventional
		return nil, errors.New("All datasources are empty")
	}

	return streamInputSource, nil

}

func extractSerialization(dataMap map[string]interface{}, output bool) streamanalytics.Serialization {

	var serialization streamanalytics.Serialization
	serializationList := dataMap["serialization"].([]interface{})
	serial := serializationList[0].(map[string]interface{})
	serialType := serial["type"].(string)

	switch serialType {
	case string(streamanalytics.TypeAvro):
		serialization = streamanalytics.AvroSerialization{
			Type: streamanalytics.TypeAvro,
		}

	case string(streamanalytics.TypeCsv):

		// TODO: check for optional parameters
		tfd := serial["field_delimiter"].(string)
		serialization = streamanalytics.CsvSerialization{
			Type: streamanalytics.TypeCsv,
			CsvSerializationProperties: &streamanalytics.CsvSerializationProperties{
				FieldDelimiter: &tfd,
				Encoding:       streamanalytics.Encoding(serial["encoding"].(string)),
			},
		}

	case string(streamanalytics.TypeJSON):

		format := streamanalytics.LineSeparated
		if output {
			if formatStr, ok := serial["format"].(string); ok && formatStr != "" {
				if formatStr == string(streamanalytics.Array) {
					format = streamanalytics.Array
				}
			}
		}

		serialization = streamanalytics.JSONSerialization{
			Type: streamanalytics.TypeJSON,
			JSONSerializationProperties: &streamanalytics.JSONSerializationProperties{
				Encoding: streamanalytics.Encoding(serial["encoding"].(string)),
				Format:   format,
			},
		}
	}

	return serialization
}

func streamAnalyticsInputfromSchema(data interface{}) (*streamanalytics.Input, error) {
	dataMap := data.(map[string]interface{})
	inputType := dataMap["type"].(string)
	inputName := dataMap["input_name"].(string)

	input := &streamanalytics.Input{
		Name: &inputName,
		Type: &streamAnalyticsInputType,
	}

	serialization := extractSerialization(dataMap, false)

	switch inputType {
	case string(streamanalytics.TypeReference):
		log.Println("[INFO] using Reference Type")

		datasource, err := extractReferenceDataSource(dataMap)

		if err != nil {
			return nil, err
		}
		inputProperties := streamanalytics.ReferenceInputProperties{
			Serialization: serialization,
			Type:          streamanalytics.TypeReference,
			Datasource:    datasource,
		}

		input.Properties = inputProperties

	case string(streamanalytics.TypeStream):
		log.Println("[INFO] using Stream Type")

		datasource, err := extractStreamDataSource(dataMap)

		if err != nil {
			return nil, err
		}
		inputProperties := streamanalytics.StreamInputProperties{
			Serialization: serialization,
			Type:          streamanalytics.TypeStream,
			Datasource:    datasource,
		}

		input.Properties = inputProperties
	default:
		return nil, errors.New("The input type not supported")

	}

	return input, nil
}

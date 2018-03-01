package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/services/streamanalytics/mgmt/2016-03-01/streamanalytics"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

var streamAnalyticsFunctionType = "Microsoft.StreamAnalytics/streamingjobs/functions"

func streamAnalyticsFunctionSchema() *schema.Schema {

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"function_name": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"input_types": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Required: true,
				},
				"output_type": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"user_defined": &schema.Schema{
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"script": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"machine_learning": &schema.Schema{
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"endpoint": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
							"api_key": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
							"batch_size": &schema.Schema{
								Type:         schema.TypeInt,
								Optional:     true,
								Default:      10,
								ValidateFunc: validation.IntBetween(10, 1000),
							},
							"input_name": &schema.Schema{
								Type:     schema.TypeString,
								Required: true,
							},
							"column_names": &schema.Schema{
								Type:     schema.TypeList,
								Required: true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"name": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"data_type": &schema.Schema{
											Type:     schema.TypeString,
											Required: true,
										},
										"map_to": &schema.Schema{
											Type:     schema.TypeInt,
											Required: true,
										},
									},
								},
							},
							"output_mapping": &schema.Schema{
								Type:     schema.TypeMap,
								Required: true,
							},
						},
					},
				},
			},
		},
	}

}

func extractColumnName(columnSchema interface{}) streamanalytics.AzureMachineLearningWebServiceInputColumn {
	columnMap := columnSchema.(map[string]interface{})
	name := columnMap["name"].(string)
	dataType := columnMap["data_type"].(string)
	mapTo := columnMap["map_to"].(int)
	mapTo32 := int32(mapTo)

	return streamanalytics.AzureMachineLearningWebServiceInputColumn{
		Name:     &name,
		DataType: &dataType,
		MapTo:    &mapTo32,
	}
}
func extractFunctionBinding(funcMap map[string]interface{}) streamanalytics.FunctionBinding {

	var functionBinding streamanalytics.FunctionBinding

	if udfSchema, ok := funcMap["user_defined"].([]interface{}); ok && len(udfSchema) != 0 {
		udfMap := udfSchema[0].(map[string]interface{})
		script := udfMap["script"].(string)
		functionBinding = streamanalytics.JavaScriptFunctionBinding{
			Type: streamanalytics.TypeMicrosoftStreamAnalyticsJavascriptUdf,
			JavaScriptFunctionBindingProperties: &streamanalytics.JavaScriptFunctionBindingProperties{
				Script: &script,
			},
		}
	} else if mlSchema, ok := funcMap["machine_learning"].([]interface{}); ok && len(mlSchema) != 0 {
		mlMap := mlSchema[0].(map[string]interface{})
		endpoint := mlMap["endpoint"].(string)
		batchSize := mlMap["batch_size"].(int)
		batchSize32 := int32(batchSize)
		apiKey := mlMap["api_key"].(string)
		inputName := mlMap["input_name"].(string)
		columnNamesList := mlMap["column_names"].([]interface{})
		outputMappingMap := mlMap["output_mapping"].(map[string]interface{})

		var columnNames []streamanalytics.AzureMachineLearningWebServiceInputColumn
		for _, columnNameIntf := range columnNamesList {
			columnNames = append(columnNames, extractColumnName(columnNameIntf))
		}
		var outputMapping []streamanalytics.AzureMachineLearningWebServiceOutputColumn
		for name, typeIntf := range outputMappingMap {
			dataType := typeIntf.(string)
			outputColumn := streamanalytics.AzureMachineLearningWebServiceOutputColumn{
				Name:     &name,
				DataType: &dataType,
			}

			outputMapping = append(outputMapping, outputColumn)
		}

		functionBinding = streamanalytics.AzureMachineLearningWebServiceFunctionBinding{
			Type: streamanalytics.TypeMicrosoftMachineLearningWebService,
			AzureMachineLearningWebServiceFunctionBindingProperties: &streamanalytics.AzureMachineLearningWebServiceFunctionBindingProperties{
				Endpoint:  &endpoint,
				BatchSize: &batchSize32,
				APIKey:    &apiKey,
				Inputs: &streamanalytics.AzureMachineLearningWebServiceInputs{
					Name:        &inputName,
					ColumnNames: &columnNames,
				},
				Outputs: &outputMapping,
			},
		}
	}

	return functionBinding
}

func streamAnalyticsFunctionFromSchema(funcSchema interface{}) streamanalytics.Function {

	funcMap := funcSchema.(map[string]interface{})
	outputType := funcMap["output_type"].(string)
	funcName := funcMap["function_name"].(string)
	inputTypeList := funcMap["input_types"].([]interface{})

	var inputTypes []streamanalytics.FunctionInput

	for _, inputTypeIntf := range inputTypeList {
		inputTypeStr := inputTypeIntf.(string)
		inputTypes = append(inputTypes, streamanalytics.FunctionInput{
			DataType: &inputTypeStr,
		})
	}

	funcBinding := extractFunctionBinding(funcMap)

	funcConf := streamanalytics.ScalarFunctionConfiguration{
		Inputs: &inputTypes,
		Output: &streamanalytics.FunctionOutput{
			DataType: &outputType,
		},
		Binding: funcBinding,
	}

	function := streamanalytics.Function{
		Name: &funcName,
		Type: &streamAnalyticsFunctionType,
		Properties: &streamanalytics.ScalarFunctionProperties{
			Type: streamanalytics.TypeScalar,
			ScalarFunctionConfiguration: &funcConf,
		},
	}
	return function

}

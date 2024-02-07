// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func schemaStreamAnalyticsStreamInputSerialization() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"type": {
					Type:     pluginsdk.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(inputs.EventSerializationTypeAvro),
						string(inputs.EventSerializationTypeCsv),
						string(inputs.EventSerializationTypeJson),
					}, false),
				},

				"field_delimiter": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						" ",
						",",
						"	",
						"|",
						";",
					}, false),
				},

				"encoding": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(inputs.EncodingUTFEight),
					}, false),
				},
			},
		},
	}
}

func expandStreamAnalyticsStreamInputSerialization(input []interface{}) (inputs.Serialization, error) {
	v := input[0].(map[string]interface{})

	inputType := v["type"].(string)
	encoding := v["encoding"].(string)
	fieldDelimiter := v["field_delimiter"].(string)

	var props interface{}

	switch inputType {
	case string(inputs.EventSerializationTypeAvro):
		return inputs.AvroSerialization{
			Properties: &props,
		}, nil

	case string(inputs.EventSerializationTypeCsv):
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `type` is set to `Csv`")
		}
		if fieldDelimiter == "" {
			return nil, fmt.Errorf("`field_delimiter` must be set when `type` is set to `Csv`")
		}
		return inputs.CsvSerialization{
			Properties: &inputs.CsvSerializationProperties{
				Encoding:       pointer.To(inputs.Encoding(encoding)),
				FieldDelimiter: pointer.To(fieldDelimiter),
			},
		}, nil

	case string(inputs.EventSerializationTypeJson):
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `type` is set to `Json`")
		}

		return inputs.JsonSerialization{
			Properties: &inputs.JsonSerializationProperties{
				Encoding: pointer.To(inputs.Encoding(encoding)),
			},
		}, nil
	}

	return nil, fmt.Errorf("Unsupported Input Type %q", inputType)
}

func expandStreamAnalyticsStreamInputSerializationTyped(serialization []Serialization) (inputs.Serialization, error) {
	v := serialization[0]

	inputType := v.Type
	encoding := v.Encoding
	fieldDelimiter := v.FieldDelimiter

	var props interface{}

	switch inputType {
	case string(inputs.EventSerializationTypeAvro):
		return inputs.AvroSerialization{
			Properties: &props,
		}, nil

	case string(inputs.EventSerializationTypeCsv):
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `type` is set to `Csv`")
		}
		if fieldDelimiter == "" {
			return nil, fmt.Errorf("`field_delimiter` must be set when `type` is set to `Csv`")
		}
		return inputs.CsvSerialization{
			Properties: &inputs.CsvSerializationProperties{
				Encoding:       pointer.To(inputs.Encoding(encoding)),
				FieldDelimiter: pointer.To(fieldDelimiter),
			},
		}, nil

	case string(inputs.EventSerializationTypeJson):
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `type` is set to `Json`")
		}

		return inputs.JsonSerialization{
			Properties: &inputs.JsonSerializationProperties{
				Encoding: pointer.To(inputs.Encoding(encoding)),
			},
		}, nil
	}

	return nil, fmt.Errorf("Unsupported Input Type %q", inputType)
}

func flattenStreamAnalyticsStreamInputSerialization(input inputs.Serialization) []interface{} {
	var encoding string
	var fieldDelimiter string
	var inputType string

	if _, ok := input.(inputs.AvroSerialization); ok {
		inputType = string(inputs.EventSerializationTypeAvro)
	}

	if csv, ok := input.(inputs.CsvSerialization); ok {
		if props := csv.Properties; props != nil {
			if v := props.Encoding; v != nil {
				encoding = string(*v)
			}

			if v := props.FieldDelimiter; v != nil {
				fieldDelimiter = *v
			}
		}

		inputType = string(inputs.EventSerializationTypeCsv)
	}

	if json, ok := input.(inputs.JsonSerialization); ok {
		if props := json.Properties; props != nil {
			if v := props.Encoding; v != nil {
				encoding = string(*v)
			}
		}

		inputType = string(inputs.EventSerializationTypeJson)
	}

	return []interface{}{
		map[string]interface{}{
			"encoding":        encoding,
			"type":            inputType,
			"field_delimiter": fieldDelimiter,
		},
	}
}

func flattenStreamAnalyticsStreamInputSerializationTyped(input inputs.Serialization) Serialization {
	var encoding string
	var fieldDelimiter string
	var inputType string

	if _, ok := input.(inputs.AvroSerialization); ok {
		inputType = string(inputs.EventSerializationTypeAvro)
	}

	if csv, ok := input.(inputs.CsvSerialization); ok {
		if props := csv.Properties; props != nil {
			if v := props.Encoding; v != nil {
				encoding = string(*v)
			}

			if v := props.FieldDelimiter; v != nil {
				fieldDelimiter = *v
			}
		}

		inputType = string(inputs.EventSerializationTypeCsv)
	}

	if json, ok := input.(inputs.JsonSerialization); ok {
		if props := json.Properties; props != nil {
			if v := props.Encoding; v != nil {
				encoding = string(*v)
			}
		}

		inputType = string(inputs.EventSerializationTypeJson)
	}

	return Serialization{
		Encoding:       encoding,
		Type:           inputType,
		FieldDelimiter: fieldDelimiter,
	}
}

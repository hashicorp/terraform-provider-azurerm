// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package streamanalytics

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func schemaStreamAnalyticsOutputSerialization() *pluginsdk.Schema {
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
						string(outputs.EventSerializationTypeAvro),
						string(outputs.EventSerializationTypeCsv),
						string(outputs.EventSerializationTypeJson),
						string(outputs.EventSerializationTypeParquet),
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
						string(outputs.EncodingUTFEight),
					}, false),
				},

				"format": {
					Type:     pluginsdk.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(outputs.JsonOutputSerializationFormatArray),
						string(outputs.JsonOutputSerializationFormatLineSeparated),
					}, false),
				},
			},
		},
	}
}

func expandStreamAnalyticsOutputSerialization(input []interface{}) (outputs.Serialization, error) {
	v := input[0].(map[string]interface{})

	outputType := v["type"].(string)
	encoding := v["encoding"].(string)
	fieldDelimiter := v["field_delimiter"].(string)
	format := v["format"].(string)

	switch outputType {
	case string(outputs.EventSerializationTypeAvro):
		if encoding != "" {
			return nil, fmt.Errorf("`encoding` cannot be set when `type` is set to `Avro`")
		}
		if fieldDelimiter != "" {
			return nil, fmt.Errorf("`field_delimiter` cannot be set when `type` is set to `Avro`")
		}
		if format != "" {
			return nil, fmt.Errorf("`format` cannot be set when `type` is set to `Avro`")
		}
		var props interface{}
		return outputs.AvroSerialization{
			Properties: &props,
		}, nil

	case string(outputs.EventSerializationTypeCsv):
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `type` is set to `Csv`")
		}
		if fieldDelimiter == "" {
			return nil, fmt.Errorf("`field_delimiter` must be set when `type` is set to `Csv`")
		}
		if format != "" {
			return nil, fmt.Errorf("`format` cannot be set when `type` is set to `Csv`")
		}
		return outputs.CsvSerialization{
			Properties: &outputs.CsvSerializationProperties{
				Encoding:       pointer.To(outputs.Encoding(encoding)),
				FieldDelimiter: pointer.To(fieldDelimiter),
			},
		}, nil

	case string(outputs.EventSerializationTypeJson):
		if encoding == "" {
			return nil, fmt.Errorf("`encoding` must be specified when `type` is set to `Json`")
		}
		if format == "" {
			return nil, fmt.Errorf("`format` must be specified when `type` is set to `Json`")
		}
		if fieldDelimiter != "" {
			return nil, fmt.Errorf("`field_delimiter` cannot be set when `type` is set to `Json`")
		}

		return outputs.JsonSerialization{
			Properties: &outputs.JsonSerializationProperties{
				Encoding: pointer.To(outputs.Encoding(encoding)),
				Format:   pointer.To(outputs.JsonOutputSerializationFormat(format)),
			},
		}, nil

	case string(outputs.EventSerializationTypeParquet):
		if encoding != "" {
			return nil, fmt.Errorf("`encoding` cannot be set when `type` is set to `Parquet`")
		}
		if fieldDelimiter != "" {
			return nil, fmt.Errorf("`field_delimiter` cannot be set when `type` is set to `Parquet`")
		}
		if format != "" {
			return nil, fmt.Errorf("`format` cannot be set when `type` is set to `Parquet`")
		}

		var props interface{}
		return outputs.ParquetSerialization{
			Properties: &props,
		}, nil
	}

	return nil, fmt.Errorf("Unsupported Output Type %q", outputType)
}

func flattenStreamAnalyticsOutputSerialization(input outputs.Serialization) []interface{} {
	var encoding string
	var outputType string
	var fieldDelimiter string
	var format string

	if _, ok := input.(outputs.AvroSerialization); ok {
		outputType = string(outputs.EventSerializationTypeAvro)
	}

	if csv, ok := input.(outputs.CsvSerialization); ok {
		if props := csv.Properties; props != nil {
			if props.Encoding != nil {
				encoding = string(*props.Encoding)
			}

			if props.FieldDelimiter != nil {
				fieldDelimiter = *props.FieldDelimiter
			}
		}

		outputType = string(outputs.EventSerializationTypeCsv)
	}

	if json, ok := input.(outputs.JsonSerialization); ok {
		if props := json.Properties; props != nil {
			if props.Encoding != nil {
				encoding = string(*props.Encoding)
			}
			if props.Format != nil {
				format = string(*props.Format)
			}
		}

		outputType = string(outputs.EventSerializationTypeJson)
	}

	if _, ok := input.(outputs.ParquetSerialization); ok {
		outputType = string(outputs.EventSerializationTypeParquet)
	}

	return []interface{}{
		map[string]interface{}{
			"encoding":        encoding,
			"type":            outputType,
			"format":          format,
			"field_delimiter": fieldDelimiter,
		},
	}
}

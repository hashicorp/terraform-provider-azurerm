// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_stream_input_eventhub_v2

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
)

func expandStreamAnalyticsStreamInputSerializationTyped(serialization []Serialization) (inputs.Serialization, error) {
	v := serialization[0]

	inputType := v.Type
	encoding := v.Encoding
	fieldDelimiter := v.FieldDelimiter

	var props any

	switch inputType {
	case string(inputs.EventSerializationTypeAvro):
		return inputs.AvroSerialization{
			Properties: &props,
		}, nil

	case string(inputs.EventSerializationTypeCsv):
		if encoding == "" {
			return nil, errors.New("`encoding` must be specified when `type` is set to `Csv`")
		}
		if fieldDelimiter == "" {
			return nil, errors.New("`field_delimiter` must be set when `type` is set to `Csv`")
		}
		return inputs.CsvSerialization{
			Properties: &inputs.CsvSerializationProperties{
				Encoding:       pointer.ToEnum[inputs.Encoding](encoding),
				FieldDelimiter: pointer.To(fieldDelimiter),
			},
		}, nil

	case string(inputs.EventSerializationTypeJson):
		if encoding == "" {
			return nil, errors.New("`encoding` must be specified when `type` is set to `Json`")
		}

		return inputs.JsonSerialization{
			Properties: &inputs.JsonSerializationProperties{
				Encoding: pointer.ToEnum[inputs.Encoding](encoding),
			},
		}, nil
	}

	return nil, fmt.Errorf("unsupported Input Type %q", inputType)
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

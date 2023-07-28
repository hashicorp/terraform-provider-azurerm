package schedule

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoMLVertical interface {
}

func unmarshalAutoMLVerticalImplementation(input []byte) (AutoMLVertical, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling AutoMLVertical into map[string]interface: %+v", err)
	}

	value, ok := temp["taskType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "Classification") {
		var out Classification
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Classification: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Forecasting") {
		var out Forecasting
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Forecasting: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImageClassification") {
		var out ImageClassification
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageClassification: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImageClassificationMultilabel") {
		var out ImageClassificationMultilabel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageClassificationMultilabel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImageInstanceSegmentation") {
		var out ImageInstanceSegmentation
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageInstanceSegmentation: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ImageObjectDetection") {
		var out ImageObjectDetection
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageObjectDetection: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Regression") {
		var out Regression
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Regression: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TextClassification") {
		var out TextClassification
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TextClassification: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TextClassificationMultilabel") {
		var out TextClassificationMultilabel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TextClassificationMultilabel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TextNER") {
		var out TextNer
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TextNer: %+v", err)
		}
		return out, nil
	}

	type RawAutoMLVerticalImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawAutoMLVerticalImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

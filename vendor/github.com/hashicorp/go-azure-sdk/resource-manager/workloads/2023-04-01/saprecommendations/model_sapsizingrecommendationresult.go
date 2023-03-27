package saprecommendations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPSizingRecommendationResult interface {
}

func unmarshalSAPSizingRecommendationResultImplementation(input []byte) (SAPSizingRecommendationResult, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SAPSizingRecommendationResult into map[string]interface: %+v", err)
	}

	value, ok := temp["deploymentType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "SingleServer") {
		var out SingleServerRecommendationResult
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SingleServerRecommendationResult: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "ThreeTier") {
		var out ThreeTierRecommendationResult
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ThreeTierRecommendationResult: %+v", err)
		}
		return out, nil
	}

	type RawSAPSizingRecommendationResultImpl struct {
		Type   string                 `json:"-"`
		Values map[string]interface{} `json:"-"`
	}
	out := RawSAPSizingRecommendationResultImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

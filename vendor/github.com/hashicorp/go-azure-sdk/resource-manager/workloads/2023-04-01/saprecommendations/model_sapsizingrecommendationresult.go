package saprecommendations

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPSizingRecommendationResult interface {
	SAPSizingRecommendationResult() BaseSAPSizingRecommendationResultImpl
}

var _ SAPSizingRecommendationResult = BaseSAPSizingRecommendationResultImpl{}

type BaseSAPSizingRecommendationResultImpl struct {
	DeploymentType SAPDeploymentType `json:"deploymentType"`
}

func (s BaseSAPSizingRecommendationResultImpl) SAPSizingRecommendationResult() BaseSAPSizingRecommendationResultImpl {
	return s
}

var _ SAPSizingRecommendationResult = RawSAPSizingRecommendationResultImpl{}

// RawSAPSizingRecommendationResultImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSAPSizingRecommendationResultImpl struct {
	sAPSizingRecommendationResult BaseSAPSizingRecommendationResultImpl
	Type                          string
	Values                        map[string]interface{}
}

func (s RawSAPSizingRecommendationResultImpl) SAPSizingRecommendationResult() BaseSAPSizingRecommendationResultImpl {
	return s.sAPSizingRecommendationResult
}

func UnmarshalSAPSizingRecommendationResultImplementation(input []byte) (SAPSizingRecommendationResult, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SAPSizingRecommendationResult into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["deploymentType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseSAPSizingRecommendationResultImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSAPSizingRecommendationResultImpl: %+v", err)
	}

	return RawSAPSizingRecommendationResultImpl{
		sAPSizingRecommendationResult: parent,
		Type:                          value,
		Values:                        temp,
	}, nil

}

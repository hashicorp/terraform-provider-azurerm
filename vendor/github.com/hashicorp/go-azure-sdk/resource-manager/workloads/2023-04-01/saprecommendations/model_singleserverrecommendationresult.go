package saprecommendations

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SAPSizingRecommendationResult = SingleServerRecommendationResult{}

type SingleServerRecommendationResult struct {
	VMSku *string `json:"vmSku,omitempty"`

	// Fields inherited from SAPSizingRecommendationResult

	DeploymentType SAPDeploymentType `json:"deploymentType"`
}

func (s SingleServerRecommendationResult) SAPSizingRecommendationResult() BaseSAPSizingRecommendationResultImpl {
	return BaseSAPSizingRecommendationResultImpl{
		DeploymentType: s.DeploymentType,
	}
}

var _ json.Marshaler = SingleServerRecommendationResult{}

func (s SingleServerRecommendationResult) MarshalJSON() ([]byte, error) {
	type wrapper SingleServerRecommendationResult
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SingleServerRecommendationResult: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SingleServerRecommendationResult: %+v", err)
	}

	decoded["deploymentType"] = "SingleServer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SingleServerRecommendationResult: %+v", err)
	}

	return encoded, nil
}

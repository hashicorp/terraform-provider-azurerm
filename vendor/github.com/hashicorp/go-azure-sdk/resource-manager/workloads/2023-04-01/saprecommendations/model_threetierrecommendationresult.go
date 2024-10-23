package saprecommendations

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SAPSizingRecommendationResult = ThreeTierRecommendationResult{}

type ThreeTierRecommendationResult struct {
	ApplicationServerInstanceCount *int64  `json:"applicationServerInstanceCount,omitempty"`
	ApplicationServerVMSku         *string `json:"applicationServerVmSku,omitempty"`
	CentralServerInstanceCount     *int64  `json:"centralServerInstanceCount,omitempty"`
	CentralServerVMSku             *string `json:"centralServerVmSku,omitempty"`
	DatabaseInstanceCount          *int64  `json:"databaseInstanceCount,omitempty"`
	DbVMSku                        *string `json:"dbVmSku,omitempty"`

	// Fields inherited from SAPSizingRecommendationResult

	DeploymentType SAPDeploymentType `json:"deploymentType"`
}

func (s ThreeTierRecommendationResult) SAPSizingRecommendationResult() BaseSAPSizingRecommendationResultImpl {
	return BaseSAPSizingRecommendationResultImpl{
		DeploymentType: s.DeploymentType,
	}
}

var _ json.Marshaler = ThreeTierRecommendationResult{}

func (s ThreeTierRecommendationResult) MarshalJSON() ([]byte, error) {
	type wrapper ThreeTierRecommendationResult
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ThreeTierRecommendationResult: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ThreeTierRecommendationResult: %+v", err)
	}

	decoded["deploymentType"] = "ThreeTier"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ThreeTierRecommendationResult: %+v", err)
	}

	return encoded, nil
}

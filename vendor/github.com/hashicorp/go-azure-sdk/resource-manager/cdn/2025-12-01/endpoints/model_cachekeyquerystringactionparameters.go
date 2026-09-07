package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = CacheKeyQueryStringActionParameters{}

type CacheKeyQueryStringActionParameters struct {
	QueryParameters     *string             `json:"queryParameters,omitempty"`
	QueryStringBehavior QueryStringBehavior `json:"queryStringBehavior"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s CacheKeyQueryStringActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = CacheKeyQueryStringActionParameters{}

func (s CacheKeyQueryStringActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper CacheKeyQueryStringActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CacheKeyQueryStringActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CacheKeyQueryStringActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleCacheKeyQueryStringBehaviorActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CacheKeyQueryStringActionParameters: %+v", err)
	}

	return encoded, nil
}

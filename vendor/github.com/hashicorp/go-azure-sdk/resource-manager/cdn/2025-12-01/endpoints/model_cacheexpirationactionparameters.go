package endpoints

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryRuleActionParameters = CacheExpirationActionParameters{}

type CacheExpirationActionParameters struct {
	CacheBehavior CacheBehavior `json:"cacheBehavior"`
	CacheDuration *string       `json:"cacheDuration,omitempty"`
	CacheType     CacheType     `json:"cacheType"`

	// Fields inherited from DeliveryRuleActionParameters

	TypeName DeliveryRuleActionParametersType `json:"typeName"`
}

func (s CacheExpirationActionParameters) DeliveryRuleActionParameters() BaseDeliveryRuleActionParametersImpl {
	return BaseDeliveryRuleActionParametersImpl{
		TypeName: s.TypeName,
	}
}

var _ json.Marshaler = CacheExpirationActionParameters{}

func (s CacheExpirationActionParameters) MarshalJSON() ([]byte, error) {
	type wrapper CacheExpirationActionParameters
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CacheExpirationActionParameters: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CacheExpirationActionParameters: %+v", err)
	}

	decoded["typeName"] = "DeliveryRuleCacheExpirationActionParameters"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CacheExpirationActionParameters: %+v", err)
	}

	return encoded, nil
}

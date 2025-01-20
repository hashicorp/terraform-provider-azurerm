package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ GatewayResponseCacheProperties = GatewayLocalResponseCachePerInstanceProperties{}

type GatewayLocalResponseCachePerInstanceProperties struct {
	Size       *string `json:"size,omitempty"`
	TimeToLive *string `json:"timeToLive,omitempty"`

	// Fields inherited from GatewayResponseCacheProperties

	ResponseCacheType string `json:"responseCacheType"`
}

func (s GatewayLocalResponseCachePerInstanceProperties) GatewayResponseCacheProperties() BaseGatewayResponseCachePropertiesImpl {
	return BaseGatewayResponseCachePropertiesImpl{
		ResponseCacheType: s.ResponseCacheType,
	}
}

var _ json.Marshaler = GatewayLocalResponseCachePerInstanceProperties{}

func (s GatewayLocalResponseCachePerInstanceProperties) MarshalJSON() ([]byte, error) {
	type wrapper GatewayLocalResponseCachePerInstanceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GatewayLocalResponseCachePerInstanceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GatewayLocalResponseCachePerInstanceProperties: %+v", err)
	}

	decoded["responseCacheType"] = "LocalCachePerInstance"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GatewayLocalResponseCachePerInstanceProperties: %+v", err)
	}

	return encoded, nil
}

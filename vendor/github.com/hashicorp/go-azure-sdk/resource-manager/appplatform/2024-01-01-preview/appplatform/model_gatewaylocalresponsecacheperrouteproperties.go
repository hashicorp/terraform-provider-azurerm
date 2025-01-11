package appplatform

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ GatewayResponseCacheProperties = GatewayLocalResponseCachePerRouteProperties{}

type GatewayLocalResponseCachePerRouteProperties struct {
	Size       *string `json:"size,omitempty"`
	TimeToLive *string `json:"timeToLive,omitempty"`

	// Fields inherited from GatewayResponseCacheProperties

	ResponseCacheType string `json:"responseCacheType"`
}

func (s GatewayLocalResponseCachePerRouteProperties) GatewayResponseCacheProperties() BaseGatewayResponseCachePropertiesImpl {
	return BaseGatewayResponseCachePropertiesImpl{
		ResponseCacheType: s.ResponseCacheType,
	}
}

var _ json.Marshaler = GatewayLocalResponseCachePerRouteProperties{}

func (s GatewayLocalResponseCachePerRouteProperties) MarshalJSON() ([]byte, error) {
	type wrapper GatewayLocalResponseCachePerRouteProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling GatewayLocalResponseCachePerRouteProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling GatewayLocalResponseCachePerRouteProperties: %+v", err)
	}

	decoded["responseCacheType"] = "LocalCachePerRoute"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling GatewayLocalResponseCachePerRouteProperties: %+v", err)
	}

	return encoded, nil
}

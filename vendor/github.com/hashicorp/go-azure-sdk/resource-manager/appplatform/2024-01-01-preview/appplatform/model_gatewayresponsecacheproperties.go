package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayResponseCacheProperties interface {
}

// RawGatewayResponseCachePropertiesImpl is returned when the Discriminated Value
// doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawGatewayResponseCachePropertiesImpl struct {
	Type   string
	Values map[string]interface{}
}

func unmarshalGatewayResponseCachePropertiesImplementation(input []byte) (GatewayResponseCacheProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling GatewayResponseCacheProperties into map[string]interface: %+v", err)
	}

	value, ok := temp["responseCacheType"].(string)
	if !ok {
		return nil, nil
	}

	if strings.EqualFold(value, "LocalCachePerInstance") {
		var out GatewayLocalResponseCachePerInstanceProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GatewayLocalResponseCachePerInstanceProperties: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LocalCachePerRoute") {
		var out GatewayLocalResponseCachePerRouteProperties
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into GatewayLocalResponseCachePerRouteProperties: %+v", err)
		}
		return out, nil
	}

	out := RawGatewayResponseCachePropertiesImpl{
		Type:   value,
		Values: temp,
	}
	return out, nil

}

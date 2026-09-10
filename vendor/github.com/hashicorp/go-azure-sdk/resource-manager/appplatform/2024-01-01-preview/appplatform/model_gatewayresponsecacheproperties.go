package appplatform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayResponseCacheProperties interface {
	GatewayResponseCacheProperties() BaseGatewayResponseCachePropertiesImpl
}

var _ GatewayResponseCacheProperties = BaseGatewayResponseCachePropertiesImpl{}

type BaseGatewayResponseCachePropertiesImpl struct {
	ResponseCacheType string `json:"responseCacheType"`
}

func (s BaseGatewayResponseCachePropertiesImpl) GatewayResponseCacheProperties() BaseGatewayResponseCachePropertiesImpl {
	return s
}

var _ GatewayResponseCacheProperties = RawGatewayResponseCachePropertiesImpl{}

// RawGatewayResponseCachePropertiesImpl is returned when the Discriminated Value doesn't match any of the defined types.
// It can also be used as a Request Payload to provide a raw JSON payload, which is useful
// for preserving arbitrary/extensible JSON properties across a round-trip.
type RawGatewayResponseCachePropertiesImpl struct {
	gatewayResponseCacheProperties BaseGatewayResponseCachePropertiesImpl
	Type                           string
	Values                         map[string]interface{}
}

func (s RawGatewayResponseCachePropertiesImpl) GatewayResponseCacheProperties() BaseGatewayResponseCachePropertiesImpl {
	return s.gatewayResponseCacheProperties
}

func (s RawGatewayResponseCachePropertiesImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Values)
}

func UnmarshalGatewayResponseCachePropertiesImplementation(input []byte) (GatewayResponseCacheProperties, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling GatewayResponseCacheProperties into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["responseCacheType"]; ok {
		value = fmt.Sprintf("%v", v)
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

	var parent BaseGatewayResponseCachePropertiesImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseGatewayResponseCachePropertiesImpl: %+v", err)
	}

	return RawGatewayResponseCachePropertiesImpl{
		gatewayResponseCacheProperties: parent,
		Type:                           value,
		Values:                         temp,
	}, nil

}

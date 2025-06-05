package namespaces

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticRoutingEnrichment interface {
	StaticRoutingEnrichment() BaseStaticRoutingEnrichmentImpl
}

var _ StaticRoutingEnrichment = BaseStaticRoutingEnrichmentImpl{}

type BaseStaticRoutingEnrichmentImpl struct {
	Key       *string                     `json:"key,omitempty"`
	ValueType StaticRoutingEnrichmentType `json:"valueType"`
}

func (s BaseStaticRoutingEnrichmentImpl) StaticRoutingEnrichment() BaseStaticRoutingEnrichmentImpl {
	return s
}

var _ StaticRoutingEnrichment = RawStaticRoutingEnrichmentImpl{}

// RawStaticRoutingEnrichmentImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawStaticRoutingEnrichmentImpl struct {
	staticRoutingEnrichment BaseStaticRoutingEnrichmentImpl
	Type                    string
	Values                  map[string]interface{}
}

func (s RawStaticRoutingEnrichmentImpl) StaticRoutingEnrichment() BaseStaticRoutingEnrichmentImpl {
	return s.staticRoutingEnrichment
}

func UnmarshalStaticRoutingEnrichmentImplementation(input []byte) (StaticRoutingEnrichment, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling StaticRoutingEnrichment into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["valueType"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "String") {
		var out StaticStringRoutingEnrichment
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into StaticStringRoutingEnrichment: %+v", err)
		}
		return out, nil
	}

	var parent BaseStaticRoutingEnrichmentImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseStaticRoutingEnrichmentImpl: %+v", err)
	}

	return RawStaticRoutingEnrichmentImpl{
		staticRoutingEnrichment: parent,
		Type:                    value,
		Values:                  temp,
	}, nil

}

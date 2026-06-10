package namespaces

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ StaticRoutingEnrichment = StaticStringRoutingEnrichment{}

type StaticStringRoutingEnrichment struct {
	Value *string `json:"value,omitempty"`

	// Fields inherited from StaticRoutingEnrichment

	Key       *string                     `json:"key,omitempty"`
	ValueType StaticRoutingEnrichmentType `json:"valueType"`
}

func (s StaticStringRoutingEnrichment) StaticRoutingEnrichment() BaseStaticRoutingEnrichmentImpl {
	return BaseStaticRoutingEnrichmentImpl{
		Key:       s.Key,
		ValueType: s.ValueType,
	}
}

var _ json.Marshaler = StaticStringRoutingEnrichment{}

func (s StaticStringRoutingEnrichment) MarshalJSON() ([]byte, error) {
	type wrapper StaticStringRoutingEnrichment
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StaticStringRoutingEnrichment: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StaticStringRoutingEnrichment: %+v", err)
	}

	decoded["valueType"] = "String"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StaticStringRoutingEnrichment: %+v", err)
	}

	return encoded, nil
}

package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryAttributeMapping = StaticDeliveryAttributeMapping{}

type StaticDeliveryAttributeMapping struct {
	Properties *StaticDeliveryAttributeMappingProperties `json:"properties,omitempty"`

	// Fields inherited from DeliveryAttributeMapping

	Name *string                      `json:"name,omitempty"`
	Type DeliveryAttributeMappingType `json:"type"`
}

func (s StaticDeliveryAttributeMapping) DeliveryAttributeMapping() BaseDeliveryAttributeMappingImpl {
	return BaseDeliveryAttributeMappingImpl{
		Name: s.Name,
		Type: s.Type,
	}
}

var _ json.Marshaler = StaticDeliveryAttributeMapping{}

func (s StaticDeliveryAttributeMapping) MarshalJSON() ([]byte, error) {
	type wrapper StaticDeliveryAttributeMapping
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling StaticDeliveryAttributeMapping: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling StaticDeliveryAttributeMapping: %+v", err)
	}

	decoded["type"] = "Static"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling StaticDeliveryAttributeMapping: %+v", err)
	}

	return encoded, nil
}

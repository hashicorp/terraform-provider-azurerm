package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ DeliveryAttributeMapping = DynamicDeliveryAttributeMapping{}

type DynamicDeliveryAttributeMapping struct {
	Properties *DynamicDeliveryAttributeMappingProperties `json:"properties,omitempty"`

	// Fields inherited from DeliveryAttributeMapping
	Name *string `json:"name,omitempty"`
}

var _ json.Marshaler = DynamicDeliveryAttributeMapping{}

func (s DynamicDeliveryAttributeMapping) MarshalJSON() ([]byte, error) {
	type wrapper DynamicDeliveryAttributeMapping
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DynamicDeliveryAttributeMapping: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DynamicDeliveryAttributeMapping: %+v", err)
	}
	decoded["type"] = "Dynamic"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DynamicDeliveryAttributeMapping: %+v", err)
	}

	return encoded, nil
}

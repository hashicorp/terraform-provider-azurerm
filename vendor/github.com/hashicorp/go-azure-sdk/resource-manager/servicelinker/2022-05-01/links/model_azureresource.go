package links

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ TargetServiceBase = AzureResource{}

type AzureResource struct {
	Id                 *string                     `json:"id,omitempty"`
	ResourceProperties AzureResourcePropertiesBase `json:"resourceProperties"`

	// Fields inherited from TargetServiceBase
}

var _ json.Marshaler = AzureResource{}

func (s AzureResource) MarshalJSON() ([]byte, error) {
	type wrapper AzureResource
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureResource: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureResource: %+v", err)
	}
	decoded["type"] = "AzureResource"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureResource: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureResource{}

func (s *AzureResource) UnmarshalJSON(bytes []byte) error {
	type alias AzureResource
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into AzureResource: %+v", err)
	}

	s.Id = decoded.Id

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureResource into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["resourceProperties"]; ok {
		impl, err := unmarshalAzureResourcePropertiesBaseImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'ResourceProperties' for 'AzureResource': %+v", err)
		}
		s.ResourceProperties = impl
	}
	return nil
}

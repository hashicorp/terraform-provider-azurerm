package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LinkedService = WebLinkedService{}

type WebLinkedService struct {
	TypeProperties WebLinkedServiceTypeProperties `json:"typeProperties"`

	// Fields inherited from LinkedService

	Annotations *[]interface{}                     `json:"annotations,omitempty"`
	ConnectVia  *IntegrationRuntimeReference       `json:"connectVia,omitempty"`
	Description *string                            `json:"description,omitempty"`
	Parameters  *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Type        string                             `json:"type"`
	Version     *string                            `json:"version,omitempty"`
}

func (s WebLinkedService) LinkedService() BaseLinkedServiceImpl {
	return BaseLinkedServiceImpl{
		Annotations: s.Annotations,
		ConnectVia:  s.ConnectVia,
		Description: s.Description,
		Parameters:  s.Parameters,
		Type:        s.Type,
		Version:     s.Version,
	}
}

var _ json.Marshaler = WebLinkedService{}

func (s WebLinkedService) MarshalJSON() ([]byte, error) {
	type wrapper WebLinkedService
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WebLinkedService: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WebLinkedService: %+v", err)
	}

	decoded["type"] = "Web"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WebLinkedService: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &WebLinkedService{}

func (s *WebLinkedService) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Annotations *[]interface{}                     `json:"annotations,omitempty"`
		ConnectVia  *IntegrationRuntimeReference       `json:"connectVia,omitempty"`
		Description *string                            `json:"description,omitempty"`
		Parameters  *map[string]ParameterSpecification `json:"parameters,omitempty"`
		Type        string                             `json:"type"`
		Version     *string                            `json:"version,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Annotations = decoded.Annotations
	s.ConnectVia = decoded.ConnectVia
	s.Description = decoded.Description
	s.Parameters = decoded.Parameters
	s.Type = decoded.Type
	s.Version = decoded.Version

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WebLinkedService into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["typeProperties"]; ok {
		impl, err := UnmarshalWebLinkedServiceTypePropertiesImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'TypeProperties' for 'WebLinkedService': %+v", err)
		}
		s.TypeProperties = impl
	}

	return nil
}

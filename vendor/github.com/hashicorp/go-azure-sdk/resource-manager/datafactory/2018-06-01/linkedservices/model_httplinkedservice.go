package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LinkedService = HTTPLinkedService{}

type HTTPLinkedService struct {
	TypeProperties HTTPLinkedServiceTypeProperties `json:"typeProperties"`

	// Fields inherited from LinkedService

	Annotations *[]interface{}                     `json:"annotations,omitempty"`
	ConnectVia  *IntegrationRuntimeReference       `json:"connectVia,omitempty"`
	Description *string                            `json:"description,omitempty"`
	Parameters  *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Type        string                             `json:"type"`
	Version     *string                            `json:"version,omitempty"`
}

func (s HTTPLinkedService) LinkedService() BaseLinkedServiceImpl {
	return BaseLinkedServiceImpl{
		Annotations: s.Annotations,
		ConnectVia:  s.ConnectVia,
		Description: s.Description,
		Parameters:  s.Parameters,
		Type:        s.Type,
		Version:     s.Version,
	}
}

var _ json.Marshaler = HTTPLinkedService{}

func (s HTTPLinkedService) MarshalJSON() ([]byte, error) {
	type wrapper HTTPLinkedService
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling HTTPLinkedService: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling HTTPLinkedService: %+v", err)
	}

	decoded["type"] = "HttpServer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling HTTPLinkedService: %+v", err)
	}

	return encoded, nil
}

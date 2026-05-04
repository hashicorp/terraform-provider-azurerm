package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LinkedService = ODataLinkedService{}

type ODataLinkedService struct {
	TypeProperties ODataLinkedServiceTypeProperties `json:"typeProperties"`

	// Fields inherited from LinkedService

	Annotations *[]interface{}                     `json:"annotations,omitempty"`
	ConnectVia  *IntegrationRuntimeReference       `json:"connectVia,omitempty"`
	Description *string                            `json:"description,omitempty"`
	Parameters  *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Type        string                             `json:"type"`
	Version     *string                            `json:"version,omitempty"`
}

func (s ODataLinkedService) LinkedService() BaseLinkedServiceImpl {
	return BaseLinkedServiceImpl{
		Annotations: s.Annotations,
		ConnectVia:  s.ConnectVia,
		Description: s.Description,
		Parameters:  s.Parameters,
		Type:        s.Type,
		Version:     s.Version,
	}
}

var _ json.Marshaler = ODataLinkedService{}

func (s ODataLinkedService) MarshalJSON() ([]byte, error) {
	type wrapper ODataLinkedService
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ODataLinkedService: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ODataLinkedService: %+v", err)
	}

	decoded["type"] = "OData"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ODataLinkedService: %+v", err)
	}

	return encoded, nil
}

package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LinkedService = XeroLinkedService{}

type XeroLinkedService struct {
	TypeProperties XeroLinkedServiceTypeProperties `json:"typeProperties"`

	// Fields inherited from LinkedService

	Annotations *[]interface{}                     `json:"annotations,omitempty"`
	ConnectVia  *IntegrationRuntimeReference       `json:"connectVia,omitempty"`
	Description *string                            `json:"description,omitempty"`
	Parameters  *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Type        string                             `json:"type"`
	Version     *string                            `json:"version,omitempty"`
}

func (s XeroLinkedService) LinkedService() BaseLinkedServiceImpl {
	return BaseLinkedServiceImpl{
		Annotations: s.Annotations,
		ConnectVia:  s.ConnectVia,
		Description: s.Description,
		Parameters:  s.Parameters,
		Type:        s.Type,
		Version:     s.Version,
	}
}

var _ json.Marshaler = XeroLinkedService{}

func (s XeroLinkedService) MarshalJSON() ([]byte, error) {
	type wrapper XeroLinkedService
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling XeroLinkedService: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling XeroLinkedService: %+v", err)
	}

	decoded["type"] = "Xero"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling XeroLinkedService: %+v", err)
	}

	return encoded, nil
}

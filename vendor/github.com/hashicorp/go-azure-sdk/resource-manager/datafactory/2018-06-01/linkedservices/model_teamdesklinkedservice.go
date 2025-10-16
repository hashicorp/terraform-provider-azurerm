package linkedservices

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ LinkedService = TeamDeskLinkedService{}

type TeamDeskLinkedService struct {
	TypeProperties TeamDeskLinkedServiceTypeProperties `json:"typeProperties"`

	// Fields inherited from LinkedService

	Annotations *[]interface{}                     `json:"annotations,omitempty"`
	ConnectVia  *IntegrationRuntimeReference       `json:"connectVia,omitempty"`
	Description *string                            `json:"description,omitempty"`
	Parameters  *map[string]ParameterSpecification `json:"parameters,omitempty"`
	Type        string                             `json:"type"`
	Version     *string                            `json:"version,omitempty"`
}

func (s TeamDeskLinkedService) LinkedService() BaseLinkedServiceImpl {
	return BaseLinkedServiceImpl{
		Annotations: s.Annotations,
		ConnectVia:  s.ConnectVia,
		Description: s.Description,
		Parameters:  s.Parameters,
		Type:        s.Type,
		Version:     s.Version,
	}
}

var _ json.Marshaler = TeamDeskLinkedService{}

func (s TeamDeskLinkedService) MarshalJSON() ([]byte, error) {
	type wrapper TeamDeskLinkedService
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TeamDeskLinkedService: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TeamDeskLinkedService: %+v", err)
	}

	decoded["type"] = "TeamDesk"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TeamDeskLinkedService: %+v", err)
	}

	return encoded, nil
}

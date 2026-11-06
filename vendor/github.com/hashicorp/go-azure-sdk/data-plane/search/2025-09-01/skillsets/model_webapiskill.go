package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = WebApiSkill{}

type WebApiSkill struct {
	AuthIdentity        SearchIndexerDataIdentity `json:"authIdentity"`
	AuthResourceId      *string                   `json:"authResourceId,omitempty"`
	BatchSize           *int64                    `json:"batchSize,omitempty"`
	DegreeOfParallelism *int64                    `json:"degreeOfParallelism,omitempty"`
	HTTPHeaders         *map[string]string        `json:"httpHeaders,omitempty"`
	HTTPMethod          *string                   `json:"httpMethod,omitempty"`
	Timeout             *string                   `json:"timeout,omitempty"`
	Uri                 string                    `json:"uri"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s WebApiSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = WebApiSkill{}

func (s WebApiSkill) MarshalJSON() ([]byte, error) {
	type wrapper WebApiSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling WebApiSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling WebApiSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Custom.WebApiSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling WebApiSkill: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &WebApiSkill{}

func (s *WebApiSkill) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		AuthResourceId      *string                   `json:"authResourceId,omitempty"`
		BatchSize           *int64                    `json:"batchSize,omitempty"`
		DegreeOfParallelism *int64                    `json:"degreeOfParallelism,omitempty"`
		HTTPHeaders         *map[string]string        `json:"httpHeaders,omitempty"`
		HTTPMethod          *string                   `json:"httpMethod,omitempty"`
		Timeout             *string                   `json:"timeout,omitempty"`
		Uri                 string                    `json:"uri"`
		Context             *string                   `json:"context,omitempty"`
		Description         *string                   `json:"description,omitempty"`
		Inputs              []InputFieldMappingEntry  `json:"inputs"`
		Name                *string                   `json:"name,omitempty"`
		OdataType           string                    `json:"@odata.type"`
		Outputs             []OutputFieldMappingEntry `json:"outputs"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.AuthResourceId = decoded.AuthResourceId
	s.BatchSize = decoded.BatchSize
	s.DegreeOfParallelism = decoded.DegreeOfParallelism
	s.HTTPHeaders = decoded.HTTPHeaders
	s.HTTPMethod = decoded.HTTPMethod
	s.Timeout = decoded.Timeout
	s.Uri = decoded.Uri
	s.Context = decoded.Context
	s.Description = decoded.Description
	s.Inputs = decoded.Inputs
	s.Name = decoded.Name
	s.OdataType = decoded.OdataType
	s.Outputs = decoded.Outputs

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling WebApiSkill into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authIdentity"]; ok {
		impl, err := UnmarshalSearchIndexerDataIdentityImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthIdentity' for 'WebApiSkill': %+v", err)
		}
		s.AuthIdentity = impl
	}

	return nil
}

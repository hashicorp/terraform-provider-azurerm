package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = AzureOpenAIEmbeddingSkill{}

type AzureOpenAIEmbeddingSkill struct {
	ApiKey       *string                   `json:"apiKey,omitempty"`
	AuthIdentity SearchIndexerDataIdentity `json:"authIdentity"`
	DeploymentId *string                   `json:"deploymentId,omitempty"`
	Dimensions   *int64                    `json:"dimensions,omitempty"`
	ModelName    *AzureOpenAIModelName     `json:"modelName,omitempty"`
	ResourceUri  *string                   `json:"resourceUri,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s AzureOpenAIEmbeddingSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = AzureOpenAIEmbeddingSkill{}

func (s AzureOpenAIEmbeddingSkill) MarshalJSON() ([]byte, error) {
	type wrapper AzureOpenAIEmbeddingSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling AzureOpenAIEmbeddingSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling AzureOpenAIEmbeddingSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.AzureOpenAIEmbeddingSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling AzureOpenAIEmbeddingSkill: %+v", err)
	}

	return encoded, nil
}

var _ json.Unmarshaler = &AzureOpenAIEmbeddingSkill{}

func (s *AzureOpenAIEmbeddingSkill) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ApiKey       *string                   `json:"apiKey,omitempty"`
		DeploymentId *string                   `json:"deploymentId,omitempty"`
		Dimensions   *int64                    `json:"dimensions,omitempty"`
		ModelName    *AzureOpenAIModelName     `json:"modelName,omitempty"`
		ResourceUri  *string                   `json:"resourceUri,omitempty"`
		Context      *string                   `json:"context,omitempty"`
		Description  *string                   `json:"description,omitempty"`
		Inputs       []InputFieldMappingEntry  `json:"inputs"`
		Name         *string                   `json:"name,omitempty"`
		OdataType    string                    `json:"@odata.type"`
		Outputs      []OutputFieldMappingEntry `json:"outputs"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ApiKey = decoded.ApiKey
	s.DeploymentId = decoded.DeploymentId
	s.Dimensions = decoded.Dimensions
	s.ModelName = decoded.ModelName
	s.ResourceUri = decoded.ResourceUri
	s.Context = decoded.Context
	s.Description = decoded.Description
	s.Inputs = decoded.Inputs
	s.Name = decoded.Name
	s.OdataType = decoded.OdataType
	s.Outputs = decoded.Outputs

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureOpenAIEmbeddingSkill into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authIdentity"]; ok {
		impl, err := UnmarshalSearchIndexerDataIdentityImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthIdentity' for 'AzureOpenAIEmbeddingSkill': %+v", err)
		}
		s.AuthIdentity = impl
	}

	return nil
}

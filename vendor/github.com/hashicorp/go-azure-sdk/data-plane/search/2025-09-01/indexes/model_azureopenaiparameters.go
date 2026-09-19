package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureOpenAIParameters struct {
	ApiKey       *string                   `json:"apiKey,omitempty"`
	AuthIdentity SearchIndexerDataIdentity `json:"authIdentity"`
	DeploymentId *string                   `json:"deploymentId,omitempty"`
	ModelName    *AzureOpenAIModelName     `json:"modelName,omitempty"`
	ResourceUri  *string                   `json:"resourceUri,omitempty"`
}

var _ json.Unmarshaler = &AzureOpenAIParameters{}

func (s *AzureOpenAIParameters) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		ApiKey       *string               `json:"apiKey,omitempty"`
		DeploymentId *string               `json:"deploymentId,omitempty"`
		ModelName    *AzureOpenAIModelName `json:"modelName,omitempty"`
		ResourceUri  *string               `json:"resourceUri,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.ApiKey = decoded.ApiKey
	s.DeploymentId = decoded.DeploymentId
	s.ModelName = decoded.ModelName
	s.ResourceUri = decoded.ResourceUri

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling AzureOpenAIParameters into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["authIdentity"]; ok {
		impl, err := UnmarshalSearchIndexerDataIdentityImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'AuthIdentity' for 'AzureOpenAIParameters': %+v", err)
		}
		s.AuthIdentity = impl
	}

	return nil
}

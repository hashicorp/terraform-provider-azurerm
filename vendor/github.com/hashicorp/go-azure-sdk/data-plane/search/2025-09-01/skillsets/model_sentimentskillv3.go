package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = SentimentSkillV3{}

type SentimentSkillV3 struct {
	DefaultLanguageCode  *string `json:"defaultLanguageCode,omitempty"`
	IncludeOpinionMining *bool   `json:"includeOpinionMining,omitempty"`
	ModelVersion         *string `json:"modelVersion,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s SentimentSkillV3) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = SentimentSkillV3{}

func (s SentimentSkillV3) MarshalJSON() ([]byte, error) {
	type wrapper SentimentSkillV3
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SentimentSkillV3: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SentimentSkillV3: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.V3.SentimentSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SentimentSkillV3: %+v", err)
	}

	return encoded, nil
}

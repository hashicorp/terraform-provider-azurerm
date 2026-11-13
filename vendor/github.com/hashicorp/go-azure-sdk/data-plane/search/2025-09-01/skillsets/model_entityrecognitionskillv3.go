package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = EntityRecognitionSkillV3{}

type EntityRecognitionSkillV3 struct {
	Categories          *[]string `json:"categories,omitempty"`
	DefaultLanguageCode *string   `json:"defaultLanguageCode,omitempty"`
	MinimumPrecision    *float64  `json:"minimumPrecision,omitempty"`
	ModelVersion        *string   `json:"modelVersion,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s EntityRecognitionSkillV3) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = EntityRecognitionSkillV3{}

func (s EntityRecognitionSkillV3) MarshalJSON() ([]byte, error) {
	type wrapper EntityRecognitionSkillV3
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EntityRecognitionSkillV3: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EntityRecognitionSkillV3: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.V3.EntityRecognitionSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EntityRecognitionSkillV3: %+v", err)
	}

	return encoded, nil
}

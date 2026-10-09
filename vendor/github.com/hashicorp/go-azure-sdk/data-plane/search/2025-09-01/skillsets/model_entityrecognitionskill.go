package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = EntityRecognitionSkill{}

type EntityRecognitionSkill struct {
	Categories              *[]EntityCategory               `json:"categories,omitempty"`
	DefaultLanguageCode     *EntityRecognitionSkillLanguage `json:"defaultLanguageCode,omitempty"`
	IncludeTypelessEntities *bool                           `json:"includeTypelessEntities,omitempty"`
	MinimumPrecision        *float64                        `json:"minimumPrecision,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s EntityRecognitionSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = EntityRecognitionSkill{}

func (s EntityRecognitionSkill) MarshalJSON() ([]byte, error) {
	type wrapper EntityRecognitionSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling EntityRecognitionSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling EntityRecognitionSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.EntityRecognitionSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling EntityRecognitionSkill: %+v", err)
	}

	return encoded, nil
}

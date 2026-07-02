package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = LanguageDetectionSkill{}

type LanguageDetectionSkill struct {
	DefaultCountryHint *string `json:"defaultCountryHint,omitempty"`
	ModelVersion       *string `json:"modelVersion,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s LanguageDetectionSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = LanguageDetectionSkill{}

func (s LanguageDetectionSkill) MarshalJSON() ([]byte, error) {
	type wrapper LanguageDetectionSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LanguageDetectionSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LanguageDetectionSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.LanguageDetectionSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LanguageDetectionSkill: %+v", err)
	}

	return encoded, nil
}

package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = TextTranslationSkill{}

type TextTranslationSkill struct {
	DefaultFromLanguageCode *TextTranslationSkillLanguage `json:"defaultFromLanguageCode,omitempty"`
	DefaultToLanguageCode   TextTranslationSkillLanguage  `json:"defaultToLanguageCode"`
	SuggestedFrom           *TextTranslationSkillLanguage `json:"suggestedFrom,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s TextTranslationSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = TextTranslationSkill{}

func (s TextTranslationSkill) MarshalJSON() ([]byte, error) {
	type wrapper TextTranslationSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling TextTranslationSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling TextTranslationSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.TranslationSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling TextTranslationSkill: %+v", err)
	}

	return encoded, nil
}

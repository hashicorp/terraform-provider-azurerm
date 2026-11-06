package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = KeyPhraseExtractionSkill{}

type KeyPhraseExtractionSkill struct {
	DefaultLanguageCode *KeyPhraseExtractionSkillLanguage `json:"defaultLanguageCode,omitempty"`
	MaxKeyPhraseCount   *int64                            `json:"maxKeyPhraseCount,omitempty"`
	ModelVersion        *string                           `json:"modelVersion,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s KeyPhraseExtractionSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = KeyPhraseExtractionSkill{}

func (s KeyPhraseExtractionSkill) MarshalJSON() ([]byte, error) {
	type wrapper KeyPhraseExtractionSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling KeyPhraseExtractionSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling KeyPhraseExtractionSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.KeyPhraseExtractionSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling KeyPhraseExtractionSkill: %+v", err)
	}

	return encoded, nil
}

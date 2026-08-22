package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = CustomEntityLookupSkill{}

type CustomEntityLookupSkill struct {
	DefaultLanguageCode            *CustomEntityLookupSkillLanguage `json:"defaultLanguageCode,omitempty"`
	EntitiesDefinitionUri          *string                          `json:"entitiesDefinitionUri,omitempty"`
	GlobalDefaultAccentSensitive   *bool                            `json:"globalDefaultAccentSensitive,omitempty"`
	GlobalDefaultCaseSensitive     *bool                            `json:"globalDefaultCaseSensitive,omitempty"`
	GlobalDefaultFuzzyEditDistance *int64                           `json:"globalDefaultFuzzyEditDistance,omitempty"`
	InlineEntitiesDefinition       *[]CustomEntity                  `json:"inlineEntitiesDefinition,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s CustomEntityLookupSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = CustomEntityLookupSkill{}

func (s CustomEntityLookupSkill) MarshalJSON() ([]byte, error) {
	type wrapper CustomEntityLookupSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling CustomEntityLookupSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling CustomEntityLookupSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.CustomEntityLookupSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling CustomEntityLookupSkill: %+v", err)
	}

	return encoded, nil
}

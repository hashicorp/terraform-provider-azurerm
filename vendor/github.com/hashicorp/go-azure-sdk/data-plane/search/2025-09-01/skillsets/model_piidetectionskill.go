package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = PIIDetectionSkill{}

type PIIDetectionSkill struct {
	DefaultLanguageCode *string                       `json:"defaultLanguageCode,omitempty"`
	Domain              *string                       `json:"domain,omitempty"`
	MaskingCharacter    *string                       `json:"maskingCharacter,omitempty"`
	MaskingMode         *PIIDetectionSkillMaskingMode `json:"maskingMode,omitempty"`
	MinimumPrecision    *float64                      `json:"minimumPrecision,omitempty"`
	ModelVersion        *string                       `json:"modelVersion,omitempty"`
	PiiCategories       *[]string                     `json:"piiCategories,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s PIIDetectionSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = PIIDetectionSkill{}

func (s PIIDetectionSkill) MarshalJSON() ([]byte, error) {
	type wrapper PIIDetectionSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling PIIDetectionSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling PIIDetectionSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.PIIDetectionSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling PIIDetectionSkill: %+v", err)
	}

	return encoded, nil
}

package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = SplitSkill{}

type SplitSkill struct {
	DefaultLanguageCode *SplitSkillLanguage `json:"defaultLanguageCode,omitempty"`
	MaximumPageLength   *int64              `json:"maximumPageLength,omitempty"`
	MaximumPagesToTake  *int64              `json:"maximumPagesToTake,omitempty"`
	PageOverlapLength   *int64              `json:"pageOverlapLength,omitempty"`
	TextSplitMode       *TextSplitMode      `json:"textSplitMode,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s SplitSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = SplitSkill{}

func (s SplitSkill) MarshalJSON() ([]byte, error) {
	type wrapper SplitSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling SplitSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling SplitSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Text.SplitSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling SplitSkill: %+v", err)
	}

	return encoded, nil
}

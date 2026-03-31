package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = ImageAnalysisSkill{}

type ImageAnalysisSkill struct {
	DefaultLanguageCode *ImageAnalysisSkillLanguage `json:"defaultLanguageCode,omitempty"`
	Details             *[]ImageDetail              `json:"details,omitempty"`
	VisualFeatures      *[]VisualFeature            `json:"visualFeatures,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s ImageAnalysisSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = ImageAnalysisSkill{}

func (s ImageAnalysisSkill) MarshalJSON() ([]byte, error) {
	type wrapper ImageAnalysisSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ImageAnalysisSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ImageAnalysisSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Vision.ImageAnalysisSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ImageAnalysisSkill: %+v", err)
	}

	return encoded, nil
}

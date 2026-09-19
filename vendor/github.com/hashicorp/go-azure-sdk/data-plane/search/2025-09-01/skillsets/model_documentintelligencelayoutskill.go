package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ SearchIndexerSkill = DocumentIntelligenceLayoutSkill{}

type DocumentIntelligenceLayoutSkill struct {
	ChunkingProperties  *DocumentIntelligenceLayoutSkillChunkingProperties  `json:"chunkingProperties,omitempty"`
	ExtractionOptions   *[]DocumentIntelligenceLayoutSkillExtractionOptions `json:"extractionOptions,omitempty"`
	MarkdownHeaderDepth *DocumentIntelligenceLayoutSkillMarkdownHeaderDepth `json:"markdownHeaderDepth,omitempty"`
	OutputFormat        *DocumentIntelligenceLayoutSkillOutputFormat        `json:"outputFormat,omitempty"`
	OutputMode          *DocumentIntelligenceLayoutSkillOutputMode          `json:"outputMode,omitempty"`

	// Fields inherited from SearchIndexerSkill

	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s DocumentIntelligenceLayoutSkill) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return BaseSearchIndexerSkillImpl{
		Context:     s.Context,
		Description: s.Description,
		Inputs:      s.Inputs,
		Name:        s.Name,
		OdataType:   s.OdataType,
		Outputs:     s.Outputs,
	}
}

var _ json.Marshaler = DocumentIntelligenceLayoutSkill{}

func (s DocumentIntelligenceLayoutSkill) MarshalJSON() ([]byte, error) {
	type wrapper DocumentIntelligenceLayoutSkill
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DocumentIntelligenceLayoutSkill: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DocumentIntelligenceLayoutSkill: %+v", err)
	}

	decoded["@odata.type"] = "#Microsoft.Skills.Util.DocumentIntelligenceLayoutSkill"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DocumentIntelligenceLayoutSkill: %+v", err)
	}

	return encoded, nil
}

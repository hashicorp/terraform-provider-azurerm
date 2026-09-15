package skillsets

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerSkill interface {
	SearchIndexerSkill() BaseSearchIndexerSkillImpl
}

var _ SearchIndexerSkill = BaseSearchIndexerSkillImpl{}

type BaseSearchIndexerSkillImpl struct {
	Context     *string                   `json:"context,omitempty"`
	Description *string                   `json:"description,omitempty"`
	Inputs      []InputFieldMappingEntry  `json:"inputs"`
	Name        *string                   `json:"name,omitempty"`
	OdataType   string                    `json:"@odata.type"`
	Outputs     []OutputFieldMappingEntry `json:"outputs"`
}

func (s BaseSearchIndexerSkillImpl) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return s
}

var _ SearchIndexerSkill = RawSearchIndexerSkillImpl{}

// RawSearchIndexerSkillImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawSearchIndexerSkillImpl struct {
	searchIndexerSkill BaseSearchIndexerSkillImpl
	Type               string
	Values             map[string]interface{}
}

func (s RawSearchIndexerSkillImpl) SearchIndexerSkill() BaseSearchIndexerSkillImpl {
	return s.searchIndexerSkill
}

func UnmarshalSearchIndexerSkillImplementation(input []byte) (SearchIndexerSkill, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling SearchIndexerSkill into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["@odata.type"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.AzureOpenAIEmbeddingSkill") {
		var out AzureOpenAIEmbeddingSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AzureOpenAIEmbeddingSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Util.ConditionalSkill") {
		var out ConditionalSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ConditionalSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.CustomEntityLookupSkill") {
		var out CustomEntityLookupSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into CustomEntityLookupSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Util.DocumentExtractionSkill") {
		var out DocumentExtractionSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DocumentExtractionSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Util.DocumentIntelligenceLayoutSkill") {
		var out DocumentIntelligenceLayoutSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DocumentIntelligenceLayoutSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.V3.EntityLinkingSkill") {
		var out EntityLinkingSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EntityLinkingSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.EntityRecognitionSkill") {
		var out EntityRecognitionSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EntityRecognitionSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.V3.EntityRecognitionSkill") {
		var out EntityRecognitionSkillV3
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EntityRecognitionSkillV3: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Vision.ImageAnalysisSkill") {
		var out ImageAnalysisSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ImageAnalysisSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.KeyPhraseExtractionSkill") {
		var out KeyPhraseExtractionSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KeyPhraseExtractionSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.LanguageDetectionSkill") {
		var out LanguageDetectionSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LanguageDetectionSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.MergeSkill") {
		var out MergeSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MergeSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Vision.OcrSkill") {
		var out OcrSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OcrSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.PIIDetectionSkill") {
		var out PIIDetectionSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into PIIDetectionSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.SentimentSkill") {
		var out SentimentSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SentimentSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.V3.SentimentSkill") {
		var out SentimentSkillV3
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SentimentSkillV3: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Util.ShaperSkill") {
		var out ShaperSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into ShaperSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.SplitSkill") {
		var out SplitSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SplitSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Text.TranslationSkill") {
		var out TextTranslationSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TextTranslationSkill: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "#Microsoft.Skills.Custom.WebApiSkill") {
		var out WebApiSkill
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebApiSkill: %+v", err)
		}
		return out, nil
	}

	var parent BaseSearchIndexerSkillImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseSearchIndexerSkillImpl: %+v", err)
	}

	return RawSearchIndexerSkillImpl{
		searchIndexerSkill: parent,
		Type:               value,
		Values:             temp,
	}, nil

}

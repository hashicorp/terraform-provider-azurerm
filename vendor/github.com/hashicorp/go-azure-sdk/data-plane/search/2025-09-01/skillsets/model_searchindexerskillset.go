package skillsets

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndexerSkillset struct {
	CognitiveServices CognitiveServicesAccount       `json:"cognitiveServices"`
	Description       *string                        `json:"description,omitempty"`
	EncryptionKey     *SearchResourceEncryptionKey   `json:"encryptionKey,omitempty"`
	IndexProjections  *SearchIndexerIndexProjections `json:"indexProjections,omitempty"`
	KnowledgeStore    *SearchIndexerKnowledgeStore   `json:"knowledgeStore,omitempty"`
	Name              string                         `json:"name"`
	OdataEtag         *string                        `json:"@odata.etag,omitempty"`
	Skills            []SearchIndexerSkill           `json:"skills"`
}

var _ json.Unmarshaler = &SearchIndexerSkillset{}

func (s *SearchIndexerSkillset) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Description      *string                        `json:"description,omitempty"`
		EncryptionKey    *SearchResourceEncryptionKey   `json:"encryptionKey,omitempty"`
		IndexProjections *SearchIndexerIndexProjections `json:"indexProjections,omitempty"`
		KnowledgeStore   *SearchIndexerKnowledgeStore   `json:"knowledgeStore,omitempty"`
		Name             string                         `json:"name"`
		OdataEtag        *string                        `json:"@odata.etag,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Description = decoded.Description
	s.EncryptionKey = decoded.EncryptionKey
	s.IndexProjections = decoded.IndexProjections
	s.KnowledgeStore = decoded.KnowledgeStore
	s.Name = decoded.Name
	s.OdataEtag = decoded.OdataEtag

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SearchIndexerSkillset into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["cognitiveServices"]; ok {
		impl, err := UnmarshalCognitiveServicesAccountImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'CognitiveServices' for 'SearchIndexerSkillset': %+v", err)
		}
		s.CognitiveServices = impl
	}

	if v, ok := temp["skills"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Skills into list []json.RawMessage: %+v", err)
		}

		output := make([]SearchIndexerSkill, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalSearchIndexerSkillImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Skills' for 'SearchIndexerSkillset': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Skills = output
	}

	return nil
}

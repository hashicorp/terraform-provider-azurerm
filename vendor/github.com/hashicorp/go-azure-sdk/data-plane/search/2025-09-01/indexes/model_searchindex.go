package indexes

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchIndex struct {
	Analyzers             *[]LexicalAnalyzer           `json:"analyzers,omitempty"`
	CharFilters           *[]CharFilter                `json:"charFilters,omitempty"`
	CorsOptions           *CorsOptions                 `json:"corsOptions,omitempty"`
	DefaultScoringProfile *string                      `json:"defaultScoringProfile,omitempty"`
	Description           *string                      `json:"description,omitempty"`
	EncryptionKey         *SearchResourceEncryptionKey `json:"encryptionKey,omitempty"`
	Fields                []SearchField                `json:"fields"`
	Name                  string                       `json:"name"`
	Normalizers           *[]LexicalNormalizer         `json:"normalizers,omitempty"`
	OdataEtag             *string                      `json:"@odata.etag,omitempty"`
	ScoringProfiles       *[]ScoringProfile            `json:"scoringProfiles,omitempty"`
	Semantic              *SemanticSettings            `json:"semantic,omitempty"`
	Similarity            Similarity                   `json:"similarity"`
	Suggesters            *[]Suggester                 `json:"suggesters,omitempty"`
	TokenFilters          *[]TokenFilter               `json:"tokenFilters,omitempty"`
	Tokenizers            *[]LexicalTokenizer          `json:"tokenizers,omitempty"`
	VectorSearch          *VectorSearch                `json:"vectorSearch,omitempty"`
}

var _ json.Unmarshaler = &SearchIndex{}

func (s *SearchIndex) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		CorsOptions           *CorsOptions                 `json:"corsOptions,omitempty"`
		DefaultScoringProfile *string                      `json:"defaultScoringProfile,omitempty"`
		Description           *string                      `json:"description,omitempty"`
		EncryptionKey         *SearchResourceEncryptionKey `json:"encryptionKey,omitempty"`
		Fields                []SearchField                `json:"fields"`
		Name                  string                       `json:"name"`
		OdataEtag             *string                      `json:"@odata.etag,omitempty"`
		ScoringProfiles       *[]ScoringProfile            `json:"scoringProfiles,omitempty"`
		Semantic              *SemanticSettings            `json:"semantic,omitempty"`
		Suggesters            *[]Suggester                 `json:"suggesters,omitempty"`
		VectorSearch          *VectorSearch                `json:"vectorSearch,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.CorsOptions = decoded.CorsOptions
	s.DefaultScoringProfile = decoded.DefaultScoringProfile
	s.Description = decoded.Description
	s.EncryptionKey = decoded.EncryptionKey
	s.Fields = decoded.Fields
	s.Name = decoded.Name
	s.OdataEtag = decoded.OdataEtag
	s.ScoringProfiles = decoded.ScoringProfiles
	s.Semantic = decoded.Semantic
	s.Suggesters = decoded.Suggesters
	s.VectorSearch = decoded.VectorSearch

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SearchIndex into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["analyzers"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Analyzers into list []json.RawMessage: %+v", err)
		}

		output := make([]LexicalAnalyzer, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalLexicalAnalyzerImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Analyzers' for 'SearchIndex': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Analyzers = &output
	}

	if v, ok := temp["charFilters"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling CharFilters into list []json.RawMessage: %+v", err)
		}

		output := make([]CharFilter, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalCharFilterImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'CharFilters' for 'SearchIndex': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.CharFilters = &output
	}

	if v, ok := temp["normalizers"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Normalizers into list []json.RawMessage: %+v", err)
		}

		output := make([]LexicalNormalizer, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalLexicalNormalizerImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Normalizers' for 'SearchIndex': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Normalizers = &output
	}

	if v, ok := temp["similarity"]; ok {
		impl, err := UnmarshalSimilarityImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'Similarity' for 'SearchIndex': %+v", err)
		}
		s.Similarity = impl
	}

	if v, ok := temp["tokenFilters"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling TokenFilters into list []json.RawMessage: %+v", err)
		}

		output := make([]TokenFilter, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalTokenFilterImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'TokenFilters' for 'SearchIndex': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.TokenFilters = &output
	}

	if v, ok := temp["tokenizers"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling Tokenizers into list []json.RawMessage: %+v", err)
		}

		output := make([]LexicalTokenizer, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalLexicalTokenizerImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'Tokenizers' for 'SearchIndex': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.Tokenizers = &output
	}

	return nil
}

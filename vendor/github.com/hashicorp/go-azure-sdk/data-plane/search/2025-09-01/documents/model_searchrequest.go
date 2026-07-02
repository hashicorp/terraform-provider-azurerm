package documents

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchRequest struct {
	Answers                       *QueryAnswerType   `json:"answers,omitempty"`
	Captions                      *QueryCaptionType  `json:"captions,omitempty"`
	Count                         *bool              `json:"count,omitempty"`
	Debug                         *QueryDebugMode    `json:"debug,omitempty"`
	Facets                        *[]string          `json:"facets,omitempty"`
	Filter                        *string            `json:"filter,omitempty"`
	Highlight                     *string            `json:"highlight,omitempty"`
	HighlightPostTag              *string            `json:"highlightPostTag,omitempty"`
	HighlightPreTag               *string            `json:"highlightPreTag,omitempty"`
	MinimumCoverage               *float64           `json:"minimumCoverage,omitempty"`
	Orderby                       *string            `json:"orderby,omitempty"`
	QueryType                     *QueryType         `json:"queryType,omitempty"`
	ScoringParameters             *[]string          `json:"scoringParameters,omitempty"`
	ScoringProfile                *string            `json:"scoringProfile,omitempty"`
	ScoringStatistics             *ScoringStatistics `json:"scoringStatistics,omitempty"`
	Search                        *string            `json:"search,omitempty"`
	SearchFields                  *string            `json:"searchFields,omitempty"`
	SearchMode                    *SearchMode        `json:"searchMode,omitempty"`
	Select                        *string            `json:"select,omitempty"`
	SemanticConfiguration         *string            `json:"semanticConfiguration,omitempty"`
	SemanticErrorHandling         *SemanticErrorMode `json:"semanticErrorHandling,omitempty"`
	SemanticMaxWaitInMilliseconds *int64             `json:"semanticMaxWaitInMilliseconds,omitempty"`
	SemanticQuery                 *string            `json:"semanticQuery,omitempty"`
	SessionId                     *string            `json:"sessionId,omitempty"`
	Skip                          *int64             `json:"skip,omitempty"`
	Top                           *int64             `json:"top,omitempty"`
	VectorFilterMode              *VectorFilterMode  `json:"vectorFilterMode,omitempty"`
	VectorQueries                 *[]VectorQuery     `json:"vectorQueries,omitempty"`
}

var _ json.Unmarshaler = &SearchRequest{}

func (s *SearchRequest) UnmarshalJSON(bytes []byte) error {
	var decoded struct {
		Answers                       *QueryAnswerType   `json:"answers,omitempty"`
		Captions                      *QueryCaptionType  `json:"captions,omitempty"`
		Count                         *bool              `json:"count,omitempty"`
		Debug                         *QueryDebugMode    `json:"debug,omitempty"`
		Facets                        *[]string          `json:"facets,omitempty"`
		Filter                        *string            `json:"filter,omitempty"`
		Highlight                     *string            `json:"highlight,omitempty"`
		HighlightPostTag              *string            `json:"highlightPostTag,omitempty"`
		HighlightPreTag               *string            `json:"highlightPreTag,omitempty"`
		MinimumCoverage               *float64           `json:"minimumCoverage,omitempty"`
		Orderby                       *string            `json:"orderby,omitempty"`
		QueryType                     *QueryType         `json:"queryType,omitempty"`
		ScoringParameters             *[]string          `json:"scoringParameters,omitempty"`
		ScoringProfile                *string            `json:"scoringProfile,omitempty"`
		ScoringStatistics             *ScoringStatistics `json:"scoringStatistics,omitempty"`
		Search                        *string            `json:"search,omitempty"`
		SearchFields                  *string            `json:"searchFields,omitempty"`
		SearchMode                    *SearchMode        `json:"searchMode,omitempty"`
		Select                        *string            `json:"select,omitempty"`
		SemanticConfiguration         *string            `json:"semanticConfiguration,omitempty"`
		SemanticErrorHandling         *SemanticErrorMode `json:"semanticErrorHandling,omitempty"`
		SemanticMaxWaitInMilliseconds *int64             `json:"semanticMaxWaitInMilliseconds,omitempty"`
		SemanticQuery                 *string            `json:"semanticQuery,omitempty"`
		SessionId                     *string            `json:"sessionId,omitempty"`
		Skip                          *int64             `json:"skip,omitempty"`
		Top                           *int64             `json:"top,omitempty"`
		VectorFilterMode              *VectorFilterMode  `json:"vectorFilterMode,omitempty"`
	}
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}

	s.Answers = decoded.Answers
	s.Captions = decoded.Captions
	s.Count = decoded.Count
	s.Debug = decoded.Debug
	s.Facets = decoded.Facets
	s.Filter = decoded.Filter
	s.Highlight = decoded.Highlight
	s.HighlightPostTag = decoded.HighlightPostTag
	s.HighlightPreTag = decoded.HighlightPreTag
	s.MinimumCoverage = decoded.MinimumCoverage
	s.Orderby = decoded.Orderby
	s.QueryType = decoded.QueryType
	s.ScoringParameters = decoded.ScoringParameters
	s.ScoringProfile = decoded.ScoringProfile
	s.ScoringStatistics = decoded.ScoringStatistics
	s.Search = decoded.Search
	s.SearchFields = decoded.SearchFields
	s.SearchMode = decoded.SearchMode
	s.Select = decoded.Select
	s.SemanticConfiguration = decoded.SemanticConfiguration
	s.SemanticErrorHandling = decoded.SemanticErrorHandling
	s.SemanticMaxWaitInMilliseconds = decoded.SemanticMaxWaitInMilliseconds
	s.SemanticQuery = decoded.SemanticQuery
	s.SessionId = decoded.SessionId
	s.Skip = decoded.Skip
	s.Top = decoded.Top
	s.VectorFilterMode = decoded.VectorFilterMode

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling SearchRequest into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["vectorQueries"]; ok {
		var listTemp []json.RawMessage
		if err := json.Unmarshal(v, &listTemp); err != nil {
			return fmt.Errorf("unmarshaling VectorQueries into list []json.RawMessage: %+v", err)
		}

		output := make([]VectorQuery, 0)
		for i, val := range listTemp {
			impl, err := UnmarshalVectorQueryImplementation(val)
			if err != nil {
				return fmt.Errorf("unmarshaling index %d field 'VectorQueries' for 'SearchRequest': %+v", i, err)
			}
			output = append(output, impl)
		}
		s.VectorQueries = &output
	}

	return nil
}

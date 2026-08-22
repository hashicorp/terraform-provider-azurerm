package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchDocumentsResult struct {
	OdataCount                          *int64                     `json:"@odata.count,omitempty"`
	OdataNextLink                       *string                    `json:"@odata.nextLink,omitempty"`
	SearchAnswers                       *[]AnswerResult            `json:"@search.answers,omitempty"`
	SearchCoverage                      *float64                   `json:"@search.coverage,omitempty"`
	SearchFacets                        *map[string][]FacetResult  `json:"@search.facets,omitempty"`
	SearchNextPageParameters            *SearchRequest             `json:"@search.nextPageParameters,omitempty"`
	SearchSemanticPartialResponseReason *SemanticErrorReason       `json:"@search.semanticPartialResponseReason,omitempty"`
	SearchSemanticPartialResponseType   *SemanticSearchResultsType `json:"@search.semanticPartialResponseType,omitempty"`
	Value                               []SearchResult             `json:"value"`
}

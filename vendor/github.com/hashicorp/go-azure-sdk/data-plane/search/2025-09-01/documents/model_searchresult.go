package documents

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchResult struct {
	SearchCaptions             *[]CaptionResult     `json:"@search.captions,omitempty"`
	SearchDocumentDebugInfo    *DocumentDebugInfo   `json:"@search.documentDebugInfo,omitempty"`
	SearchHighlights           *map[string][]string `json:"@search.highlights,omitempty"`
	SearchRerankerBoostedScore *float64             `json:"@search.rerankerBoostedScore,omitempty"`
	SearchRerankerScore        *float64             `json:"@search.rerankerScore,omitempty"`
	SearchScore                float64              `json:"@search.score"`
}

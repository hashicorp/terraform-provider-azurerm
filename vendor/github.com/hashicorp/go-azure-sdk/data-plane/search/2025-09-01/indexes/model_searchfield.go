package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchField struct {
	Analyzer            *LexicalAnalyzerName   `json:"analyzer,omitempty"`
	Dimensions          *int64                 `json:"dimensions,omitempty"`
	Facetable           *bool                  `json:"facetable,omitempty"`
	Fields              *[]SearchField         `json:"fields,omitempty"`
	Filterable          *bool                  `json:"filterable,omitempty"`
	IndexAnalyzer       *LexicalAnalyzerName   `json:"indexAnalyzer,omitempty"`
	Key                 *bool                  `json:"key,omitempty"`
	Name                string                 `json:"name"`
	Normalizer          *LexicalNormalizerName `json:"normalizer,omitempty"`
	Retrievable         *bool                  `json:"retrievable,omitempty"`
	SearchAnalyzer      *LexicalAnalyzerName   `json:"searchAnalyzer,omitempty"`
	Searchable          *bool                  `json:"searchable,omitempty"`
	Sortable            *bool                  `json:"sortable,omitempty"`
	Stored              *bool                  `json:"stored,omitempty"`
	SynonymMaps         *[]string              `json:"synonymMaps,omitempty"`
	Type                SearchFieldDataType    `json:"type"`
	VectorEncoding      *VectorEncodingFormat  `json:"vectorEncoding,omitempty"`
	VectorSearchProfile *string                `json:"vectorSearchProfile,omitempty"`
}

package indexes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalyzeRequest struct {
	Analyzer     *LexicalAnalyzerName   `json:"analyzer,omitempty"`
	CharFilters  *[]CharFilterName      `json:"charFilters,omitempty"`
	Normalizer   *LexicalNormalizerName `json:"normalizer,omitempty"`
	Text         string                 `json:"text"`
	TokenFilters *[]TokenFilterName     `json:"tokenFilters,omitempty"`
	Tokenizer    *LexicalTokenizerName  `json:"tokenizer,omitempty"`
}

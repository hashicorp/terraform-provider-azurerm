package threatintelligence

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceFilteringCriteria struct {
	Ids             *[]string                            `json:"ids,omitempty"`
	IncludeDisabled *bool                                `json:"includeDisabled,omitempty"`
	Keywords        *[]string                            `json:"keywords,omitempty"`
	MaxConfidence   *int64                               `json:"maxConfidence,omitempty"`
	MaxValidUntil   *string                              `json:"maxValidUntil,omitempty"`
	MinConfidence   *int64                               `json:"minConfidence,omitempty"`
	MinValidUntil   *string                              `json:"minValidUntil,omitempty"`
	PageSize        *int64                               `json:"pageSize,omitempty"`
	PatternTypes    *[]string                            `json:"patternTypes,omitempty"`
	SkipToken       *string                              `json:"skipToken,omitempty"`
	SortBy          *[]ThreatIntelligenceSortingCriteria `json:"sortBy,omitempty"`
	Sources         *[]string                            `json:"sources,omitempty"`
	ThreatTypes     *[]string                            `json:"threatTypes,omitempty"`
}

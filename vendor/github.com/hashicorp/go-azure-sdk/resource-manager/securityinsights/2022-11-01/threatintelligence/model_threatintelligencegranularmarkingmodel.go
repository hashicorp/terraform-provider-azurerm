package threatintelligence

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceGranularMarkingModel struct {
	Language   *string   `json:"language,omitempty"`
	MarkingRef *int64    `json:"markingRef,omitempty"`
	Selectors  *[]string `json:"selectors,omitempty"`
}

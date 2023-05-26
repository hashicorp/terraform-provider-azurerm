package threatintelligence

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceSortingCriteria struct {
	ItemKey   *string                         `json:"itemKey,omitempty"`
	SortOrder *ThreatIntelligenceSortingOrder `json:"sortOrder,omitempty"`
}

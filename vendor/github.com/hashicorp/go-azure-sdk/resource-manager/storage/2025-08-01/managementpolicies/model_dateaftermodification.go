package managementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DateAfterModification struct {
	DaysAfterCreationGreaterThan       *float64 `json:"daysAfterCreationGreaterThan,omitempty"`
	DaysAfterLastAccessTimeGreaterThan *float64 `json:"daysAfterLastAccessTimeGreaterThan,omitempty"`
	DaysAfterLastTierChangeGreaterThan *float64 `json:"daysAfterLastTierChangeGreaterThan,omitempty"`
	DaysAfterModificationGreaterThan   *float64 `json:"daysAfterModificationGreaterThan,omitempty"`
}

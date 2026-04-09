package managementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DateAfterCreation struct {
	DaysAfterCreationGreaterThan       float64  `json:"daysAfterCreationGreaterThan"`
	DaysAfterLastTierChangeGreaterThan *float64 `json:"daysAfterLastTierChangeGreaterThan,omitempty"`
}

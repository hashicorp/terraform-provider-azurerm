package managedenvironments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkloadProfileStatesProperties struct {
	CurrentCount *int64 `json:"currentCount,omitempty"`
	MaximumCount *int64 `json:"maximumCount,omitempty"`
	MinimumCount *int64 `json:"minimumCount,omitempty"`
}

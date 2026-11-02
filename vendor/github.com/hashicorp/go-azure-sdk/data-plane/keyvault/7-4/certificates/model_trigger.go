package certificates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Trigger struct {
	DaysBeforeExpiry   *int64 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int64 `json:"lifetime_percentage,omitempty"`
}

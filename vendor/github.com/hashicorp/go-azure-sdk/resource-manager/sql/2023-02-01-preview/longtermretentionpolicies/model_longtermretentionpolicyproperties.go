package longtermretentionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LongTermRetentionPolicyProperties struct {
	MakeBackupsImmutable *bool   `json:"makeBackupsImmutable,omitempty"`
	MonthlyRetention     *string `json:"monthlyRetention,omitempty"`
	WeekOfYear           *int64  `json:"weekOfYear,omitempty"`
	WeeklyRetention      *string `json:"weeklyRetention,omitempty"`
	YearlyRetention      *string `json:"yearlyRetention,omitempty"`
}

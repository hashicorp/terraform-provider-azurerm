package commitmentplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommitmentPeriod struct {
	Count     *int64           `json:"count,omitempty"`
	EndDate   *string          `json:"endDate,omitempty"`
	Quota     *CommitmentQuota `json:"quota,omitempty"`
	StartDate *string          `json:"startDate,omitempty"`
	Tier      *string          `json:"tier,omitempty"`
}

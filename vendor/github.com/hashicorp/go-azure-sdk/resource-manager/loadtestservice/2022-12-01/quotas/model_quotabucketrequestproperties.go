package quotas

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaBucketRequestProperties struct {
	CurrentQuota *int64                                  `json:"currentQuota,omitempty"`
	CurrentUsage *int64                                  `json:"currentUsage,omitempty"`
	Dimensions   *QuotaBucketRequestPropertiesDimensions `json:"dimensions,omitempty"`
	NewQuota     *int64                                  `json:"newQuota,omitempty"`
}

package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaReportFilterRequest struct {
	QuotaTarget              *string    `json:"quotaTarget,omitempty"`
	QuotaType                *QuotaType `json:"quotaType,omitempty"`
	UsageThresholdPercentage *int64     `json:"usageThresholdPercentage,omitempty"`
}

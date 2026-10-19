package volumes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QuotaReport struct {
	IsDerivedQuota        *bool      `json:"isDerivedQuota,omitempty"`
	PercentageUsed        *float64   `json:"percentageUsed,omitempty"`
	QuotaLimitTotalInKiBs *int64     `json:"quotaLimitTotalInKiBs,omitempty"`
	QuotaLimitUsedInKiBs  *int64     `json:"quotaLimitUsedInKiBs,omitempty"`
	QuotaTarget           *string    `json:"quotaTarget,omitempty"`
	QuotaType             *QuotaType `json:"quotaType,omitempty"`
}

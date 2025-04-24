package deletedworkspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceCapping struct {
	DailyQuotaGb        *float64             `json:"dailyQuotaGb,omitempty"`
	DataIngestionStatus *DataIngestionStatus `json:"dataIngestionStatus,omitempty"`
	QuotaNextResetTime  *string              `json:"quotaNextResetTime,omitempty"`
}

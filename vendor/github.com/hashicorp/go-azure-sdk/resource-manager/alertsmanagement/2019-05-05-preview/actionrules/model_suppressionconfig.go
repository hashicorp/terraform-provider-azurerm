package actionrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuppressionConfig struct {
	RecurrenceType SuppressionType      `json:"recurrenceType"`
	Schedule       *SuppressionSchedule `json:"schedule,omitempty"`
}

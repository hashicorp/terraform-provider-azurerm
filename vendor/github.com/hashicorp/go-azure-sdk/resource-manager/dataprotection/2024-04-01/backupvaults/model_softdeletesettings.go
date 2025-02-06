package backupvaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SoftDeleteSettings struct {
	RetentionDurationInDays *float64         `json:"retentionDurationInDays,omitempty"`
	State                   *SoftDeleteState `json:"state,omitempty"`
}

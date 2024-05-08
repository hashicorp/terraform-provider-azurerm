package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RangerProfile struct {
	RangerAdmin    RangerAdminSpec    `json:"rangerAdmin"`
	RangerAudit    *RangerAuditSpec   `json:"rangerAudit,omitempty"`
	RangerUsersync RangerUsersyncSpec `json:"rangerUsersync"`
}

package networksecurityperimeterconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningIssueProperties struct {
	Description *string    `json:"description,omitempty"`
	IssueType   *IssueType `json:"issueType,omitempty"`
	Severity    *Severity  `json:"severity,omitempty"`
}

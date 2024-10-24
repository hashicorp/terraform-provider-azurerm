package connectionmonitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityIssue struct {
	Context  *[]map[string]string `json:"context,omitempty"`
	Origin   *Origin              `json:"origin,omitempty"`
	Severity *Severity            `json:"severity,omitempty"`
	Type     *IssueType           `json:"type,omitempty"`
}

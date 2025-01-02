package networksecurityperimeterconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningIssueProperties struct {
	Description          *string                    `json:"description,omitempty"`
	IssueType            *ProvisioningIssueType     `json:"issueType,omitempty"`
	Severity             *ProvisioningIssueSeverity `json:"severity,omitempty"`
	SuggestedAccessRules *[]AccessRule              `json:"suggestedAccessRules,omitempty"`
	SuggestedResourceIds *[]string                  `json:"suggestedResourceIds,omitempty"`
}

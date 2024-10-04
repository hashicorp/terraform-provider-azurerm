package networksecurityperimeterconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisioningIssue struct {
	Name       *string                      `json:"name,omitempty"`
	Properties *ProvisioningIssueProperties `json:"properties,omitempty"`
}

package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SSODetailsResponse struct {
	AadDomains      *[]string  `json:"aadDomains,omitempty"`
	AdminUsers      *[]string  `json:"adminUsers,omitempty"`
	IsSsoEnabled    *SSOStatus `json:"isSsoEnabled,omitempty"`
	MetadataURL     *string    `json:"metadataUrl,omitempty"`
	SingleSignOnURL *string    `json:"singleSignOnUrl,omitempty"`
}

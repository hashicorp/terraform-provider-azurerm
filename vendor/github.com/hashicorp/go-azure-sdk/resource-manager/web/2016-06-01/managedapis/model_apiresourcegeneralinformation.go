package managedapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiResourceGeneralInformation struct {
	Description   *string `json:"description,omitempty"`
	DisplayName   *string `json:"displayName,omitempty"`
	IconURL       *string `json:"iconUrl,omitempty"`
	ReleaseTag    *string `json:"releaseTag,omitempty"`
	TermsOfUseURL *string `json:"termsOfUseUrl,omitempty"`
}

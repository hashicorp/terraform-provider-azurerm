package signupsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TermsOfServiceProperties struct {
	ConsentRequired *bool   `json:"consentRequired,omitempty"`
	Enabled         *bool   `json:"enabled,omitempty"`
	Text            *string `json:"text,omitempty"`
}

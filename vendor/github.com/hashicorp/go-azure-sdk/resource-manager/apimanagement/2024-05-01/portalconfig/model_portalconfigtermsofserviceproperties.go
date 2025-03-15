package portalconfig

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PortalConfigTermsOfServiceProperties struct {
	RequireConsent *bool   `json:"requireConsent,omitempty"`
	Text           *string `json:"text,omitempty"`
}

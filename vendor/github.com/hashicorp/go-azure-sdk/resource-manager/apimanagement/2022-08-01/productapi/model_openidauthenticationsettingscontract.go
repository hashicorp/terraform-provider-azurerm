package productapi

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OpenIdAuthenticationSettingsContract struct {
	BearerTokenSendingMethods *[]BearerTokenSendingMethods `json:"bearerTokenSendingMethods,omitempty"`
	OpenidProviderId          *string                      `json:"openidProviderId,omitempty"`
}

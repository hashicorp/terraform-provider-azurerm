package channel

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SlackChannelProperties struct {
	ClientId                *string `json:"clientId,omitempty"`
	ClientSecret            *string `json:"clientSecret,omitempty"`
	IsEnabled               bool    `json:"isEnabled"`
	IsValidated             *bool   `json:"IsValidated,omitempty"`
	LandingPageURL          *string `json:"landingPageUrl,omitempty"`
	LastSubmissionId        *string `json:"lastSubmissionId,omitempty"`
	RedirectAction          *string `json:"redirectAction,omitempty"`
	RegisterBeforeOAuthFlow *bool   `json:"registerBeforeOAuthFlow,omitempty"`
	Scopes                  *string `json:"scopes,omitempty"`
	SigningSecret           *string `json:"signingSecret,omitempty"`
	VerificationToken       *string `json:"verificationToken,omitempty"`
}

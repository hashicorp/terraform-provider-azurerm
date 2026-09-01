package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicySettings struct {
	CaptchaExpirationInMinutes             *int64                      `json:"captchaExpirationInMinutes,omitempty"`
	CustomBlockResponseBody                *string                     `json:"customBlockResponseBody,omitempty"`
	CustomBlockResponseStatusCode          *int64                      `json:"customBlockResponseStatusCode,omitempty"`
	EnabledState                           *PolicyEnabledState         `json:"enabledState,omitempty"`
	JavascriptChallengeExpirationInMinutes *int64                      `json:"javascriptChallengeExpirationInMinutes,omitempty"`
	LogScrubbing                           *PolicySettingsLogScrubbing `json:"logScrubbing,omitempty"`
	Mode                                   *PolicyMode                 `json:"mode,omitempty"`
	RedirectURL                            *string                     `json:"redirectUrl,omitempty"`
	RequestBodyCheck                       *PolicyRequestBodyCheck     `json:"requestBodyCheck,omitempty"`
}

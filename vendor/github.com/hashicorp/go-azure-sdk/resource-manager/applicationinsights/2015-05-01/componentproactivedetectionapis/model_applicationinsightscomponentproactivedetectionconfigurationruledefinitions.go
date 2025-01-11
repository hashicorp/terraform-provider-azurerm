package componentproactivedetectionapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentProactiveDetectionConfigurationRuleDefinitions struct {
	Description                *string `json:"Description,omitempty"`
	DisplayName                *string `json:"DisplayName,omitempty"`
	HelpURL                    *string `json:"HelpUrl,omitempty"`
	IsEnabledByDefault         *bool   `json:"IsEnabledByDefault,omitempty"`
	IsHidden                   *bool   `json:"IsHidden,omitempty"`
	IsInPreview                *bool   `json:"IsInPreview,omitempty"`
	Name                       *string `json:"Name,omitempty"`
	SupportsEmailNotifications *bool   `json:"SupportsEmailNotifications,omitempty"`
}

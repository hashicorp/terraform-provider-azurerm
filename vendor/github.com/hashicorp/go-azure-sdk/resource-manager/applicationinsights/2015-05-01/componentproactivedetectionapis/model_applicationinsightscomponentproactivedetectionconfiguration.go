package componentproactivedetectionapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentProactiveDetectionConfiguration struct {
	CustomEmails                   *[]string                                                                   `json:"CustomEmails,omitempty"`
	Enabled                        *bool                                                                       `json:"Enabled,omitempty"`
	LastUpdatedTime                *string                                                                     `json:"LastUpdatedTime,omitempty"`
	Name                           *string                                                                     `json:"Name,omitempty"`
	RuleDefinitions                *ApplicationInsightsComponentProactiveDetectionConfigurationRuleDefinitions `json:"RuleDefinitions,omitempty"`
	SendEmailsToSubscriptionOwners *bool                                                                       `json:"SendEmailsToSubscriptionOwners,omitempty"`
}

package componentproactivedetectionapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationInsightsComponentProactiveDetectionConfiguration struct {
	CustomEmails                   *[]string                                                                   `json:"customEmails,omitempty"`
	Enabled                        *bool                                                                       `json:"enabled,omitempty"`
	LastUpdatedTime                *string                                                                     `json:"lastUpdatedTime,omitempty"`
	Name                           *string                                                                     `json:"name,omitempty"`
	RuleDefinitions                *ApplicationInsightsComponentProactiveDetectionConfigurationRuleDefinitions `json:"ruleDefinitions,omitempty"`
	SendEmailsToSubscriptionOwners *bool                                                                       `json:"sendEmailsToSubscriptionOwners,omitempty"`
}

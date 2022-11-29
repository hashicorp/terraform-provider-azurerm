package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledQueryRuleProperties struct {
	Actions                               *Actions                    `json:"actions,omitempty"`
	AutoMitigate                          *bool                       `json:"autoMitigate,omitempty"`
	CheckWorkspaceAlertsStorageConfigured *bool                       `json:"checkWorkspaceAlertsStorageConfigured,omitempty"`
	CreatedWithApiVersion                 *string                     `json:"createdWithApiVersion,omitempty"`
	Criteria                              *ScheduledQueryRuleCriteria `json:"criteria,omitempty"`
	Description                           *string                     `json:"description,omitempty"`
	DisplayName                           *string                     `json:"displayName,omitempty"`
	Enabled                               *bool                       `json:"enabled,omitempty"`
	EvaluationFrequency                   *string                     `json:"evaluationFrequency,omitempty"`
	IsLegacyLogAnalyticsRule              *bool                       `json:"isLegacyLogAnalyticsRule,omitempty"`
	IsWorkspaceAlertsStorageConfigured    *bool                       `json:"isWorkspaceAlertsStorageConfigured,omitempty"`
	MuteActionsDuration                   *string                     `json:"muteActionsDuration,omitempty"`
	OverrideQueryTimeRange                *string                     `json:"overrideQueryTimeRange,omitempty"`
	Scopes                                *[]string                   `json:"scopes,omitempty"`
	Severity                              *AlertSeverity              `json:"severity,omitempty"`
	SkipQueryValidation                   *bool                       `json:"skipQueryValidation,omitempty"`
	TargetResourceTypes                   *[]string                   `json:"targetResourceTypes,omitempty"`
	WindowSize                            *string                     `json:"windowSize,omitempty"`
}

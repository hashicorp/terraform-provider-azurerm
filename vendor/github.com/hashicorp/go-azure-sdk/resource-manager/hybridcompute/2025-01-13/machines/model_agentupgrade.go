package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentUpgrade struct {
	CorrelationId             *string                `json:"correlationId,omitempty"`
	DesiredVersion            *string                `json:"desiredVersion,omitempty"`
	EnableAutomaticUpgrade    *bool                  `json:"enableAutomaticUpgrade,omitempty"`
	LastAttemptDesiredVersion *string                `json:"lastAttemptDesiredVersion,omitempty"`
	LastAttemptMessage        *string                `json:"lastAttemptMessage,omitempty"`
	LastAttemptStatus         *LastAttemptStatusEnum `json:"lastAttemptStatus,omitempty"`
	LastAttemptTimestamp      *string                `json:"lastAttemptTimestamp,omitempty"`
}

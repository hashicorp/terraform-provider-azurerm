package sessionhost

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SessionHostHealthCheckReport struct {
	AdditionalFailureDetails *SessionHostHealthCheckFailureDetails `json:"additionalFailureDetails,omitempty"`
	HealthCheckName          *HealthCheckName                      `json:"healthCheckName,omitempty"`
	HealthCheckResult        *HealthCheckResult                    `json:"healthCheckResult,omitempty"`
}

package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserFacingError struct {
	Code              *string            `json:"code,omitempty"`
	Details           *[]UserFacingError `json:"details,omitempty"`
	InnerError        *InnerError        `json:"innerError,omitempty"`
	IsRetryable       *bool              `json:"isRetryable,omitempty"`
	IsUserError       *bool              `json:"isUserError,omitempty"`
	Message           *string            `json:"message,omitempty"`
	Properties        *map[string]string `json:"properties,omitempty"`
	RecommendedAction *[]string          `json:"recommendedAction,omitempty"`
	Target            *string            `json:"target,omitempty"`
}

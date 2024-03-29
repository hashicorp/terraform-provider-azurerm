package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Error struct {
	Code       *string          `json:"code,omitempty"`
	Details    *[]Error         `json:"details,omitempty"`
	InnerError *ErrorInnerError `json:"innerError,omitempty"`
	Message    *string          `json:"message,omitempty"`
	Target     *string          `json:"target,omitempty"`
}

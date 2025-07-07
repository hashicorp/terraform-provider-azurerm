package account

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ErrorModel struct {
	Code    *string       `json:"code,omitempty"`
	Details *[]ErrorModel `json:"details,omitempty"`
	Message *string       `json:"message,omitempty"`
	Target  *string       `json:"target,omitempty"`
}

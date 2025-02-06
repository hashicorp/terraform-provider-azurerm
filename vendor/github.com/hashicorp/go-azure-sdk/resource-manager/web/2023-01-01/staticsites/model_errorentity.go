package staticsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ErrorEntity struct {
	Code            *string        `json:"code,omitempty"`
	Details         *[]ErrorEntity `json:"details,omitempty"`
	ExtendedCode    *string        `json:"extendedCode,omitempty"`
	InnerErrors     *[]ErrorEntity `json:"innerErrors,omitempty"`
	Message         *string        `json:"message,omitempty"`
	MessageTemplate *string        `json:"messageTemplate,omitempty"`
	Parameters      *[]string      `json:"parameters,omitempty"`
	Target          *string        `json:"target,omitempty"`
}

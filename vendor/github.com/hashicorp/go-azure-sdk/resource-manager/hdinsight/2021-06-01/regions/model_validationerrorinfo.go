package regions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidationErrorInfo struct {
	Code             *string   `json:"code,omitempty"`
	ErrorResource    *string   `json:"errorResource,omitempty"`
	Message          *string   `json:"message,omitempty"`
	MessageArguments *[]string `json:"messageArguments,omitempty"`
}

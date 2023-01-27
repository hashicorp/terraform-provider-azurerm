package workflowrunactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Expression struct {
	Error          *AzureResourceErrorInfo `json:"error,omitempty"`
	Subexpressions *[]Expression           `json:"subexpressions,omitempty"`
	Text           *string                 `json:"text,omitempty"`
	Value          *interface{}            `json:"value,omitempty"`
}

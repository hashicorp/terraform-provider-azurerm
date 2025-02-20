package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFunctionActivityTypeProperties struct {
	Body         *string                     `json:"body,omitempty"`
	FunctionName string                      `json:"functionName"`
	Headers      *map[string]string          `json:"headers,omitempty"`
	Method       AzureFunctionActivityMethod `json:"method"`
}

package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFunctionActivityTypeProperties struct {
	Body         *interface{}                `json:"body,omitempty"`
	FunctionName interface{}                 `json:"functionName"`
	Headers      *map[string]interface{}     `json:"headers,omitempty"`
	Method       AzureFunctionActivityMethod `json:"method"`
}

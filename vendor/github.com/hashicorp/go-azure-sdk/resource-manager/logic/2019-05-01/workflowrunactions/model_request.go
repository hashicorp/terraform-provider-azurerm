package workflowrunactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Request struct {
	Headers *interface{} `json:"headers,omitempty"`
	Method  *string      `json:"method,omitempty"`
	Uri     *string      `json:"uri,omitempty"`
}

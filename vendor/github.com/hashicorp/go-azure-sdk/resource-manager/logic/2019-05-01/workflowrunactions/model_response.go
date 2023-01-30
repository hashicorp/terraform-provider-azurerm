package workflowrunactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Response struct {
	BodyLink   *ContentLink `json:"bodyLink,omitempty"`
	Headers    *interface{} `json:"headers,omitempty"`
	StatusCode *int64       `json:"statusCode,omitempty"`
}

package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkflowHealth struct {
	Error *ErrorEntity        `json:"error,omitempty"`
	State WorkflowHealthState `json:"state"`
}

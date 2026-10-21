package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TaskContainerExecutionInformation struct {
	ContainerId *string `json:"containerId,omitempty"`
	Error       *string `json:"error,omitempty"`
	State       *string `json:"state,omitempty"`
}

package scriptactions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScriptActionExecutionSummary struct {
	InstanceCount *int64  `json:"instanceCount,omitempty"`
	Status        *string `json:"status,omitempty"`
}

package experiments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StepStatus struct {
	Branches *[]BranchStatus `json:"branches,omitempty"`
	Status   *string         `json:"status,omitempty"`
	StepId   *string         `json:"stepId,omitempty"`
	StepName *string         `json:"stepName,omitempty"`
}

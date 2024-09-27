package experiments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BranchStatus struct {
	Actions    *[]ActionStatus `json:"actions,omitempty"`
	BranchId   *string         `json:"branchId,omitempty"`
	BranchName *string         `json:"branchName,omitempty"`
	Status     *string         `json:"status,omitempty"`
}

package replicationrecoveryplans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RecoveryPlanProtectedItem struct {
	Id               *string `json:"id,omitempty"`
	VirtualMachineId *string `json:"virtualMachineId,omitempty"`
}

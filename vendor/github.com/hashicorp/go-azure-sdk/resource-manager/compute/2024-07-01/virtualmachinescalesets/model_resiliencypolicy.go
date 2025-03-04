package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ResiliencyPolicy struct {
	ResilientVMCreationPolicy *ResilientVMCreationPolicy `json:"resilientVMCreationPolicy,omitempty"`
	ResilientVMDeletionPolicy *ResilientVMDeletionPolicy `json:"resilientVMDeletionPolicy,omitempty"`
}

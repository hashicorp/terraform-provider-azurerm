package backupinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PolicyInfo struct {
	PolicyId         string            `json:"policyId"`
	PolicyParameters *PolicyParameters `json:"policyParameters,omitempty"`
	PolicyVersion    *string           `json:"policyVersion,omitempty"`
}

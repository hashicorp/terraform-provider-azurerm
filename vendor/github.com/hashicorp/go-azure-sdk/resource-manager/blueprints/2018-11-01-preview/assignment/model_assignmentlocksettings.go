package assignment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssignmentLockSettings struct {
	ExcludedActions    *[]string           `json:"excludedActions,omitempty"`
	ExcludedPrincipals *[]string           `json:"excludedPrincipals,omitempty"`
	Mode               *AssignmentLockMode `json:"mode,omitempty"`
}

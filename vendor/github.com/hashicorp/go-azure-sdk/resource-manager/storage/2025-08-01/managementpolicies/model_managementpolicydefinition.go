package managementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementPolicyDefinition struct {
	Actions ManagementPolicyAction  `json:"actions"`
	Filters *ManagementPolicyFilter `json:"filters,omitempty"`
}

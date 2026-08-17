package redisenterprise

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccessPolicyAssignmentProperties struct {
	AccessPolicyName  string                               `json:"accessPolicyName"`
	ProvisioningState *ProvisioningState                   `json:"provisioningState,omitempty"`
	User              AccessPolicyAssignmentPropertiesUser `json:"user"`
}

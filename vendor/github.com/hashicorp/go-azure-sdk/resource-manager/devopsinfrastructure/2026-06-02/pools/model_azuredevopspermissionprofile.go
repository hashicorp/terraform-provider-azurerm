package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureDevOpsPermissionProfile struct {
	Groups *[]string                 `json:"groups,omitempty"`
	Kind   AzureDevOpsPermissionType `json:"kind"`
	Users  *[]string                 `json:"users,omitempty"`
}

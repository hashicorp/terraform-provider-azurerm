package managementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchManagementGroupRequest struct {
	DisplayName   *string `json:"displayName,omitempty"`
	ParentGroupId *string `json:"parentGroupId,omitempty"`
}

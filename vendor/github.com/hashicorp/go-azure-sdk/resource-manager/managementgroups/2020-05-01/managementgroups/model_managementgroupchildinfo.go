package managementgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupChildInfo struct {
	Children    *[]ManagementGroupChildInfo `json:"children,omitempty"`
	DisplayName *string                     `json:"displayName,omitempty"`
	Id          *string                     `json:"id,omitempty"`
	Name        *string                     `json:"name,omitempty"`
	Type        *ManagementGroupChildType   `json:"type,omitempty"`
}

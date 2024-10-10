package rolemanagementpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserSet struct {
	Description *string   `json:"description,omitempty"`
	Id          *string   `json:"id,omitempty"`
	IsBackup    *bool     `json:"isBackup,omitempty"`
	UserType    *UserType `json:"userType,omitempty"`
}

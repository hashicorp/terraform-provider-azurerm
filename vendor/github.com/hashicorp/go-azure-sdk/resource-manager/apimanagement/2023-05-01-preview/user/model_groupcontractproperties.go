package user

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupContractProperties struct {
	BuiltIn     *bool      `json:"builtIn,omitempty"`
	Description *string    `json:"description,omitempty"`
	DisplayName string     `json:"displayName"`
	ExternalId  *string    `json:"externalId,omitempty"`
	Type        *GroupType `json:"type,omitempty"`
}

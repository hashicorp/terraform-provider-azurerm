package caches

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NfsAccessRule struct {
	Access         NfsAccessRuleAccess `json:"access"`
	AnonymousGID   *string             `json:"anonymousGID,omitempty"`
	AnonymousUID   *string             `json:"anonymousUID,omitempty"`
	Filter         *string             `json:"filter,omitempty"`
	RootSquash     *bool               `json:"rootSquash,omitempty"`
	Scope          NfsAccessRuleScope  `json:"scope"`
	SubmountAccess *bool               `json:"submountAccess,omitempty"`
	Suid           *bool               `json:"suid,omitempty"`
}

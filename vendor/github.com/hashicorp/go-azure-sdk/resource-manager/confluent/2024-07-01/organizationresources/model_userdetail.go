package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserDetail struct {
	AadEmail          *string `json:"aadEmail,omitempty"`
	EmailAddress      string  `json:"emailAddress"`
	FirstName         *string `json:"firstName,omitempty"`
	LastName          *string `json:"lastName,omitempty"`
	UserPrincipalName *string `json:"userPrincipalName,omitempty"`
}

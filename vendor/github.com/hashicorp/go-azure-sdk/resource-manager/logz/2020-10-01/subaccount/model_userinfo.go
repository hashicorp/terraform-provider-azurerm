package subaccount

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserInfo struct {
	EmailAddress *string `json:"emailAddress,omitempty"`
	FirstName    *string `json:"firstName,omitempty"`
	LastName     *string `json:"lastName,omitempty"`
	PhoneNumber  *string `json:"phoneNumber,omitempty"`
}

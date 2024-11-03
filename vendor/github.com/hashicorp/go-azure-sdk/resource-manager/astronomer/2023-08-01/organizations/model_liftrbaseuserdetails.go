package organizations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LiftrBaseUserDetails struct {
	EmailAddress string  `json:"emailAddress"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	PhoneNumber  *string `json:"phoneNumber,omitempty"`
	Upn          *string `json:"upn,omitempty"`
}

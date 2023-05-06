package orders

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Address struct {
	AddressLine1 *string `json:"addressLine1,omitempty"`
	AddressLine2 *string `json:"addressLine2,omitempty"`
	AddressLine3 *string `json:"addressLine3,omitempty"`
	City         *string `json:"city,omitempty"`
	Country      string  `json:"country"`
	PostalCode   *string `json:"postalCode,omitempty"`
	State        *string `json:"state,omitempty"`
}

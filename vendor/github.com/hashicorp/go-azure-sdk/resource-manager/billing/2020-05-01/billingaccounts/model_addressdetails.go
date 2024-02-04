package billingaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AddressDetails struct {
	AddressLine1 string  `json:"addressLine1"`
	AddressLine2 *string `json:"addressLine2,omitempty"`
	AddressLine3 *string `json:"addressLine3,omitempty"`
	City         *string `json:"city,omitempty"`
	CompanyName  *string `json:"companyName,omitempty"`
	Country      string  `json:"country"`
	District     *string `json:"district,omitempty"`
	Email        *string `json:"email,omitempty"`
	FirstName    *string `json:"firstName,omitempty"`
	LastName     *string `json:"lastName,omitempty"`
	MiddleName   *string `json:"middleName,omitempty"`
	PhoneNumber  *string `json:"phoneNumber,omitempty"`
	PostalCode   *string `json:"postalCode,omitempty"`
	Region       *string `json:"region,omitempty"`
}

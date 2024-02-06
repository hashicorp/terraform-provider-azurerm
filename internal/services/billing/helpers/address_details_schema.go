// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/billing/2020-05-01/billingaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AddressDetails struct {
	AddressLine1 string `tfschema:"address_line_1"`
	AddressLine2 string `tfschema:"address_line_2"`
	AddressLine3 string `tfschema:"address_line_3"`
	City         string `tfschema:"city"`
	CompanyName  string `tfschema:"company_name"`
	Country      string `tfschema:"country"`
	District     string `tfschema:"district"`
	Email        string `tfschema:"email"`
	FirstName    string `tfschema:"first_name"`
	LastName     string `tfschema:"last_name"`
	MiddleName   string `tfschema:"middle_name"`
	PhoneNumber  string `tfschema:"phone_number"`
	PostalCode   string `tfschema:"postal_code"`
	Region       string `tfschema:"region"`
}

func AddressDetailsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"address_line_1": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"address_line_2": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"address_line_3": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"city": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"company_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"country": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"district": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"email": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"first_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"last_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"middle_name": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"phone_number": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"postal_code": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"region": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func FlattenAddressDetails(input *billingaccounts.AddressDetails) []AddressDetails {
	if input == nil {
		return []AddressDetails{}
	}

	result := AddressDetails{}

	if props := input; props != nil {
		result.AddressLine1 = input.AddressLine1
		result.AddressLine2 = pointer.From(input.AddressLine2)
		result.AddressLine3 = pointer.From(input.AddressLine3)
		result.City = pointer.From(input.City)
		result.CompanyName = pointer.From(input.CompanyName)
		result.Country = input.Country
		result.District = pointer.From(input.District)
		result.Email = pointer.From(input.Email)
		result.FirstName = pointer.From(input.FirstName)
		result.LastName = pointer.From(input.LastName)
		result.MiddleName = pointer.From(input.MiddleName)
		result.PhoneNumber = pointer.From(input.PhoneNumber)
		result.PostalCode = pointer.From(input.PostalCode)
		result.Region = pointer.From(input.Region)
	}

	return []AddressDetails{result}
}

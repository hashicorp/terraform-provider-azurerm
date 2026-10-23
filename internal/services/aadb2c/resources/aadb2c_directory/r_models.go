// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package aadb2c_directory

type AadB2cDirectoryModel struct {
	BillingType           string            `tfschema:"billing_type"`
	CountryCode           string            `tfschema:"country_code"`
	DataResidencyLocation string            `tfschema:"data_residency_location"`
	DisplayName           string            `tfschema:"display_name"`
	DomainName            string            `tfschema:"domain_name"`
	EffectiveStartDate    string            `tfschema:"effective_start_date"`
	ResourceGroup         string            `tfschema:"resource_group_name"`
	Sku                   string            `tfschema:"sku_name"`
	Tags                  map[string]string `tfschema:"tags"`
	TenantId              string            `tfschema:"tenant_id"`
}

func (r Resource) ModelObject() any {
	return &AadB2cDirectoryModel{}
}

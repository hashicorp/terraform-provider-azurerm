// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package trusted_signing_account

type TrustedSigningAccountModel struct {
	Name              string            `tfschema:"name"`
	Location          string            `tfschema:"location"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	AccountUri        string            `tfschema:"account_uri"`
	SkuName           string            `tfschema:"sku_name"`
	Tags              map[string]string `tfschema:"tags"`
}

func (m Resource) ModelObject() any {
	return &TrustedSigningAccountModel{}
}

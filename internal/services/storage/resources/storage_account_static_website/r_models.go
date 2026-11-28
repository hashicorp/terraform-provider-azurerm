// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_static_website

type AccountStaticWebsiteResourceModel struct {
	StorageAccountId string `tfschema:"storage_account_id"`
	Error404Document string `tfschema:"error_404_document"`
	IndexDocument    string `tfschema:"index_document"`
}

func (a Resource) ModelObject() any {
	return &AccountStaticWebsiteResourceModel{}
}

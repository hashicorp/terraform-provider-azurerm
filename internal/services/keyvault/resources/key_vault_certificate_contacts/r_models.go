// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package key_vault_certificate_contacts

type KeyVaultCertificateContactsResourceModel struct {
	KeyVaultId string    `tfschema:"key_vault_id"`
	Contact    []Contact `tfschema:"contact"`
}

func (r Resource) ModelObject() any {
	return &KeyVaultCertificateContactsResourceModel{}
}

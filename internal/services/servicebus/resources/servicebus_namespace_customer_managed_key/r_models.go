// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package servicebus_namespace_customer_managed_key

type ServiceBusNamespaceCustomerManagedKeyModel struct {
	NamespaceID                     string `tfschema:"namespace_id"`
	KeyVaultKeyID                   string `tfschema:"key_vault_key_id"`
	InfrastructureEncryptionEnabled bool   `tfschema:"infrastructure_encryption_enabled"`
}

func (r Resource) ModelObject() any {
	return &ServiceBusNamespaceCustomerManagedKeyModel{}
}

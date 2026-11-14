// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_customer_managed_key

type DataFactoryCustomerManagedKeyModel struct {
	DataFactoryID          string `tfschema:"data_factory_id"`
	CustomerManagedKeyID   string `tfschema:"customer_managed_key_id"`
	UserAssignedIdentityID string `tfschema:"user_assigned_identity_id"`
}

func (r Resource) ModelObject() any {
	return &DataFactoryCustomerManagedKeyModel{}
}

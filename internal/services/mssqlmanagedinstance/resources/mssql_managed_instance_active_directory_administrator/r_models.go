// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_active_directory_administrator

type MsSqlManagedInstanceActiveDirectoryAdministratorModel struct {
	ManagedInstanceId         string `tfschema:"managed_instance_id"`
	AzureADAuthenticationOnly bool   `tfschema:"azuread_authentication_only"`
	LoginUsername             string `tfschema:"login_username"`
	ObjectId                  string `tfschema:"object_id"`
	TenantId                  string `tfschema:"tenant_id"`
}

func (r Resource) ModelObject() any {
	return &MsSqlManagedInstanceActiveDirectoryAdministratorModel{}
}

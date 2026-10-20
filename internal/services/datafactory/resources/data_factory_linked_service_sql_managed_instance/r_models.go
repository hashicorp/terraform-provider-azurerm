// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sql_managed_instance

type LinkedServiceSqlManagedInstanceModel struct {
	Name                     string                           `tfschema:"name"`
	DataFactoryID            string                           `tfschema:"data_factory_id"`
	Annotations              []string                         `tfschema:"annotations"`
	ConnectionString         string                           `tfschema:"connection_string"`
	Description              string                           `tfschema:"description"`
	IntegrationRuntimeName   string                           `tfschema:"integration_runtime_name"`
	KeyVaultConnectionString []KeyVaultConnectionStringConfig `tfschema:"key_vault_connection_string"`
	KeyVaultPassword         []KeyVaultPasswordConfig         `tfschema:"key_vault_password"`
	Parameters               map[string]any                   `tfschema:"parameters"`
	ServicePrincipalID       string                           `tfschema:"service_principal_id"`
	ServicePrincipalKey      string                           `tfschema:"service_principal_key"`
	Tenant                   string                           `tfschema:"tenant"`
}

func (r Resource) ModelObject() any {
	return &LinkedServiceSqlManagedInstanceModel{}
}

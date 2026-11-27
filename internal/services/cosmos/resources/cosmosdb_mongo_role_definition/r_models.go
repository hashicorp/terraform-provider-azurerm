// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_role_definition

type CosmosDbMongoRoleDefinitionResourceModel struct {
	CosmosMongoDatabaseId string      `tfschema:"cosmos_mongo_database_id"`
	RoleName              string      `tfschema:"role_name"`
	InheritedRoleNames    []string    `tfschema:"inherited_role_names"`
	Privileges            []Privilege `tfschema:"privilege"`
}

func (r Resource) ModelObject() any {
	return &CosmosDbMongoRoleDefinitionResourceModel{}
}

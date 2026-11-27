// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cosmosdb_mongo_user_definition

type CosmosDbMongoUserDefinitionResourceModel struct {
	CosmosMongoDatabaseId string   `tfschema:"cosmos_mongo_database_id"`
	Username              string   `tfschema:"username"`
	Password              string   `tfschema:"password"`
	InheritedRoleNames    []string `tfschema:"inherited_role_names"`
}

func (r Resource) ModelObject() any {
	return &CosmosDbMongoUserDefinitionResourceModel{}
}

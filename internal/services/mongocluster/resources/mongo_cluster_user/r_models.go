// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mongo_cluster_user

type MongoClusterUserResourceModel struct {
	IdentityProviderType string              `tfschema:"identity_provider_type"`
	MongoClusterId       string              `tfschema:"mongo_cluster_id"`
	ObjectId             string              `tfschema:"object_id"`
	PrincipalType        string              `tfschema:"principal_type"`
	Roles                []DatabaseRoleModel `tfschema:"role"`
}

type DatabaseRoleModel struct {
	Database string `tfschema:"database"`
	Name     string `tfschema:"name"`
}

func (r Resource) ModelObject() any {
	return &MongoClusterUserResourceModel{}
}

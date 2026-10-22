// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server_active_directory_administrator

type MySQLFlexibleServerAdministratorModel struct {
	ServerId   string `tfschema:"server_id"`
	IdentityId string `tfschema:"identity_id"`
	Login      string `tfschema:"login"`
	ObjectId   string `tfschema:"object_id"`
	TenantId   string `tfschema:"tenant_id"`
}

func (r Resource) ModelObject() any {
	return &MySQLFlexibleServerAdministratorModel{}
}

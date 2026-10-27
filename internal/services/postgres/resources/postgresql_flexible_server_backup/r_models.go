// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_backup

func (r Resource) ModelObject() any {
	return &PostgresqlFlexibleServerBackupResourceModel{}
}

type PostgresqlFlexibleServerBackupResourceModel struct {
	Name          string `tfschema:"name"`
	ServerId      string `tfschema:"server_id"`
	CompletedTime string `tfschema:"completed_time"`
}

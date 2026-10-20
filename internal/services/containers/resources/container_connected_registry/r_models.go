// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package container_connected_registry

type ContainerConnectedRegistryModel struct {
	Name                string                   `tfschema:"name"`
	ContainerRegistryId string                   `tfschema:"container_registry_id"`
	ParentRegistryId    string                   `tfschema:"parent_registry_id"`
	SyncTokenId         string                   `tfschema:"sync_token_id"`
	SyncSchedule        string                   `tfschema:"sync_schedule"`
	SyncMessageTTL      string                   `tfschema:"sync_message_ttl"`
	SyncWindow          string                   `tfschema:"sync_window"`
	Mode                string                   `tfschema:"mode"`
	RepoNotifications   []RepositoryNotification `tfschema:"notification"`
	ClientTokenIds      []string                 `tfschema:"client_token_ids"`
	LogLevel            string                   `tfschema:"log_level"`
	AuditLogEnabled     bool                     `tfschema:"audit_log_enabled"`
}

func (r Resource) ModelObject() any {
	return &ContainerConnectedRegistryModel{}
}

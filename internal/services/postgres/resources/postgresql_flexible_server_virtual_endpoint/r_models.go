// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server_virtual_endpoint

type PostgresqlFlexibleServerVirtualEndpointModel struct {
	Name            string `tfschema:"name"`
	SourceServerId  string `tfschema:"source_server_id"`
	ReplicaServerId string `tfschema:"replica_server_id"`
	Type            string `tfschema:"type"`
}

func (r Resource) ModelObject() any {
	return &PostgresqlFlexibleServerVirtualEndpointModel{}
}

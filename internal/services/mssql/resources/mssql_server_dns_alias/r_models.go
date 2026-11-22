// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_server_dns_alias

type ServerDNSAliasModel struct {
	MsSQLServerId string `tfschema:"mssql_server_id"`
	Name          string `tfschema:"name"`
	DNSRecord     string `tfschema:"dns_record"`
}

func (m Resource) ModelObject() any {
	return &ServerDNSAliasModel{}
}

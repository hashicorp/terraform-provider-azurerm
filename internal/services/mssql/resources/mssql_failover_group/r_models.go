// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_failover_group

type MsSqlFailoverGroupModel struct {
	Databases                            []string             `tfschema:"databases"`
	Name                                 string               `tfschema:"name"`
	PartnerServers                       []PartnerServerModel `tfschema:"partner_server"`
	ReadonlyEndpointFailurePolicyEnabled bool                 `tfschema:"readonly_endpoint_failover_policy_enabled"`
	ServerId                             string               `tfschema:"server_id"`
	Tags                                 map[string]string    `tfschema:"tags"`

	ReadWriteEndpointFailurePolicy []ReadWriteEndpointFailurePolicyModel `tfschema:"read_write_endpoint_failover_policy"`
}

type PartnerServerModel struct {
	ID       string `tfschema:"id"`
	Location string `tfschema:"location"`
	Role     string `tfschema:"role"`
}

type ReadWriteEndpointFailurePolicyModel struct {
	GraceMinutes int64  `tfschema:"grace_minutes"`
	Mode         string `tfschema:"mode"`
}

func (r Resource) ModelObject() any {
	return &MsSqlFailoverGroupModel{}
}

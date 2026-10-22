// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance_failover_group

type MsSqlManagedInstanceFailoverGroupModel struct {
	Name                                  string `tfschema:"name"`
	Location                              string `tfschema:"location"`
	ManagedInstanceId                     string `tfschema:"managed_instance_id"`
	PartnerManagedInstanceId              string `tfschema:"partner_managed_instance_id"`
	ReadOnlyEndpointFailoverPolicyEnabled bool   `tfschema:"readonly_endpoint_failover_policy_enabled"`
	SecondaryType                         string `tfschema:"secondary_type"`

	ReadWriteEndpointFailurePolicy []MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel `tfschema:"read_write_endpoint_failover_policy"`

	PartnerRegion []MsSqlManagedInstancePartnerRegionModel `tfschema:"partner_region"`
	Role          string                                   `tfschema:"role"`
}

type MsSqlManagedInstanceReadWriteEndpointFailurePolicyModel struct {
	GraceMinutes int64  `tfschema:"grace_minutes"`
	Mode         string `tfschema:"mode"`
}

type MsSqlManagedInstancePartnerRegionModel struct {
	Location string `tfschema:"location"`
	Role     string `tfschema:"role"`
}

func (r Resource) ModelObject() any {
	return &MsSqlManagedInstanceFailoverGroupModel{}
}

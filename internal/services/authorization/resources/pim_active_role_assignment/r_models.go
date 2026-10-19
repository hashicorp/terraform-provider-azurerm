// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pim_active_role_assignment

type PimActiveRoleAssignmentModel struct {
	RoleDefinitionId string                                `tfschema:"role_definition_id"`
	Scope            string                                `tfschema:"scope"`
	PrincipalId      string                                `tfschema:"principal_id"`
	PrincipalType    string                                `tfschema:"principal_type"`
	Justification    string                                `tfschema:"justification"`
	TicketInfo       []PimActiveRoleAssignmentTicketInfo   `tfschema:"ticket"`
	ScheduleInfo     []PimActiveRoleAssignmentScheduleInfo `tfschema:"schedule"`
}

func (Resource) ModelObject() any {
	return &PimActiveRoleAssignmentModel{}
}

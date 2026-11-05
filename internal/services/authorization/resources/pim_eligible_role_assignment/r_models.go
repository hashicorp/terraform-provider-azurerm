// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pim_eligible_role_assignment

type PimEligibleRoleAssignmentModel struct {
	RoleDefinitionId string                                  `tfschema:"role_definition_id"`
	Scope            string                                  `tfschema:"scope"`
	PrincipalId      string                                  `tfschema:"principal_id"`
	PrincipalType    string                                  `tfschema:"principal_type"`
	Justification    string                                  `tfschema:"justification"`
	TicketInfo       []PimEligibleRoleAssignmentTicketInfo   `tfschema:"ticket"`
	ScheduleInfo     []PimEligibleRoleAssignmentScheduleInfo `tfschema:"schedule"`
	Condition        string                                  `tfschema:"condition"`
	ConditionVersion string                                  `tfschema:"condition_version"`
}

func (Resource) ModelObject() any {
	return &PimEligibleRoleAssignmentModel{}
}

package roleassignmentschedules

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentScheduleProperties struct {
	AssignmentType                  *AssignmentType     `json:"assignmentType,omitempty"`
	Condition                       *string             `json:"condition,omitempty"`
	ConditionVersion                *string             `json:"conditionVersion,omitempty"`
	CreatedOn                       *string             `json:"createdOn,omitempty"`
	EndDateTime                     *string             `json:"endDateTime,omitempty"`
	ExpandedProperties              *ExpandedProperties `json:"expandedProperties,omitempty"`
	LinkedRoleEligibilityScheduleId *string             `json:"linkedRoleEligibilityScheduleId,omitempty"`
	MemberType                      *MemberType         `json:"memberType,omitempty"`
	PrincipalId                     *string             `json:"principalId,omitempty"`
	PrincipalType                   *PrincipalType      `json:"principalType,omitempty"`
	RoleAssignmentScheduleRequestId *string             `json:"roleAssignmentScheduleRequestId,omitempty"`
	RoleDefinitionId                *string             `json:"roleDefinitionId,omitempty"`
	Scope                           *string             `json:"scope,omitempty"`
	StartDateTime                   *string             `json:"startDateTime,omitempty"`
	Status                          *Status             `json:"status,omitempty"`
	UpdatedOn                       *string             `json:"updatedOn,omitempty"`
}

func (o *RoleAssignmentScheduleProperties) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleAssignmentScheduleProperties) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}

func (o *RoleAssignmentScheduleProperties) GetEndDateTimeAsTime() (*time.Time, error) {
	if o.EndDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleAssignmentScheduleProperties) SetEndDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDateTime = &formatted
}

func (o *RoleAssignmentScheduleProperties) GetStartDateTimeAsTime() (*time.Time, error) {
	if o.StartDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleAssignmentScheduleProperties) SetStartDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDateTime = &formatted
}

func (o *RoleAssignmentScheduleProperties) GetUpdatedOnAsTime() (*time.Time, error) {
	if o.UpdatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleAssignmentScheduleProperties) SetUpdatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedOn = &formatted
}

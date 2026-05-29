package roleassignments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleAssignmentProperties struct {
	Condition                          *string        `json:"condition,omitempty"`
	ConditionVersion                   *string        `json:"conditionVersion,omitempty"`
	CreatedBy                          *string        `json:"createdBy,omitempty"`
	CreatedOn                          *string        `json:"createdOn,omitempty"`
	DelegatedManagedIdentityResourceId *string        `json:"delegatedManagedIdentityResourceId,omitempty"`
	Description                        *string        `json:"description,omitempty"`
	PrincipalId                        string         `json:"principalId"`
	PrincipalType                      *PrincipalType `json:"principalType,omitempty"`
	RoleDefinitionId                   string         `json:"roleDefinitionId"`
	Scope                              *string        `json:"scope,omitempty"`
	UpdatedBy                          *string        `json:"updatedBy,omitempty"`
	UpdatedOn                          *string        `json:"updatedOn,omitempty"`
}

func (o *RoleAssignmentProperties) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleAssignmentProperties) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}

func (o *RoleAssignmentProperties) GetUpdatedOnAsTime() (*time.Time, error) {
	if o.UpdatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleAssignmentProperties) SetUpdatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedOn = &formatted
}

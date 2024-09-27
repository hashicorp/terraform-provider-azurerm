package roledefinitions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RoleDefinitionProperties struct {
	AssignableScopes *[]string     `json:"assignableScopes,omitempty"`
	CreatedBy        *string       `json:"createdBy,omitempty"`
	CreatedOn        *string       `json:"createdOn,omitempty"`
	Description      *string       `json:"description,omitempty"`
	Permissions      *[]Permission `json:"permissions,omitempty"`
	RoleName         *string       `json:"roleName,omitempty"`
	Type             *string       `json:"type,omitempty"`
	UpdatedBy        *string       `json:"updatedBy,omitempty"`
	UpdatedOn        *string       `json:"updatedOn,omitempty"`
}

func (o *RoleDefinitionProperties) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleDefinitionProperties) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}

func (o *RoleDefinitionProperties) GetUpdatedOnAsTime() (*time.Time, error) {
	if o.UpdatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *RoleDefinitionProperties) SetUpdatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedOn = &formatted
}

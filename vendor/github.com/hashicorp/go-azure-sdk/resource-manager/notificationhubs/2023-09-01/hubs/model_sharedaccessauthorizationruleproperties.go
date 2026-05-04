package hubs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SharedAccessAuthorizationRuleProperties struct {
	ClaimType    *string        `json:"claimType,omitempty"`
	ClaimValue   *string        `json:"claimValue,omitempty"`
	CreatedTime  *string        `json:"createdTime,omitempty"`
	KeyName      *string        `json:"keyName,omitempty"`
	ModifiedTime *string        `json:"modifiedTime,omitempty"`
	PrimaryKey   *string        `json:"primaryKey,omitempty"`
	Revision     *int64         `json:"revision,omitempty"`
	Rights       []AccessRights `json:"rights"`
	SecondaryKey *string        `json:"secondaryKey,omitempty"`
}

func (o *SharedAccessAuthorizationRuleProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SharedAccessAuthorizationRuleProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

func (o *SharedAccessAuthorizationRuleProperties) GetModifiedTimeAsTime() (*time.Time, error) {
	if o.ModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SharedAccessAuthorizationRuleProperties) SetModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ModifiedTime = &formatted
}

package resources

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenericResourceExpanded struct {
	ChangedTime       *string            `json:"changedTime,omitempty"`
	CreatedTime       *string            `json:"createdTime,omitempty"`
	Id                *string            `json:"id,omitempty"`
	Location          string             `json:"location"`
	Name              *string            `json:"name,omitempty"`
	Plan              *Plan              `json:"plan,omitempty"`
	Properties        *interface{}       `json:"properties,omitempty"`
	ProvisioningState *string            `json:"provisioningState,omitempty"`
	Tags              *map[string]string `json:"tags,omitempty"`
	Type              *string            `json:"type,omitempty"`
}

func (o *GenericResourceExpanded) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GenericResourceExpanded) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *GenericResourceExpanded) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GenericResourceExpanded) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}

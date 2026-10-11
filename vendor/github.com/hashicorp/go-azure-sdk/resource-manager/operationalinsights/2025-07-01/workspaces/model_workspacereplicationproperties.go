package workspaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceReplicationProperties struct {
	CreatedDate       *string                    `json:"createdDate,omitempty"`
	Enabled           *bool                      `json:"enabled,omitempty"`
	LastModifiedDate  *string                    `json:"lastModifiedDate,omitempty"`
	Location          *string                    `json:"location,omitempty"`
	ProvisioningState *WorkspaceReplicationState `json:"provisioningState,omitempty"`
}

func (o *WorkspaceReplicationProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkspaceReplicationProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}

func (o *WorkspaceReplicationProperties) GetLastModifiedDateAsTime() (*time.Time, error) {
	if o.LastModifiedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkspaceReplicationProperties) SetLastModifiedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedDate = &formatted
}

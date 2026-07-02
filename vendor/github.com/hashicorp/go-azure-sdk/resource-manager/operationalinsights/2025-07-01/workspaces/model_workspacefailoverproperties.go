package workspaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceFailoverProperties struct {
	LastModifiedDate *string                 `json:"lastModifiedDate,omitempty"`
	State            *WorkspaceFailoverState `json:"state,omitempty"`
}

func (o *WorkspaceFailoverProperties) GetLastModifiedDateAsTime() (*time.Time, error) {
	if o.LastModifiedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkspaceFailoverProperties) SetLastModifiedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedDate = &formatted
}

package issue

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IssueContractProperties struct {
	ApiId       *string `json:"apiId,omitempty"`
	CreatedDate *string `json:"createdDate,omitempty"`
	Description string  `json:"description"`
	State       *State  `json:"state,omitempty"`
	Title       string  `json:"title"`
	UserId      string  `json:"userId"`
}

func (o *IssueContractProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *IssueContractProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}

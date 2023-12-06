package apirelease

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiReleaseContractProperties struct {
	ApiId           *string `json:"apiId,omitempty"`
	CreatedDateTime *string `json:"createdDateTime,omitempty"`
	Notes           *string `json:"notes,omitempty"`
	UpdatedDateTime *string `json:"updatedDateTime,omitempty"`
}

func (o *ApiReleaseContractProperties) GetCreatedDateTimeAsTime() (*time.Time, error) {
	if o.CreatedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ApiReleaseContractProperties) SetCreatedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDateTime = &formatted
}

func (o *ApiReleaseContractProperties) GetUpdatedDateTimeAsTime() (*time.Time, error) {
	if o.UpdatedDateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedDateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ApiReleaseContractProperties) SetUpdatedDateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedDateTime = &formatted
}

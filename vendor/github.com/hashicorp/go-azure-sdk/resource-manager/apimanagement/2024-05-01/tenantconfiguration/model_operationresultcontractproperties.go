package tenantconfiguration

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationResultContractProperties struct {
	ActionLog  *[]OperationResultLogItemContract `json:"actionLog,omitempty"`
	Error      *ErrorResponseBody                `json:"error,omitempty"`
	Id         *string                           `json:"id,omitempty"`
	ResultInfo *string                           `json:"resultInfo,omitempty"`
	Started    *string                           `json:"started,omitempty"`
	Status     *AsyncOperationStatus             `json:"status,omitempty"`
	Updated    *string                           `json:"updated,omitempty"`
}

func (o *OperationResultContractProperties) GetStartedAsTime() (*time.Time, error) {
	if o.Started == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Started, "2006-01-02T15:04:05Z07:00")
}

func (o *OperationResultContractProperties) SetStartedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Started = &formatted
}

func (o *OperationResultContractProperties) GetUpdatedAsTime() (*time.Time, error) {
	if o.Updated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Updated, "2006-01-02T15:04:05Z07:00")
}

func (o *OperationResultContractProperties) SetUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Updated = &formatted
}

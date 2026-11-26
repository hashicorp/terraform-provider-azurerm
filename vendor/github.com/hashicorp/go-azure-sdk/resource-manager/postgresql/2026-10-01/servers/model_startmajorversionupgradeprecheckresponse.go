package servers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartMajorVersionUpgradePrecheckResponse struct {
	CreateTime *string                            `json:"createTime,omitempty"`
	Name       *string                            `json:"name,omitempty"`
	Status     *MajorVersionUpgradePrecheckStatus `json:"status,omitempty"`
}

func (o *StartMajorVersionUpgradePrecheckResponse) GetCreateTimeAsTime() (*time.Time, error) {
	if o.CreateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StartMajorVersionUpgradePrecheckResponse) SetCreateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreateTime = &formatted
}

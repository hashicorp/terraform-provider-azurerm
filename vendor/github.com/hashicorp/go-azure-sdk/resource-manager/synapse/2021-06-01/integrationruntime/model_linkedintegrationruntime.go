package integrationruntime

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LinkedIntegrationRuntime struct {
	CreateTime          *string `json:"createTime,omitempty"`
	DataFactoryLocation *string `json:"dataFactoryLocation,omitempty"`
	DataFactoryName     *string `json:"dataFactoryName,omitempty"`
	Name                *string `json:"name,omitempty"`
	SubscriptionId      *string `json:"subscriptionId,omitempty"`
}

func (o *LinkedIntegrationRuntime) GetCreateTimeAsTime() (*time.Time, error) {
	if o.CreateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *LinkedIntegrationRuntime) SetCreateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreateTime = &formatted
}

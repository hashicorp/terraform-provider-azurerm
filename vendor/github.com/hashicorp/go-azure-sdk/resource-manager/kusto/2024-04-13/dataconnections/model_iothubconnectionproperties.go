package dataconnections

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IotHubConnectionProperties struct {
	ConsumerGroup          string             `json:"consumerGroup"`
	DataFormat             *IotHubDataFormat  `json:"dataFormat,omitempty"`
	DatabaseRouting        *DatabaseRouting   `json:"databaseRouting,omitempty"`
	EventSystemProperties  *[]string          `json:"eventSystemProperties,omitempty"`
	IotHubResourceId       string             `json:"iotHubResourceId"`
	MappingRuleName        *string            `json:"mappingRuleName,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
	RetrievalStartDate     *string            `json:"retrievalStartDate,omitempty"`
	SharedAccessPolicyName string             `json:"sharedAccessPolicyName"`
	TableName              *string            `json:"tableName,omitempty"`
}

func (o *IotHubConnectionProperties) GetRetrievalStartDateAsTime() (*time.Time, error) {
	if o.RetrievalStartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RetrievalStartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *IotHubConnectionProperties) SetRetrievalStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RetrievalStartDate = &formatted
}

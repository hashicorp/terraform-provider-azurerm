package dataconnections

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
	SharedAccessPolicyName string             `json:"sharedAccessPolicyName"`
	TableName              *string            `json:"tableName,omitempty"`
}

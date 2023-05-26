package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AwsCloudTrailDataConnectorProperties struct {
	AwsRoleArn *string                              `json:"awsRoleArn,omitempty"`
	DataTypes  *AwsCloudTrailDataConnectorDataTypes `json:"dataTypes,omitempty"`
}

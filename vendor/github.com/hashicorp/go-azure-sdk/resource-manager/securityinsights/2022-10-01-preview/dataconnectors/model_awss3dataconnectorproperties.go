package dataconnectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AwsS3DataConnectorProperties struct {
	DataTypes        AwsS3DataConnectorDataTypes `json:"dataTypes"`
	DestinationTable string                      `json:"destinationTable"`
	RoleArn          string                      `json:"roleArn"`
	SqsURLs          []string                    `json:"sqsUrls"`
}

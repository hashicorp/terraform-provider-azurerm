package jobdefinitions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceEndpointProperties struct {
	AwsS3BucketId            *string `json:"awsS3BucketId,omitempty"`
	Name                     *string `json:"name,omitempty"`
	SourceEndpointResourceId *string `json:"sourceEndpointResourceId,omitempty"`
}

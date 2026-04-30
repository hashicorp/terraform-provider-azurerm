package jobsteps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStepOutput struct {
	Credential        *string            `json:"credential,omitempty"`
	DatabaseName      string             `json:"databaseName"`
	ResourceGroupName *string            `json:"resourceGroupName,omitempty"`
	SchemaName        *string            `json:"schemaName,omitempty"`
	ServerName        string             `json:"serverName"`
	SubscriptionId    *string            `json:"subscriptionId,omitempty"`
	TableName         string             `json:"tableName"`
	Type              *JobStepOutputType `json:"type,omitempty"`
}

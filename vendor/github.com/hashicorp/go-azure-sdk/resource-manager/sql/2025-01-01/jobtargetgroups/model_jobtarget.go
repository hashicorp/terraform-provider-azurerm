package jobtargetgroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobTarget struct {
	DatabaseName      *string                       `json:"databaseName,omitempty"`
	ElasticPoolName   *string                       `json:"elasticPoolName,omitempty"`
	MembershipType    *JobTargetGroupMembershipType `json:"membershipType,omitempty"`
	RefreshCredential *string                       `json:"refreshCredential,omitempty"`
	ServerName        *string                       `json:"serverName,omitempty"`
	ShardMapName      *string                       `json:"shardMapName,omitempty"`
	Type              JobTargetType                 `json:"type"`
}

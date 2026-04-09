package registeredserverresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegisteredServerCreateParametersProperties struct {
	AgentVersion      *string `json:"agentVersion,omitempty"`
	ClusterId         *string `json:"clusterId,omitempty"`
	ClusterName       *string `json:"clusterName,omitempty"`
	FriendlyName      *string `json:"friendlyName,omitempty"`
	LastHeartBeat     *string `json:"lastHeartBeat,omitempty"`
	ServerCertificate *string `json:"serverCertificate,omitempty"`
	ServerId          *string `json:"serverId,omitempty"`
	ServerOSVersion   *string `json:"serverOSVersion,omitempty"`
	ServerRole        *string `json:"serverRole,omitempty"`
}

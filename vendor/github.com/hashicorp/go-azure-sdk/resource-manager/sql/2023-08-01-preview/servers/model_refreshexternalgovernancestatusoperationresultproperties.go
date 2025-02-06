package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RefreshExternalGovernanceStatusOperationResultProperties struct {
	ErrorMessage *string `json:"errorMessage,omitempty"`
	QueuedTime   *string `json:"queuedTime,omitempty"`
	RequestId    *string `json:"requestId,omitempty"`
	RequestType  *string `json:"requestType,omitempty"`
	ServerName   *string `json:"serverName,omitempty"`
	Status       *string `json:"status,omitempty"`
}

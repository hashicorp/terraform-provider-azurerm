package databases

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImportExportOperationResultProperties struct {
	BlobUri                    *string                                   `json:"blobUri,omitempty"`
	DatabaseName               *string                                   `json:"databaseName,omitempty"`
	ErrorMessage               *string                                   `json:"errorMessage,omitempty"`
	LastModifiedTime           *string                                   `json:"lastModifiedTime,omitempty"`
	PrivateEndpointConnections *[]PrivateEndpointConnectionRequestStatus `json:"privateEndpointConnections,omitempty"`
	QueuedTime                 *string                                   `json:"queuedTime,omitempty"`
	RequestId                  *string                                   `json:"requestId,omitempty"`
	RequestType                *string                                   `json:"requestType,omitempty"`
	ServerName                 *string                                   `json:"serverName,omitempty"`
	Status                     *string                                   `json:"status,omitempty"`
}

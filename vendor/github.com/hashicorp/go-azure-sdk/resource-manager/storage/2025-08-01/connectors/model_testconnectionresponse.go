package connectors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TestConnectionResponse struct {
	StorageConnectorErrorMessage *string `json:"storageConnectorErrorMessage,omitempty"`
	StorageConnectorMethodName   string  `json:"storageConnectorMethodName"`
	StorageConnectorRequestId    string  `json:"storageConnectorRequestId"`
}

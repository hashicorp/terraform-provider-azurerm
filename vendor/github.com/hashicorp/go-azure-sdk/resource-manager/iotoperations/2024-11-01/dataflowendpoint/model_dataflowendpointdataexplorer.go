package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointDataExplorer struct {
	Authentication DataflowEndpointDataExplorerAuthentication `json:"authentication"`
	Batching       *BatchingConfiguration                     `json:"batching,omitempty"`
	Database       string                                     `json:"database"`
	Host           string                                     `json:"host"`
}

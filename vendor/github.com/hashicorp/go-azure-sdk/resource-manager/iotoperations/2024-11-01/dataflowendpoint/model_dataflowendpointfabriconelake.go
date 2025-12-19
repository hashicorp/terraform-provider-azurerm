package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointFabricOneLake struct {
	Authentication  DataflowEndpointFabricOneLakeAuthentication `json:"authentication"`
	Batching        *BatchingConfiguration                      `json:"batching,omitempty"`
	Host            string                                      `json:"host"`
	Names           DataflowEndpointFabricOneLakeNames          `json:"names"`
	OneLakePathType DataflowEndpointFabricPathType              `json:"oneLakePathType"`
}

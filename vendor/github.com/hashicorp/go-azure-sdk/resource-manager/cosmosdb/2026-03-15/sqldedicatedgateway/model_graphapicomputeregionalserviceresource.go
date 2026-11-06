package sqldedicatedgateway

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GraphAPIComputeRegionalServiceResource struct {
	GraphApiComputeEndpoint *string        `json:"graphApiComputeEndpoint,omitempty"`
	Location                *string        `json:"location,omitempty"`
	Name                    *string        `json:"name,omitempty"`
	Status                  *ServiceStatus `json:"status,omitempty"`
}

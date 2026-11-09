package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteFailoverConnectionResourceDetails struct {
	LastUpdatedTime *string                   `json:"lastUpdatedTime,omitempty"`
	Name            *string                   `json:"name,omitempty"`
	NrpResourceUri  *string                   `json:"nrpResourceUri,omitempty"`
	Status          *FailoverConnectionStatus `json:"status,omitempty"`
}

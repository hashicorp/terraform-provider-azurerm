package workspaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedVirtualNetworkSettings struct {
	AllowedAadTenantIdsForLinking     *[]string `json:"allowedAadTenantIdsForLinking,omitempty"`
	LinkedAccessCheckOnTargetResource *bool     `json:"linkedAccessCheckOnTargetResource,omitempty"`
	PreventDataExfiltration           *bool     `json:"preventDataExfiltration,omitempty"`
}

package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserAssignedIdentitiesProfile struct {
	ControlPlaneOperators  map[string]string `json:"controlPlaneOperators"`
	DataPlaneOperators     map[string]string `json:"dataPlaneOperators"`
	ServiceManagedIdentity string            `json:"serviceManagedIdentity"`
}

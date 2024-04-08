package openshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterProfile struct {
	Domain               *string               `json:"domain,omitempty"`
	FipsValidatedModules *FipsValidatedModules `json:"fipsValidatedModules,omitempty"`
	PullSecret           *string               `json:"pullSecret,omitempty"`
	ResourceGroupId      *string               `json:"resourceGroupId,omitempty"`
	Version              *string               `json:"version,omitempty"`
}

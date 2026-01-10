package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdvancedNetworking struct {
	Enabled       *bool                            `json:"enabled,omitempty"`
	Observability *AdvancedNetworkingObservability `json:"observability,omitempty"`
	Security      *AdvancedNetworkingSecurity      `json:"security,omitempty"`
}

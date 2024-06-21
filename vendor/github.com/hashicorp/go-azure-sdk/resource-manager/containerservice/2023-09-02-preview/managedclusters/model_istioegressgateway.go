package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IstioEgressGateway struct {
	Enabled      bool               `json:"enabled"`
	NodeSelector *map[string]string `json:"nodeSelector,omitempty"`
}

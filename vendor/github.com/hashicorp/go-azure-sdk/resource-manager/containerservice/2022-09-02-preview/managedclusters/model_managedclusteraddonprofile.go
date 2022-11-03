package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterAddonProfile struct {
	Config   *map[string]string    `json:"config,omitempty"`
	Enabled  bool                  `json:"enabled"`
	Identity *UserAssignedIdentity `json:"identity,omitempty"`
}

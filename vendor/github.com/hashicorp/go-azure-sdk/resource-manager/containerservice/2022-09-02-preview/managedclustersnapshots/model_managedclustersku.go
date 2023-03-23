package managedclustersnapshots

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterSKU struct {
	Name *ManagedClusterSKUName `json:"name,omitempty"`
	Tier *ManagedClusterSKUTier `json:"tier,omitempty"`
}

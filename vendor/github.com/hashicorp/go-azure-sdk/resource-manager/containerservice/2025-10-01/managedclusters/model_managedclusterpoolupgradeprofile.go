package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedClusterPoolUpgradeProfile struct {
	KubernetesVersion string                                             `json:"kubernetesVersion"`
	Name              *string                                            `json:"name,omitempty"`
	OsType            OSType                                             `json:"osType"`
	Upgrades          *[]ManagedClusterPoolUpgradeProfileUpgradesInlined `json:"upgrades,omitempty"`
}

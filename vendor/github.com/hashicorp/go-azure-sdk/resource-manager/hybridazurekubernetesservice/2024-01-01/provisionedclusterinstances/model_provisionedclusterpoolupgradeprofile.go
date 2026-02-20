package provisionedclusterinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProvisionedClusterPoolUpgradeProfile struct {
	KubernetesVersion *string                                           `json:"kubernetesVersion,omitempty"`
	OsType            *OsType                                           `json:"osType,omitempty"`
	Upgrades          *[]ProvisionedClusterPoolUpgradeProfileProperties `json:"upgrades,omitempty"`
}

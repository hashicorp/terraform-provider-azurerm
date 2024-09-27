package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutomaticOSUpgradePolicy struct {
	DisableAutomaticRollback *bool `json:"disableAutomaticRollback,omitempty"`
	EnableAutomaticOSUpgrade *bool `json:"enableAutomaticOSUpgrade,omitempty"`
	OsRollingUpgradeDeferral *bool `json:"osRollingUpgradeDeferral,omitempty"`
	UseRollingUpgradePolicy  *bool `json:"useRollingUpgradePolicy,omitempty"`
}

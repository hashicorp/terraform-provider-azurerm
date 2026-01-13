package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpgradePolicy struct {
	AutomaticOSUpgradePolicy *AutomaticOSUpgradePolicy `json:"automaticOSUpgradePolicy,omitempty"`
	Mode                     UpgradeMode               `json:"mode"`
	RollingUpgradePolicy     *RollingUpgradePolicy     `json:"rollingUpgradePolicy,omitempty"`
}

package cosmosdb

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoscaleSettingsResource struct {
	AutoUpgradePolicy   *AutoUpgradePolicyResource `json:"autoUpgradePolicy,omitempty"`
	MaxThroughput       int64                      `json:"maxThroughput"`
	TargetMaxThroughput *int64                     `json:"targetMaxThroughput,omitempty"`
}

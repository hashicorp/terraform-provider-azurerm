package goldengatedeployments

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceConfigurationType struct {
	BundleReleaseUpgradePeriodInDays   *int64 `json:"bundleReleaseUpgradePeriodInDays,omitempty"`
	InterimReleaseUpgradePeriodInDays  *int64 `json:"interimReleaseUpgradePeriodInDays,omitempty"`
	IsInterimReleaseAutoUpgradeEnabled *bool  `json:"isInterimReleaseAutoUpgradeEnabled,omitempty"`
	MajorReleaseUpgradePeriodInDays    *int64 `json:"majorReleaseUpgradePeriodInDays,omitempty"`
	SecurityPatchUpgradePeriodInDays   *int64 `json:"securityPatchUpgradePeriodInDays,omitempty"`
}

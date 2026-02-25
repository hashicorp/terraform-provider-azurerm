package defenderforstorage

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefenderForStorageSettingProperties struct {
	IsEnabled                         *bool                             `json:"isEnabled,omitempty"`
	MalwareScanning                   *MalwareScanningProperties        `json:"malwareScanning,omitempty"`
	OverrideSubscriptionLevelSettings *bool                             `json:"overrideSubscriptionLevelSettings,omitempty"`
	SensitiveDataDiscovery            *SensitiveDataDiscoveryProperties `json:"sensitiveDataDiscovery,omitempty"`
}

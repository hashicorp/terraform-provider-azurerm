package extensions

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchExtensionProperties struct {
	AutoUpgradeMinorVersion        *bool              `json:"autoUpgradeMinorVersion,omitempty"`
	ConfigurationProtectedSettings *map[string]string `json:"configurationProtectedSettings,omitempty"`
	ConfigurationSettings          *map[string]string `json:"configurationSettings,omitempty"`
	ReleaseTrain                   *string            `json:"releaseTrain,omitempty"`
	Version                        *string            `json:"version,omitempty"`
}

package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerServiceNetworkProfileKubeProxyConfig struct {
	Enabled    *bool                                                    `json:"enabled,omitempty"`
	IPvsConfig *ContainerServiceNetworkProfileKubeProxyConfigIPvsConfig `json:"ipvsConfig,omitempty"`
	Mode       *Mode                                                    `json:"mode,omitempty"`
}

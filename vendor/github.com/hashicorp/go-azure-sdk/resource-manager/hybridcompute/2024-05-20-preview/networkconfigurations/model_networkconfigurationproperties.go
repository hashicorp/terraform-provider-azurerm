package networkconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkConfigurationProperties struct {
	KeyProperties                       *KeyProperties `json:"keyProperties,omitempty"`
	Location                            *string        `json:"location,omitempty"`
	NetworkConfigurationScopeId         *string        `json:"networkConfigurationScopeId,omitempty"`
	NetworkConfigurationScopeResourceId string         `json:"networkConfigurationScopeResourceId"`
	TenantId                            *string        `json:"tenantId,omitempty"`
}

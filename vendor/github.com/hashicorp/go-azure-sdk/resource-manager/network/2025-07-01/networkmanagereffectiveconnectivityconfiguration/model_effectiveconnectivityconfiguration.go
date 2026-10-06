package networkmanagereffectiveconnectivityconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EffectiveConnectivityConfiguration struct {
	ConfigurationGroups *[]ConfigurationGroup                `json:"configurationGroups,omitempty"`
	Id                  *string                              `json:"id,omitempty"`
	Properties          *ConnectivityConfigurationProperties `json:"properties,omitempty"`
}

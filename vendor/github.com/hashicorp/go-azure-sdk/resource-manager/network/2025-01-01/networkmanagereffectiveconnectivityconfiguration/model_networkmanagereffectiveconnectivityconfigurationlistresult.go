package networkmanagereffectiveconnectivityconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkManagerEffectiveConnectivityConfigurationListResult struct {
	SkipToken *string                               `json:"skipToken,omitempty"`
	Value     *[]EffectiveConnectivityConfiguration `json:"value,omitempty"`
}

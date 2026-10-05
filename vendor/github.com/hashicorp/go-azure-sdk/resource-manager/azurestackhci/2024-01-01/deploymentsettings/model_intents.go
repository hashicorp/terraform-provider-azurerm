package deploymentsettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Intents struct {
	Adapter                             *[]string                            `json:"adapter,omitempty"`
	AdapterPropertyOverrides            *AdapterPropertyOverrides            `json:"adapterPropertyOverrides,omitempty"`
	Name                                *string                              `json:"name,omitempty"`
	OverrideAdapterProperty             *bool                                `json:"overrideAdapterProperty,omitempty"`
	OverrideQosPolicy                   *bool                                `json:"overrideQosPolicy,omitempty"`
	OverrideVirtualSwitchConfiguration  *bool                                `json:"overrideVirtualSwitchConfiguration,omitempty"`
	QosPolicyOverrides                  *QosPolicyOverrides                  `json:"qosPolicyOverrides,omitempty"`
	TrafficType                         *[]string                            `json:"trafficType,omitempty"`
	VirtualSwitchConfigurationOverrides *VirtualSwitchConfigurationOverrides `json:"virtualSwitchConfigurationOverrides,omitempty"`
}

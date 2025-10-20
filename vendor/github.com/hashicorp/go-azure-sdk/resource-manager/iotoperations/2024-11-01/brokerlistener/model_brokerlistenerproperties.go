package brokerlistener

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BrokerListenerProperties struct {
	Ports             []ListenerPort     `json:"ports"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ServiceName       *string            `json:"serviceName,omitempty"`
	ServiceType       *ServiceType       `json:"serviceType,omitempty"`
}

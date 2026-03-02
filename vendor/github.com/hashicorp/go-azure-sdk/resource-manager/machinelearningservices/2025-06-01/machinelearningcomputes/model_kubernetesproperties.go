package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KubernetesProperties struct {
	DefaultInstanceType           *string                        `json:"defaultInstanceType,omitempty"`
	ExtensionInstanceReleaseTrain *string                        `json:"extensionInstanceReleaseTrain,omitempty"`
	ExtensionPrincipalId          *string                        `json:"extensionPrincipalId,omitempty"`
	InstanceTypes                 *map[string]InstanceTypeSchema `json:"instanceTypes,omitempty"`
	Namespace                     *string                        `json:"namespace,omitempty"`
	RelayConnectionString         *string                        `json:"relayConnectionString,omitempty"`
	ServiceBusConnectionString    *string                        `json:"serviceBusConnectionString,omitempty"`
	VcName                        *string                        `json:"vcName,omitempty"`
}

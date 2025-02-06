package trafficcontrollerinterface

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficControllerProperties struct {
	Associations           *[]ResourceId      `json:"associations,omitempty"`
	ConfigurationEndpoints *[]string          `json:"configurationEndpoints,omitempty"`
	Frontends              *[]ResourceId      `json:"frontends,omitempty"`
	ProvisioningState      *ProvisioningState `json:"provisioningState,omitempty"`
}

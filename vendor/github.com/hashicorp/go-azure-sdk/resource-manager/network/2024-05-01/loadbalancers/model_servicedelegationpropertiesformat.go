package loadbalancers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceDelegationPropertiesFormat struct {
	Actions           *[]string          `json:"actions,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	ServiceName       *string            `json:"serviceName,omitempty"`
}

package integrationruntime

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedIntegrationRuntimeTypeProperties struct {
	ComputeProperties      *IntegrationRuntimeComputeProperties      `json:"computeProperties,omitempty"`
	CustomerVirtualNetwork *IntegrationRuntimeCustomerVirtualNetwork `json:"customerVirtualNetwork,omitempty"`
	SsisProperties         *IntegrationRuntimeSsisProperties         `json:"ssisProperties,omitempty"`
}

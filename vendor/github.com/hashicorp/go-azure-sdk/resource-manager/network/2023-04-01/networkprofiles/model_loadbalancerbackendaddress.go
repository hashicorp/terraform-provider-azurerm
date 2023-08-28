package networkprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerBackendAddress struct {
	Name       *string                                     `json:"name,omitempty"`
	Properties *LoadBalancerBackendAddressPropertiesFormat `json:"properties,omitempty"`
}

package containerinstance

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	ApplicationGateway *ApplicationGateway `json:"applicationGateway,omitempty"`
	LoadBalancer       *LoadBalancer       `json:"loadBalancer,omitempty"`
}

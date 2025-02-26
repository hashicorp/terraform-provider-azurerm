package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IstioComponents struct {
	IngressGateways *[]IstioIngressGateway `json:"ingressGateways,omitempty"`
}

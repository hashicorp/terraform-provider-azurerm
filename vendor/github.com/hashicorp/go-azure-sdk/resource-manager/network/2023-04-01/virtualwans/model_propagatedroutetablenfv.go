package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PropagatedRouteTableNfv struct {
	Ids    *[]RoutingConfigurationNfvSubResource `json:"ids,omitempty"`
	Labels *[]string                             `json:"labels,omitempty"`
}

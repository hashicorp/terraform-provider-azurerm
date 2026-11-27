package httprouteconfig

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPRouteTarget struct {
	ContainerApp string  `json:"containerApp"`
	Label        *string `json:"label,omitempty"`
	Revision     *string `json:"revision,omitempty"`
}

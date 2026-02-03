package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteSourceDetails struct {
	Circuit *string `json:"circuit,omitempty"`
	Pri     *string `json:"pri,omitempty"`
	Sec     *string `json:"sec,omitempty"`
}

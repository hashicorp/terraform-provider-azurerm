package expressrouteconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExpressRouteConnection struct {
	Id         *string                           `json:"id,omitempty"`
	Name       string                            `json:"name"`
	Properties *ExpressRouteConnectionProperties `json:"properties,omitempty"`
}

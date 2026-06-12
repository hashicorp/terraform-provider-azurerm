package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Parameter struct {
	AsPath      *[]string `json:"asPath,omitempty"`
	Community   *[]string `json:"community,omitempty"`
	RoutePrefix *[]string `json:"routePrefix,omitempty"`
}

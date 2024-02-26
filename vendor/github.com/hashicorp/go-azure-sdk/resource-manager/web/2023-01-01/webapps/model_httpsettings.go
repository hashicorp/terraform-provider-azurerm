package webapps

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HTTPSettings struct {
	ForwardProxy *ForwardProxy       `json:"forwardProxy,omitempty"`
	RequireHTTPS *bool               `json:"requireHttps,omitempty"`
	Routes       *HTTPSettingsRoutes `json:"routes,omitempty"`
}

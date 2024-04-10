package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayCustomError struct {
	CustomErrorPageUrl *string                                  `json:"customErrorPageUrl,omitempty"`
	StatusCode         *ApplicationGatewayCustomErrorStatusCode `json:"statusCode,omitempty"`
}

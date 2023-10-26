package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxFrontendIPConfiguration struct {
	PrivateIPAddresses *[]NginxPrivateIPAddress `json:"privateIPAddresses,omitempty"`
	PublicIPAddresses  *[]NginxPublicIPAddress  `json:"publicIPAddresses,omitempty"`
}

package nginxdeployment

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NginxNetworkProfile struct {
	FrontEndIPConfiguration       *NginxFrontendIPConfiguration       `json:"frontEndIPConfiguration,omitempty"`
	NetworkInterfaceConfiguration *NginxNetworkInterfaceConfiguration `json:"networkInterfaceConfiguration,omitempty"`
}

package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PanoramaConfig struct {
	CgName          *string `json:"cgName,omitempty"`
	ConfigString    string  `json:"configString"`
	DgName          *string `json:"dgName,omitempty"`
	HostName        *string `json:"hostName,omitempty"`
	PanoramaServer  *string `json:"panoramaServer,omitempty"`
	PanoramaServer2 *string `json:"panoramaServer2,omitempty"`
	TplName         *string `json:"tplName,omitempty"`
	VMAuthKey       *string `json:"vmAuthKey,omitempty"`
}

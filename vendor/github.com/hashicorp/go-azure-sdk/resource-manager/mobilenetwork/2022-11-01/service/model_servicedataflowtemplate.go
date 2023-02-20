package service

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceDataFlowTemplate struct {
	Direction    SdfDirection `json:"direction"`
	Ports        *[]string    `json:"ports,omitempty"`
	Protocol     []string     `json:"protocol"`
	RemoteIPList []string     `json:"remoteIpList"`
	TemplateName string       `json:"templateName"`
}

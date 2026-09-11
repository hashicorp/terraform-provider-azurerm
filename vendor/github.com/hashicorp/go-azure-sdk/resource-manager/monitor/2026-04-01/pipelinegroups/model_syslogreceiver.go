package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SyslogReceiver struct {
	AllowSkipPriHeader *bool              `json:"allowSkipPriHeader,omitempty"`
	AllowedFormats     *[]AllowedFormats  `json:"allowedFormats,omitempty"`
	Endpoint           string             `json:"endpoint"`
	TransportProtocol  *TransportProtocol `json:"transportProtocol,omitempty"`
}

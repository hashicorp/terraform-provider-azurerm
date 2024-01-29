package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FrontendSetting struct {
	BackendConfiguration  EndpointConfiguration `json:"backendConfiguration"`
	FrontendConfiguration EndpointConfiguration `json:"frontendConfiguration"`
	Name                  string                `json:"name"`
	Protocol              ProtocolType          `json:"protocol"`
}

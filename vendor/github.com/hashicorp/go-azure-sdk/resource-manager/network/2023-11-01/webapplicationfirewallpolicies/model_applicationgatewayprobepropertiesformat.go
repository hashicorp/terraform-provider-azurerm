package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayProbePropertiesFormat struct {
	Host                                *string                                     `json:"host,omitempty"`
	Interval                            *int64                                      `json:"interval,omitempty"`
	Match                               *ApplicationGatewayProbeHealthResponseMatch `json:"match,omitempty"`
	MinServers                          *int64                                      `json:"minServers,omitempty"`
	Path                                *string                                     `json:"path,omitempty"`
	PickHostNameFromBackendHTTPSettings *bool                                       `json:"pickHostNameFromBackendHttpSettings,omitempty"`
	PickHostNameFromBackendSettings     *bool                                       `json:"pickHostNameFromBackendSettings,omitempty"`
	Port                                *int64                                      `json:"port,omitempty"`
	Protocol                            *ApplicationGatewayProtocol                 `json:"protocol,omitempty"`
	ProvisioningState                   *ProvisioningState                          `json:"provisioningState,omitempty"`
	Timeout                             *int64                                      `json:"timeout,omitempty"`
	UnhealthyThreshold                  *int64                                      `json:"unhealthyThreshold,omitempty"`
}

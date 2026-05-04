package applicationgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayOnDemandProbe struct {
	BackendAddressPool                  *SubResource                                `json:"backendAddressPool,omitempty"`
	BackendHTTPSettings                 *SubResource                                `json:"backendHttpSettings,omitempty"`
	Host                                *string                                     `json:"host,omitempty"`
	Match                               *ApplicationGatewayProbeHealthResponseMatch `json:"match,omitempty"`
	Path                                *string                                     `json:"path,omitempty"`
	PickHostNameFromBackendHTTPSettings *bool                                       `json:"pickHostNameFromBackendHttpSettings,omitempty"`
	Protocol                            *ApplicationGatewayProtocol                 `json:"protocol,omitempty"`
	Timeout                             *int64                                      `json:"timeout,omitempty"`
}

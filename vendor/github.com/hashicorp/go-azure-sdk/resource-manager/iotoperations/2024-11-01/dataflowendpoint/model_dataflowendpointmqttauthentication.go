package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointMqttAuthentication struct {
	Method                                MqttAuthMethod                                               `json:"method"`
	ServiceAccountTokenSettings           *DataflowEndpointAuthenticationServiceAccountToken           `json:"serviceAccountTokenSettings,omitempty"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `json:"systemAssignedManagedIdentitySettings,omitempty"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `json:"userAssignedManagedIdentitySettings,omitempty"`
	X509CertificateSettings               *DataflowEndpointAuthenticationX509                          `json:"x509CertificateSettings,omitempty"`
}

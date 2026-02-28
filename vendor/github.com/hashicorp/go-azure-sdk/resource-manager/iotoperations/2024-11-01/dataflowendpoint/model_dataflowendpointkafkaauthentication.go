package dataflowendpoint

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataflowEndpointKafkaAuthentication struct {
	Method                                KafkaAuthMethod                                              `json:"method"`
	SaslSettings                          *DataflowEndpointAuthenticationSasl                          `json:"saslSettings,omitempty"`
	SystemAssignedManagedIdentitySettings *DataflowEndpointAuthenticationSystemAssignedManagedIdentity `json:"systemAssignedManagedIdentitySettings,omitempty"`
	UserAssignedManagedIdentitySettings   *DataflowEndpointAuthenticationUserAssignedManagedIdentity   `json:"userAssignedManagedIdentitySettings,omitempty"`
	X509CertificateSettings               *DataflowEndpointAuthenticationX509                          `json:"x509CertificateSettings,omitempty"`
}

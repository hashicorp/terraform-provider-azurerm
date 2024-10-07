package backend

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BackendServiceFabricClusterProperties struct {
	ClientCertificateId           *string                `json:"clientCertificateId,omitempty"`
	ClientCertificatethumbprint   *string                `json:"clientCertificatethumbprint,omitempty"`
	ManagementEndpoints           []string               `json:"managementEndpoints"`
	MaxPartitionResolutionRetries *int64                 `json:"maxPartitionResolutionRetries,omitempty"`
	ServerCertificateThumbprints  *[]string              `json:"serverCertificateThumbprints,omitempty"`
	ServerX509Names               *[]X509CertificateName `json:"serverX509Names,omitempty"`
}

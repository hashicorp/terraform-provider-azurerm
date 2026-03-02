package fluxconfiguration

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OCIRepositoryDefinition struct {
	Insecure              *bool                       `json:"insecure,omitempty"`
	LayerSelector         *LayerSelectorDefinition    `json:"layerSelector,omitempty"`
	LocalAuthRef          *string                     `json:"localAuthRef,omitempty"`
	RepositoryRef         *OCIRepositoryRefDefinition `json:"repositoryRef,omitempty"`
	ServiceAccountName    *string                     `json:"serviceAccountName,omitempty"`
	SyncIntervalInSeconds *int64                      `json:"syncIntervalInSeconds,omitempty"`
	TimeoutInSeconds      *int64                      `json:"timeoutInSeconds,omitempty"`
	TlsConfig             *TlsConfigDefinition        `json:"tlsConfig,omitempty"`
	Url                   *string                     `json:"url,omitempty"`
	UseWorkloadIdentity   *bool                       `json:"useWorkloadIdentity,omitempty"`
	Verify                *VerifyDefinition           `json:"verify,omitempty"`
}

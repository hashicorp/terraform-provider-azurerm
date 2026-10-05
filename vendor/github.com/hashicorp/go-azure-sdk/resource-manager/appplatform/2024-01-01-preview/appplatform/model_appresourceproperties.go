package appplatform

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppResourceProperties struct {
	AddonConfigs          *map[string]interface{}         `json:"addonConfigs,omitempty"`
	CustomPersistentDisks *[]CustomPersistentDiskResource `json:"customPersistentDisks,omitempty"`
	EnableEndToEndTLS     *bool                           `json:"enableEndToEndTLS,omitempty"`
	Fqdn                  *string                         `json:"fqdn,omitempty"`
	HTTPSOnly             *bool                           `json:"httpsOnly,omitempty"`
	IngressSettings       *IngressSettings                `json:"ingressSettings,omitempty"`
	LoadedCertificates    *[]LoadedCertificate            `json:"loadedCertificates,omitempty"`
	PersistentDisk        *PersistentDisk                 `json:"persistentDisk,omitempty"`
	ProvisioningState     *AppResourceProvisioningState   `json:"provisioningState,omitempty"`
	Public                *bool                           `json:"public,omitempty"`
	Secrets               *[]Secret                       `json:"secrets,omitempty"`
	TemporaryDisk         *TemporaryDisk                  `json:"temporaryDisk,omitempty"`
	TestEndpointAuthState *TestEndpointAuthState          `json:"testEndpointAuthState,omitempty"`
	Url                   *string                         `json:"url,omitempty"`
	VnetAddons            *AppVNetAddons                  `json:"vnetAddons,omitempty"`
	WorkloadProfileName   *string                         `json:"workloadProfileName,omitempty"`
}

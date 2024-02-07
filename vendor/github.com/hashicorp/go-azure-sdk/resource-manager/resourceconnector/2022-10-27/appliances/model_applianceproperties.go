package appliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplianceProperties struct {
	Distro               *Distro                                  `json:"distro,omitempty"`
	InfrastructureConfig *AppliancePropertiesInfrastructureConfig `json:"infrastructureConfig,omitempty"`
	ProvisioningState    *string                                  `json:"provisioningState,omitempty"`
	PublicKey            *string                                  `json:"publicKey,omitempty"`
	Status               *Status                                  `json:"status,omitempty"`
	Version              *string                                  `json:"version,omitempty"`
}

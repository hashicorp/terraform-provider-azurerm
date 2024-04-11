package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolProperties struct {
	DevBoxDefinitionName  *string                        `json:"devBoxDefinitionName,omitempty"`
	HealthStatus          *HealthStatus                  `json:"healthStatus,omitempty"`
	HealthStatusDetails   *[]HealthStatusDetail          `json:"healthStatusDetails,omitempty"`
	LicenseType           *LicenseType                   `json:"licenseType,omitempty"`
	LocalAdministrator    *LocalAdminStatus              `json:"localAdministrator,omitempty"`
	NetworkConnectionName *string                        `json:"networkConnectionName,omitempty"`
	ProvisioningState     *ProvisioningState             `json:"provisioningState,omitempty"`
	StopOnDisconnect      *StopOnDisconnectConfiguration `json:"stopOnDisconnect,omitempty"`
}

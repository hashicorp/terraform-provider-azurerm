package pools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PoolProperties struct {
	DevBoxDefinitionName  string                         `json:"devBoxDefinitionName"`
	HealthStatus          *HealthStatus                  `json:"healthStatus,omitempty"`
	HealthStatusDetails   *[]HealthStatusDetail          `json:"healthStatusDetails,omitempty"`
	LicenseType           LicenseType                    `json:"licenseType"`
	LocalAdministrator    LocalAdminStatus               `json:"localAdministrator"`
	NetworkConnectionName string                         `json:"networkConnectionName"`
	ProvisioningState     *ProvisioningState             `json:"provisioningState,omitempty"`
	StopOnDisconnect      *StopOnDisconnectConfiguration `json:"stopOnDisconnect,omitempty"`
}

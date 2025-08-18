package apigateway

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiManagementGatewayBaseProperties struct {
	Backend                 *BackendConfiguration    `json:"backend,omitempty"`
	ConfigurationApi        *GatewayConfigurationApi `json:"configurationApi,omitempty"`
	CreatedAtUtc            *string                  `json:"createdAtUtc,omitempty"`
	Frontend                *FrontendConfiguration   `json:"frontend,omitempty"`
	ProvisioningState       *string                  `json:"provisioningState,omitempty"`
	TargetProvisioningState *string                  `json:"targetProvisioningState,omitempty"`
	VirtualNetworkType      *VirtualNetworkType      `json:"virtualNetworkType,omitempty"`
}

func (o *ApiManagementGatewayBaseProperties) GetCreatedAtUtcAsTime() (*time.Time, error) {
	if o.CreatedAtUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAtUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ApiManagementGatewayBaseProperties) SetCreatedAtUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAtUtc = &formatted
}

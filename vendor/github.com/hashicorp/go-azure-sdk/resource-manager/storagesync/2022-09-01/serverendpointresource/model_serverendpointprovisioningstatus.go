package serverendpointresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerEndpointProvisioningStatus struct {
	ProvisioningStatus       *ServerProvisioningStatus               `json:"provisioningStatus,omitempty"`
	ProvisioningStepStatuses *[]ServerEndpointProvisioningStepStatus `json:"provisioningStepStatuses,omitempty"`
	ProvisioningType         *string                                 `json:"provisioningType,omitempty"`
}

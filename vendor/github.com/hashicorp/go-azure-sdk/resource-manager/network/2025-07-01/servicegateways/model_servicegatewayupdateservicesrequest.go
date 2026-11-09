package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGatewayUpdateServicesRequest struct {
	Action          *ServiceUpdateAction            `json:"action,omitempty"`
	ServiceRequests *[]ServiceGatewayServiceRequest `json:"serviceRequests,omitempty"`
}

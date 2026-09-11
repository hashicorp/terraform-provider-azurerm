package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGatewayService struct {
	Name       *string                                `json:"name,omitempty"`
	Properties *ServiceGatewayServicePropertiesFormat `json:"properties,omitempty"`
}

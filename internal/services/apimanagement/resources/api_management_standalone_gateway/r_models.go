// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_standalone_gateway

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apigateway"
)

type ApiManagementStandaloneGatewayModel struct {
	Name               string            `tfschema:"name"`
	ResourceGroupName  string            `tfschema:"resource_group_name"`
	BackendSubnetId    string            `tfschema:"backend_subnet_id"`
	Location           string            `tfschema:"location"`
	Sku                []GatewaySkuModel `tfschema:"sku"`
	Tags               map[string]string `tfschema:"tags"`
	VirtualNetworkType string            `tfschema:"virtual_network_type"`
}

type GatewaySkuModel struct {
	Capacity int64                        `tfschema:"capacity"`
	Name     apigateway.ApiGatewaySkuType `tfschema:"name"`
}

func (r Resource) ModelObject() any {
	return &ApiManagementStandaloneGatewayModel{}
}

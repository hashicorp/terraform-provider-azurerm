// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_connection

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/helpers"
)

type SpringCloudConnectorResourceModel struct {
	Name             string                     `tfschema:"name"`
	SpringCloudId    string                     `tfschema:"spring_cloud_id"`
	TargetResourceId string                     `tfschema:"target_resource_id"`
	ClientType       string                     `tfschema:"client_type"`
	AuthInfo         []helpers.AuthInfoModel    `tfschema:"authentication"`
	VnetSolution     string                     `tfschema:"vnet_solution"`
	SecretStore      []helpers.SecretStoreModel `tfschema:"secret_store"`
}

func (r Resource) ModelObject() any {
	return &SpringCloudConnectorResourceModel{}
}

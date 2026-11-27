// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package app_service_connection

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/helpers"
)

type AppServiceConnectorResourceModel struct {
	Name             string                     `tfschema:"name"`
	AppServiceId     string                     `tfschema:"app_service_id"`
	TargetResourceId string                     `tfschema:"target_resource_id"`
	ClientType       string                     `tfschema:"client_type"`
	AuthInfo         []helpers.AuthInfoModel    `tfschema:"authentication"`
	VnetSolution     string                     `tfschema:"vnet_solution"`
	SecretStore      []helpers.SecretStoreModel `tfschema:"secret_store"`
}

func (r Resource) ModelObject() any {
	return &AppServiceConnectorResourceModel{}
}

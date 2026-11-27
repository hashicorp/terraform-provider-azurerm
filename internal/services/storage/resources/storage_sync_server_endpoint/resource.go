// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_server_endpoint

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/serverendpointresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name storage_sync_server_endpoint -service-package-name storage -properties "name" -compare-values "subscription_id:storage_sync_group_id,resource_group_name:storage_sync_group_id,storage_sync_service_name:storage_sync_group_id,sync_group_name:storage_sync_group_id" -test-sequential

type Resource struct{}

var (
	_ sdk.ResourceWithUpdate   = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return serverendpointresource.ValidateServerEndpointID
}

func (r Resource) ResourceType() string {
	return "azurerm_storage_sync_server_endpoint"
}

func (r Resource) Identity() resourceids.ResourceId {
	return &serverendpointresource.ServerEndpointId{}
}

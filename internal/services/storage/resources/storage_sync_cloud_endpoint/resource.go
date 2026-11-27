// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_sync_cloud_endpoint

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagesync/2020-03-01/cloudendpointresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -resource-name storage_sync_cloud_endpoint -service-package-name storage -properties "name" -compare-values "subscription_id:storage_sync_group_id,resource_group_name:storage_sync_group_id,storage_sync_service_name:storage_sync_group_id,sync_group_name:storage_sync_group_id"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceStorageSyncCloudEndpointCreate,
		Read:   resourceStorageSyncCloudEndpointRead,
		Delete: resourceStorageSyncCloudEndpointDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&cloudendpointresource.CloudEndpointId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&cloudendpointresource.CloudEndpointId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(45 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(45 * time.Minute),
		},

		Schema: storageSyncCloudEndpointSchema(),
	}
}

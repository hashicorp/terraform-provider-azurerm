// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_volume_bucket

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/buckets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Resource struct{}

var (
	_ sdk.Resource             = Resource{}
	_ sdk.ResourceWithIdentity = Resource{}
)

func (r Resource) Identity() resourceids.ResourceId {
	return &buckets.BucketId{}
}

func (r Resource) ResourceType() string {
	return "azurerm_netapp_volume_bucket"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return buckets.ValidateBucketID
}

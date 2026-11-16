// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_metadata_schema

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/metadataschemas"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity -properties "name" -compare-values "subscription_id:api_center_service_id,resource_group_name:api_center_service_id,service_name:api_center_service_id"
var (
	_ sdk.ResourceWithUpdate        = Resource{}
	_ sdk.ResourceWithIdentity      = Resource{}
	_ sdk.ResourceWithCustomizeDiff = Resource{}
)

type Resource struct{}

func (r Resource) Identity() resourceids.ResourceId {
	return &metadataschemas.MetadataSchemaId{}
}

func (r Resource) ModelObject() any {
	return &ApiCenterMetadataSchemaModel{}
}

func (r Resource) ResourceType() string {
	return "azurerm_api_center_metadata_schema"
}

func (r Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return metadataschemas.ValidateMetadataSchemaID
}

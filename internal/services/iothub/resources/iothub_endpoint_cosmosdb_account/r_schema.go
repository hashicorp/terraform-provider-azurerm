// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_cosmosdb_account

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the endpoint. The name must be unique across endpoint types. The following names are reserved: `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubEndpointName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iothub_id": {
			Description:  "The ID of the IoT Hub to create the endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubID,
		},

		"container_name": {
			Description:  "The name of the Cosmos DB Container in the Cosmos DB Database.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"database_name": {
			Description:  "The name of the Cosmos DB Database in the Cosmos DB Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"endpoint_uri": {
			Description:  "The URI of the Cosmos DB Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"authentication_type": {
			Description:  "The type used to authenticate against the Cosmos DB Account endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(devices.AuthenticationTypeKeyBased),
			ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
		},

		"identity_id": {
			Description:   "The ID of the User Managed Identity used to authenticate against the Cosmos DB Account endpoint.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			ConflictsWith: []string{"primary_key", "secondary_key"},
		},

		"partition_key_name": {
			Description:  "The name of the partition key associated with the Cosmos DB Container.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"partition_key_template"},
		},

		"partition_key_template": {
			Description:  "The template for generating a synthetic partition key value for use within the Cosmos DB Container.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"partition_key_name"},
		},

		"primary_key": {
			Description:   "The primary key of the Cosmos DB Account.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"identity_id"},
			RequiredWith:  []string{"secondary_key"},
		},

		"secondary_key": {
			Description:   "The secondary key of the Cosmos DB Account.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"identity_id"},
			RequiredWith:  []string{"primary_key"},
		},

		"subscription_id": {
			Description: "The subscription ID for the endpoint.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C : required since this property would always be set even if it isn't specified in the tf config, otherwise it would cause a diff and break existing users
			Computed:     true,
			ValidateFunc: validation.IsUUID,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

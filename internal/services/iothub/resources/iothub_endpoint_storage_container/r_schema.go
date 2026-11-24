// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_endpoint_storage_container

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func resourceIothubEndpointStorageContainerSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the endpoint. The name must be unique across endpoint types. The following names are reserved: `events`, `operationsMonitoringEvents`, `fileNotifications` and `$default`.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IoTHubEndpointName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		// lintignore: S013
		"iothub_id": {
			Description:  "The IoTHub ID for the endpoint.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iothubValidate.IotHubID,
		},

		"container_name": {
			Description:  "The name of Storage Container in the Storage Account.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.StorageContainerName,
		},

		"file_name_format": {
			Description:  "File name format for the blob. All parameters are mandatory but can be reordered.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "{iothub}/{partition}/{YYYY}/{MM}/{DD}/{HH}/{mm}",
			ValidateFunc: iothubValidate.FileNameFormat,
		},

		"batch_frequency_in_seconds": {
			Description:  "Time interval at which blobs are written to storage.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      300,
			ValidateFunc: validation.IntBetween(60, 720),
		},

		"max_chunk_size_in_bytes": {
			Description:  "Maximum number of bytes for each blob written to storage.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      314572800,
			ValidateFunc: validation.IntBetween(10485760, 524288000),
		},

		"authentication_type": {
			Description:  "Type used to authenticate against the storage endpoint.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(devices.AuthenticationTypeKeyBased),
			ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
		},

		"identity_id": {
			Description:   "ID of the User Managed Identity used to authenticate against the storage endpoint.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ValidateFunc:  commonids.ValidateUserAssignedIdentityID,
			ConflictsWith: []string{"connection_string"},
		},

		"endpoint_uri": {
			Description:  "URI of the Storage Container endpoint. This corresponds to the `primary_blob_endpoint` of the parent Storage Account. This attribute can only be specified and is mandatory when `authentication_type` is `identityBased`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"endpoint_uri", "connection_string"},
		},

		"connection_string": {
			Description: "The connection string for the endpoint. This attribute can only be specified and is mandatory when `authentication_type` is `keyBased`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			DiffSuppressFunc: func(k, old, new string, d *pluginsdk.ResourceData) bool {
				accountKeyRegex := regexp.MustCompile("AccountKey=[^;]+")

				maskedNew := accountKeyRegex.ReplaceAllString(new, "AccountKey=****")
				return (new == d.Get(k).(string)) && (maskedNew == old)
			},
			Sensitive:     true,
			ConflictsWith: []string{"identity_id"},
			ExactlyOneOf:  []string{"endpoint_uri", "connection_string"},
		},

		"encoding": {
			Description:      "Encoding that is used to serialize messages to blobs.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			Default:          string(devices.EncodingAvro),
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validation.StringInEnumSlice(devices.PossibleEncodingValues(), true),
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

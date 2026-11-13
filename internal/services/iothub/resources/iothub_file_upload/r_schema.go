// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_file_upload

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	storageValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	devices "github.com/jackofallops/kermit/sdk/iothub/2022-04-30-preview/iothub" // azignore:AZG010 - package name does not match its path
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"iothub_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubID,
		},

		"connection_string": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: helpers.FileUploadConnectionStringDiffSuppress,
			Sensitive:        true,
			ValidateFunc:     validation.StringIsNotEmpty,
		},

		"container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: storageValidate.StorageContainerName,
		},

		"authentication_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(devices.AuthenticationTypeKeyBased),
			ValidateFunc: validation.StringInEnumSlice(devices.PossibleAuthenticationTypeValues(), false),
		},

		"default_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H",
			ValidateFunc: validation.ISO8601Duration,
		},

		"identity_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateUserAssignedIdentityID,
		},

		"lock_duration": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1M",
			ValidateFunc: validation.ISO8601Duration,
		},

		"max_delivery_count": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      10,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"notifications_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"sas_ttl": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "PT1H",
			ValidateFunc: validation.ISO8601Duration,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_factory_linked_service_sftp

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func dataFactoryLinkedServiceSftpSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},

		"data_factory_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"authentication_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(linkedservices.PossibleValuesForSftpAuthenticationType(), false),
		},

		"host": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"port": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"username": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_key_fingerprint": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"integration_runtime_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"parameters": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_key_passphrase": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"key_vault_private_key_passphrase"},
		},

		"key_vault_private_key_passphrase": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"private_key_passphrase"},
			MaxItems:      1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"password": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsNotEmpty,
			ConflictsWith: []string{"key_vault_password"},
		},

		"key_vault_password": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			ConflictsWith: []string{"password"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"private_key_content_base64": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ValidateFunc:  validation.StringIsBase64,
			ConflictsWith: []string{"key_vault_private_key_content_base64", "private_key_path"},
		},

		"key_vault_private_key_content_base64": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"private_key_content_base64", "private_key_path"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"private_key_path": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"key_vault_private_key_content_base64", "private_key_content_base64"},
		},

		"skip_host_key_validation": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

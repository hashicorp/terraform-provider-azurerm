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
			Description:  "The name of the Data Factory Linked Service. See the [Microsoft documentation](https://learn.microsoft.com/azure/data-factory/naming-rules) for all restrictions.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.LinkedServiceDatasetName,
		},

		"data_factory_id": {
			Description:  "The Data Factory ID in which to associate the Linked Service with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: factories.ValidateFactoryID,
		},

		"authentication_type": {
			Description:  "The type of authentication used to connect to the SFTP server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(linkedservices.PossibleValuesForSftpAuthenticationType(), false),
		},

		"host": {
			Description:  "The SFTP server hostname.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"port": {
			Description: "The TCP port number that the SFTP server uses to listen for client connection.",
			Type:        pluginsdk.TypeInt,
			Required:    true,
		},

		"username": {
			Description:  "The username used to log on to the SFTP server.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_properties": {
			Description: "A map of additional properties to associate with the Data Factory Linked Service.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"annotations": {
			Description: "List of tags that can be used for describing the Data Factory Linked Service.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"description": {
			Description:  "The description for the Data Factory Linked Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"host_key_fingerprint": {
			Description:  "The host key fingerprint of the SFTP server.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"integration_runtime_name": {
			Description:  "The name of the integration runtime to associate with the Data Factory Linked Service.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"parameters": {
			Description: "A map of parameters to associate with the Data Factory Linked Service.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_key_passphrase": {
			Description:   "The passphrase for the private key if the key is encrypted.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Sensitive:     true,
			ConflictsWith: []string{"key_vault_private_key_passphrase"},
		},

		"key_vault_private_key_passphrase": {
			Type:          pluginsdk.TypeList,
			Optional:      true,
			MaxItems:      1,
			ConflictsWith: []string{"private_key_passphrase"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"linked_service_name": {
						Description:  "The name of an existing Key Vault Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Description:  "The name of the secret containing the SSH private key passphrase.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"password": {
			Description:   "Password to log on to the SFTP Server for Basic Authentication.",
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
						Description:  "The name of an existing Key Vault Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Description:  "The name of the secret containing the password.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"private_key_content_base64": {
			Description:   "The Base64 encoded private key content in OpenSSH format used to log on to the SFTP server.",
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
						Description:  "The name of an existing Key Vault Data Factory Linked Service.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"secret_name": {
						Description:  "The name of the secret containing the Base64 encoded SSH private key.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"private_key_path": {
			Description:   "The absolute path to the private key file that the self-hosted integration runtime can access.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ConflictsWith: []string{"key_vault_private_key_content_base64", "private_key_content_base64"},
		},

		"skip_host_key_validation": {
			Description: "Whether to validate host key fingerprint while connecting. If set to `false`, `host_key_fingerprint` must also be set.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
	}
}

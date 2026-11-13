// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func activeDirectoryDomainServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty, // TODO: proper validation
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"domain_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DomainServiceName,
		},

		"initial_replica_set": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain_controller_ip_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"external_access_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					// location is computed here
					"location": commonschema.LocationComputed(),

					"service_status": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"subnet_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard",
				"Enterprise",
				"Premium",
			}, false),
		},

		"filtered_sync_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"notifications": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_recipients": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
					},

					"notify_dc_admins": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"notify_global_admins": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"secure_ldap": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enabled": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},

					"external_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"pfx_certificate": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsBase64,
					},

					"pfx_certificate_password": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},

					"certificate_expiry": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"certificate_thumbprint": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"public_certificate": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"security": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"kerberos_armoring_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"kerberos_rc4_encryption_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"ntlm_v1_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"sync_kerberos_passwords": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"sync_ntlm_passwords": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"sync_on_prem_passwords": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"tls_v1_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},
				},
			},
		},

		"domain_configuration_type": {
			Type:     pluginsdk.TypeString,
			ForceNew: true,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"FullySynced",
				"ResourceTrusting",
			}, false),
		},

		"tags": commonschema.Tags(),

		"deployment_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"resource_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sync_owner": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tenant_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"version": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
	}
}

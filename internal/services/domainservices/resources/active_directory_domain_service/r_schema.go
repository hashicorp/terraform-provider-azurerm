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
			Description:  "The display name for your managed Active Directory Domain Service resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty, // TODO: proper validation
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"domain_name": {
			Description:  "The Active Directory domain to use. See [official documentation](https://learn.microsoft.com/azure/active-directory-domain-services/tutorial-create-instance#create-a-managed-domain) for constraints and recommendations.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DomainServiceName,
		},

		"initial_replica_set": {
			Description: "The initial replica set inherits the same location as the Domain Service resource.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"domain_controller_ip_addresses": {
						Description: "A list of subnet IP addresses for the domain controllers in the initial replica set, typically two.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"external_access_ip_address": {
						Description: "The publicly routable IP address for the domain controllers in the initial replica set.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"id": {
						Description: "A unique ID for the replica set.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					// location is computed here
					"location": commonschema.LocationComputed(),

					"service_status": {
						Description: "The current service status for the initial replica set.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"subnet_id": {
						Description:  "The ID of the subnet in which to place the initial replica set.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},
				},
			},
		},

		"sku": {
			Description: "The SKU to use when provisioning the Domain Service resource.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"Standard",
				"Enterprise",
				"Premium",
			}, false),
		},

		"filtered_sync_enabled": {
			Description: "Whether to enable group-based filtered sync (also called scoped synchronisation).",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"notifications": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"additional_recipients": {
						Description: "A list of additional email addresses to notify when there are alerts in the managed domain.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotWhiteSpace,
						},
					},

					"notify_dc_admins": {
						Description: "Whether to notify members of the _AAD DC Administrators_ group when there are alerts in the managed domain.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"notify_global_admins": {
						Description: "Whether to notify all Global Administrators when there are alerts in the managed domain.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
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
						Description: "Whether to enable secure LDAP for the managed domain. For more information, see [official documentation on enabling LDAPS](https://learn.microsoft.com/azure/active-directory-domain-services/tutorial-configure-ldaps), paying particular attention to the section on network security to avoid unnecessarily exposing your service to Internet-borne bruteforce attacks.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
					},

					"external_access_enabled": {
						Description: "Whether to enable external access to LDAPS over the Internet.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"pfx_certificate": {
						Description:  "The certificate/private key to use for LDAPS, as a base64-encoded TripleDES-SHA1 encrypted PKCS#12 bundle (PFX file).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsBase64,
					},

					"pfx_certificate_password": {
						Description: "The password to use for decrypting the PKCS#12 bundle (PFX file).",
						Type:        pluginsdk.TypeString,
						Required:    true,
						Sensitive:   true,
					},

					"certificate_expiry": {
						Description: "The expiry time of the certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"certificate_thumbprint": {
						Description: "The thumbprint of the certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"public_certificate": {
						Description: "The public certificate.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
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
						Description: "Whether to enable Kerberos Armoring.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"kerberos_rc4_encryption_enabled": {
						Description: "Whether to enable Kerberos RC4 Encryption.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"ntlm_v1_enabled": {
						Description: "Whether to enable legacy NTLM v1 support.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"sync_kerberos_passwords": {
						Description: "Whether to synchronize Kerberos password hashes to the managed domain.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"sync_ntlm_passwords": {
						Description: "Whether to synchronize NTLM password hashes to the managed domain.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"sync_on_prem_passwords": {
						Description: "Whether to synchronize on-premises password hashes to the managed domain.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"tls_v1_enabled": {
						Description: "Whether to enable legacy TLS v1 support.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		"domain_configuration_type": {
			Description: "The configuration type of this Active Directory Domain.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				"FullySynced",
				"ResourceTrusting",
			}, false),
		},

		"tags": commonschema.Tags(),

		"deployment_id": {
			Description: "A unique ID for the managed domain deployment.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"resource_id": {
			Description: "The Azure resource ID for the domain service.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
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

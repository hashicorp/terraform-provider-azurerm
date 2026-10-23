// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package netapp_account

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func netappAccountSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AccountName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"active_directory": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dns_servers": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
					"domain": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[(\da-zA-Z-).]{1,255}$`),
							`The domain name must end with a letter or number before dot and start with a letter or number after dot and can not be longer than 255 characters in length.`,
						),
					},
					"smb_server_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[\da-zA-Z\-]{1,10}$`),
							`smb_server_name can contain a mix of numbers, upper/lowercase letters, dashes, and be no longer than 10 characters.`,
						),
					},
					"username": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"password": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"organizational_unit": {
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     "CN=Computers",
						Description: "The Organizational Unit (OU) within the Windows Active Directory where machines will be created. If blank, defaults to 'CN=Computers'",
					},
					"site_name": {
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Default:     "Default-First-Site-Name",
						Description: "The Active Directory site the service will limit Domain Controller discovery to. If blank, defaults to 'Default-First-Site-Name'",
					},
					"kerberos_ad_name": {
						Type:        pluginsdk.TypeString,
						Optional:    true,
						Description: "Name of the active directory machine. This optional parameter is used only while creating kerberos volume.",
					},
					"kerberos_kdc_ip": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv4Address,
						Description:  "IP address of the KDC server (usually same the DC). This optional parameter is used only while creating kerberos volume.",
					},
					"aes_encryption_enabled": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "If enabled, AES encryption will be enabled for SMB communication.",
					},
					"local_nfs_users_with_ldap_allowed": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "If enabled, NFS client local users can also (in addition to LDAP users) access the NFS volumes.",
					},
					"ldap_over_tls_enabled": {
						Type:         pluginsdk.TypeBool,
						Optional:     true,
						Default:      false,
						RequiredWith: []string{"active_directory.0.server_root_ca_certificate"},
						Description:  "Specifies whether or not the LDAP traffic needs to be secured via TLS.",
					},
					"server_root_ca_certificate": {
						Type:         pluginsdk.TypeString,
						Sensitive:    true,
						Optional:     true,
						RequiredWith: []string{"active_directory.0.ldap_over_tls_enabled"},
						Description:  "When LDAP over SSL/TLS is enabled, the LDAP client is required to have base64 encoded Active Directory Certificate Service's self-signed root CA certificate, this optional parameter is used only for dual protocol with LDAP user-mapping volumes.",
					},
					"ldap_signing_enabled": {
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
						Description: "Specifies whether or not the LDAP traffic needs to be signed.",
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

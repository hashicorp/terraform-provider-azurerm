// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2025-01-01/netappaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	netAppValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetAppAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetAppAccountCreate,
		Read:   resourceNetAppAccountRead,
		Update: resourceNetAppAccountUpdate,
		Delete: resourceNetAppAccountDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := netappaccounts.ParseNetAppAccountID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: netAppValidate.AccountName,
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
								ValidateFunc: validate.IPv4Address,
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
		},
	}
}

func resourceNetAppAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := netappaccounts.NewNetAppAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if d.IsNewResource() {
		existing, err := client.AccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_netapp_account", id.ID())
		}
	}

	accountParameters := netappaccounts.NetAppAccount{
		Location:   azure.NormalizeLocation(d.Get("location").(string)),
		Properties: &netappaccounts.AccountProperties{},
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	activeDirectoryRaw := d.Get("active_directory")
	if activeDirectoryRaw != nil {
		activeDirectories := activeDirectoryRaw.([]interface{})
		activeDirectoriesExpanded := expandNetAppActiveDirectories(activeDirectories)
		if len(pointer.From(activeDirectoriesExpanded)) > 0 {
			accountParameters.Properties.ActiveDirectories = activeDirectoriesExpanded
		}
	}

	anfAccountIdentityRaw := d.Get("identity")
	if anfAccountIdentityRaw != nil {
		anfAccountIdentity, ok := anfAccountIdentityRaw.([]interface{})

		if ok && len(anfAccountIdentity) > 0 {
			anfAccountIdentityExpanded, err := identity.ExpandLegacySystemAndUserAssignedMap(anfAccountIdentity)
			if err != nil {
				return err
			}
			if anfAccountIdentity != nil {
				accountParameters.Identity = anfAccountIdentityExpanded
			}
		}
	}

	if err := client.AccountsCreateOrUpdateThenPoll(ctx, id, accountParameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNetAppAccountRead(d, meta)
}

func resourceNetAppAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	update := netappaccounts.NetAppAccountPatch{
		Properties: &netappaccounts.AccountProperties{},
	}

	if d.HasChange("active_directory") {
		activeDirectoriesRaw := d.Get("active_directory").([]interface{})
		activeDirectories := expandNetAppActiveDirectories(activeDirectoriesRaw)
		update.Properties.ActiveDirectories = activeDirectories
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]interface{})
		update.Tags = tags.Expand(tagsRaw)
	}

	if d.HasChange("identity") {
		anfAccountIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}

		update.Identity = anfAccountIdentity
	}

	if err = client.AccountsUpdateThenPoll(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", id.ID(), err)
	}

	return resourceNetAppAccountRead(d, meta)
}

func resourceNetAppAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NetAppAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		if model.Identity != nil {
			anfAccountIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}

			if err := d.Set("identity", anfAccountIdentity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}
		}

		if model.Properties.ActiveDirectories != nil {
			adProps := *model.Properties.ActiveDirectories
			// response returns an array, but only 1 NetApp AD connection is allowed per the Azure platform currently
			if len(adProps) > 0 {
				// the API returns opaque('***') values for password and server_root_ca_certificate, so we pass through current state values so change detection works
				prevPassword := d.Get("active_directory.0.password").(string)
				prevCaCert := d.Get("active_directory.0.server_root_ca_certificate").(string)

				if err = d.Set("active_directory", flattenNetAppActiveDirectories(&adProps[0], &prevPassword, &prevCaCert)); err != nil {
					return fmt.Errorf("setting `active_directory`: %+v", err)
				}
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceNetAppAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NetApp.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if err := client.AccountsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandNetAppActiveDirectories(input []interface{}) *[]netappaccounts.ActiveDirectory {
	results := make([]netappaccounts.ActiveDirectory, 0)
	if input == nil {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})
		dns := strings.Join(*utils.ExpandStringSlice(v["dns_servers"].([]interface{})), ",")

		result := netappaccounts.ActiveDirectory{
			Dns:                        utils.String(dns),
			Domain:                     utils.String(v["domain"].(string)),
			OrganizationalUnit:         utils.String(v["organizational_unit"].(string)),
			Password:                   utils.String(v["password"].(string)),
			SmbServerName:              utils.String(v["smb_server_name"].(string)),
			Username:                   utils.String(v["username"].(string)),
			Site:                       utils.String(v["site_name"].(string)),
			AdName:                     utils.String(v["kerberos_ad_name"].(string)),
			KdcIP:                      utils.String(v["kerberos_kdc_ip"].(string)),
			AesEncryption:              utils.Bool(v["aes_encryption_enabled"].(bool)),
			AllowLocalNfsUsersWithLdap: utils.Bool(v["local_nfs_users_with_ldap_allowed"].(bool)),
			LdapOverTLS:                utils.Bool(v["ldap_over_tls_enabled"].(bool)),
			ServerRootCACertificate:    utils.String(v["server_root_ca_certificate"].(string)),
			LdapSigning:                utils.Bool(v["ldap_signing_enabled"].(bool)),
		}

		results = append(results, result)
	}
	return &results
}

func flattenNetAppActiveDirectories(input *netappaccounts.ActiveDirectory, prevPassword *string, prevCaCert *string) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"dns_servers":                       utils.FlattenStringSliceWithDelimiter(input.Dns, ","),
			"domain":                            input.Domain,
			"organizational_unit":               input.OrganizationalUnit,
			"password":                          prevPassword,
			"smb_server_name":                   input.SmbServerName,
			"username":                          input.Username,
			"site_name":                         input.Site,
			"kerberos_ad_name":                  input.AdName,
			"kerberos_kdc_ip":                   input.KdcIP,
			"aes_encryption_enabled":            input.AesEncryption,
			"local_nfs_users_with_ldap_allowed": input.AllowLocalNfsUsersWithLdap,
			"ldap_over_tls_enabled":             input.LdapOverTLS,
			"server_root_ca_certificate":        prevCaCert,
			"ldap_signing_enabled":              input.LdapSigning,
		},
	}
}

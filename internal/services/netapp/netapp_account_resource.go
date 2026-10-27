// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/netapp/2026-05-01/netappaccounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
						"dc_connection_encryption_enabled": {
							Type:        pluginsdk.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies if traffic between the Domain Controller and SMB server is encrypted.",
						},
						"ldap_client_preferred_server_ips": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 2,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.IsIPv4Address,
							},
							Description: "Specifies the preferred LDAP server IPs for the LDAP client.",
						},
						"ldap_search_scope": {
							Type:        pluginsdk.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Specifies the LDAP search scope for the LDAP client.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_dn": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(1, 255),
										Description:  "Specifies the group DN, which overrides the base DN for group lookups.",
										AtLeastOneOf: []string{"active_directory.0.ldap_search_scope.0.group_dn", "active_directory.0.ldap_search_scope.0.group_membership_filter", "active_directory.0.ldap_search_scope.0.user_dn"},
									},
									"group_membership_filter": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(1, 255),
										Description:  "specifies the custom LDAP search filter to be used when looking up group membership from LDAP server.",
										AtLeastOneOf: []string{"active_directory.0.ldap_search_scope.0.group_dn", "active_directory.0.ldap_search_scope.0.group_membership_filter", "active_directory.0.ldap_search_scope.0.user_dn"},
									},
									"user_dn": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(1, 255),
										Description:  "specifies the user DN, which overrides the base DN for user lookups.",
										AtLeastOneOf: []string{"active_directory.0.ldap_search_scope.0.group_dn", "active_directory.0.ldap_search_scope.0.group_membership_filter", "active_directory.0.ldap_search_scope.0.user_dn"},
									},
								},
							},
						},
						"administrator_users": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.All(
									validation.StringIsNotEmpty,
									validation.None(
										map[string]func(any, string) ([]string, []error){
											"Email Address": validation.IsEmailAddress,
											"Object ID":     validation.IsUUID,
										},
									),
								),
							},
							Description: "Domain users to be added to the built-in Administrators active directory group.",
						},
						"backup_users": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.All(
									validation.StringIsNotEmpty,
									validation.None(
										map[string]func(any, string) ([]string, []error){
											"Email Address": validation.IsEmailAddress,
											"Object ID":     validation.IsUUID,
										},
									),
								),
							},
							Description: "Domain users to be added to the built-in Backup Operators active directory group.",
						},
						"security_users": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.All(
									validation.StringIsNotEmpty,
									validation.None(
										map[string]func(any, string) ([]string, []error){
											"Email Address": validation.IsEmailAddress,
											"Object ID":     validation.IsUUID,
										},
									),
								),
							},
							Description: "Domain users to be given SeSecurityPrivilege privilege.",
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceNetAppAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.AccountClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := netappaccounts.NewNetAppAccountID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
		Location:   location.Normalize(d.Get("location").(string)),
		Properties: &netappaccounts.AccountProperties{},
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	activeDirectoryRaw := d.Get("active_directory")
	if activeDirectoryRaw != nil {
		activeDirectories := activeDirectoryRaw.([]any)
		activeDirectoriesExpanded := expandNetAppActiveDirectories(activeDirectories)
		if len(pointer.From(activeDirectoriesExpanded)) > 0 {
			accountParameters.Properties.ActiveDirectories = activeDirectoriesExpanded
		}
	}

	anfAccountIdentityRaw := d.Get("identity")
	if anfAccountIdentityRaw != nil {
		anfAccountIdentity, ok := anfAccountIdentityRaw.([]any)

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

	if err := client.AccountsCreateOrUpdateCallbackThenPoll(ctx, id, accountParameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNetAppAccountRead(d, meta)
}

func resourceNetAppAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).NetApp.AccountClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := netappaccounts.ParseNetAppAccountID(d.Id())
	if err != nil {
		return err
	}

	locks.ByID(id.ID())
	defer locks.UnlockByID(id.ID())

	existing, err := client.AccountsGet(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if ad := existing.Model.Properties.ActiveDirectories; ad != nil && len(*ad) > 0 {
		// The API doesn't return these, so we'll set these values based on config.
		// If there are no changes to the `active_directory` block this ensures we don't unintentionally wipe these values
		existingAD := (*ad)[0]
		existingAD.ServerRootCACertificate = pointer.To(d.Get("active_directory.0.server_root_ca_certificate").(string))
		existingAD.Password = pointer.To(d.Get("active_directory.0.password").(string))
		existing.Model.Properties.ActiveDirectories = pointer.To([]netappaccounts.ActiveDirectory{existingAD})
	}

	if d.HasChange("active_directory") {
		existing.Model.Properties.ActiveDirectories = expandNetAppActiveDirectories(d.Get("active_directory").([]any))
	}

	if d.HasChange("tags") {
		existing.Model.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("identity") {
		anfAccountIdentity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		existing.Model.Identity = anfAccountIdentity
	}

	if err := client.AccountsCreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id.ID(), err)
	}

	return resourceNetAppAccountRead(d, meta)
}

func resourceNetAppAccountRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[DEBUG] %s was not found - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.NetAppAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		anfAccountIdentity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}

		if err := d.Set("identity", anfAccountIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if model.Properties != nil {
			// the API doesn't return values for `active_directory.0.password` and `active_directory.0.server_root_ca_certificate`, so we pass through current state values so change detection works
			prevPassword := d.Get("active_directory.0.password").(string)
			prevCaCert := d.Get("active_directory.0.server_root_ca_certificate").(string)

			if err = d.Set("active_directory", flattenNetAppActiveDirectories(model.Properties.ActiveDirectories, &prevPassword, &prevCaCert)); err != nil {
				return fmt.Errorf("setting `active_directory`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceNetAppAccountDelete(d *pluginsdk.ResourceData, meta any) error {
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

func expandNetAppActiveDirectories(input []any) *[]netappaccounts.ActiveDirectory {
	results := make([]netappaccounts.ActiveDirectory, 0)
	if input == nil {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]any)
		dns := strings.Join(*pluginsdk.ExpandStringSlice(v["dns_servers"].([]any)), ",")

		result := netappaccounts.ActiveDirectory{
			Dns:                           pointer.To(dns),
			Domain:                        pointer.To(v["domain"].(string)),
			OrganizationalUnit:            pointer.To(v["organizational_unit"].(string)),
			Password:                      pointer.To(v["password"].(string)),
			SmbServerName:                 pointer.To(v["smb_server_name"].(string)),
			Username:                      pointer.To(v["username"].(string)),
			Site:                          pointer.To(v["site_name"].(string)),
			AdName:                        pointer.To(v["kerberos_ad_name"].(string)),
			KdcIP:                         pointer.To(v["kerberos_kdc_ip"].(string)),
			AesEncryption:                 pointer.To(v["aes_encryption_enabled"].(bool)),
			AllowLocalNfsUsersWithLdap:    pointer.To(v["local_nfs_users_with_ldap_allowed"].(bool)),
			LdapOverTLS:                   pointer.To(v["ldap_over_tls_enabled"].(bool)),
			ServerRootCACertificate:       pointer.To(v["server_root_ca_certificate"].(string)),
			LdapSigning:                   pointer.To(v["ldap_signing_enabled"].(bool)),
			EncryptDCConnections:          pointer.To(v["dc_connection_encryption_enabled"].(bool)),
			PreferredServersForLdapClient: pointer.To(strings.Join(*pluginsdk.ExpandStringSlice(v["ldap_client_preferred_server_ips"].([]any)), ",")),
		}

		var searchScope map[string]any
		if rawSearchScope := v["ldap_search_scope"].([]any); len(rawSearchScope) > 0 {
			searchScope = rawSearchScope[0].(map[string]any)
		}
		result.LdapSearchScope = expandNetAppLdapSearchScope(searchScope)

		if admins := v["administrator_users"].(*pluginsdk.Set); admins != nil {
			result.Administrators = pluginsdk.ExpandStringSlice(admins.List())
		}

		if backupUsers := v["backup_users"].(*pluginsdk.Set); backupUsers != nil {
			result.BackupOperators = pluginsdk.ExpandStringSlice(backupUsers.List())
		}

		if secUsers := v["security_users"].(*pluginsdk.Set); secUsers != nil {
			result.SecurityOperators = pluginsdk.ExpandStringSlice(secUsers.List())
		}

		results = append(results, result)
	}
	return &results
}

func expandNetAppLdapSearchScope(input map[string]any) *netappaccounts.LdapSearchScopeOpt {
	if input == nil {
		// if not supplied, set empty strings to all fields to properly update in API.
		// if LdapSearchScopeOpt is left nil or empty, the API will not update the search scope and will leave it as-is.
		return &netappaccounts.LdapSearchScopeOpt{
			GroupDN:               pointer.To(""),
			GroupMembershipFilter: pointer.To(""),
			UserDN:                pointer.To(""),
		}
	}

	result := netappaccounts.LdapSearchScopeOpt{
		GroupDN:               pointer.To(input["group_dn"].(string)),
		GroupMembershipFilter: pointer.To(input["group_membership_filter"].(string)),
		UserDN:                pointer.To(input["user_dn"].(string)),
	}

	return &result
}

func flattenNetAppActiveDirectories(input *[]netappaccounts.ActiveDirectory, prevPassword *string, prevCaCert *string) []any {
	if input == nil || len(*input) == 0 {
		return []any{}
	}

	v := (*input)[0]

	return []any{
		map[string]any{
			"dns_servers":                       pluginsdk.FlattenStringSliceWithDelimiter(v.Dns, ","),
			"domain":                            v.Domain,
			"organizational_unit":               v.OrganizationalUnit,
			"password":                          prevPassword,
			"smb_server_name":                   v.SmbServerName,
			"username":                          v.Username,
			"site_name":                         v.Site,
			"kerberos_ad_name":                  v.AdName,
			"kerberos_kdc_ip":                   v.KdcIP,
			"aes_encryption_enabled":            v.AesEncryption,
			"local_nfs_users_with_ldap_allowed": v.AllowLocalNfsUsersWithLdap,
			"ldap_over_tls_enabled":             v.LdapOverTLS,
			"server_root_ca_certificate":        prevCaCert,
			"ldap_signing_enabled":              v.LdapSigning,
			"dc_connection_encryption_enabled":  v.EncryptDCConnections,
			"ldap_client_preferred_server_ips":  pluginsdk.FlattenStringSliceWithDelimiter(v.PreferredServersForLdapClient, ","),
			"administrator_users":               pluginsdk.FlattenSlice(v.Administrators),
			"backup_users":                      pluginsdk.FlattenSlice(v.BackupOperators),
			"security_users":                    pluginsdk.FlattenSlice(v.SecurityOperators),
			"ldap_search_scope":                 flattenNetAppLdapSearchScope(v.LdapSearchScope),
		},
	}
}

func flattenNetAppLdapSearchScope(input *netappaccounts.LdapSearchScopeOpt) []any {
	if input == nil || (input.GroupDN == nil && input.GroupMembershipFilter == nil && input.UserDN == nil) {
		return []any{}
	}

	return []any{
		map[string]any{
			"group_dn":                input.GroupDN,
			"group_membership_filter": input.GroupMembershipFilter,
			"user_dn":                 input.UserDN,
		},
	}
}

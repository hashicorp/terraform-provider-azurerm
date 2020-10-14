package aadmgmt

import (
	"context"
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2017-06-01/aad"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/aadmgmt/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/aadmgmt/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmActiveDirectoryDomainService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmActiveDirectoryDomainServiceCreate,
		Read:   resourceArmActiveDirectoryDomainServiceRead,
		Update: resourceArmActiveDirectoryDomainServiceUpdate,
		Delete: resourceArmActiveDirectoryDomainServiceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ValidateDomainServiceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"filtered_sync": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"ldaps": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_access": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"ldaps": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"pfx_certificate": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: azValidate.Base64EncodedString,
						},
						"pfx_certificate_password": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
						},
						"external_access_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"notifications": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_recipients": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotWhiteSpace,
							},
						},
						"notify_dc_admins": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"notify_global_admins": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"security": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ntlm_v1": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"sync_ntlm_passwords": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"tls_v1": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			"domain_controller_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmActiveDirectoryDomainServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AadMgmt.DomainServicesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	existing, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing Domain service %q (Resource Group %q): %s", name, resourceGroup, err)
		}
	}

	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_active_directory_domain_service", *existing.ID)
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	security := d.Get("security").([]interface{})
	ldaps := d.Get("ldaps").([]interface{})
	notifications := d.Get("notifications").([]interface{})
	subnetId := d.Get("subnet_id").(string)
	filteredSync := aad.FilteredSyncDisabled
	if d.Get("filtered_sync").(bool) {
		filteredSync = aad.FilteredSyncDisabled
	}

	domainService := aad.DomainService{
		Location: &location,
		DomainServiceProperties: &aad.DomainServiceProperties{
			DomainName:             &name,
			DomainSecuritySettings: expandDomainServiceSecurity(security),
			FilteredSync:           filteredSync,
			LdapsSettings:          expandDomainServiceLdapsSettings(ldaps),
			NotificationSettings:   expandDomainServiceNotificationSettings(notifications),
			SubnetID:               &subnetId,
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, domainService)
	if err != nil {
		return fmt.Errorf("creating Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Domain Services is 2 controllers running in the cloud
	// The create API completes when the first domain controller is up and running and the DomainControllerIPAddress property will show a single IP
	// Afterwards, there will be an additional domain controller created, then the DomainControllerIPAddress property will update and show 2 addresses
	// The Azure Portal blocks users from modifying the domain service until both domain controllers are up
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      domainServiceControllerRefreshFunc(ctx, client, resourceGroup, name),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
		Timeout:      1 * time.Hour,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for both Domain Service controllers to become available: %+v", err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("null ID returned for Domain Service %q (Resource Group %q)", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmActiveDirectoryDomainServiceRead(d, meta)
}

func resourceArmActiveDirectoryDomainServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AadMgmt.DomainServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	security := d.Get("security").([]interface{})
	ldaps := d.Get("ldaps").([]interface{})
	notifications := d.Get("notifications").([]interface{})
	subnetId := d.Get("subnet_id").(string)

	filteredSync := aad.FilteredSyncDisabled
	if v, ok := d.GetOk("filtered_sync"); ok && v.(bool) {
		filteredSync = aad.FilteredSyncEnabled
	}

	domainService := aad.DomainService{
		DomainServiceProperties: &aad.DomainServiceProperties{
			DomainName:             &name,
			DomainSecuritySettings: expandDomainServiceSecurity(security),
			FilteredSync:           filteredSync,
			LdapsSettings:          expandDomainServiceLdapsSettings(ldaps),
			NotificationSettings:   expandDomainServiceNotificationSettings(notifications),
			SubnetID:               &subnetId,
		},
	}

	future, err := client.Update(ctx, resourceGroup, name, domainService)
	if err != nil {
		return fmt.Errorf("updating Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmActiveDirectoryDomainServiceRead(d, meta)
}

func resourceArmActiveDirectoryDomainServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AadMgmt.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DomainServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return err
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if domainServiceProperties := resp.DomainServiceProperties; domainServiceProperties != nil {
		if err := d.Set("domain_controller_ip_addresses", domainServiceProperties.DomainControllerIPAddress); err != nil {
			return fmt.Errorf("setting `domain_controller_ip_addresses`: %+v", err)
		}
		if err := d.Set("security", flattenDomainServiceSecurity(domainServiceProperties.DomainSecuritySettings)); err != nil {
			return fmt.Errorf("setting `security`: %+v", err)
		}
		if err := d.Set("ldaps", flattenDomainServiceLdapsSettings(domainServiceProperties.LdapsSettings)); err != nil {
			return fmt.Errorf("setting `ldaps`: %+v", err)
		}
		if err := d.Set("notifications", flattenDomainServiceNotification(domainServiceProperties.NotificationSettings)); err != nil {
			return fmt.Errorf("setting `notifications`: %+v", err)
		}

		d.Set("filtered_sync", false)
		if domainServiceProperties.FilteredSync == aad.FilteredSyncEnabled {
			d.Set("filtered_sync", true)
		}
		d.Set("subnet_id", domainServiceProperties.SubnetID)
	}

	return nil
}

func resourceArmActiveDirectoryDomainServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AadMgmt.DomainServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["domainServices"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func domainServiceControllerRefreshFunc(ctx context.Context, client *aad.DomainServicesClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Waiting for both Domain Service controllers to deploy...")
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil, "error", err
		}
		if resp.DomainControllerIPAddress == nil || len(*resp.DomainControllerIPAddress) < 2 {
			return resp, "pending", nil
		}
		return resp, "available", nil
	}
}

func expandDomainServiceSecurity(input []interface{}) *aad.DomainSecuritySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	ntlmV1 := aad.NtlmV1Disabled
	syncNtlmPasswords := aad.SyncNtlmPasswordsDisabled
	tlsV1 := aad.TLSV1Disabled

	if v["ntlm_v1"].(bool) {
		ntlmV1 = aad.NtlmV1Enabled
	}
	if v["sync_ntlm_passwords"].(bool) {
		syncNtlmPasswords = aad.SyncNtlmPasswordsEnabled
	}
	if v["tls_v1"].(bool) {
		tlsV1 = aad.TLSV1Enabled
	}

	return &aad.DomainSecuritySettings{
		NtlmV1:            ntlmV1,
		SyncNtlmPasswords: syncNtlmPasswords,
		TLSV1:             tlsV1,
	}
}

func expandDomainServiceLdapsSettings(input []interface{}) *aad.LdapsSettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	externalAccess := aad.Disabled
	ldaps := aad.LdapsDisabled
	pfxCertificate := v["pfx_certificate"].(string)
	pfxCertificatePassword := v["pfx_certificate_password"].(string)

	if v["external_access"].(bool) {
		externalAccess = aad.Enabled
	}
	if v["ldaps"].(bool) {
		ldaps = aad.LdapsEnabled
	}
	return &aad.LdapsSettings{
		ExternalAccess:         externalAccess,
		Ldaps:                  ldaps,
		PfxCertificate:         &pfxCertificate,
		PfxCertificatePassword: &pfxCertificatePassword,
	}
}

func expandDomainServiceNotificationSettings(input []interface{}) *aad.NotificationSettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	additionalRecipients := make([]string, 0)
	if ar, ok := v["additional_recipients"]; ok {
		additionalRecipients = ar.([]string)
	}

	notifyDcAdmins := aad.NotifyDcAdminsDisabled
	if n, ok := v["notify_dc_admins"]; ok && n.(bool) {
		notifyDcAdmins = aad.NotifyDcAdminsEnabled
	}

	notifyGlobalAdmins := aad.NotifyGlobalAdminsDisabled
	if n, ok := v["notify_global_admins"]; ok && n.(bool) {
		notifyGlobalAdmins = aad.NotifyGlobalAdminsEnabled
	}


	return &aad.NotificationSettings{
		AdditionalRecipients: &additionalRecipients,
		NotifyDcAdmins:       notifyDcAdmins,
		NotifyGlobalAdmins:   notifyGlobalAdmins,
	}
}

func flattenDomainServiceSecurity(input *aad.DomainSecuritySettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := map[string]bool{
		"ntlm_v1":             false,
		"sync_ntlm_passwords": false,
		"tls_v1":              false,
	}
	if input.NtlmV1 == aad.NtlmV1Enabled {
		result["ntlm_v1"] = true
	}
	if input.SyncNtlmPasswords == aad.SyncNtlmPasswordsEnabled {
		result["sync_ntlm_passwords"] = true
	}
	if input.TLSV1 == aad.TLSV1Enabled {
		result["tls_v1"] = true
	}

	return []interface{}{result}
}

func flattenDomainServiceLdapsSettings(input *aad.LdapsSettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := map[string]interface{}{
		"external_access": false,
		"ldaps":           false,
	}

	if input.ExternalAccess == aad.Enabled {
		result["external_access"] = true
	}
	if input.Ldaps == aad.LdapsEnabled {
		result["ldaps"] = true
	}
	if pfxCertificate := input.PfxCertificate; pfxCertificate != nil {
		result["pfx_certificate"] = *pfxCertificate
	}
	if pfxCertificatePassword := input.PfxCertificatePassword; pfxCertificatePassword != nil {
		result["pfx_certificate_password"] = *pfxCertificatePassword
	}
	if externalAccessIPAddress := input.ExternalAccessIPAddress; externalAccessIPAddress != nil {
		result["external_access_ip_address"] = *externalAccessIPAddress
	}

	return []interface{}{result}
}

func flattenDomainServiceNotification(input *aad.NotificationSettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := map[string]interface{}{
		"notify_dc_admins":     false,
		"notify_global_admins": false,
	}

	result["additional_recipients"] = make([]string, 0)
	if input.AdditionalRecipients != nil && len(*input.AdditionalRecipients) > 0 {
		result["additional_recipients"] = *input.AdditionalRecipients
	}
	if input.NotifyDcAdmins == aad.NotifyDcAdminsEnabled {
		result["notify_dc_admins"] = true
	}
	if input.NotifyGlobalAdmins == aad.NotifyGlobalAdminsEnabled {
		result["notify_global_admins"] = true
	}

	return []interface{}{result}
}

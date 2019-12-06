package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2017-06-01/aad"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDomainService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDomainServiceCreate,
		Read:   resourceArmDomainServiceRead,
		Update: resourceArmDomainServiceUpdate,
		Delete: resourceArmDomainServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"domain_security_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ntlm_v1": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.NtlmV1Enabled),
								string(aad.NtlmV1Disabled),
							}, false),
							Default: string(aad.NtlmV1Enabled),
						},
						"sync_ntlm_passwords": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.SyncNtlmPasswordsEnabled),
								string(aad.SyncNtlmPasswordsDisabled),
							}, false),
							Default: string(aad.SyncNtlmPasswordsEnabled),
						},
						"tls_v1": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.TLSV1Enabled),
								string(aad.TLSV1Disabled),
							}, false),
							Default: string(aad.TLSV1Enabled),
						},
					},
				},
			},

			"filtered_sync": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(aad.FilteredSyncEnabled),
					string(aad.FilteredSyncDisabled),
				}, false),
				Default: string(aad.FilteredSyncEnabled),
			},

			"ldaps_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_access": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.Enabled),
								string(aad.Disabled),
							}, false),
							Default: string(aad.Enabled),
						},
						"ldaps": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.LdapsEnabled),
								string(aad.LdapsDisabled),
							}, false),
							Default: string(aad.LdapsEnabled),
						},
						"pfx_certificate": {
							Type:         schema.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validate.Base64String(),
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

			"notification_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_recipients": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"notify_dc_admins": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.NotifyDcAdminsEnabled),
								string(aad.NotifyDcAdminsDisabled),
							}, false),
							Default: string(aad.NotifyDcAdminsEnabled),
						},
						"notify_global_admins": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(aad.NotifyGlobalAdminsEnabled),
								string(aad.NotifyGlobalAdminsDisabled),
							}, false),
							Default: string(aad.NotifyGlobalAdminsEnabled),
						},
					},
				},
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"domain_controller_ip_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDomainServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DomainServices.DomainServicesClient
	vnetClient := meta.(*ArmClient).Network.VnetClient
	ctx, cancel := timeouts.ForCreate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Domain Service %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_domain_service", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	domainSecuritySettings := d.Get("domain_security_settings").([]interface{})
	filteredSync := d.Get("filtered_sync").(string)
	ldapsSettings := d.Get("ldaps_settings").([]interface{})
	notificationSettings := d.Get("notification_settings").([]interface{})
	subnetId := d.Get("subnet_id").(string)
	t := d.Get("tags").(map[string]interface{})

	domainService := aad.DomainService{
		Location: &location,
		DomainServiceProperties: &aad.DomainServiceProperties{
			DomainName:             &name,
			DomainSecuritySettings: expandArmDomainServiceDomainSecuritySettings(domainSecuritySettings),
			FilteredSync:           aad.FilteredSync(filteredSync),
			LdapsSettings:          expandArmDomainServiceLdapsSettings(ldapsSettings),
			NotificationSettings:   expandArmDomainServiceNotificationSettings(notificationSettings),
			SubnetID:               &subnetId,
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, domainService)
	if err != nil {
		return fmt.Errorf("Error creating Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// Domain Services is 2 controllers running in the cloud
	// the create api completes Once the first domain controller is up and running, 1 ip address will show up
	// Afterwards, there will be an additional domain controller creating, then the ip address will update and become 2
	// the Azure Portal blocks users from modifying the domain service until both domain controllers are up
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      domainServiceControllerRefreshFunc(ctx, client, resourceGroup, name),
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}

	if features.SupportsCustomTimeouts() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = 1 * time.Hour
	}

	if domainControllerIPAddress, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for both Domain Service Controller up, err: %+v", err)
	} else {
		// Update DNS server settings of the virtual network
		parsedId, err := azure.ParseAzureResourceID(subnetId)
		if err != nil {
			return err
		}
		resourceGroupName := parsedId.ResourceGroup
		virtualNetworkName := parsedId.Path["virtualNetworks"]

		resp, err := vnetClient.Get(ctx, resourceGroupName, virtualNetworkName, "")
		if err != nil {
			return fmt.Errorf("Error readding Virtual Network %q (Resource Group %q): %+v", virtualNetworkName, resourceGroupName, err)
		}
		dns := domainControllerIPAddress.([]string)
		resp.DhcpOptions = &network.DhcpOptions{
			DNSServers: &dns,
		}
		if _, err := vnetClient.CreateOrUpdate(ctx, resourceGroupName, virtualNetworkName, resp); err != nil {
			return fmt.Errorf("Error updating DNS server to %+v for Virtual Network %q (Resource Group %q): %+v", domainControllerIPAddress, virtualNetworkName, resourceGroupName, err)
		}
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Domain Service %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmDomainServiceRead(d, meta)
}

func resourceArmDomainServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["domainServices"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Domain Services %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if domainServiceProperties := resp.DomainServiceProperties; domainServiceProperties != nil {
		d.Set("domain_controller_ip_address", utils.FlattenStringSlice(domainServiceProperties.DomainControllerIPAddress))
		if err := d.Set("domain_security_settings", flattenArmDomainServiceDomainSecuritySettings(domainServiceProperties.DomainSecuritySettings)); err != nil {
			return fmt.Errorf("Error setting `domain_security_settings`: %+v", err)
		}
		d.Set("filtered_sync", string(domainServiceProperties.FilteredSync))
		if err := d.Set("ldaps_settings", flattenArmDomainServiceLdapsSettings(domainServiceProperties.LdapsSettings)); err != nil {
			return fmt.Errorf("Error setting `ldaps_settings`: %+v", err)
		}
		if err := d.Set("notification_settings", flattenArmDomainServiceNotificationSettings(domainServiceProperties.NotificationSettings)); err != nil {
			return fmt.Errorf("Error setting `notification_settings`: %+v", err)
		}
		d.Set("subnet_id", domainServiceProperties.SubnetID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDomainServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	domainSecuritySettings := d.Get("domain_security_settings").([]interface{})
	filteredSync := d.Get("filtered_sync").(string)
	ldapsSettings := d.Get("ldaps_settings").([]interface{})
	notificationSettings := d.Get("notification_settings").([]interface{})
	subnetId := d.Get("subnet_id").(string)
	t := d.Get("tags").(map[string]interface{})

	domainService := aad.DomainService{
		DomainServiceProperties: &aad.DomainServiceProperties{
			DomainName:             &name,
			DomainSecuritySettings: expandArmDomainServiceDomainSecuritySettings(domainSecuritySettings),
			FilteredSync:           aad.FilteredSync(filteredSync),
			LdapsSettings:          expandArmDomainServiceLdapsSettings(ldapsSettings),
			NotificationSettings:   expandArmDomainServiceNotificationSettings(notificationSettings),
			SubnetID:               &subnetId,
		},
		Tags: tags.Expand(t),
	}

	future, err := client.Update(ctx, resourceGroup, name, domainService)
	if err != nil {
		return fmt.Errorf("Error updating Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return resourceArmDomainServiceRead(d, meta)
}

func resourceArmDomainServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["domainServices"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deleting Domain Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func domainServiceControllerRefreshFunc(ctx context.Context, client *aad.DomainServicesClient, resourceGroup, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("waiting for both Domain Service Controller up...")
		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil, "pending", err
		}
		if len(*resp.DomainControllerIPAddress) < 2 {
			return *resp.DomainControllerIPAddress, "pending", nil
		}
		return *resp.DomainControllerIPAddress, "available", nil
	}
}

func expandArmDomainServiceDomainSecuritySettings(input []interface{}) *aad.DomainSecuritySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	ntlmV1 := v["ntlm_v1"].(string)
	tlsV1 := v["tls_v1"].(string)
	syncNtlmPasswords := v["sync_ntlm_passwords"].(string)

	result := aad.DomainSecuritySettings{
		NtlmV1:            aad.NtlmV1(ntlmV1),
		SyncNtlmPasswords: aad.SyncNtlmPasswords(syncNtlmPasswords),
		TLSV1:             aad.TLSV1(tlsV1),
	}
	return &result
}

func expandArmDomainServiceLdapsSettings(input []interface{}) *aad.LdapsSettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	ldaps := v["ldaps"].(string)
	pfxCertificate := v["pfx_certificate"].(string)
	pfxCertificatePassword := v["pfx_certificate_password"].(string)
	externalAccess := v["external_access"].(string)

	result := aad.LdapsSettings{
		ExternalAccess:         aad.ExternalAccess(externalAccess),
		Ldaps:                  aad.Ldaps(ldaps),
		PfxCertificate:         &pfxCertificate,
		PfxCertificatePassword: &pfxCertificatePassword,
	}
	return &result
}

func expandArmDomainServiceNotificationSettings(input []interface{}) *aad.NotificationSettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	notifyGlobalAdmins := v["notify_global_admins"].(string)
	notifyDcAdmins := v["notify_dc_admins"].(string)
	additionalRecipients := v["additional_recipients"].([]interface{})

	result := aad.NotificationSettings{
		AdditionalRecipients: utils.ExpandStringSlice(additionalRecipients),
		NotifyDcAdmins:       aad.NotifyDcAdmins(notifyDcAdmins),
		NotifyGlobalAdmins:   aad.NotifyGlobalAdmins(notifyGlobalAdmins),
	}
	return &result
}

func flattenArmDomainServiceDomainSecuritySettings(input *aad.DomainSecuritySettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["ntlm_v1"] = string(input.NtlmV1)
	result["sync_ntlm_passwords"] = string(input.SyncNtlmPasswords)
	result["tls_v1"] = string(input.TLSV1)

	return []interface{}{result}
}

func flattenArmDomainServiceLdapsSettings(input *aad.LdapsSettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["external_access"] = string(input.ExternalAccess)
	result["ldaps"] = string(input.Ldaps)
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

func flattenArmDomainServiceNotificationSettings(input *aad.NotificationSettings) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["additional_recipients"] = utils.FlattenStringSlice(input.AdditionalRecipients)
	result["notify_dc_admins"] = string(input.NotifyDcAdmins)
	result["notify_global_admins"] = string(input.NotifyGlobalAdmins)

	return []interface{}{result}
}

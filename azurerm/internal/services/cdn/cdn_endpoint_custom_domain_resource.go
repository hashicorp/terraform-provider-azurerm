package cdn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/validate"
	keyvaultValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/keyvault/validate"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCdnEndpointCustomDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCdnEndpointCustomDomainCreate,
		Update: resourceArmCdnEndpointCustomDomainUpdate,
		Read:   resourceArmCdnEndpointCustomDomainRead,
		Delete: resourceArmCdnEndpointCustomDomainDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.EndpointCustomDomainID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Hour),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Hour),
			Delete: schema.DefaultTimeout(10 * time.Hour),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CdnEndpointCustomDomainName(),
			},

			"cdn_endpoint_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.EndpointID,
			},

			"host_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// Ipv6 (with square enclosed), Ipv4, domain name is allowed
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"cdn_managed_https_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.Shared),
								string(cdn.Dedicated),
							}, false),
						},
					},
				},
				ConflictsWith: []string{"user_managed_https_settings"},
			},

			"user_managed_https_settings": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsUUID,
						},
						"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),
						"vault_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: keyvaultValidate.VaultName,
						},
						"secret_name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: keyvaultValidate.NestedItemName,
						},
						"secret_version": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ConflictsWith: []string{"cdn_managed_https_settings"},
			},
		},
	}
}

func resourceArmCdnEndpointCustomDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	epid := d.Get("cdn_endpoint_id").(string)

	cdnEndpointId, err := parse.EndpointID(epid)
	if err != nil {
		return fmt.Errorf("parsing CDN Endpoint ID %q: %+v", epid, err)
	}

	existing, err := client.Get(ctx, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q): %+v",
				name, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, err)
		}
	}
	if existing.ID != nil && *existing.ID != "" {
		return tf.ImportAsExistsError("azurerm_cdn_endpoint_custom_domain", *existing.ID)
	}

	props := cdn.CustomDomainParameters{
		CustomDomainPropertiesParameters: &cdn.CustomDomainPropertiesParameters{
			HostName: utils.String(d.Get("host_name").(string)),
		},
	}

	future, err := client.Create(ctx, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, name, props)
	if err != nil {
		return fmt.Errorf("creating Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q): %+v",
			name, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q): %+v",
			name, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, err)
	}

	resp, err := client.Get(ctx, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q): %+v",
			name, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q): %+v",
			name, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, err)
	}

	// Enable https if specified
	var params cdn.BasicCustomDomainHTTPSParameters
	if v, ok := d.GetOk("user_managed_https_settings"); ok {
		// User managed certificate is only available for Azure CDN from Microsoft and Azure CDN from Verizon profiles.
		// https://docs.microsoft.com/en-us/azure/cdn/cdn-custom-ssl?tabs=option-2-enable-https-with-your-own-certificate#tlsssl-certificates
		pfClient := meta.(*clients.Client).Cdn.ProfilesClient
		cdnEndpointResp, err := pfClient.Get(ctx, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName)
		if err != nil {
			return fmt.Errorf("retrieving Cdn Profile %q (Resource Group %q): %+v",
				cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, err)
		}
		if cdnEndpointResp.Sku != nil && (cdnEndpointResp.Sku.Name != cdn.StandardMicrosoft && cdnEndpointResp.Sku.Name != cdn.StandardVerizon) {
			return errors.New("User managed HTTPS certificate is only available for Azure CDN from Microsoft or Azure CDN from Verizon profiles")
		}
		params = expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(v.([]interface{}))
	} else if v, ok := d.GetOk("cdn_managed_https_settings"); ok {
		params = expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(v.([]interface{}))
	}
	if params != nil {
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, name, params); err != nil {
			enableErr := err
			// Rollback the creation
			future, err := client.Delete(ctx, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, name)
			if err != nil {
				return fmt.Errorf("%+v. Addtionally, failed to delete (rollback) Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q): %+v",
					enableErr, name, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("%+v. Additionally, failed to wait for deletion (rollback) of Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q): %+v",
					enableErr, name, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, err)
			}
			return enableErr
		}
	}

	d.SetId(*resp.ID)

	return resourceArmCdnEndpointCustomDomainRead(d, meta)
}

func resourceArmCdnEndpointCustomDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName)
	if err != nil {
		return fmt.Errorf("retrieving Cdn Endpoint Custom Domain %q: %+v",
			id, err)
	}

	switch {
	case d.HasChange("cdn_managed_https_settings") && d.HasChange("user_managed_https_settings"):
		// One is turned on and the other is turned off
		return fmt.Errorf("in-place update on enabled HTTPS settings is not supported on Cdn Endpoint Custom Domain %q", id)
	case d.HasChange("cdn_managed_https_settings"):
		props := resp.CustomDomainProperties
		if props == nil {
			return errors.New("unexpected nil of `CustomDomainProperties` in response")
		}
		if props.CustomHTTPSParameters == nil {
			// disabled -> enabled
			if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName,
				expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(d.Get("cdn_managed_https_settings").([]interface{}))); err != nil {
				return fmt.Errorf("enable HTTPS on Cdn Endpoint Custom Domain %q: %+v", id, err)
			}
		} else {
			params := expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(d.Get("cdn_managed_https_settings").([]interface{}))
			if params == nil {
				// enabled -> disabled
				if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName); err != nil {
					return fmt.Errorf("disable HTTPS on Cdn Endpoint Custom Domain %q: %+v", id, err)
				}
			} else {
				return fmt.Errorf("in-place update on enabled HTTPS settings is not supported on Cdn Endpoint Custom Domain %q", id)
			}
		}
	case d.HasChange("user_managed_https_settings"):
		props := resp.CustomDomainProperties
		if props == nil {
			return errors.New("unexpected nil of `CustomDomainProperties` in response")
		}
		if props.CustomHTTPSParameters == nil {
			// disabled -> enabled
			if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName,
				expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(d.Get("user_managed_https_settings").([]interface{}))); err != nil {
				return fmt.Errorf("enable HTTPS on Cdn Endpoint Custom Domain %q: %+v", id, err)
			}
		} else {
			params := expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(d.Get("user_managed_https_settings").([]interface{}))
			if params == nil {
				// enabled -> disabled
				if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName); err != nil {
					return fmt.Errorf("disable HTTPS on Cdn Endpoint Custom Domain %q: %+v", id, err)
				}
			} else {
				return fmt.Errorf("in-place update on enabled HTTPS settings is not supported on Cdn Endpoint Custom Domain %q", id)
			}
		}
	}
	return nil
}

func resourceArmCdnEndpointCustomDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	epClient := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Cdn Endpoint Custom Domain %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Cdn Endpoint Custom Domain %q: %+v", id, err)
	}
	cdnEndpointResp, err := epClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName)
	if err != nil {
		return fmt.Errorf("retrieving Cdn Endpoint %q (Resource Group %q / Profile %q): %+v",
			id.EndpointName, id.ResourceGroup, id.ProfileName, err)
	}

	d.Set("name", resp.Name)
	d.Set("cdn_endpoint_id", cdnEndpointResp.ID)
	if props := resp.CustomDomainProperties; props != nil {
		d.Set("host_name", props.HostName)

		switch params := props.CustomHTTPSParameters.(type) {
		case *cdn.ManagedHTTPSParameters:
			if err := d.Set("cdn_managed_https_settings", flattenArmCdnEndpointCustomDomainCdnManagedHttpsSettings(params)); err != nil {
				return fmt.Errorf("setting `cdn_managed_https_settings`: %+v", err)
			}
		case *cdn.UserManagedHTTPSParameters:
			if err := d.Set("user_managed_https_settings", flattenArmCdnEndpointCustomDomainUserManagedHttpsSettings(params)); err != nil {
				return fmt.Errorf("setting `user_managed_https_settings`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmCdnEndpointCustomDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName)
	if err != nil {
		return fmt.Errorf("deleting Cdn Endpoint Custom Domain %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of Cdn Endpoint Custom Domain %q: %+v", id, err)
	}

	return nil
}

func expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(input []interface{}) cdn.BasicCustomDomainHTTPSParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &cdn.ManagedHTTPSParameters{
		CertificateSourceParameters: &cdn.CertificateSourceParameters{
			OdataType:       utils.String("#Microsoft.Azure.Cdn.Models.CdnCertificateSourceParameters"),
			CertificateType: cdn.CertificateType(raw["certificate_type"].(string)),
		},
		CertificateSource: cdn.CertificateSourceCdn,
		ProtocolType:      cdn.IPBased,
		MinimumTLSVersion: cdn.None,
	}

	return output
}

func expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(input []interface{}) cdn.BasicCustomDomainHTTPSParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &cdn.UserManagedHTTPSParameters{
		CertificateSourceParameters: &cdn.KeyVaultCertificateSourceParameters{
			OdataType:         utils.String("#Microsoft.Azure.Cdn.Models.KeyVaultCertificateSourceParameters"),
			SubscriptionID:    utils.String(raw["subscription_id"].(string)),
			ResourceGroupName: utils.String(raw["resource_group_name"].(string)),
			VaultName:         utils.String(raw["vault_name"].(string)),
			SecretName:        utils.String(raw["secret_name"].(string)),
			SecretVersion:     utils.String(raw["secret_version"].(string)),
			UpdateRule:        utils.String("NoAction"),
			DeleteRule:        utils.String("NoAction"),
		},
		CertificateSource: cdn.CertificateSourceAzureKeyVault,
		ProtocolType:      cdn.ServerNameIndication,
		MinimumTLSVersion: cdn.None,
	}

	return output
}

func flattenArmCdnEndpointCustomDomainCdnManagedHttpsSettings(input *cdn.ManagedHTTPSParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	certificateType := ""
	if params := input.CertificateSourceParameters; params != nil {
		certificateType = string(params.CertificateType)
	}

	return []interface{}{
		map[string]interface{}{
			"certificate_type": certificateType,
		},
	}
}

func flattenArmCdnEndpointCustomDomainUserManagedHttpsSettings(input *cdn.UserManagedHTTPSParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var (
		subscriptionId    string
		resourceGroupName string
		vaultName         string
		secretName        string
		secretVersion     string
	)
	if params := input.CertificateSourceParameters; params != nil {
		if params.SubscriptionID != nil {
			subscriptionId = *params.SubscriptionID
		}
		if params.ResourceGroupName != nil {
			resourceGroupName = *params.ResourceGroupName
		}
		if params.VaultName != nil {
			vaultName = *params.VaultName
		}
		if params.SecretName != nil {
			secretName = *params.SecretName
		}
		if params.SecretVersion != nil {
			secretVersion = *params.SecretVersion
		}
	}

	return []interface{}{
		map[string]interface{}{
			"subscription_id":     subscriptionId,
			"resource_group_name": resourceGroupName,
			"vault_name":          vaultName,
			"secret_name":         secretName,
			"secret_version":      secretVersion,
		},
	}
}

func enableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, resgrp, profile, endpoint, name string, params cdn.BasicCustomDomainHTTPSParameters) error {
	_, err := client.EnableCustomHTTPS(ctx, resgrp, profile, endpoint, name, &params)
	if err != nil {
		return fmt.Errorf("sending enable request: %+v", err)
	}

	log.Printf("[DEBUG] Waiting for HTTPS to enable on Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q)",
		name, resgrp, profile, endpoint)
	deadline, _ := ctx.Deadline()
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(cdn.Enabling)},
		Target:     []string{string(cdn.Enabled), string(cdn.Failed)},
		Refresh:    cdnEndpointCustomDomainHttpsRefreshFunc(ctx, client, resgrp, profile, endpoint, name),
		MinTimeout: 10 * time.Second,
		Timeout:    time.Until(deadline),
	}

	state, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for HTTPS provision state: %+v", err)
	}
	if state == cdn.Failed {
		return errors.New("HTTPS provision state is Failed")
	}

	return nil
}

func disableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, resgrp, profile, endpoint, name string) error {
	_, err := client.DisableCustomHTTPS(ctx, resgrp, profile, endpoint, name)
	if err != nil {
		return fmt.Errorf("sending disable request: %+v", err)
	}

	log.Printf("[DEBUG] Waiting for HTTPS to disable on Cdn Endpoint Custom Domain %q (Resource Group %q / Profile %q / Endpoint %q)",
		name, resgrp, profile, endpoint)
	deadline, _ := ctx.Deadline()
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(cdn.Disabling)},
		Target:     []string{string(cdn.Disabled), string(cdn.Failed)},
		Refresh:    cdnEndpointCustomDomainHttpsRefreshFunc(ctx, client, resgrp, profile, endpoint, name),
		MinTimeout: 10 * time.Second,
		Timeout:    time.Until(deadline),
	}

	state, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for HTTPS provision state: %+v", err)
	}
	if state == cdn.Failed {
		return errors.New("HTTPS provision state is Failed")
	}

	return nil
}

func cdnEndpointCustomDomainHttpsRefreshFunc(ctx context.Context, client *cdn.CustomDomainsClient, resgrp string, profile string, endpoint string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resgrp, profile, endpoint, name)
		if err != nil {
			return nil, "", fmt.Errorf("retrieving HTTPS provisioning state: %+v", err)
		}

		props := res.CustomDomainProperties
		if props == nil {
			return nil, "", errors.New("unexpected nil of `CustomDomainProperties` in response")
		}

		return props.CustomHTTPSProvisioningState, string(props.CustomHTTPSProvisioningState), nil
	}
}

package cdn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"

	keyvaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyvaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmCdnEndpointCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmCdnEndpointCustomDomainCreate,
		Read:   resourceArmCdnEndpointCustomDomainRead,
		Update: resourceArmCdnEndpointCustomDomainUpdate,
		Delete: resourceArmCdnEndpointCustomDomainDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CustomDomainID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(20 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(20 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(20 * time.Hour),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CdnEndpointCustomDomainName(),
			},

			"cdn_endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.EndpointID,
			},

			"host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"cdn_managed_https_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"certificate_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.Shared),
								string(cdn.Dedicated),
							}, false),
						},
						"protocol_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.ServerNameIndication),
								string(cdn.IPBased),
							}, false),
						},
					},
				},
				ConflictsWith: []string{"user_managed_https_settings"},
			},

			"user_managed_https_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyvaultValidate.VaultID,
						},
						"secret_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: keyvaultValidate.NestedItemName,
						},
						"secret_version": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
				ConflictsWith: []string{"cdn_managed_https_settings"},
			},
		},
		CustomizeDiff: func(ctx context.Context, diff *pluginsdk.ResourceDiff, _ interface{}) error {
			if settings, ok := diff.GetOk("cdn_managed_https_settings"); ok {
				settings := settings.([]interface{})[0].(map[string]interface{})
				cert, protocol := settings["certificate_type"].(string), settings["protocol_type"].(string)
				if cert == string(cdn.Shared) && protocol != string(cdn.ServerNameIndication) {
					return fmt.Errorf("`certificate_type = Shared` has to be used together with `protocol_type = ServerNameIndication`")
				}
				if cert == string(cdn.Dedicated) && protocol != string(cdn.IPBased) {
					return fmt.Errorf("`certificate_type = Dedicated` has to be used together with `protocol_type = IPBased`")
				}
			}
			return nil
		},
	}
}

func resourceArmCdnEndpointCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	epid := d.Get("cdn_endpoint_id").(string)

	cdnEndpointId, err := parse.EndpointID(epid)
	if err != nil {
		return err
	}

	id := parse.NewCustomDomainID(cdnEndpointId.SubscriptionId, cdnEndpointId.ResourceGroup, cdnEndpointId.ProfileName, cdnEndpointId.Name, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %q: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_endpoint_custom_domain", id.ID())
	}

	props := cdn.CustomDomainParameters{
		CustomDomainPropertiesParameters: &cdn.CustomDomainPropertiesParameters{
			HostName: utils.String(d.Get("host_name").(string)),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name, props)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	// Enable https if specified
	var params cdn.BasicCustomDomainHTTPSParameters
	if v, ok := d.GetOk("user_managed_https_settings"); ok {
		// User managed certificate is only available for Azure CDN from Microsoft and Azure CDN from Verizon profiles.
		// https://docs.microsoft.com/en-us/azure/cdn/cdn-custom-ssl?tabs=option-2-enable-https-with-your-own-certificate#tlsssl-certificates
		pfClient := meta.(*clients.Client).Cdn.ProfilesClient
		cdnEndpointResp, err := pfClient.Get(ctx, id.ResourceGroup, id.ProfileName)
		if err != nil {
			return fmt.Errorf("retrieving Cdn Profile %q (Resource Group %q): %+v",
				id.ResourceGroup, id.ProfileName, err)
		}
		if cdnEndpointResp.Sku != nil && (cdnEndpointResp.Sku.Name != cdn.StandardMicrosoft && cdnEndpointResp.Sku.Name != cdn.StandardVerizon) {
			return errors.New("user managed HTTPS certificate is only available for Azure CDN from Microsoft or Azure CDN from Verizon profiles")
		}
		params, err = expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(v.([]interface{}))
		if err != nil {
			return err
		}
	} else if v, ok := d.GetOk("cdn_managed_https_settings"); ok {
		params = expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(v.([]interface{}))
	}

	if params != nil {
		if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, id, params); err != nil {
			enableErr := err
			// Rollback the creation
			future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
			if err != nil {
				return fmt.Errorf("%+v. Addtionally, failed to delete (rollback) Cdn Endpoint Custom Domain %q: %+v",
					enableErr, id, err)
			}
			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("%+v. Additionally, failed to wait for deletion (rollback) of Cdn Endpoint Custom Domain %q: %+v",
					enableErr, id, err)
			}
			return enableErr
		}
	}

	d.SetId(id.ID())

	return resourceArmCdnEndpointCustomDomainRead(d, meta)
}

func resourceArmCdnEndpointCustomDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	switch {
	case d.HasChange("cdn_managed_https_settings") && d.HasChange("user_managed_https_settings"):
		// One is turned on, and the other is turned off
		return fmt.Errorf("in-place update on enabled HTTPS settings is not supported on %q", id)
	case d.HasChange("cdn_managed_https_settings"):
		props := resp.CustomDomainProperties
		if props == nil {
			return errors.New("unexpected nil of `CustomDomainProperties` in response")
		}
		if props.CustomHTTPSParameters == nil {
			// disabled -> enabled
			if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, *id,
				expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(d.Get("cdn_managed_https_settings").([]interface{}))); err != nil {
				return fmt.Errorf("enable HTTPS on %q: %+v", id, err)
			}
		} else {
			params := expandArmCdnEndpointCustomDomainCdnManagedHttpsSettings(d.Get("cdn_managed_https_settings").([]interface{}))
			if params == nil {
				// enabled -> disabled
				if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, *id); err != nil {
					return fmt.Errorf("disable HTTPS on %q: %+v", id, err)
				}
			} else {
				return fmt.Errorf("in-place update on enabled HTTPS settings is not supported on %q", id)
			}
		}
	case d.HasChange("user_managed_https_settings"):
		props := resp.CustomDomainProperties
		if props == nil {
			return errors.New("unexpected nil of `CustomDomainProperties` in response")
		}
		param, err := expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(d.Get("user_managed_https_settings").([]interface{}))
		if err != nil {
			return err
		}
		if props.CustomHTTPSParameters == nil {
			// disabled -> enabled
			if err := enableArmCdnEndpointCustomDomainHttps(ctx, client, *id, param); err != nil {
				return fmt.Errorf("enable HTTPS on %q: %+v", id, err)
			}
		} else {
			if param == nil {
				// enabled -> disabled
				if err := disableArmCdnEndpointCustomDomainHttps(ctx, client, *id); err != nil {
					return fmt.Errorf("disable HTTPS on %q: %+v", id, err)
				}
			} else {
				return fmt.Errorf("in-place update on enabled HTTPS settings is not supported on %q", id)
			}
		}
	}
	return resourceArmCdnEndpointCustomDomainRead(d, meta)
}

func resourceArmCdnEndpointCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	cdnEndpointId := parse.NewEndpointID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.EndpointName)

	d.Set("name", resp.Name)
	d.Set("cdn_endpoint_id", cdnEndpointId.ID())
	if props := resp.CustomDomainProperties; props != nil {
		d.Set("host_name", props.HostName)
		switch params := props.CustomHTTPSParameters.(type) {
		case cdn.ManagedHTTPSParameters:
			if err := d.Set("cdn_managed_https_settings", flattenArmCdnEndpointCustomDomainCdnManagedHttpsSettings(params)); err != nil {
				return fmt.Errorf("setting `cdn_managed_https_settings`: %+v", err)
			}
		case cdn.UserManagedHTTPSParameters:
			if err := d.Set("user_managed_https_settings", flattenArmCdnEndpointCustomDomainUserManagedHttpsSettings(params)); err != nil {
				return fmt.Errorf("setting `user_managed_https_settings`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmCdnEndpointCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.CustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CustomDomainID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
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
		ProtocolType:      cdn.ProtocolType(raw["protocol_type"].(string)),
		MinimumTLSVersion: cdn.MinimumTLSVersionNone,
	}

	return output
}

func expandArmCdnEndpointCustomDomainUserManagedHttpsSettings(input []interface{}) (cdn.BasicCustomDomainHTTPSParameters, error) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})

	keyVaultId, err := keyvaultParse.VaultID(raw["key_vault_id"].(string))
	if err != nil {
		return nil, err
	}

	output := &cdn.UserManagedHTTPSParameters{
		CertificateSourceParameters: &cdn.KeyVaultCertificateSourceParameters{
			OdataType:         utils.String("#Microsoft.Azure.Cdn.Models.KeyVaultCertificateSourceParameters"),
			SubscriptionID:    &keyVaultId.SubscriptionId,
			ResourceGroupName: &keyVaultId.ResourceGroup,
			VaultName:         &keyVaultId.Name,
			SecretName:        utils.String(raw["secret_name"].(string)),
			SecretVersion:     utils.String(raw["secret_version"].(string)),
			UpdateRule:        utils.String("NoAction"),
			DeleteRule:        utils.String("NoAction"),
		},
		CertificateSource: cdn.CertificateSourceAzureKeyVault,
		ProtocolType:      cdn.ServerNameIndication,
		MinimumTLSVersion: cdn.MinimumTLSVersionNone,
	}

	return output, nil
}

func flattenArmCdnEndpointCustomDomainCdnManagedHttpsSettings(input cdn.ManagedHTTPSParameters) []interface{} {
	certificateType := ""
	if params := input.CertificateSourceParameters; params != nil {
		certificateType = string(params.CertificateType)
	}

	return []interface{}{
		map[string]interface{}{
			"certificate_type": certificateType,
			"protocol_type":    string(input.ProtocolType),
		},
	}
}

func flattenArmCdnEndpointCustomDomainUserManagedHttpsSettings(input cdn.UserManagedHTTPSParameters) []interface{} {
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
			"key_vault_id":   keyvaultParse.NewVaultID(subscriptionId, resourceGroupName, vaultName).ID(),
			"secret_name":    secretName,
			"secret_version": secretVersion,
		},
	}
}

func enableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, id parse.CustomDomainId, params cdn.BasicCustomDomainHTTPSParameters) error {
	_, err := client.EnableCustomHTTPS(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name, &params)
	if err != nil {
		return fmt.Errorf("sending enable request: %+v", err)
	}

	log.Printf("[DEBUG] Waiting for HTTPS to enable on %q", id)
	deadline, _ := ctx.Deadline()
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(cdn.CustomHTTPSProvisioningStateEnabling)},
		Target:     []string{string(cdn.CustomHTTPSProvisioningStateEnabled), string(cdn.CustomHTTPSProvisioningStateFailed), string(cdn.CustomHTTPSProvisioningStateDisabled)},
		Refresh:    cdnEndpointCustomDomainHttpsRefreshFunc(ctx, client, id),
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
	if state == cdn.Disabled {
		return errors.New("HTTPS provision state is back to Disabled")
	}

	return nil
}

func disableArmCdnEndpointCustomDomainHttps(ctx context.Context, client *cdn.CustomDomainsClient, id parse.CustomDomainId) error {
	_, err := client.DisableCustomHTTPS(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return fmt.Errorf("sending disable request: %+v", err)
	}

	log.Printf("[DEBUG] Waiting for HTTPS to disable on %q", id)
	deadline, _ := ctx.Deadline()
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(cdn.CustomHTTPSProvisioningStateDisabling)},
		Target:     []string{string(cdn.CustomHTTPSProvisioningStateDisabled), string(cdn.CustomHTTPSProvisioningStateFailed), string(cdn.CustomHTTPSProvisioningStateEnabled)},
		Refresh:    cdnEndpointCustomDomainHttpsRefreshFunc(ctx, client, id),
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
	if state == cdn.Enabled {
		return errors.New("HTTPS provision state is back to Enabled")
	}

	return nil
}

func cdnEndpointCustomDomainHttpsRefreshFunc(ctx context.Context, client *cdn.CustomDomainsClient, id parse.CustomDomainId) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
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

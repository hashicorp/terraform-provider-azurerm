package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdcustomdomains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontdoorCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorCustomDomainCreate,
		Read:   resourceFrontdoorCustomDomainRead,
		Update: resourceFrontdoorCustomDomainUpdate,
		Delete: resourceFrontdoorCustomDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontdoorCustomDomainID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ProfileID,
			},

			"azure_dns_zone_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"deployment_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"domain_validation_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
			},

			"pre_validated_custom_domain_resource_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tls_settings": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"certificate_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							Default:  string(cdn.AfdCertificateTypeCustomerCertificate),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.AfdCertificateTypeCustomerCertificate),
								string(cdn.AfdCertificateTypeManagedCertificate),
							}, false),
						},

						"minimum_tls_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.AfdMinimumTLSVersionTLS12),
							ValidateFunc: validation.StringInSlice([]string{
								string(cdn.AfdMinimumTLSVersionTLS10),
								string(cdn.AfdMinimumTLSVersionTLS12),
							}, false),
						},

						"secret_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"validation_properties": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"expiration_date": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"validation_token": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceFrontdoorCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	// In the SDK the namespace is Microsoft.CDN in Terraform the namespace is Microsoft.Cdn
	sdkId := afdcustomdomains.NewCustomDomainID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))
	id := parse.NewFrontdoorCustomDomainID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, sdkId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_frontdoor_custom_domain", id.ID())
		}
	}

	props := afdcustomdomains.AFDDomain{
		Properties: &afdcustomdomains.AFDDomainProperties{
			AzureDnsZone:                       expandCustomDomainResourceReference(d.Get("azure_dns_zone").(string)),
			HostName:                           d.Get("host_name").(string),
			PreValidatedCustomDomainResourceId: expandCustomDomainResourceReference(d.Get("pre_validated_custom_domain_resource_id").(string)),
			TlsSettings:                        expandCustomDomainAFDDomainHttpsParameters(d.Get("tls_settings").([]interface{})),
		},
	}
	if err := client.CreateThenPoll(ctx, sdkId, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorCustomDomainRead(d, meta)
}

func resourceFrontdoorCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdcustomdomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *sdkId)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CustomDomainName)

	d.Set("frontdoor_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {

			if err := d.Set("azure_dns_zone", flattenCustomDomainResourceReference(props.AzureDnsZone)); err != nil {
				return fmt.Errorf("setting `azure_dns_zone`: %+v", err)
			}
			d.Set("deployment_status", props.DeploymentStatus)
			d.Set("domain_validation_state", props.DomainValidationState)
			d.Set("host_name", props.HostName)

			if err := d.Set("pre_validated_custom_domain_resource_id", flattenCustomDomainResourceReference(props.PreValidatedCustomDomainResourceId)); err != nil {
				return fmt.Errorf("setting `pre_validated_custom_domain_resource_id`: %+v", err)
			}
			d.Set("profile_name", props.ProfileName)
			d.Set("provisioning_state", props.ProvisioningState)

			if err := d.Set("tls_settings", flattenCustomDomainAFDDomainHttpsParameters(props.TlsSettings)); err != nil {
				return fmt.Errorf("setting `tls_settings`: %+v", err)
			}

			if err := d.Set("validation_properties", flattenCustomDomainDomainValidationProperties(props.ValidationProperties)); err != nil {
				return fmt.Errorf("setting `validation_properties`: %+v", err)
			}
		}
	}
	return nil
}

func resourceFrontdoorCustomDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdcustomdomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	props := afdcustomdomains.AFDDomainUpdateParameters{
		Properties: &afdcustomdomains.AFDDomainUpdatePropertiesParameters{
			AzureDnsZone:                       expandCustomDomainResourceReference(d.Get("azure_dns_zone").([]interface{})),
			PreValidatedCustomDomainResourceId: expandCustomDomainResourceReference(d.Get("pre_validated_custom_domain_resource_id").([]interface{})),
			TlsSettings:                        expandCustomDomainAFDDomainHttpsParameters(d.Get("tls_settings").([]interface{})),
		},
	}
	if err := client.UpdateThenPoll(ctx, *sdkId, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorCustomDomainRead(d, meta)
}

func resourceFrontdoorCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sdkId, err := afdcustomdomains.ParseCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *sdkId); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandCustomDomainResourceReference(input interface{}) *afdcustomdomains.ResourceReference {
	if input == nil {
		return nil
	}

	return &afdcustomdomains.ResourceReference{
		Id: utils.String(input.(string)),
	}
}

func expandCustomDomainAFDDomainHttpsParameters(input []interface{}) *afdcustomdomains.AFDDomainHttpsParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	certificateTypeValue := afdcustomdomains.AfdCertificateType(v["certificate_type"].(string))
	minimumTlsVersionValue := afdcustomdomains.AfdMinimumTlsVersion(v["minimum_tls_version"].(string))
	return &afdcustomdomains.AFDDomainHttpsParameters{
		CertificateType:   certificateTypeValue,
		MinimumTlsVersion: &minimumTlsVersionValue,
		Secret:            expandCustomDomainResourceReference(v["secret_id"].(string)),
	}
}

func flattenCustomDomainAFDDomainHttpsParameters(input *afdcustomdomains.AFDDomainHttpsParameters) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})
	result["certificate_type"] = input.CertificateType

	if input.MinimumTlsVersion != nil {
		result["minimum_tls_version"] = *input.MinimumTlsVersion
	}

	result["secret_id"] = *flattenCustomDomainResourceReference(input.Secret)
	return append(results, result)
}

func flattenCustomDomainDomainValidationProperties(input *afdcustomdomains.DomainValidationProperties) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	if input.ExpirationDate != nil {
		result["expiration_date"] = *input.ExpirationDate
	}

	if input.ValidationToken != nil {
		result["validation_token"] = *input.ValidationToken
	}
	return append(results, result)
}

func flattenCustomDomainResourceReference(input *afdcustomdomains.ResourceReference) *string {
	if input == nil {
		return nil
	}

	return input.Id
}

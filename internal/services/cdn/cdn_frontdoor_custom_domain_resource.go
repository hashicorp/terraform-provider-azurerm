package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	track1 "github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorCustomDomainCreate,
		Read:   resourceCdnFrontdoorCustomDomainRead,
		Update: resourceCdnFrontdoorCustomDomainUpdate,
		Delete: resourceCdnFrontdoorCustomDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			// TODO: these timeouts need adjusting..?
			Create: pluginsdk.DefaultTimeout(12 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(24 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(12 * time.Hour),
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
				// TODO: missing validation
			},

			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorProfileID,
			},

			"dns_zone_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// TODO: this should be validating the DNS Zone ID - if this could be both Public or Private we can validate that
				ValidateFunc: azure.ValidateResourceID,
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

			"pre_validated_cdn_frontdoor_custom_domain_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.FrontdoorCustomDomainID,
			},

			"tls": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"certificate_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(cdn.AfdCertificateTypeManagedCertificate),
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

						"cdn_frontdoor_secret_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.FrontdoorSecretID,
						},
					},
				},
			},

			// TODO: can this be removed?
			"cdn_frontdoor_profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"validation_properties": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						// TODO: we should probably make these top-levle fields
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

func resourceCdnFrontdoorCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontdoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorCustomDomainID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
		}
	}

	props := track1.AFDDomain{
		AFDDomainProperties: &track1.AFDDomainProperties{
			HostName:                           utils.String(d.Get("host_name").(string)),
			AzureDNSZone:                       expandResourceReference(d.Get("dns_zone_id").(string)),
			PreValidatedCustomDomainResourceID: expandResourceReference(d.Get("pre_validated_cdn_frontdoor_custom_domain_id").(string)),
			TLSSettings:                        expandCustomDomainAFDDomainHttpsParameters(d.Get("tls").([]interface{})),
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontdoorCustomDomainRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CustomDomainName)
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontdoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.AFDDomainProperties; props != nil {
		d.Set("domain_validation_state", props.DomainValidationState)
		d.Set("host_name", props.HostName)
		d.Set("cdn_frontdoor_profile_name", id.ProfileName) // TODO: this either needs to be returned from the API or removed

		if err := d.Set("dns_zone_id", flattenResourceReference(props.AzureDNSZone)); err != nil {
			return fmt.Errorf("setting `dns_zone_id`: %+v", err)
		}

		if err := d.Set("pre_validated_cdn_frontdoor_custom_domain_id", flattenResourceReference(props.PreValidatedCustomDomainResourceID)); err != nil {
			return fmt.Errorf("setting `pre_validated_cdn_frontdoor_custom_domain_id`: %+v", err)
		}

		if err := d.Set("tls", flattenCustomDomainAFDDomainHttpsParameters(props.TLSSettings)); err != nil {
			return fmt.Errorf("setting `tls`: %+v", err)
		}

		if err := d.Set("validation_properties", flattenCustomDomainDomainValidationProperties(props.ValidationProperties)); err != nil {
			return fmt.Errorf("setting `validation_properties`: %+v", err)
		}
	}

	return nil
}

func resourceCdnFrontdoorCustomDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	props := track1.AFDDomainUpdateParameters{
		AFDDomainUpdatePropertiesParameters: &track1.AFDDomainUpdatePropertiesParameters{
			AzureDNSZone:                       expandResourceReference(d.Get("dns_zone_id").(string)),
			PreValidatedCustomDomainResourceID: expandResourceReference(d.Get("pre_validated_cdn_frontdoor_custom_domain_id").(string)),
			TLSSettings:                        expandCustomDomainAFDDomainHttpsParameters(d.Get("tls").([]interface{})),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontdoorCustomDomainRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandCustomDomainAFDDomainHttpsParameters(input []interface{}) *track1.AFDDomainHTTPSParameters {
	if len(input) == 0 || input[0] == nil {
		// TODO: why is this returning nil and not an empty object?
		return nil
	}

	v := input[0].(map[string]interface{})

	return &track1.AFDDomainHTTPSParameters{
		CertificateType:   track1.AfdCertificateType(v["certificate_type"].(string)),
		MinimumTLSVersion: track1.AfdMinimumTLSVersion(v["minimum_tls_version"].(string)),
		Secret:            expandResourceReference(v["cdn_frontdoor_secret_id"].(string)),
	}
}

func flattenCustomDomainAFDDomainHttpsParameters(input *track1.AFDDomainHTTPSParameters) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"cdn_frontdoor_secret_id": flattenResourceReference(input.Secret),
			"certificate_type":        string(input.CertificateType),
			"minimum_tls_version":     string(input.MinimumTLSVersion),
		},
	}
}

func flattenCustomDomainDomainValidationProperties(input *track1.DomainValidationProperties) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	expirationDate := ""
	if input.ExpirationDate != nil {
		expirationDate = *input.ExpirationDate
	}
	validationToken := ""
	if input.ValidationToken != nil {
		validationToken = *input.ValidationToken
	}
	return []interface{}{
		map[string]interface{}{
			"expiration_date":  expirationDate,
			"validation_token": validationToken,
		},
	}
}

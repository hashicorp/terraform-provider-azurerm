package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	webValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorCustomDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainCreate,
		Read:   resourceCdnFrontDoorCustomDomainRead,
		Update: resourceCdnFrontDoorCustomDomainUpdate,
		Delete: resourceCdnFrontDoorCustomDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			// NOTE: These timeouts are extreamly long due to the manual
			// step of approving the private link if defined.
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorCustomDomainName,
			},

			// NOTE: I need this fake field because I need the profile name during the create operation
			"cdn_frontdoor_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorProfileID,
			},

			"dns_zone_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: dnsValidate.ValidateDnsZoneID,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
			},

			"pre_validated_custom_domain_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: webValidate.StaticSiteID,
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
								string(cdn.AfdCertificateTypeAzureFirstPartyManagedCertificate),
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

						// NOTE: If the secret is managed by FrontDoor this will cause a perpetual diff,
						// so this has to be an optional computed field.
						"cdn_frontdoor_secret_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validate.FrontDoorSecretID,
						},
					},
				},
			},

			"expiration_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"validation_token": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdnFrontDoorCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontDoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontdoorCustomDomainID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
	}

	preValidatedDomain := d.Get("pre_validated_custom_domain_id").(string)
	dnsZone := d.Get("dns_zone_id").(string)
	tls := d.Get("tls").([]interface{})

	// DNS Zone field is not supported for a pre-validated custom domain
	if preValidatedDomain != "" && dnsZone != "" {
		return fmt.Errorf("the 'dns_zone_id' field is not supported if the 'pre_validated_custom_domain_id' is passed")
	}

	props := cdn.AFDDomain{
		AFDDomainProperties: &cdn.AFDDomainProperties{
			HostName: utils.String(d.Get("host_name").(string)),
		},
	}

	// TODO: Still figuring out the preValidatedDomain bit...
	// Validate and set TLS settings
	if preValidatedDomain != "" {
		props.AFDDomainProperties.PreValidatedCustomDomainResourceID = expandResourceReference(preValidatedDomain)

		tlsSettings, err := expandTlsParameters(tls, true)
		if err != nil {
			return err
		}

		props.AFDDomainProperties.TLSSettings = tlsSettings
	} else {
		if dnsZone != "" {
			props.AFDDomainProperties.AzureDNSZone = expandResourceReference(dnsZone)
		}

		tlsSettings, err := expandTlsParameters(tls, false)
		if err != nil {
			return err
		}

		props.AFDDomainProperties.TLSSettings = tlsSettings
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("cdn_frontdoor_profile_id", parse.NewFrontDoorProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID())

	if props := resp.AFDDomainProperties; props != nil {
		d.Set("host_name", props.HostName)

		if err := d.Set("dns_zone_id", flattenResourceReference(props.AzureDNSZone)); err != nil {
			return fmt.Errorf("setting `dns_zone_id`: %+v", err)
		}

		if err := d.Set("pre_validated_custom_domain_id", flattenResourceReference(props.PreValidatedCustomDomainResourceID)); err != nil {
			return fmt.Errorf("setting `pre_validated_custom_domain_id`: %+v", err)
		}

		if err := d.Set("tls", flattenCustomDomainAFDDomainHttpsParameters(props.TLSSettings)); err != nil {
			return fmt.Errorf("setting `tls`: %+v", err)
		}

		if validationProps := props.ValidationProperties; validationProps != nil {
			d.Set("expiration_date", validationProps.ExpirationDate)
			d.Set("validation_token", validationProps.ValidationToken)
		}
	}

	return nil
}

func resourceCdnFrontDoorCustomDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontdoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	dnsZone := d.Get("dns_zone_id").(string)
	tls := d.Get("tls").([]interface{})
	preValidatedDomain := d.Get("pre_validated_custom_domain_id").(string)

	// DNS Zone field is not supported for a pre-validated custom domain
	if preValidatedDomain != "" && dnsZone != "" {
		return fmt.Errorf("the 'dns_zone_id' field is not supported if the 'pre_validated_custom_domain_id' is passed")
	}

	props := cdn.AFDDomainUpdateParameters{
		AFDDomainUpdatePropertiesParameters: &cdn.AFDDomainUpdatePropertiesParameters{},
	}

	// TODO: Resturcture update based on feed back in PR:
	// since there's 3 of these and this is a separate update method, can we conditionally set these (using if d.HasChange("dns_zone_id") { .. } etc

	// Validate and set TLS settings
	if preValidatedDomain != "" {
		props.AFDDomainUpdatePropertiesParameters.PreValidatedCustomDomainResourceID = expandResourceReference(preValidatedDomain)

		props.AFDDomainUpdatePropertiesParameters.TLSSettings, err = expandTlsParameters(tls, true)
		if err != nil {
			return err
		}
	} else {
		if dnsZone != "" {
			props.AFDDomainUpdatePropertiesParameters.AzureDNSZone = expandResourceReference(dnsZone)
		}

		props.AFDDomainUpdatePropertiesParameters.TLSSettings, err = expandTlsParameters(tls, false)
		if err != nil {
			return err
		}
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", *id, err)
	}

	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

func expandTlsParameters(input []interface{}, isPreValidatedDomain bool) (*cdn.AFDDomainHTTPSParameters, error) {
	if len(input) == 0 || input[0] == nil {
		// NOTE: With the Frontdoor service, they do not treat an empty object like an empty object
		// if it is not nil they assume it is fully defined and then end up throwing errors when they
		// attempt to get a value from one of the fields.
		return nil, nil
	}

	v := input[0].(map[string]interface{})

	certType := v["certificate_type"].(string)
	secret := v["cdn_frontdoor_secret_id"].(string)
	minTlsVersion := v["minimum_tls_version"].(string)

	tls := cdn.AFDDomainHTTPSParameters{}

	// NOTE: If this is a pre-validated domain you cannot pass a secret?
	if isPreValidatedDomain {
		if secret != "" {
			return nil, fmt.Errorf("the 'cdn_frontdoor_secret_id' field is not supported if the 'pre_validated_custom_domain_id' is passed")
		}
	} else {
		if tls.CertificateType == cdn.AfdCertificateTypeCustomerCertificate && secret == "" {
			return nil, fmt.Errorf("the 'cdn_frontdoor_secret_id' field must be set if the 'certificate_type' is 'CustomerCertificate'")
		} else if tls.CertificateType == cdn.AfdCertificateTypeManagedCertificate && secret != "" {
			return nil, fmt.Errorf("the 'cdn_frontdoor_secret_id' field is not supported if the 'certificate_type' is 'ManagedCertificate'")
		}

		if secret != "" {
			tls.Secret = expandResourceReference(secret)
		}
	}

	// NOTE: Minimum TLS Version is required in both pre-validated and not pre-validated
	// custom domains and the schema defaults the value to TLS 1.2
	tls.CertificateType = cdn.AfdCertificateType(certType)
	tls.MinimumTLSVersion = cdn.AfdMinimumTLSVersion(minTlsVersion)

	return &tls, nil
}

func flattenCustomDomainAFDDomainHttpsParameters(input *cdn.AFDDomainHTTPSParameters) []interface{} {
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

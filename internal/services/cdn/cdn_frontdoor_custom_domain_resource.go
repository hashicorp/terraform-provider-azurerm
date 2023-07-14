// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorCustomDomain() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainCreate,
		Read:   resourceCdnFrontDoorCustomDomainRead,
		Update: resourceCdnFrontDoorCustomDomainUpdate,
		Delete: resourceCdnFrontDoorCustomDomainDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			// NOTE: These timeouts are extremely long due to the manual
			// step of approving the private link if defined.
			Create: pluginsdk.DefaultTimeout(12 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(24 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(12 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FrontDoorCustomDomainID(id)
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

	return resource
}

func resourceCdnFrontDoorCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontDoorProfileID(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorCustomDomainID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
	}

	dnsZone := d.Get("dns_zone_id").(string)
	tls := d.Get("tls").([]interface{})

	props := cdn.AFDDomain{
		AFDDomainProperties: &cdn.AFDDomainProperties{
			HostName: utils.String(d.Get("host_name").(string)),
		},
	}

	if dnsZone != "" {
		props.AFDDomainProperties.AzureDNSZone = expandResourceReference(dnsZone)
	}

	tlsSettings, err := expandTlsParameters(tls)
	if err != nil {
		return err
	}

	props.AFDDomainProperties.TLSSettings = tlsSettings

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

	id, err := parse.FrontDoorCustomDomainID(d.Id())
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

		dnsZoneId, err := flattenDNSZoneResourceReference(props.AzureDNSZone)
		if err != nil {
			return fmt.Errorf("flattening `dns_zone_id`: %+v", err)
		}

		if err := d.Set("dns_zone_id", dnsZoneId); err != nil {
			return fmt.Errorf("setting `dns_zone_id`: %+v", err)
		}

		tls, err := flattenCustomDomainAFDDomainHttpsParameters(props.TLSSettings)
		if err != nil {
			return fmt.Errorf("flattening `tls`: %+v", err)
		}

		if err := d.Set("tls", tls); err != nil {
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

	id, err := parse.FrontDoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	props := cdn.AFDDomainUpdateParameters{
		AFDDomainUpdatePropertiesParameters: &cdn.AFDDomainUpdatePropertiesParameters{},
	}

	if d.HasChange("dns_zone_id") {
		if dnsZone := d.Get("dns_zone_id").(string); dnsZone != "" {
			props.AFDDomainUpdatePropertiesParameters.AzureDNSZone = expandResourceReference(dnsZone)
		}
	}

	if d.HasChange("tls") {
		tls := &cdn.AFDDomainHTTPSParameters{}
		tlsSettings := d.Get("tls").([]interface{})
		v := tlsSettings[0].(map[string]interface{})
		secretRaw := v["cdn_frontdoor_secret_id"].(string)

		// NOTE: Cert type has to always be passed in the update else you will get a
		// "AfdDomain.TlsSettings.CertificateType' is required but it was not set" error
		tls.CertificateType = cdn.AfdCertificateType(v["certificate_type"].(string))

		// NOTE: Secret always needs to be passed if it is defined else you will
		// receive a 500 Internal Server Error
		if secretRaw != "" {
			secret, err := parse.FrontDoorSecretID(secretRaw)
			if err != nil {
				return err
			}

			tls.Secret = expandResourceReference(secret.ID())
		}

		if d.HasChange("tls.0.minimum_tls_version") {
			tls.MinimumTLSVersion = cdn.AfdMinimumTLSVersion(v["minimum_tls_version"].(string))
		}

		if tls.CertificateType == cdn.AfdCertificateTypeCustomerCertificate && secretRaw == "" {
			return fmt.Errorf("the 'cdn_frontdoor_secret_id' field must be set if the 'certificate_type' is 'CustomerCertificate'")
		} else if tls.CertificateType == cdn.AfdCertificateTypeManagedCertificate && secretRaw != "" {
			return fmt.Errorf("the 'cdn_frontdoor_secret_id' field is not supported if the 'certificate_type' is 'ManagedCertificate'")
		}

		props.AFDDomainUpdatePropertiesParameters.TLSSettings = tls
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

	id, err := parse.FrontDoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	// delete the custom domain...
	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return nil
}

func expandTlsParameters(input []interface{}) (*cdn.AFDDomainHTTPSParameters, error) {
	// NOTE: With the Frontdoor service, they do not treat an empty object like an empty object
	// if it is not nil they assume it is fully defined and then end up throwing errors when they
	// attempt to get a value from one of the fields.
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	v := input[0].(map[string]interface{})

	certType := v["certificate_type"].(string)
	secretRaw := v["cdn_frontdoor_secret_id"].(string)
	minTlsVersion := v["minimum_tls_version"].(string)

	tls := cdn.AFDDomainHTTPSParameters{}

	if tls.CertificateType == cdn.AfdCertificateTypeCustomerCertificate && secretRaw == "" {
		return nil, fmt.Errorf("the 'cdn_frontdoor_secret_id' field must be set if the 'certificate_type' is 'CustomerCertificate'")
	} else if tls.CertificateType == cdn.AfdCertificateTypeManagedCertificate && secretRaw != "" {
		return nil, fmt.Errorf("the 'cdn_frontdoor_secret_id' field is not supported if the 'certificate_type' is 'ManagedCertificate'")
	}

	if secretRaw != "" {
		secret, err := parse.FrontDoorSecretID(secretRaw)
		if err != nil {
			return nil, err
		}

		tls.Secret = expandResourceReference(secret.ID())
	}

	// NOTE: Minimum TLS Version is required in both pre-validated and not pre-validated
	// custom domains and the schema defaults the value to TLS 1.2
	tls.CertificateType = cdn.AfdCertificateType(certType)
	tls.MinimumTLSVersion = cdn.AfdMinimumTLSVersion(minTlsVersion)

	return &tls, nil
}

func flattenCustomDomainAFDDomainHttpsParameters(input *cdn.AFDDomainHTTPSParameters) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	secretId, err := flattenSecretResourceReference(input.Secret)
	if err != nil {
		return []interface{}{}, err
	}

	return []interface{}{
		map[string]interface{}{
			"cdn_frontdoor_secret_id": secretId,
			"certificate_type":        string(input.CertificateType),
			"minimum_tls_version":     string(input.MinimumTLSVersion),
		},
	}, nil
}

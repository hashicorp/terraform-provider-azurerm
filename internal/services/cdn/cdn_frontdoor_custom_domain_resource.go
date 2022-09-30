package cdn

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
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

			// NOTE: I need the reference to ensure the correct destroy order
			// this will also offload the task of maintaining the custom domain association
			// IDs in the route resource to the custom domains themselves...
			"associate_with_cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.FrontDoorRouteID,
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

	if route := d.Get("associate_with_cdn_frontdoor_route_id").(string); route != "" {
		routeId, err := parse.FrontDoorRouteID(route)
		if err != nil {
			return err
		}
		if id.ProfileName != routeId.ProfileName {
			return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain: the configuration is invalid, the Front Door Custom Domain(Name: %q, Profile: %q) and the Front Door Route(Name: %q, Profile: %q) must belong to the same Front Door Profile", id.CustomDomainName, id.ProfileName, routeId.RouteName, routeId.ProfileName)
		}
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
	}

	// preValidatedDomain := d.Get("pre_validated_custom_domain_id").(string)
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

	tlsSettings, err := expandTlsParameters(tls, false)
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

	// Now that the Custom Domain has been created we need to associate the custom domain with the
	// route, if that field was passed...
	associateWithRoute := d.Get("associate_with_cdn_frontdoor_route_id").(string)
	if associateWithRoute != "" {
		// Lock the route for update...
		routeId, err := parse.FrontDoorRouteID(associateWithRoute)
		if err != nil {
			return err
		}

		locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
		defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

		// add the association to the route...
		writeFieldToState, err := addCustomDomainAssociationToRoute(d, meta, routeId, &id)
		if err != nil {
			return err
		}

		if writeFieldToState {
			d.Set("associate_with_cdn_frontdoor_route_id", associateWithRoute)
		}
	}

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

	id, err := parse.FrontDoorCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	if route := d.Get("associate_with_cdn_frontdoor_route_id").(string); route != "" {
		routeId, err := parse.FrontDoorRouteID(route)
		if err != nil {
			return err
		}
		if id.ProfileName != routeId.ProfileName {
			return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain: the configuration is invalid, the Front Door Custom Domain(Name: %q, Profile: %q) and the Front Door Route(Name: %q, Profile: %q) must belong to the same Front Door Profile", id.CustomDomainName, id.ProfileName, routeId.RouteName, routeId.ProfileName)
		}
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

		secret := v["cdn_frontdoor_secret_id"].(string)

		// NOTE: Cert type has to always be passed in the update else you will get a
		// "AfdDomain.TlsSettings.CertificateType' is required but it was not set" error
		tls.CertificateType = cdn.AfdCertificateType(v["certificate_type"].(string))

		// NOTE: Secret always needs to be passed if it is defined else you will
		// receive a 500 Internal Server Error
		if secret != "" {
			tls.Secret = expandResourceReference(secret)
		}

		if d.HasChange("tls.0.minimum_tls_version") {
			tls.MinimumTLSVersion = cdn.AfdMinimumTLSVersion(v["minimum_tls_version"].(string))
		}

		if tls.CertificateType == cdn.AfdCertificateTypeCustomerCertificate && secret == "" {
			return fmt.Errorf("the 'cdn_frontdoor_secret_id' field must be set if the 'certificate_type' is 'CustomerCertificate'")
		} else if tls.CertificateType == cdn.AfdCertificateTypeManagedCertificate && secret != "" {
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

	// Now that the Custom Domain has been updated we need to ensure the referential integrity of the route resource
	// and associate/unassociate the custom domain with the route, if that field was defined/removed...
	if d.HasChange("associate_with_cdn_frontdoor_route_id") {
		var writeFieldToState bool

		old, new := d.GetChange("associate_with_cdn_frontdoor_route_id")
		oldRouteValue := old.(string)
		newRouteValue := new.(string)

		// If the old value was "" and the new value is something we are adding an association (lock only new value route)
		// if the old value was something and it isn't the same as the new value we are associating the custom domain with a different route (lock both new and old routes)
		// if the old value was something and the new value is "" we are removing the association with the route (lock only the old value route)
		// the only other possibility here is that the old and new value are not empty and are the same value, which is a no op so do nothing...

		switch {
		case (oldRouteValue == "" && newRouteValue != ""):
			// Add
			newRouteId, err := parse.FrontDoorRouteID(newRouteValue)
			if err != nil {
				return err
			}

			// lock the route resource for update...
			locks.ByName(newRouteId.RouteName, cdnFrontDoorRouteResourceName)
			defer locks.UnlockByName(newRouteId.RouteName, cdnFrontDoorRouteResourceName)

			// add the association to the route...
			writeFieldToState, err = addCustomDomainAssociationToRoute(d, meta, newRouteId, id)
			if err != nil {
				return err
			}

		case (oldRouteValue != "" && newRouteValue == ""):
			// Remove
			oldRouteId, err := parse.FrontDoorRouteID(oldRouteValue)
			if err != nil {
				return err
			}

			// lock the route resource for update...
			locks.ByName(oldRouteId.RouteName, cdnFrontDoorRouteResourceName)
			defer locks.UnlockByName(oldRouteId.RouteName, cdnFrontDoorRouteResourceName)

			// remove the association from the route..
			err = removeCustomDomainAssociationFromRoute(d, meta, oldRouteId, id)
			if err != nil {
				return err
			}

		case (oldRouteValue != "" && newRouteValue != "" && !strings.EqualFold(oldRouteValue, newRouteValue)):
			// Swap
			oldRouteId, err := parse.FrontDoorRouteID(oldRouteValue)
			if err != nil {
				return err
			}

			newRouteId, err := parse.FrontDoorRouteID(newRouteValue)
			if err != nil {
				return err
			}

			// lock the route resources for update...
			locks.ByName(oldRouteId.RouteName, cdnFrontDoorRouteResourceName)
			defer locks.UnlockByName(oldRouteId.RouteName, cdnFrontDoorRouteResourceName)

			locks.ByName(newRouteId.RouteName, cdnFrontDoorRouteResourceName)
			defer locks.UnlockByName(newRouteId.RouteName, cdnFrontDoorRouteResourceName)

			// remove the association from the old route..
			err = removeCustomDomainAssociationFromRoute(d, meta, oldRouteId, id)
			if err != nil {
				return err
			}

			// add the association to the new route...
			writeFieldToState, err = addCustomDomainAssociationToRoute(d, meta, newRouteId, id)
			if err != nil {
				return err
			}
		}

		// Now that all of the routes have been processed and updated correctly
		// write the value to the state file if we need too...
		if writeFieldToState {
			d.Set("associate_with_cdn_frontdoor_route_id", newRouteValue)
		}
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

	// NOTE: If the custom domain is still associated with a route you cannot delete it
	// you must first update the route to remove the association with the custom domain...
	disassociateRoute := d.Get("associate_with_cdn_frontdoor_route_id").(string)
	if disassociateRoute != "" {
		routeId, err := parse.FrontDoorRouteID(disassociateRoute)
		if err != nil {
			return err
		}

		// NOTE: cdnFrontDoorRouteResourceName is defined in the "cdn_frontdoor_route_unlink_default_domain_resource.go" file
		locks.ByName(routeId.RouteName, cdnFrontDoorRouteResourceName)
		defer locks.UnlockByName(routeId.RouteName, cdnFrontDoorRouteResourceName)

		// remove the association from the route..
		err = removeCustomDomainAssociationFromRoute(d, meta, routeId, id)
		if err != nil {
			return err
		}
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

	if tls.CertificateType == cdn.AfdCertificateTypeCustomerCertificate && secret == "" {
		return nil, fmt.Errorf("the 'cdn_frontdoor_secret_id' field must be set if the 'certificate_type' is 'CustomerCertificate'")
	} else if tls.CertificateType == cdn.AfdCertificateTypeManagedCertificate && secret != "" {
		return nil, fmt.Errorf("the 'cdn_frontdoor_secret_id' field is not supported if the 'certificate_type' is 'ManagedCertificate'")
	}

	if secret != "" {
		tls.Secret = expandResourceReference(secret)
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

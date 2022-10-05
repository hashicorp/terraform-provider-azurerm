package cdn

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	dnsValidate "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/zones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
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

			// NOTE: Was afraid of this "Error: Cycle: azurerm_cdn_frontdoor_custom_domain.wodan, azurerm_cdn_frontdoor_route.wodan"
			// have to create a 'azurerm_cdn_frontdoor_custom_domain_routes_association' resource to manage the associations between
			// the routes and the custom domains for delete/destroy operation...
			"cdn_frontdoor_route_ids": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				ConflictsWith: func() []string {
					if !features.FourPointOhBeta() {
						return []string{"associate_with_cdn_frontdoor_route_id"}
					}
					return []string{}
				}(),

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorRouteID,
				},
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

	if !features.FourPointOhBeta() {
		resource.Schema["associate_with_cdn_frontdoor_route_id"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Optional:   true,
			Deprecated: "'associate_with_cdn_frontdoor_route_id' is no longer used and and will be removed in version 4.0 of the AzureRM provider. Please use the 'cdn_frontdoor_custom_domains_ids' field in the 'cdn_frontdoor_route' resource to control the CDN FrontDoor Custom Domain association(s) with the CDN FrontDoor Route",
			ConflictsWith: func() []string {
				if !features.FourPointOhBeta() {
					return []string{"cdn_frontdoor_route_ids"}
				}
				return []string{}
			}(),
			ValidateFunc: validate.FrontDoorRouteID,
		}
	}

	return resource
}

func resourceCdnFrontDoorCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := parse.FrontDoorProfileIDInsensitively(d.Get("cdn_frontdoor_profile_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorCustomDomainID(profileId.SubscriptionId, profileId.ResourceGroup, profileId.ProfileName, d.Get("name").(string))

	// validate the routes...
	routes := d.Get("cdn_frontdoor_route_ids").([]interface{})
	if len(routes) != 0 {
		if err := validateCustomDomainRoutes(routes, &id); err != nil {
			return err
		}
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
	d.Set("cdn_frontdoor_route_ids", routes)

	// Associate the custom domain with the route, if that field was passed...
	// if !features.FourPointOhBeta() {
	// 	if route := d.Get("associate_with_cdn_frontdoor_route_id").(string); route != "" {
	// 		routes := []interface{}{route}

	// 		if err = addCustomDomainAssociationToRoutes(d, meta, routes, &id); err != nil {
	// 			return err
	// 		} else {
	// 			d.Set("associate_with_cdn_frontdoor_route_id", route)
	// 		}
	// 	}
	// }

	// // NOTE: To associate a custom domain with more than one route, you must also make sure that all of the routes point to the same endpoint
	// // NOTE: I also now have to wait until all of the front door custom domains are created befor I move on here?
	// if routes := d.Get("associate_with_cdn_frontdoor_route_ids").([]interface{}); len(routes) != 0 {
	// 	if err = addCustomDomainAssociationToRoutes(d, meta, routes, &id); err != nil {
	// 		return err
	// 	} else {
	// 		d.Set("associate_with_cdn_frontdoor_route_ids", routes)
	// 	}
	// }

	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := customDomainIDWithErrorTxt(d.Id())
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

	d.Set("cdn_frontdoor_route_ids", d.Get("cdn_frontdoor_route_ids").([]interface{}))

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

	// validate the routes...
	if d.HasChange("cdn_frontdoor_route_ids") {
		routes := d.Get("cdn_frontdoor_route_ids").([]interface{})

		if len(routes) != 0 {
			if err := validateCustomDomainRoutes(routes, id); err != nil {
				return fmt.Errorf("azurerm_cdn_frontdoor_custom_domain: the 'cdn_frontdoor_route_ids' field is invalid: %+v", err)
			}
		}

		d.Set("cdn_frontdoor_route_ids", routes)
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
	// and associate/unassociate the custom domain with the route(s), if that field was defined/removed...

	// if !features.FourPointOhBeta() {
	// 	if route := d.Get("associate_with_cdn_frontdoor_route_id").(string); route != "" {
	// 		if d.HasChange("associate_with_cdn_frontdoor_route_id") {
	// 			old, new := d.GetChange("associate_with_cdn_frontdoor_route_id")
	// 			oldRoute := []interface{}{old.(string)}
	// 			newRoute := []interface{}{new.(string)}

	// 			switch {
	// 			case (len(oldRoute) == 0 && len(newRoute) != 0):
	// 				// Add
	// 				if err := addCustomDomainAssociationToRoutes(d, meta, newRoute, id); err != nil {
	// 					return err
	// 				}

	// 				d.Set("associate_with_cdn_frontdoor_route_id", new.(string))

	// 			case (len(oldRoute) != 0 && len(newRoute) == 0):
	// 				// Remove
	// 				if err := removeCustomDomainAssociationFromRoutes(d, meta, newRoute, id); err != nil {
	// 					return err
	// 				}

	// 				d.Set("associate_with_cdn_frontdoor_route_id", "")

	// 			case (len(oldRoute) != 0 && len(newRoute) != 0 && !strings.EqualFold(old.(string), new.(string))):
	// 				// Swap
	// 				if err = removeCustomDomainAssociationFromRoutes(d, meta, oldRoute, id); err != nil {
	// 					return err
	// 				}

	// 				d.Set("associate_with_cdn_frontdoor_route_id", "")

	// 				// add the association to the new route...
	// 				if err = addCustomDomainAssociationToRoutes(d, meta, newRoute, id); err != nil {
	// 					return err
	// 				}

	// 				d.Set("associate_with_cdn_frontdoor_route_id", new.(string))
	// 			}
	// 		}
	// 	}
	// }

	// if d.HasChange("cdn_frontdoor_route_ids") {
	// 	if routes := d.Get("cdn_frontdoor_route_ids").([]interface{}); len(routes) != 0 {
	// 		old, new := d.GetChange("associate_with_cdn_frontdoor_route_ids")
	// 		oldRoutes := old.([]interface{})
	// 		newRoutes := new.([]interface{})

	// 		switch {
	// 		case (len(oldRoutes) == 0 && len(newRoutes) != 0):
	// 			// Add
	// 			if err = addCustomDomainAssociationToRoutes(d, meta, newRoutes, id); err != nil {
	// 				return err
	// 			}

	// 			d.Set("associate_with_cdn_frontdoor_route_ids", newRoutes)

	// 		case (len(oldRoutes) != 0 && len(newRoutes) == 0):
	// 			// Remove
	// 			if err = removeCustomDomainAssociationFromRoutes(d, meta, oldRoutes, id); err != nil {
	// 				return err
	// 			}

	// 			d.Set("associate_with_cdn_frontdoor_route_ids", []interface{}{})

	// 		case (len(oldRoutes) != 0 && len(newRoutes) != 0):
	// 			// Swap
	// 			if removeRoutes, sharedRoutes := getRemoveRoutesDelta(oldRoutes, newRoutes); len(removeRoutes) != 0 {
	// 				if err = removeCustomDomainAssociationFromRoutes(d, meta, removeRoutes, id); err != nil {
	// 					return err
	// 				}

	// 				d.Set("associate_with_cdn_frontdoor_route_ids", sharedRoutes)
	// 			}

	// 			if err = addCustomDomainAssociationToRoutes(d, meta, newRoutes, id); err != nil {
	// 				return err
	// 			}

	// 			d.Set("associate_with_cdn_frontdoor_route_ids", newRoutes)
	// 		}
	// 	}
	// }

	return resourceCdnFrontDoorCustomDomainRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorCustomDomainIDInsensitively(d.Id())
	if err != nil {
		return err
	}

	// NOTE: If the custom domain is still associated with a route you cannot delete the custom domain
	// without the service throwing an error. You must first update the route to remove the association...
	if routes := d.Get("cdn_frontdoor_route_ids").([]interface{}); len(routes) != 0 {
		// remove the association from the old route..
		if err := removeCustomDomainAssociationFromRoutes(d, meta, routes, id); err != nil {
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

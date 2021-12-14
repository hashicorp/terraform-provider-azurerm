package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdCustomDomains() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdCustomDomainCreate,
		Read:   resourceAfdCustomDomainRead,
		Update: resourceAfdCustomDomainUpdate,
		Delete: resourceAfdCustomDomainDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.CdnEndpointV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AfdCustomDomainID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.CdnEndpointCustomDomainName(),
			},

			"profile_id": {
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

			"tls": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"certificate_type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
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
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAfdCustomDomainCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId := d.Get("profile_id").(string)
	customDomainName := d.Get("name").(string)
	tlsSettings := d.Get("tls").([]interface{})

	profile, err := parse.ProfileID(profileId)
	if err != nil {
		return err
	}

	id := parse.NewAfdCustomDomainID(profile.SubscriptionId, profile.ResourceGroup, profile.Name, customDomainName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %q: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_cdn_frontdoor_custom_domain", id.ID())
	}

	domain := cdn.AFDDomain{
		AFDDomainProperties: &cdn.AFDDomainProperties{
			HostName:    utils.String(d.Get("host_name").(string)),
			TLSSettings: expandTlsSettings(tlsSettings),
			//AzureDNSZone
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName, domain)
	if err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceAfdCustomDomainRead(d, meta)
}

func expandTlsSettings(input []interface{}) *cdn.AFDDomainHTTPSParameters {
	if len(input) == 0 {
		return nil
	}

	config := input[0].(map[string]interface{})
	certificateType := config["certificate_type"].(string)
	minimumTlsVersion := config["minimum_tls_version"].(string)
	secretId := config["secret_id"].(string)

	parameters := cdn.AFDDomainHTTPSParameters{}

	switch certificateType {
	case "CustomerCertificate":
		parameters.CertificateType = cdn.AfdCertificateTypeCustomerCertificate
	case "ManagedCertificate":
		parameters.CertificateType = cdn.AfdCertificateTypeManagedCertificate
	default:
		parameters.CertificateType = cdn.AfdCertificateTypeManagedCertificate
	}

	switch minimumTlsVersion {
	case "TLS10":
		parameters.MinimumTLSVersion = cdn.AfdMinimumTLSVersionTLS10
	case "TLS12":
		parameters.MinimumTLSVersion = cdn.AfdMinimumTLSVersionTLS12
	default:
		parameters.MinimumTLSVersion = cdn.AfdMinimumTLSVersionTLS12
	}

	parameters.Secret.ID = &secretId

	return &parameters
}

func resourceAfdCustomDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("host_name", resp.HostName)

	return nil
}

func resourceAfdCustomDomainUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.AfdCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAfdCustomDomainRead(d, meta)
}

func resourceAfdCustomDomainDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDCustomDomainsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdCustomDomainID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return nil
}

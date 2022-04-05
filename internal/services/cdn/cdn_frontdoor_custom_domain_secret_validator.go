package cdn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	dnsParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorCustomDomainSecretValidator() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorCustomDomainSecretValidatorCreate,
		Read:   resourceCdnFrontdoorCustomDomainSecretValidatorRead,
		Delete: resourceCdnFrontdoorCustomDomainSecretValidatorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(24 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			// TODO: Make an importer
			_, err := parse.FrontdoorCustomDomainSecretID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_custom_domain_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorCustomDomainID,
			},

			"cdn_frontdoor_custom_domain_txt_validator_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorCustomDomainTxtID,
			},
		},
	}
}

func resourceCdnFrontdoorCustomDomainSecretValidatorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomainClient := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	client := meta.(*clients.Client).Cdn.FrontdoorSecretsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	customDomainId, err := parse.FrontdoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	customDomainTxtValidatorId, err := parse.FrontdoorCustomDomainTxtID(d.Get("cdn_frontdoor_custom_domain_txt_validator_id").(string))
	if err != nil {
		return err
	}

	// I need to get the custom domain so I can look at the tls properties to get the secret id
	customDomainResp, err := customDomainClient.Get(ctx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(customDomainResp.Response) {
			return fmt.Errorf("checking for existing %s: %+v", *customDomainId, err)
		}

		return fmt.Errorf("retrieving up  %s: %+v", customDomainId, err)
	}

	// I now have the custom domain, now get the secret id from the tls settings
	tlsSecret := ""
	if customDomainResp.AFDDomainProperties != nil {
		tlsSecret = *customDomainResp.AFDDomainProperties.TLSSettings.Secret.ID
	}

	// This is ok, and shouldn't error as the secret has not propagated yet
	// Might need to add a hard link to the azurerm_cdn_frontdoor_custom_domain_route_association.test to wait
	// for the association to be made, or I can loop here until something comes back in the Secret ID...
	// "unable to prase Frontdoor Secret ID: ID was missing the `profiles` element", this might be something else on my part...
	secretId, err := parse.FrontdoorSecretID(tlsSecret)
	if err != nil {
		return fmt.Errorf("unable to prase Frontdoor Secret ID: %+v", err)
	}

	id := parse.NewFrontdoorCustomDomainSecretID(customDomainId.SubscriptionId, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName, secretId.SecretName)

	// Make sure the secret exists
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, secretId.SecretName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing %s: %+v", *secretId, err)
		}

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// "Failed", "InProgress", "NotStarted", "Succeeded"
	log.Printf("[DEBUG] Waiting for %q:%q Secret to become %q", "cdn_frontdoor_custom_domain_id", secretId.SecretName, "Succeeded")
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"InProgress", "NotStarted"},
		Target:                    []string{"Succeeded"},
		Refresh:                   cdnFrontdoorCustomDomainSecretRefreshFunc(ctx, client, secretId),
		MinTimeout:                30 * time.Second,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		ContinuousTargetOccurence: 3,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the %q:%q (Resource Group: %q) deployment state to become %q: %+v", "azurerm_cdn_frontdoor_custom_domain_secret_validator", id.SecretName, id.ResourceGroup, "Succeeded", err)
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_custom_domain_id", customDomainId.ID())
	d.Set("cdn_frontdoor_custom_domain_txt_validator_id", customDomainTxtValidatorId.ID())
	return resourceCdnFrontdoorCustomDomainSecretValidatorRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainSecretValidatorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	txtId, err := dnsParse.TxtRecordID(d.Get("dns_txt_record_id").(string))
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorCustomDomainTxtID(d.Id())
	if err != nil {
		return err
	}

	customDomainId := parse.NewFrontdoorCustomDomainID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.CustomDomainName)

	resp, err := client.Get(ctx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing %s: %+v", customDomainId, err)
		}

		return fmt.Errorf("retrieving %s: %+v", customDomainId, err)
	}

	if props := resp.AFDDomainProperties; props != nil {
		d.Set("dns_txt_record_id", txtId.ID())
		d.Set("cdn_frontdoor_custom_domain_id", customDomainId.ID())
		d.Set("cdn_frontdoor_custom_domain_validation_state", props.DomainValidationState)
	}

	return nil
}

func resourceCdnFrontdoorCustomDomainSecretValidatorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// TODO: Delete doesn't really make sense since this is a fake resource I need to think about this...

	d.SetId("")
	return nil
}

func cdnFrontdoorCustomDomainSecretRefreshFunc(ctx context.Context, client *track1.SecretsClient, id *parse.FrontdoorSecretId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if CDN Frontdoor Secret %q (Resource Group: %q) is available...", id.SecretName, id.ResourceGroup)

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecretName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving CDN Frontdoor Secret %q (Resource Group: %q) returned 404", id.SecretName, id.ResourceGroup)
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the Domain Validation State of the CDN Frontdoor Secret %q (Resource Group: %q): %+v", id.SecretName, id.ResourceGroup, err)
		}

		state := track1.DeploymentStatusFailed
		if props := resp.SecretProperties; props != nil {
			// "Failed", "InProgress", "NotStarted", "Succeeded"
			if props.DeploymentStatus != "" {
				state = props.DeploymentStatus
			}
		}

		if state == track1.DeploymentStatusFailed {
			log.Printf("[DEBUG] CDN Frontdoor Secret %q (Resource Group: %q) returned Deployment Status: %q", id.SecretName, id.ResourceGroup, state)
			return nil, string(state), nil
		}

		return resp, string(state), nil
	}
}

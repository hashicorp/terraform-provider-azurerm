package cdn

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
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
			"cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorRouteID,
			},

			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontdoorCustomDomainID,
				},
			},

			"cdn_frontdoor_custom_domain_txt_validator_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontdoorCustomDomainTxtID,
				},
			},

			"cdn_frontdoor_custom_domain_secrets_state": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"cdn_frontdoor_custom_domain_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"cdn_frontdoor_secret_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"cdn_frontdoor_secret_provisioning_state": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCdnFrontdoorCustomDomainSecretValidatorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	customDomainClient := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	client := meta.(*clients.Client).Cdn.FrontDoorSecretsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeId, err := parse.FrontdoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for the %q: %+v", "azurerm_cdn_frontdoor_custom_domain_secret_validator", err)
	}

	id := parse.NewFrontdoorCustomDomainSecretID(routeId.SubscriptionId, routeId.ResourceGroup, routeId.ProfileName, "secretValidator", uuid)

	customDomainIds := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})

	for _, customDomain := range customDomainIds {
		customDomainId, err := parse.FrontdoorCustomDomainID(customDomain.(string))
		if err != nil {
			return err
		}

		// I need to wait for the Custom Domains to be associated by the route before I can grab the TLS Settings
		// else the TLS Settings will be nil causing a panic...
		log.Printf("[DEBUG] Waiting for Custom Domain %q TLS Settings to become %q", customDomainId.CustomDomainName, "Succeeded")
		customDomainStateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"Pending"},
			Target:                    []string{"Succeeded"},
			Refresh:                   cdnFrontdoorCustomDomainTLSSettingsRefreshFunc(ctx, customDomainClient, customDomainId),
			MinTimeout:                30 * time.Second,
			Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
			ContinuousTargetOccurence: 1,
		}

		if _, err = customDomainStateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for the %q:%q (Resource Group: %q) TLS Settings to become %q: %+v", "azurerm_cdn_frontdoor_custom_domain_secret_validator", id.SecretName, id.ResourceGroup, "Succeeded", err)
		}

		// Now that I know they are there I can grab them...
		customDomainResp, err := customDomainClient.Get(ctx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
		if err != nil {
			if utils.ResponseWasNotFound(customDomainResp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", *customDomainId, err)
			}

			return fmt.Errorf("retrieving Frontdoor Custom Domain %s: %+v", customDomainId, err)
		}

		tlsSecret := ""
		if customDomainResp.AFDDomainProperties != nil {
			tlsSecret = *customDomainResp.AFDDomainProperties.TLSSettings.Secret.ID
		}

		secretId, err := parse.FrontdoorSecretIDInsensitively(tlsSecret)
		if err != nil {
			return fmt.Errorf("unable to prase Frontdoor Secret ID(%q): %+v", tlsSecret, err)
		}

		resp, err := client.Get(ctx, secretId.ResourceGroup, secretId.ProfileName, secretId.SecretName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", *secretId, err)
			}

			return fmt.Errorf("creating %s: %+v", id, err)
		}

		// NOTE: Per the service team: DeploymentStatus would be the correct check ultimately, however deployment tracking is not yet rolled out.
		// We are targeting to roll it out by end of next week(4/15/2022).
		log.Printf("[DEBUG] Waiting for Custom Domain %q secret %q to become %q", customDomainId.CustomDomainName, secretId.SecretName, "Succeeded")
		stateConf := &pluginsdk.StateChangeConf{
			Pending:                   []string{"InProgress", "NotStarted", "Updating", "Creating"},
			Target:                    []string{"Succeeded"},
			Refresh:                   cdnFrontdoorCustomDomainSecretRefreshFunc(ctx, client, secretId),
			MinTimeout:                30 * time.Second,
			Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
			ContinuousTargetOccurence: 1,
		}

		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for the %q:%q (Resource Group: %q) deployment state to become %q: %+v", "azurerm_cdn_frontdoor_custom_domain_secret_validator", id.SecretName, id.ResourceGroup, "Succeeded", err)
		}
	}

	validatorIds := d.Get("cdn_frontdoor_custom_domain_txt_validator_ids").([]interface{})

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_custom_domain_ids", customDomainIds)
	d.Set("cdn_frontdoor_route_id", routeId.ID())
	d.Set("cdn_frontdoor_custom_domain_txt_validator_ids", validatorIds)
	return resourceCdnFrontdoorCustomDomainSecretValidatorRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainSecretValidatorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	secretsClient := meta.(*clients.Client).Cdn.FrontDoorSecretsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeId, err := parse.FrontdoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	validatorIds := d.Get("cdn_frontdoor_custom_domain_txt_validator_ids").([]interface{})
	customDomainIds := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})
	customDomainsSecretsState := make([]interface{}, 0)

	for _, domain := range customDomainIds {
		secretState := make(map[string]interface{})

		id, err := parse.FrontdoorCustomDomainID(domain.(string))
		if err != nil {
			return err
		}

		secretState["cdn_frontdoor_custom_domain_id"] = id.ID()

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		tlsSecret := ""
		if props := resp.AFDDomainProperties; props != nil {
			tlsSecret = *props.TLSSettings.Secret.ID
		}

		if tlsSecret != "" {
			secretId, err := parse.FrontdoorSecretIDInsensitively(tlsSecret)
			if err != nil {
				return err
			}

			secretResp, err := secretsClient.Get(ctx, id.ResourceGroup, id.ProfileName, secretId.SecretName)
			if err != nil {
				if utils.ResponseWasNotFound(secretResp.Response) {
					return fmt.Errorf("checking for existing %s: %+v", secretId, err)
				}

				return fmt.Errorf("retrieving %s: %+v", secretId, err)
			}

			if props := secretResp.SecretProperties; props != nil {
				secretState["cdn_frontdoor_secret_id"] = secretId.ID()

				if props.ProvisioningState != "" {
					secretState["cdn_frontdoor_secret_provisioning_state"] = props.ProvisioningState
				}
			}
		}

		customDomainsSecretsState = append(customDomainsSecretsState, secretState)
	}

	d.Set("cdn_frontdoor_custom_domain_ids", customDomainIds)
	d.Set("cdn_frontdoor_route_id", routeId.ID())
	d.Set("cdn_frontdoor_custom_domain_secrets_state", customDomainsSecretsState)
	d.Set("cdn_frontdoor_custom_domain_txt_validator_ids", validatorIds)

	return nil
}

func resourceCdnFrontdoorCustomDomainSecretValidatorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func cdnFrontdoorCustomDomainTLSSettingsRefreshFunc(ctx context.Context, client *cdn.AFDCustomDomainsClient, id *parse.FrontdoorCustomDomainId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if CDN Frontdoor TLS Settings %q (Resource Group: %q) are available...", id.CustomDomainName, id.ResourceGroup)

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving CDN Frontdoor TLS Settings %q (Resource Group: %q) returned an error: %+v", id.CustomDomainName, id.ResourceGroup, err)
				return nil, "", fmt.Errorf("polling for the CDN Frontdoor TLS Settings %q (Resource Group: %q): %+v", id.CustomDomainName, id.ResourceGroup, err)
			}

			// The route and custom domain may not have been associated yet so a 404 is acceptable
			// keep polling until it shows up
			return resp, "Pending", nil
		}

		if props := resp.AFDDomainProperties; props != nil {
			validationState := props.DomainValidationState

			// First lets check the validation state of the custom domain
			if validationState == cdn.DomainValidationStateApproved {
				// Are the Domains TLS Settings available yet?
				if props.TLSSettings != nil && props.TLSSettings.Secret != nil && props.TLSSettings.Secret.ID != nil {
					return resp, "Succeeded", nil
				} else {
					// Nope.
					return resp, "Pending", nil
				}
			} else if validationState == cdn.DomainValidationStatePending || validationState == cdn.DomainValidationStateSubmitting {
				return resp, "Pending", nil
			}

			// It's not Approved, Pending or Submitting... we are in a bad state return an error...
			return nil, "", fmt.Errorf("the custom domain %q (resource group: %q) has returned the validation state of %q, which indicates that an error has occurred or that the custom domain is in an unsupported state", id.CustomDomainName, id.ResourceGroup, validationState)
		}

		// Default to pending because the Domain Properties may not be available yet...
		return resp, "Pending", nil
	}
}

func cdnFrontdoorCustomDomainSecretRefreshFunc(ctx context.Context, client *cdn.SecretsClient, id *parse.FrontdoorSecretId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if CDN Frontdoor Secret %q (Resource Group: %q) is available...", id.SecretName, id.ResourceGroup)

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.SecretName)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving CDN Frontdoor Secret %q (Resource Group: %q) returned an error: %+v", id.SecretName, id.ResourceGroup, err)
				return nil, "", fmt.Errorf("polling for the Domain Validation State of the CDN Frontdoor Secret %q (Resource Group: %q): %+v", id.SecretName, id.ResourceGroup, err)
			}

			// The secret may not have been provisioned yet so a 404 is acceptable
			// keep polling until it shows up
			return resp, "", nil
		}

		var out string
		provisioningState := cdn.AfdProvisioningStateFailed
		deploymentState := cdn.DeploymentStatusFailed
		if props := resp.SecretProperties; props != nil {
			if props.ProvisioningState != "" {
				provisioningState = props.ProvisioningState
				deploymentState = props.DeploymentStatus
			}
		}

		// Due to deployment tracking not being currently implemented in the service
		// I am first going to check the DeploymentStatus, if I get a NotStarted
		// I will fall back and use the provisioningState instead. That way when
		// they do implement the DeploymentStatus, this resource will be checking the
		// correct field once it goes live. But, for now ProvisioningState is all we
		// have to go with.

		if deploymentState == cdn.DeploymentStatusNotStarted {
			out = string(provisioningState)
		} else {
			out = string(deploymentState)
		}

		if strings.EqualFold(out, "failed") {
			log.Printf("[DEBUG] CDN Frontdoor Secret %q (Resource Group: %q) returned Deployment Status: %q", id.SecretName, id.ResourceGroup, out)
			return nil, "", fmt.Errorf("deployment state: %q", out)
		}

		return resp, out, nil
	}
}

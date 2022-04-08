package cdn

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
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
			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontdoorCustomDomainID,
				},
			},

			"cdn_frontdoor_custom_domain_route_association_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorCustomDomainRouteID,
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
	client := meta.(*clients.Client).Cdn.FrontdoorSecretsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	customDomainRouteId, err := parse.FrontdoorCustomDomainRouteID(d.Get("cdn_frontdoor_custom_domain_route_association_id").(string))
	if err != nil {
		return err
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for the %q: %+v", "azurerm_cdn_frontdoor_custom_domain_secret_validator", err)
	}

	id := parse.NewFrontdoorCustomDomainSecretID(customDomainRouteId.SubscriptionId, customDomainRouteId.ResourceGroup, customDomainRouteId.ProfileName, "secretValidator", uuid)

	customDomainIds := d.Get("cdn_frontdoor_custom_domain_ids").([]interface{})

	for _, customDomain := range customDomainIds {
		customDomainId, err := parse.FrontdoorCustomDomainID(customDomain.(string))
		if err != nil {
			return err
		}

		customDomainResp, err := customDomainClient.Get(ctx, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName)
		if err != nil {
			if utils.ResponseWasNotFound(customDomainResp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", *customDomainId, err)
			}

			return fmt.Errorf("retrieving up  %s: %+v", customDomainId, err)
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
			ContinuousTargetOccurence: 3,
		}

		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return fmt.Errorf("waiting for the %q:%q (Resource Group: %q) deployment state to become %q: %+v", "azurerm_cdn_frontdoor_custom_domain_secret_validator", id.SecretName, id.ResourceGroup, "Succeeded", err)
		}
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_custom_domain_ids", customDomainIds)
	d.Set("cdn_frontdoor_custom_domain_route_association_id", customDomainRouteId.ID())
	return resourceCdnFrontdoorCustomDomainSecretValidatorRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainSecretValidatorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	secretsClient := meta.(*clients.Client).Cdn.FrontdoorSecretsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	customDomainRouteId, err := parse.FrontdoorCustomDomainRouteID(d.Get("cdn_frontdoor_custom_domain_route_association_id").(string))
	if err != nil {
		return err
	}

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
				return fmt.Errorf("checking for existing %s: %+v", id, err)
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
	d.Set("cdn_frontdoor_custom_domain_route_association_id", customDomainRouteId.ID())
	d.Set("cdn_frontdoor_custom_domain_secrets_state", customDomainsSecretsState)

	return nil
}

func resourceCdnFrontdoorCustomDomainSecretValidatorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func cdnFrontdoorCustomDomainSecretRefreshFunc(ctx context.Context, client *track1.SecretsClient, id *parse.FrontdoorSecretId) pluginsdk.StateRefreshFunc {
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
		provisioningState := track1.AfdProvisioningStateFailed
		deploymentState := track1.DeploymentStatusFailed
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

		if deploymentState == track1.DeploymentStatusNotStarted {
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

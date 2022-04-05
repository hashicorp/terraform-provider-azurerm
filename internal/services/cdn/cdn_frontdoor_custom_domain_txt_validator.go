package cdn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	dnsParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/parse"
	dnsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontdoorCustomDomainTxtValidator() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorCustomDomainTxtValidatorCreate,
		Read:   resourceCdnFrontdoorCustomDomainTxtValidatorRead,
		Delete: resourceCdnFrontdoorCustomDomainTxtValidatorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(24 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			// TODO: Make an importer
			_, err := parse.FrontdoorCustomDomainTxtID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"dns_txt_record_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: dnsValidate.TxtRecordID,
			},

			"cdn_frontdoor_custom_domain_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorCustomDomainID,
			},

			"cdn_frontdoor_custom_domain_validation_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCdnFrontdoorCustomDomainTxtValidatorCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorCustomDomainsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	txtId, err := dnsParse.TxtRecordID(d.Get("dns_txt_record_id").(string))
	if err != nil {
		return err
	}

	customDomainId, err := parse.FrontdoorCustomDomainID(d.Get("cdn_frontdoor_custom_domain_id").(string))
	if err != nil {
		return err
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for the %q: %+v", "azurerm_cdn_frontdoor_custom_domain_txt_validator", err)
	}

	id := parse.NewFrontdoorCustomDomainTxtID(customDomainId.SubscriptionId, customDomainId.ResourceGroup, customDomainId.ProfileName, customDomainId.CustomDomainName, uuid)

	// Make sure the custom domain exists
	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("checking for existing %s: %+v", *customDomainId, err)
		}

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// DomainValidationStates: "Approved", "InternalError", "Pending", "PendingRevalidation",
	// "RefreshingValidationToken", Rejected", "Submitting", "TimedOut", "Unknown"
	log.Printf("[DEBUG] Waiting for %q:%q to become %q", "cdn_frontdoor_custom_domain_id", customDomainId.CustomDomainName, "Approved")
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Unknown", "Submitting", "Pending"},
		Target:                    []string{"Approved"},
		Refresh:                   cdnFrontdoorCustomDomainTxtRefreshFunc(ctx, client, customDomainId),
		MinTimeout:                30 * time.Second,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		ContinuousTargetOccurence: 3,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for the %q:%q (Resource Group: %q) domain validation state to become %q: %+v", "azurerm_cdn_frontdoor_custom_domain", id.CustomDomainName, id.ResourceGroup, "Approved", err)
	}

	d.SetId(id.ID())
	d.Set("dns_txt_record_id", txtId.ID())
	return resourceCdnFrontdoorCustomDomainTxtValidatorRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainTxtValidatorRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceCdnFrontdoorCustomDomainTxtValidatorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// TODO: Delete doesn't really make sense since this is a fake resource I need to think about this...

	d.SetId("")
	return nil
}

func cdnFrontdoorCustomDomainTxtRefreshFunc(ctx context.Context, client *track1.AFDCustomDomainsClient, id *parse.FrontdoorCustomDomainId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if CDN Frontdoor Custom Domain %q (Resource Group: %q) is available...", id.CustomDomainName, id.ResourceGroup)

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving CDN Frontdoor Custom Domain %q (Resource Group: %q) returned 404", id.CustomDomainName, id.ResourceGroup)
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the Domain Validation State of the CDN Frontdoor Custom Domain %q (Resource Group: %q): %+v", id.CustomDomainName, id.ResourceGroup, err)
		}

		state := track1.DomainValidationStateUnknown
		if props := resp.AFDDomainProperties; props != nil {
			// 'DomainValidationStateUnknown', 'DomainValidationStateSubmitting', 'DomainValidationStatePending',
			// 'DomainValidationStateRejected', 'DomainValidationStateTimedOut', 'DomainValidationStatePendingRevalidation',
			// 'DomainValidationStateApproved', 'DomainValidationStateRefreshingValidationToken',
			// 'DomainValidationStateInternalError'
			if props.DomainValidationState != "" {
				state = props.DomainValidationState
			}
		}

		if state == track1.DomainValidationStateRejected || state == track1.DomainValidationStateTimedOut || state == track1.DomainValidationStateInternalError {
			log.Printf("[DEBUG] CDN Frontdoor Custom Domain %q (Resource Group: %q) returned Domain Validation State: %q", id.CustomDomainName, id.ResourceGroup, state)
			return nil, string(state), nil
		}

		// not sure what to do here since they regenerated the DNS TXT record value... I think it is safe to ignore and keep polling...
		if state == track1.DomainValidationStateRefreshingValidationToken || state == track1.DomainValidationStatePendingRevalidation {
			log.Printf("[DEBUG] CDN Frontdoor Custom Domain %q (Resource Group: %q) returned Domain Validation State: %q", id.CustomDomainName, id.ResourceGroup, state)
			return resp, string(state), nil
		}

		return resp, string(state), nil
	}
}

package cdn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

var cdnFrontDoorCustomDomainResourceName = "azurerm_cdn_frontdoor_custom_domain"
var cdnFrontDoorRouteResourceName = "azurerm_cdn_frontdoor_route"

func resourceCdnFrontDoorCustomDomainAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorCustomDomainAssociationCreate,
		Read:   resourceCdnFrontDoorCustomDomainAssociationRead,
		Update: resourceCdnFrontDoorCustomDomainAssociationUpdate,
		Delete: resourceCdnFrontDoorCustomDomainAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.FrontDoorCustomDomainAssociationID(id)
			return err
		}, importCdnFrontDoorCustomDomainAssociation()),

		Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorRouteID,
			},

			"cdn_frontdoor_custom_domain_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorCustomDomainID,
				},
			},

			"link_to_default_domain": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceCdnFrontDoorCustomDomainAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for CDN FrontDoor Route <-> CDN FrontDoor Custom Domain Association creation")

	locks.ByID(cdnFrontDoorRouteResourceName)
	defer locks.UnlockByID(cdnFrontDoorRouteResourceName)

	rId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewFrontDoorCustomDomainAssociationID(rId.SubscriptionId, rId.ResourceGroup, rId.ProfileName, rId.AfdEndpointName, rId.RouteName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", id, err)
	}

	if existing.RouteProperties == nil {
		return fmt.Errorf("retrieving existing %s: 'properties' was nil", id)
	}

	customDomains := d.Get("cdn_frontdoor_custom_domain_ids").(*pluginsdk.Set).List()
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)

	if !linkToDefaultDomain && len(customDomains) == 0 {
		return fmt.Errorf("it is invalid to disable the 'LinkToDefaultDomain' for the associated CDN Front Door Route(Name: %s) since the route will not have any CDN Front Door Custom Domains associated with it", id.AssociationName)
	} else if len(customDomains) != 0 {
		if err := validateRoutesCustomDomainProfile(customDomains, id.ProfileName); err != nil {
			return err
		}
	}

	if props := existing.RouteProperties; props != nil {
		updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
			CacheConfiguration: props.CacheConfiguration,
		}

		updateProps.CustomDomains = expandCustomDomainActivatedResourceArray(customDomains)
		updateProps.LinkToDefaultDomain = expandEnabledBoolToLinkToDefaultDomain(linkToDefaultDomain)

		updateParams := azuresdkhacks.RouteUpdateParameters{
			RouteUpdatePropertiesParameters: pointer.To(updateProps),
		}

		future, err := workaroundsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName, updateParams)
		if err != nil {
			return fmt.Errorf("creating %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
	d.Set("cdn_frontdoor_route_id", rId.ID())
	d.Set("link_to_default_domain", linkToDefaultDomain)

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if resp.RouteProperties == nil {
		return fmt.Errorf("retrieving existing %s: 'properties' was nil", id)
	}

	if props := resp.RouteProperties; props != nil {
		customDomains, err := flattenCustomDomainActivatedResourceArray(props.CustomDomains)
		if err != nil {
			return err
		}

		// Need to normalize it here since the RP may have changed the casing...
		rId, err := parse.FrontDoorRouteIDInsensitively(pointer.From(resp.ID))
		if err != nil {
			return err
		}

		d.Set("cdn_frontdoor_route_id", rId.ID())
		d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
		d.Set("link_to_default_domain", flattenLinkToDefaultDomainToBool(props.LinkToDefaultDomain))
	}

	return nil
}

func resourceCdnFrontDoorCustomDomainAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.HasChange("cdn_frontdoor_custom_domain_ids") {
		id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
		if err != nil {
			return err
		}

		locks.ByID(cdnFrontDoorRouteResourceName)
		defer locks.UnlockByID(cdnFrontDoorRouteResourceName)

		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
		if err != nil {
			if utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("updating %s: was not found %+v", id, err)
			}

			return fmt.Errorf("updating %s: %+v", id, err)
		}

		if existing.RouteProperties == nil {
			return fmt.Errorf("retrieving existing %s: 'properties' was nil", id)
		}

		customDomains := d.Get("cdn_frontdoor_custom_domain_ids").(*pluginsdk.Set).List()
		linkToDefaultDomain := d.Get("link_to_default_domain").(bool)

		if !linkToDefaultDomain && len(customDomains) == 0 {
			return fmt.Errorf("it is invalid to disable the 'LinkToDefaultDomain' for the associated CDN Front Door Route(Name: %s) since the route will not have any CDN Front Door Custom Domains associated with it", id.AssociationName)
		} else if len(customDomains) != 0 {
			if err := validateRoutesCustomDomainProfile(customDomains, id.ProfileName); err != nil {
				return err
			}
		}

		if props := existing.RouteProperties; props != nil {
			updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
				CacheConfiguration:  props.CacheConfiguration,
				CustomDomains:       props.CustomDomains,
				LinkToDefaultDomain: props.LinkToDefaultDomain,
			}

			if d.HasChange("cdn_frontdoor_custom_domain_ids") {
				updateProps.CustomDomains = expandCustomDomainActivatedResourceArray(customDomains)
			}

			if d.HasChange("link_to_default_domain") {
				updateProps.LinkToDefaultDomain = expandEnabledBoolToLinkToDefaultDomain(linkToDefaultDomain)
			}

			updateParams := azuresdkhacks.RouteUpdateParameters{
				RouteUpdatePropertiesParameters: pointer.To(updateProps),
			}

			future, err := workaroundsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName, updateParams)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the update of %s: %+v", id, err)
			}
		}

		// Don't need to update Route ID here since it's force new...
		d.Set("cdn_frontdoor_custom_domain_ids", customDomains)
		d.Set("link_to_default_domain", linkToDefaultDomain)
	}

	return resourceCdnFrontDoorCustomDomainAssociationRead(d, meta)
}

func resourceCdnFrontDoorCustomDomainAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	workaroundsClient := azuresdkhacks.NewCdnFrontDoorRoutesWorkaroundClient(client)
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorCustomDomainAssociationID(d.Id())
	if err != nil {
		return err
	}

	// TODO: Need to poll the resource to see if another operation is in progress else you will get a conflict error returned from the service...
	// need to check to see if the route is currently being modified, if so wait...
	log.Printf("[DEBUG] Waiting for %s to become ready", id)

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context had no deadline")
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Creating", "Updating", "Deleting"},
		Target:                    []string{"Succeeded", "NotFound"},
		Refresh:                   frontDoorRouteRefreshFunc(ctx, client, id),
		MinTimeout:                15 * time.Second,
		ContinuousTargetOccurence: 2,
		Timeout:                   time.Until(deadline),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become ready: %+v", id, err)
	}

	locks.ByID(cdnFrontDoorRouteResourceName)
	defer locks.UnlockByID(cdnFrontDoorRouteResourceName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			// The route was already deleted...
			d.SetId("")
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if existing.RouteProperties == nil {
		return fmt.Errorf("retrieving existing %s: 'properties' was nil", id)
	}

	// You must set the LinkToDefaultDomain to 'true' here
	// else the route will be in an invalid state...
	if props := existing.RouteProperties; props != nil {
		updateProps := azuresdkhacks.RouteUpdatePropertiesParameters{
			CacheConfiguration:  props.CacheConfiguration,
			LinkToDefaultDomain: cdn.LinkToDefaultDomainEnabled,
		}

		updateProps.CustomDomains = nil

		updateParams := azuresdkhacks.RouteUpdateParameters{
			RouteUpdatePropertiesParameters: pointer.To(updateProps),
		}

		future, err := workaroundsClient.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName, updateParams)
		if err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the deletion of %s: %+v", id, err)
		}
	}

	d.SetId("")

	return nil
}

func frontDoorRouteRefreshFunc(ctx context.Context, client *cdn.RoutesClient, id *parse.FrontDoorCustomDomainAssociationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Checking to see if CDN Front Door Route %q (Resource Group: %q) is available...", id.AssociationName, id.ResourceGroup)

		resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				log.Printf("[DEBUG] Retrieving CDN Front Door Route %q (Resource Group: %q) returned 404.", id.AssociationName, id.ResourceGroup)
				return nil, "NotFound", nil
			}

			return nil, "", fmt.Errorf("polling for the state of the CDN Front Door Route %q (Resource Group: %q): %+v", id.AssociationName, id.ResourceGroup, err)
		}

		state := ""
		if props := resp.RouteProperties; props != nil {
			if props.ProvisioningState != "" {
				state = string(props.ProvisioningState)
			}
		}

		return resp, state, nil
	}
}

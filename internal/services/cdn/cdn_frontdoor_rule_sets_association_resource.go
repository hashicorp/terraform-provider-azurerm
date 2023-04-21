package cdn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const cdnFrontDoorRuleSetAssociationResourceType string = "azurerm_cdn_frontdoor_rule_sets_association"

func resourceCdnFrontDoorRuleSetAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorRuleSetAssociationCreate,
		Read:   resourceCdnFrontDoorRuleSetAssociationRead,
		Delete: resourceCdnFrontDoorRuleSetAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.FrontDoorRuleSetAssociationID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
			id, err := parse.FrontDoorRuleSetAssociationID(d.Id())
			if err != nil {
				return []*pluginsdk.ResourceData{}, err
			}

			client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
			resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
			if err != nil {
				return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if resp.RouteProperties == nil {
				return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: 'RouteProperties' are 'nil'", id)
			}
			ruleSets := flattenRuleSetResourceArray(resp.RouteProperties.RuleSets)

			d.SetId(id.ID())
			d.Set("cdn_frontdoor_route_id", parse.NewFrontDoorRouteID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName).ID())
			d.Set("cdn_frontdoor_rule_set_ids", ruleSets)

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorRouteID,
			},

			"cdn_frontdoor_rule_set_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.FrontDoorRuleSetID,
				},
			},
		},
	}
}

func resourceCdnFrontDoorRuleSetAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	// NOTE: The association name is the name of the route the resource is being associated with...
	// e.g. subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/associations/route1
	id := parse.NewFrontDoorRuleSetAssociationID(routeId.SubscriptionId, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)

	// lock the route, association and rule set resource types for create...
	locks.ByID(cdnFrontDoorRouteResourceType)
	defer locks.UnlockByID(cdnFrontDoorRouteResourceType)

	locks.ByID(cdnFrontDoorRuleSetAssociationResourceType)
	defer locks.UnlockByID(cdnFrontDoorRuleSetAssociationResourceType)

	locks.ByID(cdnFrontDoorRuleSetResourceType)
	defer locks.UnlockByID(cdnFrontDoorRuleSetResourceType)

	log.Printf("[INFO] preparing arguments for CDN FrontDoor Route <-> CDN FrontDoor Rule Set Association creation")

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", routeId, err)
	}

	props := existing.RouteProperties
	if props == nil {
		return fmt.Errorf("retrieving existing %s: 'properties' was nil", id)
	}

	ruleSetsRAW := d.Get("cdn_frontdoor_rule_set_ids").(*pluginsdk.Set).List()
	ruleSets, err := expandRuleSetIds(ruleSetsRAW)
	if err != nil {
		return err
	}

	props.RuleSets = expandRuleSetReferenceArray(ruleSets)

	routeProps := cdn.Route{
		RouteProperties: props,
	}

	// NOTE: Calling Create intentionally to avoid having to use the azuresdkhacks for the Update (PATCH) call...
	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName, routeProps)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creating of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceCdnFrontDoorRuleSetAssociationRead(d, meta)
}

func resourceCdnFrontDoorRuleSetAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRuleSetAssociationID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	props := existing.RouteProperties
	if props == nil {
		return fmt.Errorf("retrieving existing %s: 'properties' was nil", id)
	}

	d.Set("cdn_frontdoor_route_id", parse.NewFrontDoorRouteID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName).ID())

	if err := d.Set("cdn_frontdoor_rule_set_ids", flattenRuleSetResourceArray(props.RuleSets)); err != nil {
		return fmt.Errorf("setting `cdn_frontdoor_rule_set_ids`: %+v", err)
	}

	return nil
}

func resourceCdnFrontDoorRuleSetAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRuleSetAssociationID(d.Id())
	if err != nil {
		return err
	}

	// lock the route, association and rule set resource types for delete...
	locks.ByID(cdnFrontDoorRouteResourceType)
	defer locks.UnlockByID(cdnFrontDoorRouteResourceType)

	locks.ByID(cdnFrontDoorRuleSetAssociationResourceType)
	defer locks.UnlockByID(cdnFrontDoorRuleSetAssociationResourceType)

	locks.ByID(cdnFrontDoorRuleSetResourceType)
	defer locks.UnlockByID(cdnFrontDoorRuleSetResourceType)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	props := existing.RouteProperties
	if props == nil {
		return fmt.Errorf("retrieving existing %s: 'properties' was nil", id)
	}

	// remove all rule set associations from the route
	props.RuleSets = nil

	routeProps := cdn.Route{
		RouteProperties: props,
	}

	// NOTE: Calling Create intentionally to avoid having to use the azuresdkhacks for the Update (PATCH) call..
	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName, routeProps)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	d.SetId("")

	return nil
}

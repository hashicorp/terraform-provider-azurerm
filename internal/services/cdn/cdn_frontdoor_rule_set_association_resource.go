package cdn

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnFrontDoorRuleSetAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontDoorRuleSetAssociationCreate,
		Read:   resourceCdnFrontDoorRuleSetAssociationRead,
		Update: resourceCdnFrontDoorRuleSetAssociationUpdate,
		Delete: resourceCdnFrontDoorRuleSetAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := parse.FrontDoorRuleSetAssociationID(id)
			return err
		}, importCdnFrontDoorRuleSetAssociation()),

		Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontDoorRouteID,
			},

			"cdn_frontdoor_rule_set_ids": {
				Type:     pluginsdk.TypeList,
				Required: true,

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

	log.Printf("[INFO] preparing arguments for CDN FrontDoor Route <-> CDN FrontDoor Rule Set Association creation")

	// subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/associations/assoc1
	rId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	// NOTE: The association name is the name of the route the resource is being associated with...
	id := parse.NewFrontDoorRuleSetAssociationID(rId.SubscriptionId, rId.ResourceGroup, rId.ProfileName, rId.AfdEndpointName, rId.RouteName)

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("creating %s: %s was not found", id, rId)
		}

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// make sure the rule set exist and is associated with the route...
	ruleSets, err := validateRuleSetsTwo(d, meta, &id)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_route_id", rId.ID())
	d.Set("cdn_frontdoor_rule_set_ids", ruleSets)

	return resourceCdnFrontDoorRuleSetAssociationRead(d, meta)
}

func resourceCdnFrontDoorRuleSetAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FrontDoorRuleSetAssociationID(d.Id())
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

	return nil
}

func resourceCdnFrontDoorRuleSetAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontDoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if d.HasChange("cdn_frontdoor_rule_set_ids") {
		rId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
		if err != nil {
			return err
		}

		id := parse.NewFrontDoorRuleSetAssociationID(rId.SubscriptionId, rId.ResourceGroup, rId.ProfileName, rId.AfdEndpointName, rId.RouteName)

		existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
		if err != nil {
			if utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("updating %s: %s was not found", id, rId)
			}

			return fmt.Errorf("updating %s: %+v", id, err)
		}

		// make sure the route exist and the rule sets are associated with it...
		ruleSets, err := validateRuleSetsTwo(d, meta, &id)
		if err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}

		d.Set("cdn_frontdoor_rule_set_ids", ruleSets)
	}

	return resourceCdnFrontDoorRuleSetAssociationRead(d, meta)
}

func resourceCdnFrontDoorRuleSetAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// since you are deleting the resource you cannot grab the value from the config
	// because it will be empty, you have to get it from the states old value...
	oldRouteId, _ := d.GetChange("cdn_frontdoor_route_id")

	routeId, err := parse.FrontDoorRouteID(oldRouteId.(string))
	if err != nil {
		return err
	}

	id, err := parse.FrontDoorRuleSetAssociationID(d.Id())
	if err != nil {
		return err
	}

	oldRuleSetIds, _ := d.GetChange("cdn_frontdoor_rule_set_ids")
	oldRuleSet := oldRuleSetIds.([]interface{})

	x, v, err := expandRuleSets(oldRuleSet)
	if err != nil {
		return err
	}

	// I need to do a delta here, between what was in the resource and what is actually set on the route...
	newRouteRuleSetAssociations := routeRuleSetDelta(x, nil)
	if err = removeRuleSetAssociationsFromRouteTwo(d, meta, newRouteRuleSetAssociations, &newRouteRuleSetAssociations); err != nil {
		return out, err
	}

	// if there were rule sets associated with the route when you deleted
	// the association resource you need to remove those rule set associations
	// from the route...
	if err := removeRuleSetAssociationsFromRouteTwo(d, meta, &v, routeId); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	d.SetId("")

	return nil
}

func validateRuleSetsTwo(d *pluginsdk.ResourceData, meta interface{}, id *parse.FrontDoorRuleSetAssociationId) ([]interface{}, error) {
	out := make([]interface{}, 0)
	o, n := d.GetChange("cdn_frontdoor_rule_set_ids")
	oldRuleSets := o.([]interface{})
	newRuleSets := n.([]interface{})

	routeId := parse.NewFrontDoorRouteID(id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)

	if len(newRuleSets) == 0 || newRuleSets == nil || id == nil {
		return out, nil
	}

	oldID, _, err := expandRuleSets(oldRuleSets)
	if err != nil {
		return out, err
	}

	newID, result, err := expandRuleSets(newRuleSets)
	if err != nil {
		return out, err
	}

	// validate the new rule sets...
	if len(*newID) != 0 {
		notAssociated := make([]string, 0)

		for _, v := range *newID {
			// Make sure the route exists and get the routes rule set association list...
			routeRuleSettAssociations, _, err := getRouteRuleSetPropertiesTwo(d, meta, id, "cdn_frontdoor_rule_set_association")
			if err != nil {
				return out, err
			}

			// Make sure the rule set is in the routes rule set association list
			if len(routeRuleSettAssociations) == 0 || !sliceContainsString(routeRuleSettAssociations, v.ID()) {
				notAssociated = append(notAssociated, v.RuleSetName)
			}
		}

		if len(notAssociated) != 0 {
			return out, fmt.Errorf("the CDN FrontDoor Route(Name: %q) is currently not associated with the CDN FrontDoor Rule Sets: %s. Please remove the CDN FrontDoor Rule Sets from your 'cdn_frontdoor_rule_set_association' code block", id.AssociationName, strings.Join(notAssociated, ", "))
		}

		// check for dupe entries in the resources rule set list...
		if err := ruleSetSliceHasDuplicates(newID, "cdn_frontdoor_rule_set_association"); err != nil {
			return out, err
		}
	}

	// to get get the delta between the old and the new list, which should be the routes new rule set
	// associations, we need to compare the old list and new list and only record the rule set id if
	// they are in both the old and the new lists...
	newRouteRuleSetAssociations := routeRuleSetDelta(oldID, newID)
	if err = removeRuleSetAssociationsFromRouteTwo(d, meta, newRouteRuleSetAssociations, &routeId); err != nil {
		return out, err
	}

	return result, nil
}

// func validateRuleSets(d *pluginsdk.ResourceData, meta interface{}, id *parse.FrontDoorRuleSetAssociationId) ([]interface{}, error) {
// 	out := make([]interface{}, 0)

// 	routeId, err := parse.FrontDoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
// 	if err != nil {
// 		return out, err
// 	}

// 	resourceRuleSets := d.Get("cdn_frontdoor_rule_set_ids").([]interface{})

// 	ids, _, err := expandRuleSets(resourceRuleSets)
// 	if err != nil {
// 		return out, err
// 	}

// 	// Make sure the route exists and get the routes rule set association list...
// 	associations, _, err := getRouteRuleSetProperties(d, meta, routeId, "cdn_frontdoor_rule_set_association")
// 	if err != nil {
// 		return out, err
// 	}

// 	if len(associations) == 0 {
// 		return out, fmt.Errorf("the CDN FrontDoor Route(Name: %q) currently does not have any CDN FrontDoor Rule Sets associated with it. Please remove the 'cdn_frontdoor_rule_set_association' code block from your configuration", routeId.RouteName)
// 	}

// 	// validate the rule sets...
// 	if len(*ids) != 0 {
// 		notAssociated := make([]string, 0)

// 		for _, v := range *ids {
// 			// Make sure the rule set is in the routes rule set association list
// 			if !sliceContainsString(associations, v.ID()) {
// 				notAssociated = append(notAssociated, v.RuleSetName)
// 			} else {
// 				// the rule set is associated with the route
// 				out = append(out, v.ID())
// 			}
// 		}

// 		if len(notAssociated) > 0 {
// 			return out, fmt.Errorf("the CDN FrontDoor Route(Name: %q) is currently not associated with the following CDN FrontDoor Rule Sets: %s. Please remove the CDN FrontDoor Rule Sets from your 'cdn_frontdoor_rule_set_association' code block", routeId.RouteName, strings.Join(notAssociated, ", "))
// 		}

// 		// This is wrong, I need to remove the rule set association in the route if it does not exist in the association resource...
// 		// now make sure all of the routes rule sets are being referenced by the rule set association resource
// 		for _, v := range associations {
// 			routeRuleSet, err := parse.FrontDoorRuleSetID(v.(string))
// 			if err != nil {
// 				return out, err
// 			}

// 			// Make sure the routes rule set is in the resources rule set list
// 			if !sliceContainsString(resourceRuleSets, routeRuleSet.ID()) {
// 				notAssociated = append(notAssociated, routeRuleSet.RuleSetName)
// 			}
// 		}

// 		if len(notAssociated) > 0 {
// 			return out, fmt.Errorf("the CDN FrontDoor Route(Name: %q) is currently associated with the following CDN FrontDoor Rule Sets: %s. Please add these CDN FrontDoor Rule Sets to your 'cdn_frontdoor_rule_set_association' code block", routeId.RouteName, strings.Join(notAssociated, ", "))
// 		}
// 	} else {
// 		if len(associations) > 0 {
// 			associatedRuleSetName := make([]string, 0)

// 			for _, v := range associations {
// 				id, err := parse.FrontDoorRuleSetID(v.(string))
// 				if err != nil {
// 					return out, err
// 				}

// 				associatedRuleSetName = append(associatedRuleSetName, id.RuleSetName)
// 			}

// 			return out, fmt.Errorf("the CDN FrontDoor Route(Name: %q) is currently associated with the following CDN FrontDoor Rule Sets: %s. Please add these CDN FrontDoor Rule Sets to your 'cdn_frontdoor_rule_set_association' code block", routeId.RouteName, strings.Join(associatedRuleSetName, ", "))
// 		}
// 	}

// 	return out, nil
// }

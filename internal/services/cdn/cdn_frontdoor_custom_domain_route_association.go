package cdn

import (
	"fmt"
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

func resourceCdnFrontdoorCustomDomainRouteAssociation() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnFrontdoorCustomDomainRouteAssociationCreate,
		Read:   resourceCdnFrontdoorCustomDomainRouteAssociationRead,
		Update: resourceCdnFrontdoorCustomDomainRouteAssociationUpdate,
		Delete: resourceCdnFrontdoorCustomDomainRouteAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			// TODO: Make an importer
			_, err := parse.FrontdoorCustomDomainID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"cdn_frontdoor_route_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorRouteID,
			},

			"cdn_frontdoor_custom_domain_txt_validator_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FrontdoorCustomDomainTxtID,
			},

			"custom_domains": {
				Type:     pluginsdk.TypeList,
				Optional: true,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						// NOTE: I am using the Insensitively here because Portal lowercases everything
						"id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: ValidateFrontdoorCustomDomainIDInsensitively,
						},

						"active": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCdnFrontdoorCustomDomainRouteAssociationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeId, err := parse.FrontdoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return fmt.Errorf("generating UUID for the %q: %+v", "azurerm_cdn_frontdoor_custom_domain_association", err)
	}

	id := parse.NewFrontdoorCustomDomainRouteID(routeId.SubscriptionId, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName, uuid)

	// Make sure the Route exists
	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", routeId, err)
		}

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	props := track1.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &track1.RouteUpdatePropertiesParameters{
			CustomDomains: expandRouteActivatedResourceReferenceArray(d.Get("custom_domains").([]interface{})),
		},
	}

	// You must pass the Cache Configuration if it exist else you will remove it and disable compression if enabled
	if existing.RouteProperties.CacheConfiguration != nil {
		props.CacheConfiguration = existing.RouteProperties.CacheConfiguration
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, props)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("cdn_frontdoor_route_id", routeId.ID())
	d.Set("cdn_frontdoor_custom_domain_txt_validator_id", d.Get("cdn_frontdoor_custom_domain_txt_validator_id").(string))
	return resourceCdnFrontdoorCustomDomainRouteAssociationRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainRouteAssociationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeId, err := parse.FrontdoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	validatorId, err := parse.FrontdoorCustomDomainTxtID(d.Get("cdn_frontdoor_custom_domain_txt_validator_id").(string))
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorCustomDomainRouteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, routeId.ResourceGroup, routeId.ProfileName, routeId.AfdEndpointName, routeId.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("cdn_frontdoor_route_id", routeId.ID())
	d.Set("cdn_frontdoor_custom_domain_txt_validator_id", validatorId.ID())

	if props := resp.RouteProperties; props != nil {
		if err := d.Set("custom_domains", flattenRouteActivatedResourceReferenceArray(props.CustomDomains)); err != nil {
			return fmt.Errorf("setting `custom_domains`: %+v", err)
		}
	}

	return nil
}

func resourceCdnFrontdoorCustomDomainRouteAssociationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeId, err := parse.FrontdoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	id, err := parse.FrontdoorCustomDomainRouteID(d.Id())
	if err != nil {
		return err
	}

	// Make sure the Route exists
	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for existing %s: %+v", routeId, err)
		}

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	props := track1.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &track1.RouteUpdatePropertiesParameters{
			CustomDomains: expandRouteActivatedResourceReferenceArray(d.Get("custom_domains").([]interface{})),
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, props)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of %s: %+v", id, err)
	}

	return resourceCdnFrontdoorCustomDomainRouteAssociationRead(d, meta)
}

func resourceCdnFrontdoorCustomDomainRouteAssociationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorRoutesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: Need to Delete just the custom domain associations on the route
	id, err := parse.FrontdoorRouteID(d.Get("cdn_frontdoor_route_id").(string))
	if err != nil {
		return err
	}

	// NOTE: I had to change the SDK MarshalJSON resource to allow nil
	props := track1.RouteUpdateParameters{
		RouteUpdatePropertiesParameters: &track1.RouteUpdatePropertiesParameters{
			CustomDomains: nil,
		},
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, props)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	d.SetId("")
	return nil
}

func expandRouteResourceReferenceArray(input []interface{}) *[]track1.ResourceReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	results := make([]track1.ResourceReference, 0)

	for _, item := range input {
		results = append(results, track1.ResourceReference{
			ID: utils.String(item.(string)),
		})
	}

	return &results
}

func expandRouteActivatedResourceReferenceArray(input []interface{}) *[]track1.ActivatedResourceReference {
	results := make([]track1.ActivatedResourceReference, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, track1.ActivatedResourceReference{
			ID: utils.String(v["id"].(string)),
		})
	}
	return &results
}

func flattenRouteActivatedResourceReferenceArray(inputs *[]track1.ActivatedResourceReference) []interface{} {
	results := make([]interface{}, 0)
	if inputs == nil {
		return results
	}

	for _, input := range *inputs {
		result := make(map[string]interface{})

		if input.ID != nil {
			result["id"] = *input.ID
		}

		if input.IsActive != nil {
			result["active"] = *input.IsActive
		}
		results = append(results, result)
	}

	return results
}

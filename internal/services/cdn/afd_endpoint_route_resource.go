package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdEndpointRoutes() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdEndpointRouteCreate,
		Read:   resourceAfdEndpointRouteRead,
		Update: resourceAfdEndpointRouteUpdate,
		Delete: resourceAfdEndpointRouteDelete,

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
			_, err := parse.EndpointID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AfdEndpointsID,
			},

			"origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AfdOriginGroupsID,
			},

			"forwarding_protocol": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.ForwardingProtocolHTTPOnly),
					string(cdn.ForwardingProtocolHTTPSOnly),
					string(cdn.ForwardingProtocolMatchRequest),
				}, false),
			},

			"supported_protocols": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 2,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(cdn.AFDEndpointProtocolsHTTP),
						string(cdn.AFDEndpointProtocolsHTTPS),
					}, false),
				},
			},

			"link_to_default_domain": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"custom_domains": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceAfdEndpointRouteCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointRouteClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routeName := d.Get("name").(string)

	// parse endpoint_id
	endpointId := d.Get("endpoint_id").(string)
	endpoint, err := parse.AfdEndpointsID(endpointId)
	if err != nil {
		return err
	}

	// parse origin_group_id
	originGroupId := d.Get("origin_group_id").(string)
	originGroupRef := &cdn.ResourceReference{
		ID: &originGroupId,
	}

	var enabledState cdn.EnabledState = cdn.EnabledStateEnabled
	if !d.Get("enabled").(bool) {
		enabledState = cdn.EnabledStateDisabled
	} else {
		enabledState = cdn.EnabledStateEnabled
	}

	id := parse.NewAfdEndpointRouteID(endpoint.SubscriptionId, endpoint.ResourceGroup, endpoint.ProfileName, endpoint.AfdEndpointName, routeName)

	// link to default domain
	var linkToDefault cdn.LinkToDefaultDomain
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)
	if linkToDefaultDomain {
		linkToDefault = cdn.LinkToDefaultDomainEnabled
	} else {
		linkToDefault = cdn.LinkToDefaultDomainDisabled
	}

	// parse custom domains (TypeList)
	customDomains := d.Get("custom_domains").([]interface{})
	if len(customDomains) == 0 || customDomains[0] == nil {
		return nil
	}
	// create an Array of ResourceReferences per custom domain
	customDomainsArray := make([]cdn.ResourceReference, 0)
	for _, v := range customDomains {

		resourceId := v.(string)
		resourceReference := cdn.ResourceReference{
			ID: &resourceId,
		}
		customDomainsArray = append(customDomainsArray, resourceReference)
	}

	// forwarding protocol
	forwardingProtocol := d.Get("forwarding_protocol").(string)

	// supported protocols
	supportedProtocols := d.Get("supported_protocols").([]interface{})
	if len(supportedProtocols) == 0 || supportedProtocols[0] == nil {
		return nil
	}
	// create an Array of ResourceReferences for supported protocols
	supportedProtocolsArray := make([]cdn.AFDEndpointProtocols, 0)
	for _, v := range supportedProtocols {

		protocol := v.(string)
		var supportedProtocol cdn.AFDEndpointProtocols
		switch protocol {
		case "Http":
			supportedProtocol = cdn.AFDEndpointProtocolsHTTP
		case "Https":
			supportedProtocol = cdn.AFDEndpointProtocolsHTTPS
		}

		supportedProtocolsArray = append(supportedProtocolsArray, supportedProtocol)
	}

	route := cdn.Route{
		Name: &routeName,
		RouteProperties: &cdn.RouteProperties{
			CustomDomains:       &customDomainsArray,
			OriginGroup:         originGroupRef,
			EnabledState:        enabledState,
			SupportedProtocols:  &supportedProtocolsArray,
			ForwardingProtocol:  cdn.ForwardingProtocol(forwardingProtocol),
			LinkToDefaultDomain: linkToDefault,
		},
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, routeName, route)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAfdEndpointRouteRead(d, meta)
}

func resourceAfdEndpointRouteRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointRouteClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdEndpointRouteID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %q: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("enabled", resp.EnabledState)
	d.Set("custom_domains", resp.CustomDomains)
	d.Set("id", resp.ID)

	return nil
}
func resourceAfdEndpointRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	id, err := parse.AfdEndpointRouteID(d.Id())
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceAfdEndpointRouteRead(d, meta)
}
func resourceAfdEndpointRouteDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointRouteClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdEndpointRouteID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}

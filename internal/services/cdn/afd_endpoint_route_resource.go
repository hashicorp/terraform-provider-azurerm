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

			// CustomDomains - Domains referenced by this endpoint.
			"custom_domains": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			// OriginGroup - A reference to the origin group.
			"origin_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AfdOriginGroupsID,
			},

			// OriginPath - A directory path on the origin that AzureFrontDoor can use to retrieve content from, e.g. contoso.cloudapp.net/originpath.
			"origin_path": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// RuleSets - rule sets referenced by this endpoint.
			"rule_sets": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			// SupportedProtocols - List of supported protocols for this route.
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

			// PatternsToMatch - The route patterns of the rule.
			"patterns_to_match": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			// CompressionSettings - compression settings
			// IsCompressionEnabled - Indicates whether content compression is enabled on AzureFrontDoor. Default value is false. If compression is enabled, content will be served as compressed if user requests for a compressed version. Content won't be compressed on AzureFrontDoor when requested content is smaller than 1 byte or larger than 1 MB.
			"enable_caching": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			// ContentTypesToCompress - List of content types on which compression applies. The value should be a valid MIME type.
			"content_types_to_compress": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},

			// QueryStringCachingBehavior - Defines how CDN caches requests that include query strings. You can ignore any query strings when caching, bypass caching to prevent requests that contain query strings from being cached, or cache every request with a unique URL. Possible values include: 'AfdQueryStringCachingBehaviorIgnoreQueryString', 'AfdQueryStringCachingBehaviorUseQueryString', 'AfdQueryStringCachingBehaviorNotSet'
			"query_string_caching_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  cdn.AfdQueryStringCachingBehaviorNotSet,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.AfdQueryStringCachingBehaviorIgnoreQueryString),
					string(cdn.AfdQueryStringCachingBehaviorNotSet),
					string(cdn.AfdQueryStringCachingBehaviorUseQueryString),
				}, false),
			},

			// ForwardingProtocol - Protocol this rule will use when forwarding traffic to backends. Possible values include: 'ForwardingProtocolHTTPOnly', 'ForwardingProtocolHTTPSOnly', 'ForwardingProtocolMatchRequest'
			"forwarding_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  cdn.ForwardingProtocolMatchRequest,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.ForwardingProtocolHTTPOnly),
					string(cdn.ForwardingProtocolHTTPSOnly),
					string(cdn.ForwardingProtocolMatchRequest),
				}, false),
			},

			// LinkToDefaultDomain - whether this route will be linked to the default endpoint domain.
			"link_to_default_domain": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			// HTTPSRedirect - Whether to automatically redirect HTTP traffic to HTTPS traffic. Note that this is a easy way to set up this rule and it will be the first rule that gets executed.
			"https_redirect": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  cdn.HTTPSRedirectDisabled,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.HTTPSRedirectDisabled),
					string(cdn.HTTPSRedirectEnabled),
				}, false),
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

	// caching
	cachingEnabled := d.Get("enable_caching").(bool)
	contentTypesToCompress := d.Get("content_types_to_compress").([]interface{})
	// create an array of content types
	var compressionSettings cdn.CompressionSettings
	contentTypesToCompressArray := make([]string, 0)
	for _, contentType := range contentTypesToCompress {
		pattern := contentType.(string)
		contentTypesToCompressArray = append(contentTypesToCompressArray, pattern)
	}

	if cachingEnabled {
		compressionSettings.IsCompressionEnabled = &cachingEnabled
		compressionSettings.ContentTypesToCompress = &contentTypesToCompressArray
	}

	// endpoint route enabled state
	var enabledState cdn.EnabledState = cdn.EnabledStateEnabled
	if !d.Get("enabled").(bool) {
		enabledState = cdn.EnabledStateDisabled
	} else {
		enabledState = cdn.EnabledStateEnabled
	}

	id := parse.NewAfdEndpointRouteID(endpoint.SubscriptionId, endpoint.ResourceGroup, endpoint.ProfileName, endpoint.AfdEndpointName, routeName)

	// parse custom_domains (TypeList)
	customDomains := d.Get("custom_domains").([]interface{})
	// create an Array of ResourceReferences per custom domain
	customDomainsArray := make([]cdn.ResourceReference, 0)
	for _, v := range customDomains {
		resourceId := v.(string)
		resourceReference := cdn.ResourceReference{
			ID: &resourceId,
		}
		customDomainsArray = append(customDomainsArray, resourceReference)
	}

	// link to default domain
	var linkToDefault cdn.LinkToDefaultDomain
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)
	if linkToDefaultDomain {
		linkToDefault = cdn.LinkToDefaultDomainEnabled
	} else {
		linkToDefault = cdn.LinkToDefaultDomainDisabled
	}

	// parse rule_sets (TypeList)
	ruleSets := d.Get("rule_sets").([]interface{})
	// create an Array of ResourceReferences per custom domain
	ruleSetsArray := make([]cdn.ResourceReference, 0)
	for _, r := range ruleSets {
		ruleSetId := r.(string)
		resourceReference := cdn.ResourceReference{
			ID: &ruleSetId,
		}
		ruleSetsArray = append(ruleSetsArray, resourceReference)
	}

	// forwarding protocol
	forwardingProtocol := d.Get("forwarding_protocol").(string)

	// HTTPSRedirect
	var httpsRedirectSet cdn.HTTPSRedirect
	httpsRedirect := d.Get("https_redirect").(string)
	if httpsRedirect == "Enabled" {
		httpsRedirectSet = cdn.HTTPSRedirectEnabled
	} else {
		httpsRedirectSet = cdn.HTTPSRedirectDisabled
	}

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

	// patterns_to_match
	patternsToMatch := d.Get("patterns_to_match").([]interface{})
	// create an Array of ResourceReferences per custom domain
	patternsToMatchArray := make([]string, 0)
	for _, p := range patternsToMatch {
		pattern := p.(string)
		patternsToMatchArray = append(patternsToMatchArray, pattern)
	}

	route := cdn.Route{
		RouteProperties: &cdn.RouteProperties{
			OriginGroup:         originGroupRef,
			EnabledState:        enabledState,
			SupportedProtocols:  &supportedProtocolsArray,
			ForwardingProtocol:  cdn.ForwardingProtocol(forwardingProtocol),
			LinkToDefaultDomain: linkToDefault,
			RuleSets:            &ruleSetsArray,
			PatternsToMatch:     &patternsToMatchArray,
			HTTPSRedirect:       httpsRedirectSet,
		},
	}

	// originPath
	originPath := d.Get("origin_path").(string)
	if originPath != "" {
		route.OriginPath = &originPath
	}

	// query_string_caching_behavior
	queryStringCachingBehavior := d.Get("query_string_caching_behavior").(string)
	switch queryStringCachingBehavior {
	case "IgnoreQueryString":
		route.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorIgnoreQueryString
	case "NotSet":
		route.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorNotSet
	case "UseQueryString":
		route.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorUseQueryString
	default:
		route.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorNotSet
	}

	// use route.CompressionSettings only when caching is enabled
	if cachingEnabled == true {
		route.CompressionSettings = compressionSettings
	} else {
		route.CompressionSettings = nil
	}

	// custom domains can only be set when link_to_default_domain is false
	if linkToDefaultDomain == true {
		route.CustomDomains = nil
	} else {
		route.CustomDomains = &customDomainsArray
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

	d.Set("query_string_caching_behavior", resp.QueryStringCachingBehavior)
	d.Set("supported_protocols", resp.SupportedProtocols)
	d.Set("patterns_to_match", resp.PatternsToMatch)
	d.Set("forwarding_protocol", resp.ForwardingProtocol)
	d.Set("custom_domains", resp.CustomDomains)

	if resp.EnabledState == cdn.EnabledStateEnabled {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	d.Set("https_redirect", resp.HTTPSRedirect)

	return nil
}
func resourceAfdEndpointRouteUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDEndpointRouteClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdEndpointRouteID(d.Id())
	if err != nil {
		return err
	}

	cachingEnabled := d.Get("enable_caching").(bool)
	contentTypesToCompress := d.Get("content_types_to_compress").([]interface{})
	queryStringCachingBehavior := d.Get("query_string_caching_behavior").(string)
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)
	forwardingProtocol := d.Get("forwarding_protocol").(string)
	patternsToMatch := d.Get("patterns_to_match").([]interface{})
	customDomains := d.Get("custom_domains").([]interface{})

	// create an array of content types
	contentTypesToCompressArray := make([]string, 0)
	for _, contentType := range contentTypesToCompress {
		pattern := contentType.(string)
		contentTypesToCompressArray = append(contentTypesToCompressArray, pattern)
	}

	var routeUpdate cdn.RouteUpdateParameters
	var routeUpdateProperties cdn.RouteUpdatePropertiesParameters

	// patterns_to_match
	if d.HasChange("patterns_to_match") {
		log.Printf("[DEBUG] Updating patterns_to_match for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)

		patternsToMatchArray := make([]string, 0)
		for _, p := range patternsToMatch {
			pattern := p.(string)
			patternsToMatchArray = append(patternsToMatchArray, pattern)
		}

		routeUpdateProperties.PatternsToMatch = &patternsToMatchArray
	}

	// forwarding_protocol
	if d.HasChange("forwarding_protocol") {
		log.Printf("[DEBUG] Updating forwarding_protocol for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)
		routeUpdateProperties.ForwardingProtocol = cdn.ForwardingProtocol(forwardingProtocol)
	}

	// linkToDefaultDomain
	if d.HasChange("link_to_default_domain") || d.HasChange("custom_domains") {
		log.Printf("[DEBUG] Updating link_to_default_domain configuration for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)

		customDomainsArray := make([]cdn.ResourceReference, 0)
		for _, v := range customDomains {
			resourceId := v.(string)
			resourceReference := cdn.ResourceReference{
				ID: &resourceId,
			}
			customDomainsArray = append(customDomainsArray, resourceReference)
		}

		if linkToDefaultDomain == true {
			routeUpdateProperties.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled
			routeUpdateProperties.CustomDomains = nil
		} else {
			routeUpdateProperties.LinkToDefaultDomain = cdn.LinkToDefaultDomainDisabled
			routeUpdateProperties.CustomDomains = &customDomainsArray
		}
	}

	// update caching & compression settings
	if cachingEnabled == true && (d.HasChange("query_string_caching_behavior") || d.HasChange("content_types_to_compress")) {

		log.Printf("[DEBUG] Updating query_string_caching_behavior configuration for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)

		switch queryStringCachingBehavior {
		case "IgnoreQueryString":
			routeUpdateProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorIgnoreQueryString
		case "NotSet":
			routeUpdateProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorNotSet
		case "UseQueryString":
			routeUpdateProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorUseQueryString
		default:
			routeUpdateProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorNotSet
		}

		routeUpdateProperties.CompressionSettings = cdn.CompressionSettings{
			IsCompressionEnabled:   &cachingEnabled,
			ContentTypesToCompress: &contentTypesToCompressArray,
		}

	}

	if cachingEnabled == false {
		routeUpdateProperties.CompressionSettings = nil
		//routeUpdateProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorNotSet
	}

	if d.HasChange("custom_domains") || d.HasChange("link_to_default_domain") {
		log.Printf("[DEBUG] Updating custom domains for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)

		// link to default domain
		var linkToDefault cdn.LinkToDefaultDomain

		if linkToDefaultDomain {
			linkToDefault = cdn.LinkToDefaultDomainEnabled
		} else {
			linkToDefault = cdn.LinkToDefaultDomainDisabled
		}

		// parse custom_domains (TypeList)
		customDomains := d.Get("custom_domains").([]interface{})
		// create an Array of ResourceReferences per custom domain
		customDomainsArray := make([]cdn.ResourceReference, 0)
		for _, v := range customDomains {
			resourceId := v.(string)
			resourceReference := cdn.ResourceReference{
				ID: &resourceId,
			}
			customDomainsArray = append(customDomainsArray, resourceReference)
		}
		if linkToDefaultDomain {
			routeUpdateProperties.LinkToDefaultDomain = linkToDefault
		} else {
			routeUpdateProperties.LinkToDefaultDomain = linkToDefault
			routeUpdateProperties.CustomDomains = &customDomainsArray
		}
	}

	routeUpdate.RouteUpdatePropertiesParameters = &routeUpdateProperties

	future, err := client.Update(ctx, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, routeUpdate)

	if err != nil {
		return fmt.Errorf("updating Front Door Endpoint Route %q (Resource Group %q): %+v", id.RouteName, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of Front Door Route %q (Resource Group %q): %+v", id.RouteName, id.ResourceGroup, err)
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

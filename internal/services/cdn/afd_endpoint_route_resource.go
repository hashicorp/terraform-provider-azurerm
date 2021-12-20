package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
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

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AfdEndpointRouteID(id)
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
				Default:  true,
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
					ValidateFunc: validate.AfdCustomDomainID,
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
				Optional:     true,
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
				Optional: true,
				MinItems: 1,
				MaxItems: 2,
				Elem: &pluginsdk.Schema{
					Type:    pluginsdk.TypeString,
					Default: string(cdn.AFDEndpointProtocolsHTTP),
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
				Default:  true,
			},

			// HTTPSRedirect - Whether to automatically redirect HTTP traffic to HTTPS traffic. Note that this is a easy way to set up this rule and it will be the first rule that gets executed.
			"https_redirect": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
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

	id := parse.NewAfdEndpointRouteID(endpoint.SubscriptionId, endpoint.ResourceGroup, endpoint.ProfileName, endpoint.AfdEndpointName, routeName)

	route := cdn.Route{}
	routeProperties := cdn.RouteProperties{}

	// caching
	cachingEnabled := d.Get("enable_caching").(bool)
	contentTypesToCompress := d.Get("content_types_to_compress").([]interface{})
	// create an array of content types
	contentTypesToCompressArray := make([]string, 0)
	for _, contentType := range contentTypesToCompress {
		pattern := contentType.(string)
		contentTypesToCompressArray = append(contentTypesToCompressArray, pattern)
	}

	compressionSettings := cdn.CompressionSettings{
		IsCompressionEnabled:   utils.Bool(cachingEnabled),
		ContentTypesToCompress: &contentTypesToCompressArray,
	}

	routeProperties.CompressionSettings = compressionSettings

	// endpoint route enabled state
	enabledState := d.Get("enabled").(bool)
	if enabledState {
		routeProperties.EnabledState = cdn.EnabledStateEnabled
	} else {
		routeProperties.EnabledState = cdn.EnabledStateDisabled
	}

	// parse origin_group_id
	originGroupId := d.Get("origin_group_id").(string)
	originGroupRef := &cdn.ResourceReference{
		ID: &originGroupId,
	}
	routeProperties.OriginGroup = originGroupRef

	// link to default domain
	linkToDefaultDomain := d.Get("link_to_default_domain").(bool)
	if linkToDefaultDomain {
		routeProperties.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled
	} else {
		routeProperties.LinkToDefaultDomain = cdn.LinkToDefaultDomainDisabled
	}

	// parse custom_domains (TypeList)
	customDomains := d.Get("custom_domains").([]interface{})
	if customDomains != nil {
		customDomainsArray := make([]cdn.ResourceReference, 0)
		for _, v := range customDomains {
			resourceId := v.(string)
			resourceReference := cdn.ResourceReference{
				ID: &resourceId,
			}
			customDomainsArray = append(customDomainsArray, resourceReference)
		}
		routeProperties.CustomDomains = &customDomainsArray
	}

	// parse rule_sets (TypeList)
	ruleSets := d.Get("rule_sets").([]interface{})
	if ruleSets != nil {
		ruleSetsArray := make([]cdn.ResourceReference, 0)
		for _, r := range ruleSets {
			ruleSetId := r.(string)
			resourceReference := cdn.ResourceReference{
				ID: &ruleSetId,
			}
			ruleSetsArray = append(ruleSetsArray, resourceReference)
		}
		routeProperties.RuleSets = &ruleSetsArray
	}

	// forwarding protocol
	forwardingProtocol := d.Get("forwarding_protocol").(string)
	if forwardingProtocol != "" {
		routeProperties.ForwardingProtocol = cdn.ForwardingProtocol(forwardingProtocol)
	}

	// supported protocols
	supportedProtocols := d.Get("supported_protocols").([]interface{})
	if supportedProtocols != nil {
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
		routeProperties.SupportedProtocols = &supportedProtocolsArray
	}

	// patterns_to_match
	patternsToMatch := d.Get("patterns_to_match").([]interface{})
	patternsToMatchArray := make([]string, 0)
	for _, p := range patternsToMatch {
		pattern := p.(string)
		patternsToMatchArray = append(patternsToMatchArray, pattern)
	}
	routeProperties.PatternsToMatch = &patternsToMatchArray

	// HTTPSRedirect
	httpsRedirect := d.Get("https_redirect").(bool)
	if httpsRedirect {
		routeProperties.HTTPSRedirect = cdn.HTTPSRedirectEnabled
	} else {
		routeProperties.HTTPSRedirect = cdn.HTTPSRedirectDisabled
	}

	// enabledState
	if enabledState {
		routeProperties.EnabledState = cdn.EnabledStateEnabled
	} else {
		routeProperties.EnabledState = cdn.EnabledStateDisabled
	}

	// originPath
	originPath := d.Get("origin_path").(string)
	if originPath != "" {
		routeProperties.OriginPath = &originPath
	}

	// query_string_caching_behavior
	queryStringCachingBehavior := d.Get("query_string_caching_behavior").(string)
	switch queryStringCachingBehavior {
	case "IgnoreQueryString":
		routeProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorIgnoreQueryString
	case "NotSet":
		routeProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorNotSet
	case "UseQueryString":
		routeProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorUseQueryString
	}

	// use route.CompressionSettings only when caching is enabled
	if cachingEnabled {
		routeProperties.CompressionSettings = compressionSettings
	} else {
		routeProperties.CompressionSettings = nil
	}

	route.RouteProperties = &routeProperties

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

	if resp.HTTPSRedirect == cdn.HTTPSRedirectEnabled {
		d.Set("https_redirect", true)
	} else {
		d.Set("https_redirect", false)
	}

	if resp.EnabledState == cdn.EnabledStateEnabled {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

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

	if d.HasChange("enabled") {
		log.Printf("[DEBUG] Updating enabled for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)

		enabledState := d.Get("enabled").(bool)
		if enabledState {
			routeUpdateProperties.EnabledState = cdn.EnabledStateEnabled
		} else {
			routeUpdateProperties.EnabledState = cdn.EnabledStateDisabled
		}
	}

	if d.HasChange("https_redirect") {
		log.Printf("[DEBUG] Updating https redirect for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)

		httpsRedirect := d.Get("https_redirect").(bool)
		if httpsRedirect {
			routeUpdateProperties.HTTPSRedirect = cdn.HTTPSRedirectEnabled
		} else {
			routeUpdateProperties.HTTPSRedirect = cdn.HTTPSRedirectDisabled
		}
	}

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
	if d.HasChange("custom_domains") {
		log.Printf("[DEBUG] Updating link_to_default_domain configuration for route %s on endpoint %s", id.RouteName, id.AfdEndpointName)

		customDomainsArray := make([]cdn.ResourceReference, 0)
		for _, v := range customDomains {
			resourceId := v.(string)
			resourceReference := cdn.ResourceReference{
				ID: &resourceId,
			}
			customDomainsArray = append(customDomainsArray, resourceReference)
		}
		routeUpdateProperties.CustomDomains = &customDomainsArray

	}

	if d.HasChange("link_to_default_domain") {
		if linkToDefaultDomain {
			routeUpdateProperties.LinkToDefaultDomain = cdn.LinkToDefaultDomainEnabled
		} else {
			routeUpdateProperties.LinkToDefaultDomain = cdn.LinkToDefaultDomainDisabled
		}
	}

	// update caching & compression settings
	if cachingEnabled && (d.HasChange("query_string_caching_behavior") || d.HasChange("content_types_to_compress")) {

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

		compressionSettings := cdn.CompressionSettings{
			IsCompressionEnabled:   &cachingEnabled,
			ContentTypesToCompress: &contentTypesToCompressArray,
		}

		routeUpdateProperties.CompressionSettings = compressionSettings

	}

	if !cachingEnabled {
		routeUpdateProperties.QueryStringCachingBehavior = cdn.AfdQueryStringCachingBehaviorNotSet
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

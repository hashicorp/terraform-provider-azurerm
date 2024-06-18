// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCdnOriginGroups() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCdnOriginGroupsCreate,
		Read:   resourceCdnOriginGroupsRead,
		Update: resourceCdnOriginGroupsUpdate,
		Delete: resourceCdnOriginGroupsDelete,

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
			"endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.EndpointID,
			},

			"origin_group": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.OriginGroupName,
						},

						"health_probe": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MaxItems: 1,

							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"protocol": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.ProbeProtocolHTTP),
											string(cdn.ProbeProtocolHTTPS),
										}, false),
									},

									"request_type": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(cdn.HealthProbeRequestTypeHEAD),
										ValidateFunc: validation.StringInSlice([]string{
											string(cdn.HealthProbeRequestTypeGET),
											string(cdn.HealthProbeRequestTypeHEAD),
										}, false),
									},

									"interval_in_seconds": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(5, 31536000),
									},

									"path": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										Default:      "/",
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},

						"origins": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							MinItems: 1,

							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.OriginID,
							},
						},

						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceCdnOriginGroupsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM CDN Origin Groups creation...")

	id, err := parse.EndpointID(d.Get("endpoint_id").(string))
	if err != nil {
		return err
	}

	existing, err := endpointsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	// No import error for this resource as it does not make sense...

	if existing.Origins == nil {
		return fmt.Errorf("'origins' is required field but was not found in the referenced `azurerm_cdn_endpoint` %s", id)
	}

	endpoint := cdn.Endpoint{
		Location: existing.Location,
		EndpointProperties: &cdn.EndpointProperties{
			Origins:                          existing.Origins,
			OriginPath:                       existing.OriginPath,
			ContentTypesToCompress:           existing.ContentTypesToCompress,
			OriginHostHeader:                 existing.OriginHostHeader,
			IsCompressionEnabled:             existing.IsCompressionEnabled,
			IsHTTPAllowed:                    existing.IsHTTPAllowed,
			IsHTTPSAllowed:                   existing.IsHTTPSAllowed,
			QueryStringCachingBehavior:       existing.QueryStringCachingBehavior,
			OptimizationType:                 existing.OptimizationType,
			ProbePath:                        existing.ProbePath,
			GeoFilters:                       existing.GeoFilters,
			DefaultOriginGroup:               existing.DefaultOriginGroup,
			URLSigningKeys:                   existing.URLSigningKeys,
			DeliveryPolicy:                   existing.DeliveryPolicy,
			WebApplicationFirewallPolicyLink: existing.WebApplicationFirewallPolicyLink,
		},
		Tags: existing.Tags,
	}

	originGroups := expandAzureRMCdnEndpointOriginGroups(d.Get("origin_group").(*pluginsdk.Set).List())
	endpoint.EndpointProperties.OriginGroups = originGroups

	log.Printf("\n\n\n***************************************************************************************************")
	log.Printf("==> resourceCdnOriginGroupsCreate <==")
	log.Printf("***************************************************************************************************\n\n")
	log.Printf("originGroups: %+v", endpoint.EndpointProperties.OriginGroups)

	// now we have to fix up the origins values and swap them out for default values else the resource groups will not be created...
	for _, originGroupOrigins := range *originGroups {

		for _, originGroup := range *originGroupOrigins.Origins {

			originID, err := parse.OriginID(*originGroup.ID)
			if err != nil {
				return err
			}

			for _, origin := range *endpoint.Origins {
				if originID.Name == *origin.Name {
					origin.Priority = pointer.To(int32(1))
					origin.Weight = pointer.To(int32(1000))

					log.Printf("Priority: %d", origin.Priority)
					log.Printf("Weight: %d", origin.Weight)
				}

			}
		}
	}

	log.Printf("***************************************************************************************************")
	log.Printf("***************************************************************************************************\n\n\n")

	future, err := endpointsClient.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for the creation of `cdn_origin_groups` for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnOriginGroupsRead(d, meta)
}

func resourceCdnOriginGroupsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	endpointsClient := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM CDN Origin Groups update...")

	id, err := parse.EndpointID(d.Get("endpoint_id").(string))
	if err != nil {
		return err
	}

	existing, err := endpointsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if existing.Origins == nil {
		return fmt.Errorf("'origins' is required field but was not found in the referenced `azurerm_cdn_endpoint` %s", id)
	}

	endpoint := cdn.Endpoint{
		Location: existing.Location,
		EndpointProperties: &cdn.EndpointProperties{
			Origins:                          existing.Origins,
			OriginPath:                       existing.OriginPath,
			ContentTypesToCompress:           existing.ContentTypesToCompress,
			OriginHostHeader:                 existing.OriginHostHeader,
			IsCompressionEnabled:             existing.IsCompressionEnabled,
			IsHTTPAllowed:                    existing.IsHTTPAllowed,
			IsHTTPSAllowed:                   existing.IsHTTPSAllowed,
			QueryStringCachingBehavior:       existing.QueryStringCachingBehavior,
			OptimizationType:                 existing.OptimizationType,
			ProbePath:                        existing.ProbePath,
			GeoFilters:                       existing.GeoFilters,
			DefaultOriginGroup:               existing.DefaultOriginGroup,
			URLSigningKeys:                   existing.URLSigningKeys,
			DeliveryPolicy:                   existing.DeliveryPolicy,
			WebApplicationFirewallPolicyLink: existing.WebApplicationFirewallPolicyLink,
		},
		Tags: existing.Tags,
	}

	originGroups := expandAzureRMCdnEndpointOriginGroups(d.Get("origin_group").(*pluginsdk.Set).List())
	endpoint.EndpointProperties.OriginGroups = originGroups

	log.Printf("\n\n\n***************************************************************************************************")
	log.Printf("==> resourceCdnOriginGroupsUpdate <==")
	log.Printf("***************************************************************************************************\n\n")
	log.Printf("originGroups: %+v", endpoint.EndpointProperties.OriginGroups)

	// now we have to fix up the origins values and swap them out for default values else the resource groups will not be created...
	for _, originGroupOrigins := range *originGroups {

		for _, originGroup := range *originGroupOrigins.Origins {

			originID, err := parse.OriginID(*originGroup.ID)
			if err != nil {
				return err
			}

			for _, origin := range *endpoint.Origins {
				if originID.Name == *origin.Name {
					// if endpoint.DefaultOriginGroup == nil {
					// 	endpoint.DefaultOriginGroup = pointer.To(originGroup)
					// }
					if origin.OriginHostHeader == nil {
						origin.OriginHostHeader = origin.HostName
					}
					origin.Priority = pointer.To(int32(1))
					origin.Weight = pointer.To(int32(1000))

					log.Printf("Priority: %d", origin.Priority)
					log.Printf("Weight: %d", origin.Weight)
					log.Printf("OriginHostHeader: %s", *origin.OriginHostHeader)
				}

			}
		}
	}

	log.Printf("defaultOriginGroup: %s", *endpoint.DefaultOriginGroup.ID)
	log.Printf("***************************************************************************************************")
	log.Printf("***************************************************************************************************\n\n\n")

	future, err := endpointsClient.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("updating `cdn_origin_groups` for %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, endpointsClient.Client); err != nil {
		return fmt.Errorf("waiting for update of `cdn_origin_groups` for %s: %+v", *id, err)
	}

	return resourceCdnOriginGroupsRead(d, meta)
}

func resourceCdnOriginGroupsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("endpoint_id", id.ID())

	if props := existing.EndpointProperties; props != nil {
		originGroups := flattenAzureRMCdnEndpointOriginGroups(props.OriginGroups, id)
		if err := d.Set("origin_group", originGroups); err != nil {
			return fmt.Errorf("setting `origin_group`: %+v", err)
		}
	}

	return nil
}

func resourceCdnOriginGroupsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	endpoint := cdn.Endpoint{
		Location: existing.Location,
		EndpointProperties: &cdn.EndpointProperties{
			Origins:                          existing.Origins,
			OriginPath:                       existing.OriginPath,
			ContentTypesToCompress:           existing.ContentTypesToCompress,
			OriginHostHeader:                 existing.OriginHostHeader,
			IsCompressionEnabled:             existing.IsCompressionEnabled,
			IsHTTPAllowed:                    existing.IsHTTPAllowed,
			IsHTTPSAllowed:                   existing.IsHTTPSAllowed,
			QueryStringCachingBehavior:       existing.QueryStringCachingBehavior,
			OptimizationType:                 existing.OptimizationType,
			ProbePath:                        existing.ProbePath,
			GeoFilters:                       existing.GeoFilters,
			DefaultOriginGroup:               existing.DefaultOriginGroup,
			URLSigningKeys:                   existing.URLSigningKeys,
			DeliveryPolicy:                   existing.DeliveryPolicy,
			WebApplicationFirewallPolicyLink: existing.WebApplicationFirewallPolicyLink,
		},
		Tags: existing.Tags,
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("deleting `cdn_origin_groups` for %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of `cdn_origin_groups` for %s: %+v", *id, err)
	}

	return nil
}

func expandAzureRMCdnEndpointOriginGroups(input []interface{}) *[]cdn.DeepCreatedOriginGroup {
	results := make([]cdn.DeepCreatedOriginGroup, 0)
	if len(input) == 0 {
		return nil
	}

	log.Printf("\n\n\n***************************************************************************************************")
	log.Printf("==> expandAzureRMCdnEndpointOriginGroups <==")
	log.Printf("***************************************************************************************************\n\n")

	for _, v := range input {
		data := v.(map[string]interface{})

		log.Printf("data: %+v", data)

		result := cdn.DeepCreatedOriginGroup{
			DeepCreatedOriginGroupProperties: &cdn.DeepCreatedOriginGroupProperties{
				HealthProbeSettings: &cdn.HealthProbeParameters{},
				Origins:             &[]cdn.ResourceReference{},
			},
		}

		result.Name = pointer.To(data["name"].(string))

		healthProbeRaw := data["health_probe"].(*pluginsdk.Set).List()
		healthProbe := healthProbeRaw[0].(map[string]interface{})

		result.HealthProbeSettings.ProbeIntervalInSeconds = pointer.To(int32(healthProbe["interval_in_seconds"].(int)))
		result.HealthProbeSettings.ProbePath = pointer.To(healthProbe["path"].(string))
		result.HealthProbeSettings.ProbeRequestType = cdn.HealthProbeRequestType(healthProbe["request_type"].(string))

		result.HealthProbeSettings.ProbeProtocol = cdn.ProbeProtocolHTTP
		if strings.EqualFold(healthProbe["protocol"].(string), string(cdn.ProbeProtocolHTTPS)) {
			result.HealthProbeSettings.ProbeProtocol = cdn.ProbeProtocolHTTPS
		}

		log.Printf("interval_in_seconds: %d", *result.HealthProbeSettings.ProbeIntervalInSeconds)
		log.Printf("path: %s", *result.HealthProbeSettings.ProbePath)
		log.Printf("request_type: %+s", result.HealthProbeSettings.ProbeRequestType)
		log.Printf("protocol: %s", result.HealthProbeSettings.ProbeProtocol)

		origins := make([]cdn.ResourceReference, 0)

		for i, id := range data["origins"].(*pluginsdk.Set).List() {
			resourceReference := cdn.ResourceReference{
				ID: pointer.To(id.(string)),
			}

			log.Printf("origins[%d]: %s", i, *resourceReference.ID)

			origins = append(origins, resourceReference)
		}

		log.Printf("origins count: %d", len(origins))

		result.DeepCreatedOriginGroupProperties.Origins = pointer.To(origins)

		results = append(results, result)
	}

	return &results
}

func flattenAzureRMCdnEndpointOriginGroups(input *[]cdn.DeepCreatedOriginGroup, endpoint *parse.EndpointId) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, v := range *input {
		result := make(map[string]interface{}, 0)

		originGroupName := ""
		if v.Name != nil {
			originGroupName = *v.Name
			result["name"] = originGroupName
		}

		result["id"] = parse.NewOriginGroupID(endpoint.SubscriptionId, endpoint.ResourceGroup, endpoint.ProfileName, endpoint.Name, originGroupName).ID()

		if props := v.DeepCreatedOriginGroupProperties; props != nil {
			healthProbeSettings := make(map[string]interface{}, 0)
			if v := props.HealthProbeSettings; v != nil {
				healthProbeSettings["interval_in_seconds"] = *v.ProbeIntervalInSeconds
				healthProbeSettings["path"] = *v.ProbeIntervalInSeconds
				healthProbeSettings["request_type"] = v.ProbeRequestType
				healthProbeSettings["protocol"] = v.ProbeProtocol
			}

			result["health_probe"] = healthProbeSettings

			var origins []string
			for _, origin := range *props.Origins {
				origins = append(origins, *origin.ID)
			}

			result["origins"] = origins
		}

		results = append(results, result)
	}

	return results
}

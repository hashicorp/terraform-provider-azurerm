// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceCdnOrigins() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceCdnOriginsCreate,
		Read:   resourceCdnOriginsRead,
		Update: resourceCdnOriginsUpdate,
		Delete: resourceCdnOriginsDelete,

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
			"endpoint_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.EndpointID,
			},

			// {
			//         "enabled": true,
			// },

			"origins": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.OriginName,
						},

						"host_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"http_port": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      80,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"https_port": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      443,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"origin_host_header": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.Any(validation.IsIPv6Address, validation.IsIPv4Address, validation.StringIsNotEmpty),
						},

						"priority": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(1, 5),
						},

						"weight": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      1000,
							ValidateFunc: validation.IntBetween(1, 1000),
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

	return resource
}

func resourceCdnOriginsCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM CDN Origins creation...")

	id, err := parse.EndpointID(d.Get("endpoint_id").(string))
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}

	if existing.Origins != nil {
		return tf.ImportAsExistsError("azurerm_cdn_origins", id.ID())
	}

	endpoint := cdn.Endpoint{
		Location: existing.Location,
		EndpointProperties: &cdn.EndpointProperties{
			IsHTTPAllowed:              existing.IsHTTPAllowed,
			IsHTTPSAllowed:             existing.IsHTTPSAllowed,
			QueryStringCachingBehavior: existing.QueryStringCachingBehavior,
		},
		Tags: existing.Tags,
	}

	origins := expandAzureRmCdnOriginsOrigins(d)
	originCount := len(origins)
	endpoint.EndpointProperties.Origins = &origins

	if originCount > 1 && existing.DefaultOriginGroup == nil {
		return fmt.Errorf("%s creating more than one 'origin' is not allowed if the Default Origin Group has not been set", id)
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceCdnOriginsRead(d, meta)
}

func resourceCdnOriginsUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM CDN Origins update...")

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
			IsHTTPAllowed:              existing.IsHTTPAllowed,
			IsHTTPSAllowed:             existing.IsHTTPSAllowed,
			QueryStringCachingBehavior: existing.QueryStringCachingBehavior,
		},
		Tags: existing.Tags,
	}

	origins := expandAzureRmCdnOriginsOrigins(d)
	if len(origins) > 0 {
		endpoint.EndpointProperties.Origins = &origins
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	return resourceCdnOriginsRead(d, meta)
}

func resourceCdnOriginsRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if props := resp.EndpointProperties; props != nil {
		origins := flattenAzureRMCdnOriginsOrigins(props.Origins, subscriptionId, id)
		if err := d.Set("origins", origins); err != nil {
			return fmt.Errorf("setting `origins`: %+v", err)
		}
	}

	return nil
}

func resourceCdnOriginsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.EndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM CDN Origins delete...")

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
			IsHTTPAllowed:              existing.IsHTTPAllowed,
			IsHTTPSAllowed:             existing.IsHTTPSAllowed,
			QueryStringCachingBehavior: existing.QueryStringCachingBehavior,
		},
		Tags: existing.Tags,
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.ProfileName, id.Name, endpoint)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return resourceCdnOriginsRead(d, meta)
}

func expandAzureRmCdnOriginsOrigins(d *pluginsdk.ResourceData) []cdn.DeepCreatedOrigin {
	input := d.Get("origins").(*pluginsdk.Set).List()
	origins := make([]cdn.DeepCreatedOrigin, 0)

	for _, v := range input {
		data := v.(map[string]interface{})

		name := data["name"].(string)
		hostName := data["host_name"].(string)
		httpPort := data["http_port"].(int32)
		httpsPort := data["https_port"].(int32)
		originHostHeader := data["origin_host_header"].(string)
		priority := data["priority"].(int32)
		weight := data["weight"].(int32)

		origin := cdn.DeepCreatedOrigin{
			Name: utils.String(name),
			DeepCreatedOriginProperties: &cdn.DeepCreatedOriginProperties{
				HostName:         pointer.To(hostName),
				HTTPPort:         pointer.To(httpPort),
				HTTPSPort:        pointer.To(httpsPort),
				OriginHostHeader: pointer.To(originHostHeader),
				Priority:         pointer.To(priority),
				Weight:           pointer.To(weight),
				Enabled:          pointer.To(true),
			},
		}

		if v, ok := data["https_port"]; ok {
			port := v.(int)
			origin.DeepCreatedOriginProperties.HTTPSPort = utils.Int32(int32(port))
		}

		if v, ok := data["http_port"]; ok {
			port := v.(int)
			origin.DeepCreatedOriginProperties.HTTPPort = utils.Int32(int32(port))
		}

		origins = append(origins, origin)
	}

	return origins
}

func flattenAzureRMCdnOriginsOrigins(input *[]cdn.DeepCreatedOrigin, subscriptionId string, endpointId *parse.EndpointId) []interface{} {
	results := make([]interface{}, 0)

	if list := input; list != nil {
		for _, i := range *list {
			name := ""
			if i.Name != nil {
				name = *i.Name
			}

			id := parse.NewOriginID(subscriptionId, endpointId.ResourceGroup, endpointId.ProfileName, endpointId.Name, name)

			var hostName string
			var httpPort int32
			var httpsPort int32
			var originHostHeader string
			var priority int32
			var weight int32

			if props := i.DeepCreatedOriginProperties; props != nil {
				if v := props.HostName; v != nil {
					hostName = pointer.From(v)
				}

				if v := props.HTTPPort; v != nil {
					httpPort = pointer.From(v)
				}

				if v := props.HTTPSPort; v != nil {
					httpsPort = pointer.From(v)
				}

				if v := props.OriginHostHeader; v != nil {
					originHostHeader = pointer.From(v)
				}

				if v := props.Priority; v != nil {
					priority = pointer.From(v)
				}

				if v := props.Weight; v != nil {
					weight = pointer.From(v)
				}

			}

			results = append(results, map[string]interface{}{
				"name":               name,
				"host_name":          hostName,
				"http_port":          httpPort,
				"https_port":         httpsPort,
				"origin_host_header": originHostHeader,
				"priority":           priority,
				"weight":             weight,
				"id":                 id.ID(),
			})
		}
	}

	return results
}

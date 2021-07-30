package signalr

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/sdk/signalr"
	signalrValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/signalr/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSignalRService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmSignalRServiceCreate,
		Read:   resourceArmSignalRServiceRead,
		Update: resourceArmSignalRServiceUpdate,
		Delete: resourceArmSignalRServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := signalr.ParseSignalRID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Free_F1",
								"Standard_S1",
							}, false),
						},

						"capacity": {
							Type:         pluginsdk.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 5, 10, 20, 50, 100}),
						},
					},
				},
			},

			"features": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"flag": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(signalr.FeatureFlagsEnableConnectivityLogs),
								string(signalr.FeatureFlagsEnableMessagingLogs),
								string(signalr.FeatureFlagsServiceMode),
							}, false),
						},

						"value": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"upstream_endpoint": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"category_pattern": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"event_pattern": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"hub_pattern": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"url_template": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: signalrValidate.UrlTemplate,
						},
					},
				},
			},

			"cors": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"allowed_origins": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"server_port": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmSignalRServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	location := azure.NormalizeLocation(d.Get("location").(string))

	id := signalr.NewSignalRID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_signalr_service", id.ID())
	}

	sku := d.Get("sku").([]interface{})
	featureFlags := d.Get("features").(*pluginsdk.Set).List()
	cors := d.Get("cors").([]interface{})
	upstreamSettings := d.Get("upstream_endpoint").(*pluginsdk.Set).List()

	expandedFeatures := expandSignalRFeatures(featureFlags)

	// Upstream configurations are only allowed when the SignalR service is in `Serverless` mode
	if len(upstreamSettings) > 0 && !signalRIsInServerlessMode(expandedFeatures) {
		return fmt.Errorf("Upstream configurations are only allowed when the SignalR Service is in `Serverless` mode")
	}

	resourceType := signalr.SignalRResource{
		Location: utils.String(location),
		Properties: &signalr.SignalRProperties{
			Cors:     expandSignalRCors(cors),
			Features: expandedFeatures,
			Upstream: expandUpstreamSettings(upstreamSettings),
		},
		Sku:  expandSignalRServiceSku(sku),
		Tags: expandTags(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, resourceType); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmSignalRServiceUpdate(d, meta)
}

func resourceArmSignalRServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keys, err := client.ListKeys(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing keys for %s: %+v", *id, err)
	}

	d.Set("name", id.SignalRName)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if err = d.Set("sku", flattenSignalRServiceSku(model.Sku)); err != nil {
			return fmt.Errorf("setting `sku`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("hostname", props.HostName)
			d.Set("ip_address", props.ExternalIP)
			d.Set("public_port", props.PublicPort)
			d.Set("server_port", props.ServerPort)

			if err := d.Set("features", flattenSignalRFeatures(props.Features)); err != nil {
				return fmt.Errorf("setting `features`: %+v", err)
			}

			if err := d.Set("cors", flattenSignalRCors(props.Cors)); err != nil {
				return fmt.Errorf("setting `cors`: %+v", err)
			}

			if err := d.Set("upstream_endpoint", flattenUpstreamSettings(props.Upstream)); err != nil {
				return fmt.Errorf("setting `upstream_endpoint`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, flattenTags(model.Tags)); err != nil {
				return err
			}
		}
	}

	if model := keys.Model; model != nil {
		d.Set("primary_access_key", model.PrimaryKey)
		d.Set("primary_connection_string", model.PrimaryConnectionString)
		d.Set("secondary_access_key", model.SecondaryKey)
		d.Set("secondary_connection_string", model.SecondaryConnectionString)
	}

	return nil
}

func resourceArmSignalRServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	resourceType := signalr.SignalRResource{}

	if d.HasChanges("cors", "features", "upstream_endpoint") {
		resourceType.Properties = &signalr.SignalRProperties{}

		if d.HasChange("cors") {
			corsRaw := d.Get("cors").([]interface{})
			resourceType.Properties.Cors = expandSignalRCors(corsRaw)
		}

		if d.HasChange("features") {
			featuresRaw := d.Get("features").(*pluginsdk.Set).List()
			resourceType.Properties.Features = expandSignalRFeatures(featuresRaw)
		}

		if d.HasChange("upstream_endpoint") {
			featuresRaw := d.Get("upstream_endpoint").(*pluginsdk.Set).List()
			resourceType.Properties.Upstream = expandUpstreamSettings(featuresRaw)
		}
	}

	if d.HasChange("sku") {
		sku := d.Get("sku").([]interface{})
		resourceType.Sku = expandSignalRServiceSku(sku)
	}

	if d.HasChange("tags") {
		tagsRaw := d.Get("tags").(map[string]interface{})
		resourceType.Tags = expandTags(tagsRaw)
	}

	if err := client.UpdateThenPoll(ctx, *id, resourceType); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmSignalRServiceRead(d, meta)
}

func resourceArmSignalRServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.Client
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func signalRIsInServerlessMode(features *[]signalr.SignalRFeature) bool {
	if features == nil {
		return false
	}

	for _, feature := range *features {
		if feature.Flag == signalr.FeatureFlagsServiceMode {
			return strings.EqualFold(feature.Value, "Serverless")
		}
	}

	return false
}

func expandSignalRFeatures(input []interface{}) *[]signalr.SignalRFeature {
	features := make([]signalr.SignalRFeature, 0)
	for _, featureValue := range input {
		value := featureValue.(map[string]interface{})

		feature := signalr.SignalRFeature{
			Flag:  signalr.FeatureFlags(value["flag"].(string)),
			Value: value["value"].(string),
		}

		features = append(features, feature)
	}

	return &features
}

func flattenSignalRFeatures(features *[]signalr.SignalRFeature) []interface{} {
	if features == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0)
	for _, feature := range *features {
		result = append(result, map[string]interface{}{
			"flag":  string(feature.Flag),
			"value": feature.Value,
		})
	}
	return result
}

func expandUpstreamSettings(input []interface{}) *signalr.ServerlessUpstreamSettings {
	upstreamTemplates := make([]signalr.UpstreamTemplate, 0)

	for _, upstreamSetting := range input {
		setting := upstreamSetting.(map[string]interface{})

		upstreamTemplate := signalr.UpstreamTemplate{
			HubPattern:      utils.String(strings.Join(*utils.ExpandStringSlice(setting["hub_pattern"].([]interface{})), ",")),
			EventPattern:    utils.String(strings.Join(*utils.ExpandStringSlice(setting["event_pattern"].([]interface{})), ",")),
			CategoryPattern: utils.String(strings.Join(*utils.ExpandStringSlice(setting["category_pattern"].([]interface{})), ",")),
			UrlTemplate:     setting["url_template"].(string),
		}

		upstreamTemplates = append(upstreamTemplates, upstreamTemplate)
	}

	return &signalr.ServerlessUpstreamSettings{
		Templates: &upstreamTemplates,
	}
}

func flattenUpstreamSettings(upstreamSettings *signalr.ServerlessUpstreamSettings) []interface{} {
	result := make([]interface{}, 0)
	if upstreamSettings == nil || upstreamSettings.Templates == nil {
		return result
	}

	for _, settings := range *upstreamSettings.Templates {
		categoryPattern := make([]interface{}, 0)
		if settings.CategoryPattern != nil {
			categoryPatterns := strings.Split(*settings.CategoryPattern, ",")
			categoryPattern = utils.FlattenStringSlice(&categoryPatterns)
		}

		eventPattern := make([]interface{}, 0)
		if settings.EventPattern != nil {
			eventPatterns := strings.Split(*settings.EventPattern, ",")
			eventPattern = utils.FlattenStringSlice(&eventPatterns)
		}

		hubPattern := make([]interface{}, 0)
		if settings.HubPattern != nil {
			hubPatterns := strings.Split(*settings.HubPattern, ",")
			hubPattern = utils.FlattenStringSlice(&hubPatterns)
		}

		result = append(result, map[string]interface{}{
			"url_template":     settings.UrlTemplate,
			"hub_pattern":      hubPattern,
			"event_pattern":    eventPattern,
			"category_pattern": categoryPattern,
		})
	}
	return result
}

func expandSignalRCors(input []interface{}) *signalr.SignalRCorsSettings {
	corsSettings := signalr.SignalRCorsSettings{}

	if len(input) == 0 || input[0] == nil {
		return &corsSettings
	}

	setting := input[0].(map[string]interface{})
	origins := setting["allowed_origins"].(*pluginsdk.Set).List()

	allowedOrigins := make([]string, 0)
	for _, param := range origins {
		allowedOrigins = append(allowedOrigins, param.(string))
	}

	corsSettings.AllowedOrigins = &allowedOrigins

	return &corsSettings
}

func flattenSignalRCors(input *signalr.SignalRCorsSettings) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})

	allowedOrigins := make([]interface{}, 0)
	if s := input.AllowedOrigins; s != nil {
		for _, v := range *s {
			allowedOrigins = append(allowedOrigins, v)
		}
	}
	result["allowed_origins"] = pluginsdk.NewSet(pluginsdk.HashString, allowedOrigins)

	return append(results, result)
}

func expandSignalRServiceSku(input []interface{}) *signalr.ResourceSku {
	v := input[0].(map[string]interface{})
	return &signalr.ResourceSku{
		Name:     v["name"].(string),
		Capacity: utils.Int64(int64(v["capacity"].(int))),
	}
}

func flattenSignalRServiceSku(input *signalr.ResourceSku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	capacity := 0
	if input.Capacity != nil {
		capacity = int(*input.Capacity)
	}

	return []interface{}{
		map[string]interface{}{
			"capacity": capacity,
			"name":     input.Name,
		},
	}
}

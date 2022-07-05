package signalr

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/sdk/2020-05-01/signalr"
	signalrValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ServiceV0ToV1{},
		}),
		SchemaVersion: 1,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := signalr.ParseSignalRID(id)
			return err
		}),

		Schema: resourceArmSignalRServiceSchema(),
	}
}

func resourceArmSignalRServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
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
	connectivityLogsEnabled := false
	if v, ok := d.GetOk("connectivity_logs_enabled"); ok {
		connectivityLogsEnabled = v.(bool)
	}
	messagingLogsEnabled := false
	if v, ok := d.GetOk("messaging_logs_enabled"); ok {
		messagingLogsEnabled = v.(bool)
	}
	liveTraceEnabled := false
	if v, ok := d.GetOk("live_trace_enabled"); ok {
		liveTraceEnabled = v.(bool)
	}
	serviceMode := "Default"
	if v, ok := d.GetOk("service_mode"); ok {
		serviceMode = v.(string)
	}

	cors := d.Get("cors").([]interface{})
	upstreamSettings := d.Get("upstream_endpoint").(*pluginsdk.Set).List()

	expandedFeatures := make([]signalr.SignalRFeature, 0)
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsEnableConnectivityLogs, strconv.FormatBool(connectivityLogsEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsEnableMessagingLogs, strconv.FormatBool(messagingLogsEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature("EnableLiveTrace", strconv.FormatBool(liveTraceEnabled)))
	expandedFeatures = append(expandedFeatures, signalRFeature(signalr.FeatureFlagsServiceMode, serviceMode))

	// Upstream configurations are only allowed when the SignalR service is in `Serverless` mode
	if len(upstreamSettings) > 0 && !signalRIsInServerlessMode(&expandedFeatures) {
		return fmt.Errorf("Upstream configurations are only allowed when the SignalR Service is in `Serverless` mode")
	}

	resourceType := signalr.SignalRResource{
		Location: utils.String(location),
		Properties: &signalr.SignalRProperties{
			Cors:     expandSignalRCors(cors),
			Features: &expandedFeatures,
			Upstream: expandUpstreamSettings(upstreamSettings),
		},
		Sku:  expandSignalRServiceSku(sku),
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, resourceType); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceArmSignalRServiceUpdate(d, meta)
}

func resourceArmSignalRServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
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

	d.Set("name", id.ResourceName)
	d.Set("resource_group_name", id.ResourceGroupName)

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

			connectivityLogsEnabled := false
			messagingLogsEnabled := false
			liveTraceEnabled := false
			serviceMode := "Default"
			for _, feature := range *props.Features {
				if feature.Flag == signalr.FeatureFlagsEnableConnectivityLogs {
					connectivityLogsEnabled = strings.EqualFold(feature.Value, "True")
				}
				if feature.Flag == signalr.FeatureFlagsEnableMessagingLogs {
					messagingLogsEnabled = strings.EqualFold(feature.Value, "True")
				}
				if feature.Flag == "EnableLiveTrace" {
					liveTraceEnabled = strings.EqualFold(feature.Value, "True")
				}
				if feature.Flag == signalr.FeatureFlagsServiceMode {
					serviceMode = feature.Value
				}
			}
			d.Set("connectivity_logs_enabled", connectivityLogsEnabled)
			d.Set("messaging_logs_enabled", messagingLogsEnabled)
			d.Set("live_trace_enabled", liveTraceEnabled)
			d.Set("service_mode", serviceMode)

			if err := d.Set("cors", flattenSignalRCors(props.Cors)); err != nil {
				return fmt.Errorf("setting `cors`: %+v", err)
			}

			if err := d.Set("upstream_endpoint", flattenUpstreamSettings(props.Upstream)); err != nil {
				return fmt.Errorf("setting `upstream_endpoint`: %+v", err)
			}

			if err := tags.FlattenAndSet(d, model.Tags); err != nil {
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
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	resourceType := signalr.SignalRResource{}

	if d.HasChanges("cors", "features", "upstream_endpoint", "connectivity_logs_enabled", "messaging_logs_enabled", "service_mode", "live_trace_enabled") {
		resourceType.Properties = &signalr.SignalRProperties{}

		if d.HasChange("cors") {
			corsRaw := d.Get("cors").([]interface{})
			resourceType.Properties.Cors = expandSignalRCors(corsRaw)
		}

		if d.HasChanges("connectivity_logs_enabled", "messaging_logs_enabled", "service_mode", "live_trace_enabled") {
			features := make([]signalr.SignalRFeature, 0)
			if d.HasChange("connectivity_logs_enabled") {
				connectivityLogsEnabled := false
				if v, ok := d.GetOk("connectivity_logs_enabled"); ok {
					connectivityLogsEnabled = v.(bool)
				}
				features = append(features, signalRFeature(signalr.FeatureFlagsEnableConnectivityLogs, strconv.FormatBool(connectivityLogsEnabled)))
			}

			if d.HasChange("messaging_logs_enabled") {
				messagingLogsEnabled := false
				if v, ok := d.GetOk("messaging_logs_enabled"); ok {
					messagingLogsEnabled = v.(bool)
				}
				features = append(features, signalRFeature(signalr.FeatureFlagsEnableMessagingLogs, strconv.FormatBool(messagingLogsEnabled)))
			}

			if d.HasChange("live_trace_enabled") {
				liveTraceEnabled := false
				if v, ok := d.GetOk("live_trace_enabled"); ok {
					liveTraceEnabled = v.(bool)
				}
				features = append(features, signalRFeature("EnableLiveTrace", strconv.FormatBool(liveTraceEnabled)))
			}

			if d.HasChange("service_mode") {
				serviceMode := "Default"
				if v, ok := d.GetOk("service_mode"); ok {
					serviceMode = v.(string)
				}
				features = append(features, signalRFeature(signalr.FeatureFlagsServiceMode, serviceMode))
			}
			resourceType.Properties.Features = &features
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
		resourceType.Tags = tags.Expand(tagsRaw)
	}

	if err := client.UpdateThenPoll(ctx, *id, resourceType); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmSignalRServiceRead(d, meta)
}

func resourceArmSignalRServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SignalR.SignalRClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := signalr.ParseSignalRID(d.Id())
	if err != nil {
		return err
	}

	// @tombuildsstuff: we can't use DeleteThenPoll here since the API returns a 404 on the Future in time
	future, err := client.Delete(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err := future.Poller.PollUntilDone(); err != nil {
		if !response.WasNotFound(future.Poller.HttpResponse) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
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

func signalRFeature(featureFlag signalr.FeatureFlags, value string) signalr.SignalRFeature {
	return signalr.SignalRFeature{
		Flag:  featureFlag,
		Value: value,
	}
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

func resourceArmSignalRServiceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.NoZeroValues,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

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
							"Premium_P1",
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

		"connectivity_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"messaging_logs_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"live_trace_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"service_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "Default",
			ValidateFunc: validation.StringInSlice([]string{
				"Serverless",
				"Classic",
				"Default",
			}, false),
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

		"tags": commonschema.Tags(),
	}
}

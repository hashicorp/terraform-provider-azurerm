package springcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/appplatform/mgmt/2022-05-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSpringCloudGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSpringCloudGatewayCreateUpdate,
		Read:   resourceSpringCloudGatewayRead,
		Update: resourceSpringCloudGatewayCreateUpdate,
		Delete: resourceSpringCloudGatewayDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SpringCloudGatewayID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
				}, false),
			},

			"spring_cloud_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SpringCloudServiceID,
			},

			"api_metadata": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"documentation_url": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},

						"server_url": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsURLWithHTTPorHTTPS,
						},

						"title": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"version": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"cors": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"credentials_allowed": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
						},

						"allowed_headers": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"allowed_methods": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"DELETE",
									"GET",
									"HEAD",
									"MERGE",
									"POST",
									"OPTIONS",
									"PUT",
								}, false),
							},
						},

						"allowed_origins": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"exposed_headers": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"max_age_seconds": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"https_only": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"instance_count": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 500),
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"quota": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"cpu": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "1",
							// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"memory": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "2Gi",
							// NOTE: we're intentionally not validating this field since additional values are possible when enabled by the service team
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"sso": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"client_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"client_secret": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"issuer_uri": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"scope": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},

			"url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}
func resourceSpringCloudGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.GatewayClient
	servicesClient := meta.(*clients.Client).AppPlatform.ServicesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	springId, err := parse.SpringCloudServiceID(d.Get("spring_cloud_service_id").(string))
	if err != nil {
		return err
	}
	id := parse.NewSpringCloudGatewayID(subscriptionId, springId.ResourceGroup, springId.SpringName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_spring_cloud_gateway", id.ID())
		}
	}

	service, err := servicesClient.Get(ctx, springId.ResourceGroup, springId.SpringName)
	if err != nil {
		return fmt.Errorf("checking for presence of existing Spring Cloud Service %q (Resource Group %q): %+v", springId.SpringName, springId.ResourceGroup, err)
	}
	if service.Sku == nil || service.Sku.Name == nil || service.Sku.Tier == nil {
		return fmt.Errorf("invalid `sku` for Spring Cloud Service %q (Resource Group %q)", springId.SpringName, springId.ResourceGroup)
	}

	gatewayResource := appplatform.GatewayResource{
		Properties: &appplatform.GatewayProperties{
			APIMetadataProperties: expandGatewayGatewayAPIMetadataProperties(d.Get("api_metadata").([]interface{})),
			CorsProperties:        expandGatewayGatewayCorsProperties(d.Get("cors").([]interface{})),
			HTTPSOnly:             utils.Bool(d.Get("https_only").(bool)),
			Public:                utils.Bool(d.Get("public_network_access_enabled").(bool)),
			ResourceRequests:      expandGatewayGatewayResourceRequests(d.Get("quota").([]interface{})),
			SsoProperties:         expandGatewaySsoProperties(d.Get("sso").([]interface{})),
		},
		Sku: &appplatform.Sku{
			Name:     service.Sku.Name,
			Tier:     service.Sku.Tier,
			Capacity: utils.Int32(int32(d.Get("instance_count").(int))),
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.GatewayName, gatewayResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSpringCloudGatewayRead(d, meta)
}

func resourceSpringCloudGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.GatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] appplatform %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}
	d.Set("name", id.GatewayName)
	d.Set("spring_cloud_service_id", parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID())
	if resp.Sku != nil {
		d.Set("instance_count", resp.Sku.Capacity)
	}
	if props := resp.Properties; props != nil {
		if err := d.Set("api_metadata", flattenGatewayGatewayAPIMetadataProperties(props.APIMetadataProperties)); err != nil {
			return fmt.Errorf("setting `api_metadata`: %+v", err)
		}
		if err := d.Set("cors", flattenGatewayGatewayCorsProperties(props.CorsProperties)); err != nil {
			return fmt.Errorf("setting `cors`: %+v", err)
		}
		d.Set("https_only", props.HTTPSOnly)
		d.Set("public_network_access_enabled", props.Public)
		if err := d.Set("quota", flattenGatewayGatewayResourceRequests(props.ResourceRequests)); err != nil {
			return fmt.Errorf("setting `resource_requests`: %+v", err)
		}
		if err := d.Set("sso", flattenGatewaySsoProperties(props.SsoProperties, d.Get("sso").([]interface{}))); err != nil {
			return fmt.Errorf("setting `sso`: %+v", err)
		}
		d.Set("url", props.URL)
	}
	return nil
}

func resourceSpringCloudGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AppPlatform.GatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SpringCloudGatewayID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.GatewayName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}
	return nil
}

func expandGatewayGatewayAPIMetadataProperties(input []interface{}) *appplatform.GatewayAPIMetadataProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &appplatform.GatewayAPIMetadataProperties{
		Title:         utils.String(v["title"].(string)),
		Description:   utils.String(v["description"].(string)),
		Documentation: utils.String(v["documentation_url"].(string)),
		Version:       utils.String(v["version"].(string)),
		ServerURL:     utils.String(v["server_url"].(string)),
	}
}

func expandGatewayGatewayCorsProperties(input []interface{}) *appplatform.GatewayCorsProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &appplatform.GatewayCorsProperties{
		AllowedOrigins:   utils.ExpandStringSlice(v["allowed_origins"].(*pluginsdk.Set).List()),
		AllowedMethods:   utils.ExpandStringSlice(v["allowed_methods"].(*pluginsdk.Set).List()),
		AllowedHeaders:   utils.ExpandStringSlice(v["allowed_headers"].(*pluginsdk.Set).List()),
		MaxAge:           utils.Int32(int32(v["max_age_seconds"].(int))),
		AllowCredentials: utils.Bool(v["credentials_allowed"].(bool)),
		ExposedHeaders:   utils.ExpandStringSlice(v["exposed_headers"].(*pluginsdk.Set).List()),
	}
}

func expandGatewayGatewayResourceRequests(input []interface{}) *appplatform.GatewayResourceRequests {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &appplatform.GatewayResourceRequests{
		CPU:    utils.String(v["cpu"].(string)),
		Memory: utils.String(v["memory"].(string)),
	}
}

func expandGatewaySsoProperties(input []interface{}) *appplatform.SsoProperties {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &appplatform.SsoProperties{
		Scope:        utils.ExpandStringSlice(v["scope"].(*pluginsdk.Set).List()),
		ClientID:     utils.String(v["client_id"].(string)),
		ClientSecret: utils.String(v["client_secret"].(string)),
		IssuerURI:    utils.String(v["issuer_uri"].(string)),
	}
}

func flattenGatewayGatewayAPIMetadataProperties(input *appplatform.GatewayAPIMetadataProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var description string
	if input.Description != nil {
		description = *input.Description
	}
	var documentation string
	if input.Documentation != nil {
		documentation = *input.Documentation
	}
	var serverUrl string
	if input.ServerURL != nil {
		serverUrl = *input.ServerURL
	}
	var title string
	if input.Title != nil {
		title = *input.Title
	}
	var version string
	if input.Version != nil {
		version = *input.Version
	}
	return []interface{}{
		map[string]interface{}{
			"description":       description,
			"documentation_url": documentation,
			"server_url":        serverUrl,
			"title":             title,
			"version":           version,
		},
	}
}

func flattenGatewayGatewayCorsProperties(input *appplatform.GatewayCorsProperties) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var allowCredentials bool
	if input.AllowCredentials != nil {
		allowCredentials = *input.AllowCredentials
	}
	var maxAge int32
	if input.MaxAge != nil {
		maxAge = *input.MaxAge
	}
	return []interface{}{
		map[string]interface{}{
			"credentials_allowed": allowCredentials,
			"allowed_headers":     utils.FlattenStringSlice(input.AllowedHeaders),
			"allowed_methods":     utils.FlattenStringSlice(input.AllowedMethods),
			"allowed_origins":     utils.FlattenStringSlice(input.AllowedOrigins),
			"exposed_headers":     utils.FlattenStringSlice(input.ExposedHeaders),
			"max_age_seconds":     maxAge,
		},
	}
}

func flattenGatewayGatewayResourceRequests(input *appplatform.GatewayResourceRequests) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var cpu string
	if input.CPU != nil {
		cpu = *input.CPU
	}
	var memory string
	if input.Memory != nil {
		memory = *input.Memory
	}
	return []interface{}{
		map[string]interface{}{
			"cpu":    cpu,
			"memory": memory,
		},
	}
}

func flattenGatewaySsoProperties(input *appplatform.SsoProperties, old []interface{}) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	oldItems := make(map[string]map[string]interface{})
	for _, item := range old {
		v := item.(map[string]interface{})
		if name, ok := v["issuer_uri"]; ok {
			oldItems[name.(string)] = v
		}
	}

	var issuerUri string
	if input.IssuerURI != nil {
		issuerUri = *input.IssuerURI
	}
	var clientId string
	var clientSecret string
	if oldItem, ok := oldItems[issuerUri]; ok {
		if value, ok := oldItem["client_id"]; ok {
			clientId = value.(string)
		}
		if value, ok := oldItem["client_secret"]; ok {
			clientSecret = value.(string)
		}
	}
	return []interface{}{
		map[string]interface{}{
			"client_id":     clientId,
			"client_secret": clientSecret,
			"issuer_uri":    issuerUri,
			"scope":         utils.FlattenStringSlice(input.Scope),
		},
	}
}

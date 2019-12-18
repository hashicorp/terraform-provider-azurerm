package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/signalr/mgmt/2018-10-01/signalr"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSignalRService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSignalRServiceCreateUpdate,
		Read:   resourceArmSignalRServiceRead,
		Update: resourceArmSignalRServiceCreateUpdate,
		Delete: resourceArmSignalRServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Free_F1",
								"Standard_S1",
							}, false),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntInSlice([]int{1, 2, 5, 10, 20, 50, 100}),
						},
					},
				},
			},

			"features": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"flag": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(signalr.EnableConnectivityLogs),
								string(signalr.ServiceMode),
							}, false),
						},

						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"cors": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"server_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"primary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmSignalRServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SignalR.Client
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	sku := d.Get("sku").([]interface{})
	t := d.Get("tags").(map[string]interface{})
	featureFlags := d.Get("features").([]interface{})
	cors := d.Get("cors").([]interface{})
	expandedTags := tags.Expand(t)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_signalr_service", *existing.ID)
		}
	}

	properties := &signalr.CreateOrUpdateProperties{
		Cors:     expandSignalRCors(cors),
		Features: expandSignalRFeatures(featureFlags),
	}

	parameters := &signalr.CreateParameters{
		Location:   utils.String(location),
		Sku:        expandSignalRServiceSku(sku),
		Tags:       expandedTags,
		Properties: properties,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating or updating SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the result of creating or updating SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("SignalR %q (Resource Group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmSignalRServiceRead(d, meta)
}

func resourceArmSignalRServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SignalR.Client
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["SignalR"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] SignalR %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	keys, err := client.ListKeys(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error getting keys of SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err = d.Set("sku", flattenSignalRServiceSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if properties := resp.Properties; properties != nil {
		d.Set("hostname", properties.HostName)
		d.Set("ip_address", properties.ExternalIP)
		d.Set("public_port", properties.PublicPort)
		d.Set("server_port", properties.ServerPort)

		if err := d.Set("features", flattenSignalRFeatures(properties.Features)); err != nil {
			return fmt.Errorf("Error setting `features`: %+v", err)
		}

		if err := d.Set("cors", flattenSignalRCors(properties.Cors)); err != nil {
			return fmt.Errorf("Error setting `cors`: %+v", err)
		}
	}

	d.Set("primary_access_key", keys.PrimaryKey)
	d.Set("primary_connection_string", keys.PrimaryConnectionString)
	d.Set("secondary_access_key", keys.SecondaryKey)
	d.Set("secondary_connection_string", keys.SecondaryConnectionString)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmSignalRServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).SignalR.Client
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["SignalR"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		return nil
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandSignalRFeatures(input []interface{}) *[]signalr.Feature {
	features := make([]signalr.Feature, 0)
	for _, featureValue := range input {
		value := featureValue.(map[string]interface{})

		feature := signalr.Feature{
			Flag:  signalr.FeatureFlags(value["flag"].(string)),
			Value: utils.String(value["value"].(string)),
		}

		features = append(features, feature)
	}

	return &features
}

func flattenSignalRFeatures(features *[]signalr.Feature) []interface{} {
	result := make([]interface{}, 0)
	if features != nil {
		for _, feature := range *features {
			value := ""
			if feature.Value != nil {
				value = *feature.Value
			}

			result = append(result, map[string]interface{}{
				"flag":  string(feature.Flag),
				"value": value,
			})
		}
	}
	return result
}

func expandSignalRCors(input []interface{}) *signalr.CorsSettings {
	corsSettings := signalr.CorsSettings{}

	if len(input) == 0 {
		return &corsSettings
	}

	setting := input[0].(map[string]interface{})
	origins := setting["allowed_origins"].(*schema.Set).List()

	allowedOrigins := make([]string, 0)
	for _, param := range origins {
		allowedOrigins = append(allowedOrigins, param.(string))
	}

	corsSettings.AllowedOrigins = &allowedOrigins

	return &corsSettings
}

func flattenSignalRCors(input *signalr.CorsSettings) []interface{} {
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
	result["allowed_origins"] = schema.NewSet(schema.HashString, allowedOrigins)

	return append(results, result)
}

func expandSignalRServiceSku(input []interface{}) *signalr.ResourceSku {
	v := input[0].(map[string]interface{})
	return &signalr.ResourceSku{
		Name:     utils.String(v["name"].(string)),
		Capacity: utils.Int32(int32(v["capacity"].(int))),
	}
}

func flattenSignalRServiceSku(input *signalr.ResourceSku) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if input.Name != nil {
		result["name"] = *input.Name
	}

	if input.Capacity != nil {
		result["capacity"] = *input.Capacity
	}

	return []interface{}{result}
}

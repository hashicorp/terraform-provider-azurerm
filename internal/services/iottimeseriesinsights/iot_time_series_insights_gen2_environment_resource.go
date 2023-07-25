// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIoTTimeSeriesInsightsGen2Environment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIoTTimeSeriesInsightsGen2EnvironmentCreateUpdate,
		Read:   resourceIoTTimeSeriesInsightsGen2EnvironmentRead,
		Update: resourceIoTTimeSeriesInsightsGen2EnvironmentCreateUpdate,
		Delete: resourceIoTTimeSeriesInsightsGen2EnvironmentDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := environments.ParseEnvironmentID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w\._\(\)]+$`),
					"IoT Time Series Insights Gen2 Environment name must contain only word characters, periods, underscores, and parentheses.",
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(environments.SkuNameLOne),
				}, false),
			},

			"warm_store_data_retention_time": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azValidate.ISO8601Duration,
			},
			"id_properties": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			"storage": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"key": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"data_access_fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceIoTTimeSeriesInsightsGen2EnvironmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.Environments
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := environments.NewEnvironmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, environments.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_gen2_environment", id.ID())
		}
	}

	sku := convertEnvironmentSkuName(d.Get("sku_name").(string))
	payload := environments.Gen2EnvironmentCreateOrUpdateParameters{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku:      sku,
		Properties: environments.Gen2EnvironmentCreationProperties{
			TimeSeriesIdProperties: expandIdProperties(d.Get("id_properties").([]interface{})),
			StorageConfiguration:   expandStorage(d.Get("storage").([]interface{})),
		},
	}

	if v, ok := d.GetOk("warm_store_data_retention_time"); ok {
		payload.Properties.WarmStoreConfiguration = &environments.WarmStoreConfigurationProperties{
			DataRetention: v.(string),
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIoTTimeSeriesInsightsGen2EnvironmentRead(d, meta)
}

func resourceIoTTimeSeriesInsightsGen2EnvironmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.Environments
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := environments.ParseEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, environments.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.EnvironmentName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		environment, ok := (*model).(environments.Gen2EnvironmentResource)
		if !ok {
			return fmt.Errorf("retrieving %s: expected a Gen2EnvironmentResource but got: %+v", *id, *model)
		}

		d.Set("sku_name", string(environment.Sku.Name))
		d.Set("location", location.Normalize(environment.Location))
		d.Set("data_access_fqdn", environment.Properties.DataAccessFqdn)
		if err := d.Set("id_properties", flattenIdProperties(environment.Properties.TimeSeriesIdProperties)); err != nil {
			return fmt.Errorf("setting `id_properties`: %+v", err)
		}
		if props := environment.Properties.WarmStoreConfiguration; props != nil {
			d.Set("warm_store_data_retention_time", props.DataRetention)
		}
		if err := d.Set("storage", flattenIoTTimeSeriesGen2EnvironmentStorage(environment.Properties.StorageConfiguration, d.Get("storage.0.key").(string))); err != nil {
			return fmt.Errorf("setting `storage`: %+v", err)
		}

		if err := tags.FlattenAndSet(d, environment.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceIoTTimeSeriesInsightsGen2EnvironmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.Environments
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := environments.ParseEnvironmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func convertEnvironmentSkuName(skuName string) environments.Sku {
	return environments.Sku{
		Name: environments.SkuName(skuName),
		// Gen2 cannot set capacity manually but SDK requires capacity
		Capacity: int64(1),
	}
}

func expandStorage(input []interface{}) environments.Gen2StorageConfigurationInput {
	storageMap := input[0].(map[string]interface{})
	return environments.Gen2StorageConfigurationInput{
		AccountName:   storageMap["name"].(string),
		ManagementKey: storageMap["key"].(string),
	}
}

func expandIdProperties(input []interface{}) []environments.TimeSeriesIdProperty {
	result := make([]environments.TimeSeriesIdProperty, 0)
	for _, item := range input {
		result = append(result, environments.TimeSeriesIdProperty{
			Name: pointer.To(item.(string)),
			Type: pointer.To(environments.PropertyTypeString),
		})
	}
	return result
}

func flattenIdProperties(input []environments.TimeSeriesIdProperty) []string {
	output := make([]string, 0)

	for _, v := range input {
		if v.Name != nil {
			output = append(output, *v.Name)
		}
	}

	return output
}

func flattenIoTTimeSeriesGen2EnvironmentStorage(input environments.Gen2StorageConfigurationOutput, key string) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"key":  key,
			"name": input.AccountName,
		},
	}
}

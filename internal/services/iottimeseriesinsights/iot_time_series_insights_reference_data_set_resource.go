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
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/referencedatasets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIoTTimeSeriesInsightsReferenceDataSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIoTTimeSeriesInsightsReferenceDataSetCreateUpdate,
		Read:   resourceIoTTimeSeriesInsightsReferenceDataSetRead,
		Update: resourceIoTTimeSeriesInsightsReferenceDataSetCreateUpdate,
		Delete: resourceIoTTimeSeriesInsightsReferenceDataSetDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := referencedatasets.ParseReferenceDataSetID(id)
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
					regexp.MustCompile(`^[A-Za-z0-9]{3,63}`),
					"IoT Time Series Insights Reference Data Set name must contain only alphanumeric characters and be between 3 and 63 characters.",
				),
			},

			"time_series_insights_environment_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: environments.ValidateEnvironmentID,
			},

			"data_string_comparison_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(referencedatasets.DataStringComparisonBehaviorOrdinal),
				ValidateFunc: validation.StringInSlice([]string{
					string(referencedatasets.DataStringComparisonBehaviorOrdinal),
					string(referencedatasets.DataStringComparisonBehaviorOrdinalIgnoreCase),
				}, false),
			},

			"key_property": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(referencedatasets.ReferenceDataKeyPropertyTypeBool),
								string(referencedatasets.ReferenceDataKeyPropertyTypeDateTime),
								string(referencedatasets.ReferenceDataKeyPropertyTypeDouble),
								string(referencedatasets.ReferenceDataKeyPropertyTypeString),
							}, false),
						},
					},
				},
			},

			"location": commonschema.Location(),

			"tags": commonschema.Tags(),
		},
	}
}

func resourceIoTTimeSeriesInsightsReferenceDataSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSets
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	environmentId, err := environments.ParseEnvironmentID(d.Get("time_series_insights_environment_id").(string))
	if err != nil {
		return err
	}
	id := referencedatasets.NewReferenceDataSetID(environmentId.SubscriptionId, environmentId.ResourceGroupName, environmentId.EnvironmentName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_reference_data_set", id.ID())
		}
	}

	payload := referencedatasets.ReferenceDataSetCreateOrUpdateParameters{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: referencedatasets.ReferenceDataSetCreationProperties{
			DataStringComparisonBehavior: pointer.To(referencedatasets.DataStringComparisonBehavior(d.Get("data_string_comparison_behavior").(string))),
			KeyProperties:                expandIoTTimeSeriesInsightsReferenceDataSetKeyProperties(d.Get("key_property").(*pluginsdk.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIoTTimeSeriesInsightsReferenceDataSetRead(d, meta)
}

func resourceIoTTimeSeriesInsightsReferenceDataSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := referencedatasets.ParseReferenceDataSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ReferenceDataSetName)
	d.Set("time_series_insights_environment_id", environments.NewEnvironmentID(id.SubscriptionId, id.ResourceGroupName, id.EnvironmentName).ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			dataStringComparisonBehavior := ""
			if props.DataStringComparisonBehavior != nil {
				dataStringComparisonBehavior = string(*props.DataStringComparisonBehavior)
			}
			d.Set("data_string_comparison_behavior", dataStringComparisonBehavior)

			if err := d.Set("key_property", flattenIoTTimeSeriesInsightsReferenceDataSetKeyProperties(props.KeyProperties)); err != nil {
				return fmt.Errorf("setting `key_property`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceIoTTimeSeriesInsightsReferenceDataSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := referencedatasets.ParseReferenceDataSetID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandIoTTimeSeriesInsightsReferenceDataSetKeyProperties(input []interface{}) []referencedatasets.ReferenceDataSetKeyProperty {
	properties := make([]referencedatasets.ReferenceDataSetKeyProperty, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		attr := v.(map[string]interface{})

		properties = append(properties, referencedatasets.ReferenceDataSetKeyProperty{
			Type: pointer.To(referencedatasets.ReferenceDataKeyPropertyType(attr["type"].(string))),
			Name: utils.String(attr["name"].(string)),
		})
	}

	return properties
}

func flattenIoTTimeSeriesInsightsReferenceDataSetKeyProperties(input []referencedatasets.ReferenceDataSetKeyProperty) []interface{} {
	properties := make([]interface{}, 0)
	for _, property := range input {
		propertyName := ""
		if property.Name != nil {
			propertyName = *property.Name
		}

		propertyType := ""
		if property.Type != nil {
			propertyType = string(*property.Type)
		}

		properties = append(properties, map[string]interface{}{
			"name": propertyName,
			"type": propertyType,
		})
	}

	return properties
}

package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/timeseriesinsights/mgmt/2020-05-15/timeseriesinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iottimeseriesinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iottimeseriesinsights/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			_, err := parse.ReferenceDataSetID(id)
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
				ValidateFunc: validate.TimeSeriesInsightsEnvironmentID,
			},

			"data_string_comparison_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(timeseriesinsights.Ordinal),
				ValidateFunc: validation.StringInSlice([]string{
					string(timeseriesinsights.Ordinal),
					string(timeseriesinsights.OrdinalIgnoreCase),
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
								string(timeseriesinsights.ReferenceDataKeyPropertyTypeBool),
								string(timeseriesinsights.ReferenceDataKeyPropertyTypeDateTime),
								string(timeseriesinsights.ReferenceDataKeyPropertyTypeDouble),
								string(timeseriesinsights.ReferenceDataKeyPropertyTypeString),
							}, false),
						},
					},
				},
			},

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),
		},
	}
}

func resourceIoTTimeSeriesInsightsReferenceDataSetCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	environmentID := d.Get("time_series_insights_environment_id").(string)
	envId, err := parse.EnvironmentID(environmentID)
	if err != nil {
		return err
	}
	id := parse.NewReferenceDataSetID(envId.SubscriptionId, envId.ResourceGroup, envId.Name, d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_reference_data_set", id.ID())
		}
	}

	dataset := timeseriesinsights.ReferenceDataSetCreateOrUpdateParameters{
		Location: &location,
		Tags:     tags.Expand(t),
		ReferenceDataSetCreationProperties: &timeseriesinsights.ReferenceDataSetCreationProperties{
			DataStringComparisonBehavior: timeseriesinsights.DataStringComparisonBehavior(d.Get("data_string_comparison_behavior").(string)),
			KeyProperties:                expandIoTTimeSeriesInsightsReferenceDataSetKeyProperties(d.Get("key_property").(*pluginsdk.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.EnvironmentName, id.Name, dataset); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIoTTimeSeriesInsightsReferenceDataSetRead(d, meta)
}

func resourceIoTTimeSeriesInsightsReferenceDataSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ReferenceDataSetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil || resp.ID == nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving IoT Time Series Insights Reference Data Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("time_series_insights_environment_id", strings.Split(d.Id(), "/referenceDataSets")[0])
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.ReferenceDataSetResourceProperties; props != nil {
		d.Set("data_string_comparison_behavior", string(props.DataStringComparisonBehavior))
		if err := d.Set("key_property", flattenIoTTimeSeriesInsightsReferenceDataSetKeyProperties(props.KeyProperties)); err != nil {
			return fmt.Errorf("setting `key_property`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceIoTTimeSeriesInsightsReferenceDataSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.ReferenceDataSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ReferenceDataSetID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting IoT Time Series Insights Reference Data Set %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandIoTTimeSeriesInsightsReferenceDataSetKeyProperties(input []interface{}) *[]timeseriesinsights.ReferenceDataSetKeyProperty {
	properties := make([]timeseriesinsights.ReferenceDataSetKeyProperty, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		attr := v.(map[string]interface{})

		properties = append(properties, timeseriesinsights.ReferenceDataSetKeyProperty{
			Type: timeseriesinsights.ReferenceDataKeyPropertyType(attr["type"].(string)),
			Name: utils.String(attr["name"].(string)),
		})
	}

	return &properties
}

func flattenIoTTimeSeriesInsightsReferenceDataSetKeyProperties(input *[]timeseriesinsights.ReferenceDataSetKeyProperty) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	properties := make([]interface{}, 0)
	for _, property := range *input {
		attr := make(map[string]interface{})
		attr["type"] = string(property.Type)
		if name := property.Name; name != nil {
			attr["name"] = *property.Name
		}
		properties = append(properties, attr)
	}

	return properties
}

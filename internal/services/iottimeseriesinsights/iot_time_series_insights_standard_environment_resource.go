// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/timeseriesinsights/2020-05-15/environments"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIoTTimeSeriesInsightsStandardEnvironment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIoTTimeSeriesInsightsStandardEnvironmentCreateUpdate,
		Read:   resourceIoTTimeSeriesInsightsStandardEnvironmentRead,
		Update: resourceIoTTimeSeriesInsightsStandardEnvironmentCreateUpdate,
		Delete: resourceIoTTimeSeriesInsightsStandardEnvironmentDelete,
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
					"IoT Time Series Insights Standard Environment name must contain only word characters, periods, underscores, and parentheses.",
				),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"S1_1",
					"S1_2",
					"S1_3",
					"S1_4",
					"S1_5",
					"S1_6",
					"S1_7",
					"S1_8",
					"S1_9",
					"S1_10",
					"S2_1",
					"S2_2",
					"S2_3",
					"S2_4",
					"S2_5",
					"S2_6",
					"S2_7",
					"S2_8",
					"S2_9",
					"S2_10",
				}, false),
			},

			"data_retention_time": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ISO8601Duration,
			},

			"storage_limit_exceeded_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(environments.StorageLimitExceededBehaviorPurgeOldData),
				ValidateFunc: validation.StringInSlice([]string{
					string(environments.StorageLimitExceededBehaviorPurgeOldData),
					string(environments.StorageLimitExceededBehaviorPauseIngress),
				}, false),
			},

			"partition_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceIoTTimeSeriesInsightsStandardEnvironmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.Environments
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := environments.NewEnvironmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, environments.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_environment", id.ID())
		}
	}

	sku, err := expandEnvironmentSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku: %+v", err)
	}

	environment := environments.Gen1EnvironmentCreateOrUpdateParameters{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Sku:      *sku,
		Properties: environments.Gen1EnvironmentCreationProperties{
			StorageLimitExceededBehavior: pointer.To(environments.StorageLimitExceededBehavior(d.Get("storage_limit_exceeded_behavior").(string))),
			DataRetentionTime:            d.Get("data_retention_time").(string),
		},
	}

	if v, ok := d.GetOk("partition_key"); ok {
		environment.Properties.PartitionKeyProperties = &[]environments.TimeSeriesIdProperty{
			{
				Name: utils.String(v.(string)),
				Type: pointer.To(environments.PropertyTypeString),
			},
		}
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, environment); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIoTTimeSeriesInsightsStandardEnvironmentRead(d, meta)
}

func resourceIoTTimeSeriesInsightsStandardEnvironmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		environment, ok := (*model).(environments.Gen1EnvironmentResource)
		if !ok {
			return fmt.Errorf("retrieving %s: expected a Gen1EnvironmentResource but got: %+v", *id, *model)
		}

		d.Set("location", location.Normalize(environment.Location))
		d.Set("sku_name", flattenEnvironmentSkuName(environment.Sku))

		d.Set("data_retention_time", environment.Properties.DataRetentionTime)
		storageLimitExceededBehavior := ""
		if environment.Properties.StorageLimitExceededBehavior != nil {
			storageLimitExceededBehavior = string(*environment.Properties.StorageLimitExceededBehavior)
		}
		d.Set("storage_limit_exceeded_behavior", storageLimitExceededBehavior)

		partitionKey := ""
		if partition := environment.Properties.PartitionKeyProperties; partition != nil {
			for _, v := range *partition {
				if v.Name == nil {
					continue
				}

				partitionKey = *v.Name
			}
		}
		d.Set("partition_key", partitionKey)

		if err := tags.FlattenAndSet(d, environment.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	return nil
}

func resourceIoTTimeSeriesInsightsStandardEnvironmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

func expandEnvironmentSkuName(skuName string) (*environments.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("sku_name (%s) has the worng number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var name environments.SkuName
	switch parts[0] {
	case "S1":
		name = environments.SkuNameSOne
	case "S2":
		name = environments.SkuNameSTwo
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capacity %s to int", skuName, parts[2])
	}

	return &environments.Sku{
		Name:     name,
		Capacity: int64(capacity),
	}, nil
}

func flattenEnvironmentSkuName(input environments.Sku) string {
	return fmt.Sprintf("%s_%d", string(input.Name), input.Capacity)
}

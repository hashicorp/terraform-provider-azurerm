package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/timeseriesinsights/mgmt/2020-05-15/timeseriesinsights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iottimeseriesinsights/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
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
			_, err := parse.EnvironmentID(id)
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				ValidateFunc: azValidate.ISO8601Duration,
			},

			"storage_limit_exceeded_behavior": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(timeseriesinsights.PurgeOldData),
				ValidateFunc: validation.StringInSlice([]string{
					string(timeseriesinsights.PurgeOldData),
					string(timeseriesinsights.PauseIngress),
				}, false),
			},

			"partition_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceIoTTimeSeriesInsightsStandardEnvironmentCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewEnvironmentID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	sku, err := expandEnvironmentSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku: %+v", err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if existing.Value != nil {
			environment, ok := existing.Value.AsGen1EnvironmentResource()
			if !ok {
				return fmt.Errorf("exisiting resource was not a %s", id)
			}

			if environment.ID != nil && *environment.ID != "" {
				return tf.ImportAsExistsError("azurerm_iot_time_series_insights_environment", *environment.ID)
			}
		}
	}

	environment := timeseriesinsights.Gen1EnvironmentCreateOrUpdateParameters{
		Location: &location,
		Tags:     tags.Expand(t),
		Sku:      sku,
		Gen1EnvironmentCreationProperties: &timeseriesinsights.Gen1EnvironmentCreationProperties{
			StorageLimitExceededBehavior: timeseriesinsights.StorageLimitExceededBehavior(d.Get("storage_limit_exceeded_behavior").(string)),
			DataRetentionTime:            utils.String(d.Get("data_retention_time").(string)),
		},
	}

	if v, ok := d.GetOk("partition_key"); ok {
		partition := make([]timeseriesinsights.TimeSeriesIDProperty, 1)
		partition[0] = timeseriesinsights.TimeSeriesIDProperty{
			Name: utils.String(v.(string)),
			Type: timeseriesinsights.String,
		}
		environment.Gen1EnvironmentCreationProperties.PartitionKeyProperties = &partition
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, environment)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceIoTTimeSeriesInsightsStandardEnvironmentRead(d, meta)
}

func resourceIoTTimeSeriesInsightsStandardEnvironmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EnvironmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil || resp.Value == nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving IoT Time Series Insights Standard Environment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	environment, ok := resp.Value.AsGen1EnvironmentResource()
	if !ok {
		return fmt.Errorf("exisiting resource was not a standard IoT Time Series Insights Standard Environment %q (Resource Group %q)", id.Name, id.ResourceGroup)
	}

	d.Set("name", environment.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("sku_name", flattenEnvironmentSkuName(environment.Sku))
	if location := environment.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := environment.Gen1EnvironmentResourceProperties; props != nil {
		d.Set("storage_limit_exceeded_behavior", string(props.StorageLimitExceededBehavior))
		d.Set("data_retention_time", props.DataRetentionTime)

		if partition := props.PartitionKeyProperties; partition != nil && len(*partition) > 0 {
			for _, v := range *partition {
				d.Set("partition_key", v.Name)
			}
		}
	}

	return tags.FlattenAndSet(d, environment.Tags)
}

func resourceIoTTimeSeriesInsightsStandardEnvironmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.EnvironmentID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting IoT Time Series Insights Standard Environment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandEnvironmentSkuName(skuName string) (*timeseriesinsights.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("sku_name (%s) has the worng number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var name timeseriesinsights.SkuName
	switch parts[0] {
	case "S1":
		name = timeseriesinsights.S1
	case "S2":
		name = timeseriesinsights.S2
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &timeseriesinsights.Sku{
		Name:     name,
		Capacity: utils.Int32(int32(capacity)),
	}, nil
}

func flattenEnvironmentSkuName(input *timeseriesinsights.Sku) string {
	if input == nil || input.Capacity == nil {
		return ""
	}

	return fmt.Sprintf("%s_%d", string(input.Name), *input.Capacity)
}

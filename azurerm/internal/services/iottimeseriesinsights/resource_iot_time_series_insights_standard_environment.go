package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/timeseriesinsights/mgmt/2018-08-15-preview/timeseriesinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	azValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIoTTimeSeriesInsightsStandardEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIoTTimeSeriesInsightsStandardEnvironmentCreateUpdate,
		Read:   resourceArmIoTTimeSeriesInsightsStandardEnvironmentRead,
		Update: resourceArmIoTTimeSeriesInsightsStandardEnvironmentCreateUpdate,
		Delete: resourceArmIoTTimeSeriesInsightsStandardEnvironmentDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TimeSeriesInsightsEnvironmentID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeString,
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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.ISO8601Duration,
			},

			"storage_limit_exceeded_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(timeseriesinsights.PurgeOldData),
				ValidateFunc: validation.StringInSlice([]string{
					string(timeseriesinsights.PurgeOldData),
					string(timeseriesinsights.PauseIngress),
				}, false),
			},

			"partition_key": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmIoTTimeSeriesInsightsStandardEnvironmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})
	sku, err := expandEnvironmentSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku: %+v", err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IoT Time Series Insights Standard Environment %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.Value != nil {
			environment, ok := existing.Value.AsStandardEnvironmentResource()
			if !ok {
				return fmt.Errorf("exisiting resource was not a standard IoT Time Series Insights Standard Environment %q (Resource Group %q)", name, resourceGroup)
			}

			if environment.ID != nil && *environment.ID != "" {
				return tf.ImportAsExistsError("azurerm_iot_time_series_insights_environment", *environment.ID)
			}
		}
	}

	environment := timeseriesinsights.StandardEnvironmentCreateOrUpdateParameters{
		Location: &location,
		Tags:     tags.Expand(t),
		Sku:      sku,
		StandardEnvironmentCreationProperties: &timeseriesinsights.StandardEnvironmentCreationProperties{
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
		environment.StandardEnvironmentCreationProperties.PartitionKeyProperties = &partition
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, environment)
	if err != nil {
		return fmt.Errorf("creating/updating IoT Time Series Insights Standard Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of IoT Time Series Insights Standard Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("retrieving IoT Time Series Insights Standard Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resource, ok := resp.Value.AsStandardEnvironmentResource()
	if !ok {
		return fmt.Errorf("resource was not a standard IoT Time Series Insights Standard Environment %q (Resource Group %q)", name, resourceGroup)
	}

	if resource.ID == nil || *resource.ID == "" {
		return fmt.Errorf("cannot read IoT Time Series Insights Standard Environment %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resource.ID)

	return resourceArmIoTTimeSeriesInsightsStandardEnvironmentRead(d, meta)
}

func resourceArmIoTTimeSeriesInsightsStandardEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TimeSeriesInsightsEnvironmentID(d.Id())
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

	environment, ok := resp.Value.AsStandardEnvironmentResource()
	if !ok {
		return fmt.Errorf("exisiting resource was not a standard IoT Time Series Insights Standard Environment %q (Resource Group %q)", id.Name, id.ResourceGroup)
	}

	d.Set("name", environment.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("sku_name", flattenEnvironmentSkuName(environment.Sku))
	if location := environment.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := environment.StandardEnvironmentResourceProperties; props != nil {
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

func resourceArmIoTTimeSeriesInsightsStandardEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TimeSeriesInsightsEnvironmentID(d.Id())
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

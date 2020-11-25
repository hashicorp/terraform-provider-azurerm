package iottimeseriesinsights

import (
	"fmt"
	"regexp"
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

func resourceArmIoTTimeSeriesInsightsGen2Environment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIoTTimeSeriesInsightsGen2EnvironmentCreateUpdate,
		Read:   resourceArmIoTTimeSeriesInsightsGen2EnvironmentRead,
		Update: resourceArmIoTTimeSeriesInsightsGen2EnvironmentCreateUpdate,
		Delete: resourceArmIoTTimeSeriesInsightsGen2EnvironmentDelete,
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
					"IoT Time Series Insights Gen2 Environment name must contain only word characters, periods, underscores, and parentheses.",
				),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"L1",
				}, false),
			},

			"data_retention_time": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.ISO8601Duration,
			},
			"property": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							Required: true,
						},
					},
				},
			},
			"storage": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmIoTTimeSeriesInsightsGen2EnvironmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.EnvironmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})
	sku, err := convertEnvironmentSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku: %+v", err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IoT Time Series Insights Gen2 Environment %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.Value != nil {
			environment, ok := existing.Value.AsLongTermEnvironmentResource()
			if !ok {
				return fmt.Errorf("exisiting resource was not IoT Time Series Insights Gen2 Environment %q (Resource Group %q)", name, resourceGroup)
			}

			if environment.ID != nil && *environment.ID != "" {
				return tf.ImportAsExistsError("azurerm_iot_time_series_insights_environment", *environment.ID)
			}
		}
	}

	environment := timeseriesinsights.LongTermEnvironmentCreateOrUpdateParameters{
		Location: &location,
		Tags:     tags.Expand(t),
		Sku:      sku,
		LongTermEnvironmentCreationProperties: &timeseriesinsights.LongTermEnvironmentCreationProperties{
			TimeSeriesIDProperties: expandIdProperties(d.Get("property").([]interface{})),
			WarmStoreConfiguration: &timeseriesinsights.WarmStoreConfigurationProperties{
				DataRetention: utils.String(d.Get("data_retention_time").(string)),
			},
			StorageConfiguration: expandStorage(d.Get("storage").([]interface{})),
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, environment)
	if err != nil {
		return fmt.Errorf("creating/updating IoT Time Series Gen2 Standard Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of IoT Time Series Insights Gen2 Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("retrieving IoT Time Series Insights Gen2 Environment %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resource, ok := resp.Value.AsLongTermEnvironmentResource()
	if !ok {
		return fmt.Errorf("resource was not IoT Time Series Insights Gen2 Environment %q (Resource Group %q)", name, resourceGroup)
	}

	if resource.ID == nil || *resource.ID == "" {
		return fmt.Errorf("cannot read IoT Time Series Insights Gen2 Environment %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resource.ID)

	return resourceArmIoTTimeSeriesInsightsGen2EnvironmentRead(d, meta)
}

func resourceArmIoTTimeSeriesInsightsGen2EnvironmentRead(d *schema.ResourceData, meta interface{}) error {
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

	environment, ok := resp.Value.AsLongTermEnvironmentResource()
	if !ok {
		return fmt.Errorf("exisiting resource was not a standard IoT Time Series Insights Standard Environment %q (Resource Group %q)", id.Name, id.ResourceGroup)
	}

	d.Set("name", environment.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("sku_name", environment.Sku.Name)
	if location := environment.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	return tags.FlattenAndSet(d, environment.Tags)
}

func resourceArmIoTTimeSeriesInsightsGen2EnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("deleting IoT Time Series Insights Gen2 Environment %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func convertEnvironmentSkuName(skuName string) (*timeseriesinsights.Sku, error) {
	var name timeseriesinsights.SkuName
	switch skuName {
	case "L1":
		name = timeseriesinsights.L1
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, skuName)
	}

	//Gen2 cannot set capacity manually but SDK requires capacity
	capacity := utils.Int32(1)

	return &timeseriesinsights.Sku{
		Name:     name,
		Capacity: capacity,
	}, nil
}

func expandStorage(input []interface{}) *timeseriesinsights.LongTermStorageConfigurationInput {
	storageMap := input[0].(map[string]interface{})
	accountName := storageMap["name"].(string)
	managementKey := storageMap["key"].(string)

	return &timeseriesinsights.LongTermStorageConfigurationInput{
		AccountName:   &accountName,
		ManagementKey: &managementKey,
	}
}

func expandIdProperties(input []interface{}) *[]timeseriesinsights.TimeSeriesIDProperty {
	result := make([]timeseriesinsights.TimeSeriesIDProperty, 0)
	propertiesMap := input[0].(map[string]interface{})["ids"].([]interface{})
	for _, item := range propertiesMap {
		result = append(result, timeseriesinsights.TimeSeriesIDProperty{
			Name: utils.String(item.(string)),
			Type: "String",
		})
	}
	return &result
}

package costmanagement

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2019-10-01/costmanagement"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/costmanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/costmanagement/validate"
	resourceValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceCostManagementExportResourceGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceCostManagementExportResourceGroupCreateUpdate,
		Read:   resourceCostManagementExportResourceGroupRead,
		Update: resourceCostManagementExportResourceGroupCreateUpdate,
		Delete: resourceCostManagementExportResourceGroupDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CostManagementExportResourceGroupID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ExportName,
			},

			"resource_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceValidate.ResourceGroupID,
			},

			"active": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"recurrence_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(costmanagement.RecurrenceTypeDaily),
					string(costmanagement.RecurrenceTypeWeekly),
					string(costmanagement.RecurrenceTypeMonthly),
					string(costmanagement.RecurrenceTypeAnnually),
				}, false),
			},

			"recurrence_period_start": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"recurrence_period_end": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"delivery_info": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"storage_account_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"container_name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.ExportContainerName,
						},
						"root_folder_path": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"query": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"time_frame": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(costmanagement.Custom),
								string(costmanagement.MonthToDate),
								string(costmanagement.TheLastMonth),
								string(costmanagement.TheLastWeek),
								string(costmanagement.TheLastYear),
								string(costmanagement.WeekToDate),
								string(costmanagement.MonthToDate),
							}, false),
						},
					},
				},
			},
		},
	}
}

func resourceCostManagementExportResourceGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CostManagement.ExportClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_id").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Cost Management Export Resource Group %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_cost_management_export_resource_group", *existing.ID)
		}
	}

	from, _ := time.Parse(time.RFC3339, d.Get("recurrence_period_start").(string))
	to, _ := time.Parse(time.RFC3339, d.Get("recurrence_period_end").(string))

	status := costmanagement.Active
	if v := d.Get("active"); !v.(bool) {
		status = costmanagement.Inactive
	}

	properties := &costmanagement.ExportProperties{
		Schedule: &costmanagement.ExportSchedule{
			Recurrence: costmanagement.RecurrenceType(d.Get("recurrence_type").(string)),
			RecurrencePeriod: &costmanagement.ExportRecurrencePeriod{
				From: &date.Time{Time: from},
				To:   &date.Time{Time: to},
			},
			Status: status,
		},
		DeliveryInfo: expandExportDeliveryInfo(d.Get("delivery_info").([]interface{})),
		Format:       costmanagement.Csv,
		Definition:   expandExportQuery(d.Get("query").([]interface{})),
	}

	account := costmanagement.Export{
		ExportProperties: properties,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, account); err != nil {
		return fmt.Errorf("creating/updating Cost Management Export Resource Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Cost Management Export Resource Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("cannot read Cost Management Export Resource Group %q (Resource Group %q) ID", name, resourceGroup)
	}

	id := *resp.ID
	// The ID is missing the prefix `/` which causes our uri parse to fail
	if !strings.HasPrefix(id, "/") {
		id = fmt.Sprintf("/%s", id)
	}

	d.SetId(id)

	return resourceCostManagementExportResourceGroupRead(d, meta)
}

func resourceCostManagementExportResourceGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CostManagement.ExportClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CostManagementExportResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceId, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Cost Management Export Resource Group %q (Resource Group %q): %+v", id.Name, id.ResourceId, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_id", id.ResourceId)

	if schedule := resp.Schedule; schedule != nil {
		if recurrencePeriod := schedule.RecurrencePeriod; recurrencePeriod != nil {
			d.Set("recurrence_period_start", recurrencePeriod.From.Format(time.RFC3339))
			d.Set("recurrence_period_end", recurrencePeriod.To.Format(time.RFC3339))
		}
		status := false
		if schedule.Status == costmanagement.Active {
			status = true
		}
		d.Set("active", status)
		d.Set("recurrence_type", schedule.Recurrence)
	}
	if err := d.Set("delivery_info", flattenExportDeliveryInfo(resp.DeliveryInfo)); err != nil {
		return fmt.Errorf("setting `delivery_info`: %+v", err)
	}

	if err := d.Set("query", flattenExportQuery(resp.Definition)); err != nil {
		return fmt.Errorf("setting `query`: %+v", err)
	}

	return nil
}

func resourceCostManagementExportResourceGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CostManagement.ExportClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CostManagementExportResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceId, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting Cost Management Export Resource Group %q (Resource Group %q): %+v", id.Name, id.ResourceId, err)
		}
	}

	return nil
}

func expandExportDeliveryInfo(input []interface{}) *costmanagement.ExportDeliveryInfo {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	deliveryInfo := &costmanagement.ExportDeliveryInfo{
		Destination: &costmanagement.ExportDeliveryDestination{
			ResourceID:     utils.String(attrs["storage_account_id"].(string)),
			Container:      utils.String(attrs["container_name"].(string)),
			RootFolderPath: utils.String(attrs["root_folder_path"].(string)),
		},
	}

	return deliveryInfo
}

func flattenExportDeliveryInfo(input *costmanagement.ExportDeliveryInfo) []interface{} {
	if input == nil || input.Destination == nil {
		return []interface{}{}
	}

	destination := input.Destination
	attrs := make(map[string]interface{})
	if resourceID := destination.ResourceID; resourceID != nil {
		attrs["storage_account_id"] = *resourceID
	}
	if containerName := destination.Container; containerName != nil {
		attrs["container_name"] = *containerName
	}
	if rootFolderPath := destination.RootFolderPath; rootFolderPath != nil {
		attrs["root_folder_path"] = *rootFolderPath
	}

	return []interface{}{attrs}
}

func expandExportQuery(input []interface{}) *costmanagement.QueryDefinition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attrs := input[0].(map[string]interface{})
	definitionInfo := &costmanagement.QueryDefinition{
		Type:      utils.String(attrs["type"].(string)),
		Timeframe: costmanagement.TimeframeType(attrs["time_frame"].(string)),
	}

	return definitionInfo
}

func flattenExportQuery(input *costmanagement.QueryDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	attrs := make(map[string]interface{})
	if queryType := input.Type; queryType != nil {
		attrs["type"] = *queryType
	}
	attrs["time_frame"] = string(input.Timeframe)

	return []interface{}{attrs}
}

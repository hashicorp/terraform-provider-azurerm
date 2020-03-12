package costmanagement

import (
	"fmt"
	"github.com/Azure/go-autorest/autorest/date"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/costmanagement/mgmt/2019-10-01/costmanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/costmanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmCostManagementExportResourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmCostManagementExportResourceGroupCreate,
		Read:   resourceArmCostManagementExportResourceGroupRead,
		Delete: resourceArmCostManagementExportResourceGroupDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.CostManagementExportResourceGroupID(id)
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
				ValidateFunc: validateCostManagementExportName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"recurrence_type": {
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(costmanagement.RecurrenceTypeDaily),
					string(costmanagement.RecurrenceTypeWeekly),
					string(costmanagement.RecurrenceTypeMonthly),
					string(costmanagement.RecurrenceTypeAnnually),
				},false),
			},

			"recurrence_period_start": {
				Type: schema.TypeString,
				Required: true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"recurrence_period_end": {
				Type: schema.TypeString,
				Optional: true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"tags": tags.ForceNewSchema(),
		},
	}
}

func resourceArmCostManagementExportResourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CostManagement.ExportClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
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

	t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema
	props.Recurrence.EndTime = &date.Time{Time: t}

	recurrencePeriod := costmanagement.ExportRecurrencePeriod{
		From: time.Parse(time.RFC3339, d.Get("recurrence_period_from").(string)),
		To:   nil,
	}

	schedule := &costmanagement.ExportSchedule{
		Recurrence: costmanagement.RecurrenceType(d.Get("recurrence_type").(string)),
	}

	properties := &costmanagement.ExportProperties{
		Schedule: schedule,
	}

	account := costmanagement.Export{
		Tags: tags.Expand(t),
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

	d.SetId(*resp.ID)

	return resourceArmCostManagementExportResourceGroupRead(d, meta)
}

func resourceArmCostManagementExportResourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CostManagement.ExportClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CostManagementExportResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Cost Management Export Resource Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmCostManagementExportResourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).CostManagement.ExportClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CostManagementExportResourceGroupID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting Cost Management Export Resource Group %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func validateCostManagementExportName(v interface{}, k string) (warnings []string, errors []error) {
	name := v.(string)

	if regexp.MustCompile(`^[\s]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q must not consist of whitespace", k))
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(name) {
		errors = append(errors, fmt.Errorf("%q may only contain letters and digits: %q", k, name))
	}

	if len(name) < 3 || len(name) > 24 {
		errors = append(errors, fmt.Errorf("%q must be (inclusive) between 3 and 24 characters long but is %d", k, len(name)))
	}

	return warnings, errors
}
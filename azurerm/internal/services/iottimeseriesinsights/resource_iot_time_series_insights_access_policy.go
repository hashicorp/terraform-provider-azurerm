package iottimeseriesinsights

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/timeseriesinsights/mgmt/2018-08-15-preview/timeseriesinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmIoTTimeSeriesInsightsAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmIoTTimeSeriesInsightsAccessPolicyCreateUpdate,
		Read:   resourceArmIoTTimeSeriesInsightsAccessPolicyRead,
		Update: resourceArmIoTTimeSeriesInsightsAccessPolicyCreateUpdate,
		Delete: resourceArmIoTTimeSeriesInsightsAccessPolicyDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.TimeSeriesInsightsAccessPolicyID(id)
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
					"IoT Time Series Insights Access Policy name must contain only word characters, periods, underscores, and parentheses.",
				),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"environment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w\._\(\)]+$`),
					"IoT Time Series Insights Environment name must contain only word characters, periods, underscores, and parentheses.",
				),
			},

			"principal_object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"roles": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(timeseriesinsights.Contributor),
						string(timeseriesinsights.Reader),
					}, false),
				},
			},
		},
	}
}

func resourceArmIoTTimeSeriesInsightsAccessPolicyCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.AccessPoliciesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	environmentName := d.Get("environment_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, environmentName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IoT Time Series Insights Access Policy %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_access_policy", *existing.ID)
		}
	}

	policy := timeseriesinsights.AccessPolicyCreateOrUpdateParameters{
		&timeseriesinsights.AccessPolicyResourceProperties{
			Description:       utils.String(d.Get("description").(string)),
			PrincipalObjectID: utils.String(d.Get("principal_object_id").(string)),
			Roles:             expandIoTTimeSeriesInsightsAccessPolicyRoles(d.Get("roles").(*schema.Set).List()),
		},
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, environmentName, name, policy)
	if err != nil {
		return fmt.Errorf("creating/updating IoT Time Series Insights Access Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, environmentName, name)
	if err != nil {
		return fmt.Errorf("retrieving IoT Time Series Insights Access Policy %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read IoT Time Series Insights Access Policy %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceArmIoTTimeSeriesInsightsAccessPolicyRead(d, meta)
}

func resourceArmIoTTimeSeriesInsightsAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.AccessPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TimeSeriesInsightsAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving IoT Time Series Insights Access Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("environment_name", id.EnvironmentName)

	if props := resp.AccessPolicyResourceProperties; props != nil {
		d.Set("description", props.Description)
		d.Set("principal_object_id", props.PrincipalObjectID)
		d.Set("roles", flattenIoTTimeSeriesInsightsAccessPolicyRoles(resp.Roles))
	}

	return nil
}

func resourceArmIoTTimeSeriesInsightsAccessPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.AccessPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.TimeSeriesInsightsAccessPolicyID(d.Id())
	if err != nil {
		return err
	}

	response, err := client.Delete(ctx, id.ResourceGroup, id.EnvironmentName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(response) {
			return fmt.Errorf("deleting IoT Time Series Insights Access Policy %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func expandIoTTimeSeriesInsightsAccessPolicyRoles(input []interface{}) *[]timeseriesinsights.AccessPolicyRole {
	roles := make([]timeseriesinsights.AccessPolicyRole, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		roles = append(roles, timeseriesinsights.AccessPolicyRole(v.(string)))
	}

	return &roles
}

func flattenIoTTimeSeriesInsightsAccessPolicyRoles(input *[]timeseriesinsights.AccessPolicyRole) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, string(item))
		}
	}
	return result
}

package iottimeseriesinsights

import (
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/timeseriesinsights/mgmt/2020-05-15/timeseriesinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/migration"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iottimeseriesinsights/validate"
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
			_, err := parse.AccessPolicyID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			migration.TimeSeriesInsightsAccessPolicyV0(),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[-\w\._\(\)]+$`),
					"IoT Time Series Insights Access Policy name must contain only word characters, periods, underscores, hyphens, and parentheses.",
				),
			},

			"time_series_insights_environment_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.TimeSeriesInsightsEnvironmentID,
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	environmentId, err := parse.TimeSeriesInsightsEnvironmentID(d.Get("time_series_insights_environment_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewAccessPolicyID(environmentId.ResourceGroup, environmentId.Name, name)
	if d.IsNewResource() {
		existing, err := client.Get(ctx, environmentId.ResourceGroup, environmentId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing IoT Time Series Insights Access Policy %q (Resource Group %q): %s", name, environmentId.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_iot_time_series_insights_access_policy", id.ID(subscriptionId))
		}
	}

	policy := timeseriesinsights.AccessPolicyCreateOrUpdateParameters{
		AccessPolicyResourceProperties: &timeseriesinsights.AccessPolicyResourceProperties{
			Description:       utils.String(d.Get("description").(string)),
			PrincipalObjectID: utils.String(d.Get("principal_object_id").(string)),
			Roles:             expandIoTTimeSeriesInsightsAccessPolicyRoles(d.Get("roles").(*schema.Set).List()),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, environmentId.ResourceGroup, environmentId.Name, name, policy); err != nil {
		return fmt.Errorf("creating/updating IoT Time Series Insights Access Policy %q (Resource Group %q): %+v", name, environmentId.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, environmentId.ResourceGroup, environmentId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving IoT Time Series Insights Access Policy %q (Resource Group %q): %+v", name, environmentId.ResourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("cannot read IoT Time Series Insights Access Policy %q (Resource Group %q) ID", name, environmentId.ResourceGroup)
	}

	d.SetId(id.ID(subscriptionId))

	return resourceArmIoTTimeSeriesInsightsAccessPolicyRead(d, meta)
}

func resourceArmIoTTimeSeriesInsightsAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).IoTTimeSeriesInsights.AccessPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AccessPolicyID(d.Id())
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

	environmentId := parse.NewTimeSeriesInsightsEnvironmentID(id.ResourceGroup, id.EnvironmentName)

	d.Set("name", resp.Name)
	d.Set("time_series_insights_environment_id", environmentId.ID(subscriptionId))

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

	id, err := parse.AccessPolicyID(d.Id())
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

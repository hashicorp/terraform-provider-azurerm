package portal

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/portal/mgmt/2019-01-01-preview/portal"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDashboardCreateUpdate,
		Read:   resourceArmDashboardRead,
		Update: resourceArmDashboardCreateUpdate,
		Delete: resourceArmDashboardDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateDashboardName,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"location":            azure.SchemaLocation(),
			"tags":                tags.Schema(),
			"dashboard_properties": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				StateFunc: utils.NormalizeJson,
			},
		},
	}
}

func resourceArmDashboardCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	t := d.Get("tags").(map[string]interface{})
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	dashboardProps := d.Get("dashboard_properties").(string)

	dashboard := portal.Dashboard{
		Location: &location,
		Tags:     tags.Expand(t),
	}

	var dashboardProperties portal.DashboardProperties

	if err := json.Unmarshal([]byte(dashboardProps), &dashboardProperties); err != nil {
		return fmt.Errorf("Error parsing JSON: %+v", err)
	}
	dashboard.DashboardProperties = &dashboardProperties

	_, err := client.CreateOrUpdate(ctx, resourceGroup, name, dashboard)
	if err != nil {
		return fmt.Errorf("Error creating/updating Dashboard %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	// get it back again to set the props
	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making Read request for Dashboard %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmDashboardRead(d, meta)
}

func resourceArmDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, parseErr := azure.ParseAzureResourceID(d.Id())
	if parseErr != nil {
		return parseErr
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["dashboards"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request for Dashboard %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*resp.Location))
	}

	props, jsonErr := json.Marshal(resp.DashboardProperties)
	if jsonErr != nil {
		return fmt.Errorf("Error parsing DashboardProperties JSON: %+v", jsonErr)
	}
	d.Set("dashboard_properties", string(props))

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Portal.DashboardsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, parseErr := azure.ParseAzureResourceID(d.Id())
	if parseErr != nil {
		return parseErr
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["dashboards"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error retrieving Key Vault %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func validateDashboardName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)

	if len(value) > 64 {
		errors = append(errors, fmt.Errorf("%q may not exceed 64 characters in length", k))
	}

	// only alpanumeric and hyphens
	if matched := regexp.MustCompile(`^[-\w]+$`).Match([]byte(value)); !matched {
		errors = append(errors, fmt.Errorf("%q may only contain alphanumeric and hyphen characters", k))
	}

	return warnings, errors
}

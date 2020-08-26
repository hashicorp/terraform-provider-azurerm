package loganalytics

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsSavedSearch() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsSavedSearchCreateUpdate,
		Read:   resourceArmLogAnalyticsSavedSearchRead,
		Update: resourceArmLogAnalyticsSavedSearchCreateUpdate,
		Delete: resourceArmLogAnalyticsSavedSearchDelete,

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
			"log_analytics_workspace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"category": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"function_alias": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"function_parameters": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
				// TODO Test less than 2 and waaaay higher than 2
				// ValidateFunc: validation.IntAtLeast(2),
			},
		},
	}
}

func resourceArmLogAnalyticsSavedSearchCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Saved Search creation.")

	name := d.Get("name").(string)
	workspaceID := d.Get("log_analytics_workspace_id").(string)
	id, err := parse.LogAnalyticsWorkspaceID(workspaceID)
	if err != nil {
		return err
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Log Analytics Saved Search %q (WorkSpace %q / Resource Group %q): %s", name, id.Name, id.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_saved_search", *existing.ID)
		}
	}

	parameters := operationalinsights.SavedSearch{
		SavedSearchProperties: &operationalinsights.SavedSearchProperties{
			Category:           utils.String(d.Get("category").(string)),
			DisplayName:        utils.String(d.Get("display_name").(string)),
			Query:              utils.String(d.Get("query").(string)),
			FunctionAlias:      utils.String(d.Get("function_alias").(string)),
			FunctionParameters: utils.String(d.Get("function_parameters").(string)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, name, parameters); err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name, name)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read Log Analytics Saved Search %q (WorkSpace %q / Resource Group %q): %s", name, id.Name, id.ResourceGroup, err)
	}

	d.SetId(*read.ID)

	return resourceArmLogAnalyticsSavedSearchRead(d, meta)
}

func resourceArmLogAnalyticsSavedSearchRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	workspaceClient := meta.(*clients.Client).LogAnalytics.WorkspacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	workspace, err := workspaceClient.Get(ctx, id.ResourceGroup, id.WorkspaceName)
	if err != nil {
		if utils.ResponseWasNotFound(workspace.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf(" making Read request on AzureRM Log Analytics workspaces %q: %+v", id.Name, err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics Saved Search %q (WorkSpace %q / Resource Group %q): %s", id.Name, id.WorkspaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("log_analytics_workspace_id", workspace.ID)

	if props := resp.SavedSearchProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("category", props.Category)
		d.Set("query", props.Query)
		d.Set("function_alias", props.FunctionAlias)
		d.Set("function_parameters", props.FunctionParameters)
		d.Set("version", props.Version)
	}

	return nil
}

func resourceArmLogAnalyticsSavedSearchDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.SavedSearchesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := parse.LogAnalyticsSavedSearchID(d.Id())
	if err != nil {
		return err
	}

	_, err = client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return fmt.Errorf("issuing AzureRM delete request for Log Analytics Saved Search %q (WorkSpace %q / Resource Group %q): %s", id.Name, id.WorkspaceName, id.ResourceGroup, err)
	}

	return nil
}

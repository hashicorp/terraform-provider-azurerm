package loganalytics

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsLinkedService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsLinkedServiceCreateUpdate,
		Read:   resourceArmLogAnalyticsLinkedServiceRead,
		Update: resourceArmLogAnalyticsLinkedServiceCreateUpdate,
		Delete: resourceArmLogAnalyticsLinkedServiceDelete,

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
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"workspace_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
			},

			"linked_service_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Automation",
				ValidateFunc: validation.StringInSlice([]string{
					"Automation",
					"Cluster",
				}, false),
			},

			"read_access_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				ExactlyOneOf: []string{"read_access_id", "write_access_id"},
			},

			"write_access_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				ExactlyOneOf: []string{"read_access_id", "write_access_id"},
			},

			// Exported properties
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmLogAnalyticsLinkedServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Linked Services creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	workspaceId := d.Get("workspace_id").(string)
	serviceType := d.Get("linked_service_type").(string)
	t := d.Get("tags").(map[string]interface{})

	workspace, err := parse.LogAnalyticsWorkspaceID(workspaceId)
	if err != nil {
		return fmt.Errorf("Linked Service %q (Workspace %q / Resource Group %q): %+v", serviceType, workspace.Name, resourceGroup, err)
	}

	id := parse.NewLogAnalyticsLinkedServiceID(resourceGroup, serviceType, workspace.Name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, workspace.Name, serviceType)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Linked Service %q (Workspace %q / Resource Group %q): %+v", serviceType, workspace.Name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_log_analytics_linked_service", *existing.ID)
		}
	}

	parameters := operationalinsights.LinkedService{
		LinkedServiceProperties: &operationalinsights.LinkedServiceProperties{},
		Tags:                    tags.Expand(t),
	}

	if d.Get("read_access_id") != "" {
		parameters.LinkedServiceProperties.ResourceID = utils.String(d.Get("read_access_id").(string))
	}

	if d.Get("write_access_id") != "" {
		parameters.LinkedServiceProperties.WriteAccessResourceID = utils.String(d.Get("write_access_id").(string))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, workspace.Name, serviceType, parameters)
	if err != nil {
		return fmt.Errorf("creating Linked Service %q (Workspace %q / Resource Group %q): %+v", serviceType, workspace.Name, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Linked Service %q (Workspace %q / Resource Group %q): %+v", serviceType, workspace.Name, resourceGroup, err)
	}

	_, err = client.Get(ctx, resourceGroup, workspace.Name, serviceType)
	if err != nil {
		return fmt.Errorf("retrieving Linked Service %q (Workspace %q / Resource Group %q): %+v", serviceType, workspace.Name, resourceGroup, err)
	}

	d.SetId(id.ID(subscriptionId))

	return resourceArmLogAnalyticsLinkedServiceRead(d, meta)
}

func resourceArmLogAnalyticsLinkedServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	workspaceName := id.Path["workspaces"]
	serviceType := id.Path["linkedServices"]
	workspace := parse.NewLogAnalyticsWorkspaceID(workspaceName, resourceGroup)

	resp, err := client.Get(ctx, resourceGroup, workspaceName, serviceType)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics Linked Service '%s': %+v", serviceType, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("workspace_id", workspace.ID(subscriptionId))
	d.Set("linked_service_type", serviceType)

	if props := resp.LinkedServiceProperties; props != nil {
		d.Set("read_access_id", props.ResourceID)
		d.Set("write_access_id", props.WriteAccessResourceID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmLogAnalyticsLinkedServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	workspaceName := id.Path["workspaces"]
	serviceType := id.Path["linkedServices"]

	future, err := client.Delete(ctx, resourceGroup, workspaceName, serviceType)
	if err != nil {
		return fmt.Errorf("deleting Log Analytics Linked Service %q (Workspace %q / Resource Group %q): %+v", serviceType, workspaceName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Log Analytics Linked Service %q (Workspace %q / Resource Group %q): %+v", serviceType, workspaceName, resourceGroup, err)
		}
	}

	return nil
}

package loganalytics

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
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

			// TODO: Add Force New if changed
			"workspace_name": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.LogAnalyticsWorkspaceName,
				ConflictsWith:    []string{"workspace_id"},
				Deprecated:       "This field has been deprecated in favour of `workspace_id` and will be removed in a future version of the provider",
			},

			// TODO: Add Force New if changed
			"workspace_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
				ConflictsWith:    []string{"workspace_name"},
			},

			// TODO: Add Defualt value to Automation if empty
			"linked_service_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"automation",
					"cluster",
				}, false),
				Deprecated: "This field has been deprecated and will be removed in a future version of the provider",
			},

			"resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ConflictsWith: []string{"read_access_id"},
				Deprecated:    "This field has been deprecated in favour of `read_access_id` and will be removed in a future version of the provider",
			},

			"read_access_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  azure.ValidateResourceID,
				ExactlyOneOf:  []string{"read_access_id", "write_access_id", "resource_id"},
				ConflictsWith: []string{"resource_id"},
			},

			"write_access_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				ExactlyOneOf: []string{"read_access_id", "write_access_id", "resource_id"},
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

	// TODO: Add legacy attributes here
	// Convert workspace name to workspace id
	// Convert linked_service_name to id.LinkedServiceName which is actually linked service type
	// Convert resource_id to read_access_id

	resourceGroup := d.Get("resource_group_name").(string)
	workspaceId := d.Get("workspace_id").(string)
	readAccess := d.Get("read_access_id").(string)
	writeAccess := d.Get("write_access_id").(string)
	t := d.Get("tags").(map[string]interface{})

	workspace, err := parse.LogAnalyticsWorkspaceID(workspaceId)
	if err != nil {
		return fmt.Errorf("Linked Service (Resource Group %q) unable to parse workspace id: %+v", resourceGroup, err)
	}

	id := parse.NewLogAnalyticsLinkedServiceID(subscriptionId, resourceGroup, workspace.WorkspaceName, LogAnalyticsLinkedServiceType(readAccess))

	if strings.EqualFold(id.LinkedServiceName, "Cluster") && writeAccess == "" {
		return fmt.Errorf("Linked Service '%s/%s' (Resource Group %q): A linked Log Analytics Cluster requires the 'write_access_id' attribute to be set", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup)
	}

	if strings.EqualFold(id.LinkedServiceName, "Automation") && readAccess == "" {
		return fmt.Errorf("Linked Service '%s/%s' (Resource Group %q): A linked Automation Account requires the 'read_access_id' attribute to be set", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, workspace.WorkspaceName, id.LinkedServiceName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
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

	if id.LinkedServiceName == "Automation" {
		parameters.LinkedServiceProperties.ResourceID = utils.String(readAccess)
	}

	if id.LinkedServiceName == "Cluster" {
		parameters.LinkedServiceProperties.WriteAccessResourceID = utils.String(writeAccess)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, workspace.WorkspaceName, id.LinkedServiceName, parameters)
	if err != nil {
		return fmt.Errorf("creating Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating future for Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
	}

	_, err = client.Get(ctx, resourceGroup, workspace.WorkspaceName, id.LinkedServiceName)
	if err != nil {
		return fmt.Errorf("retrieving Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
	}

	d.SetId(id.ID())

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
	workspace := parse.NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, workspaceName)

	resp, err := client.Get(ctx, resourceGroup, workspaceName, serviceType)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on AzureRM Log Analytics Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, serviceType, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("workspace_id", workspace.ID())

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
		return fmt.Errorf("deleting Log Analytics Linked Service '%s/%s' (Resource Group %q): %+v", workspaceName, serviceType, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Log Analytics Linked Service '%s/%s' (Resource Group %q): %+v", workspaceName, serviceType, resourceGroup, err)
		}
	}

	// (@WodansSon) - This is a bug in the service API, it returns instantly from the delete call with a 200
	// so we must wait for the state to change before we return from the delete function
	deleteWait := logAnalyticsLinkedServiceDeleteWaitForState(ctx, meta, d.Timeout(schema.TimeoutDelete), resourceGroup, workspaceName, serviceType)

	if _, err := deleteWait.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Log Analytics Cluster to finish deleting '%s/%s' (Resource Group %q): %+v", workspaceName, serviceType, resourceGroup, err)
	}

	return nil
}

func LogAnalyticsLinkedServiceType(readAccessId string) string {
	if readAccessId != "" {
		return "Automation"
	}

	return "Cluster"
}

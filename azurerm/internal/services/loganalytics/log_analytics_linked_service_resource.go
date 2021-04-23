package loganalytics

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/operationalinsights/mgmt/2020-08-01/operationalinsights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	validateAuto "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/automation/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceLogAnalyticsLinkedService() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogAnalyticsLinkedServiceCreateUpdate,
		Read:   resourceLogAnalyticsLinkedServiceRead,
		Update: resourceLogAnalyticsLinkedServiceCreateUpdate,
		Delete: resourceLogAnalyticsLinkedServiceDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			// TODO: Remove in 3.0
			"workspace_name": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     validate.LogAnalyticsWorkspaceName,
				ExactlyOneOf:     []string{"workspace_name", "workspace_id"},
				Deprecated:       "This field has been deprecated in favour of `workspace_id` and will be removed in a future version of the provider",
			},

			"workspace_id": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc:     azure.ValidateResourceID,
				ExactlyOneOf:     []string{"workspace_name", "workspace_id"},
			},

			// TODO: Remove in 3.0
			"linked_service_name": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					"automation",
					"cluster",
				}, false),
				Deprecated: "This field has been deprecated and will be removed in a future version of the provider",
			},

			// TODO: Remove in 3.0
			"resource_id": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				ExactlyOneOf: []string{"read_access_id", "write_access_id", "resource_id"},
				Deprecated:   "This field has been deprecated in favour of `read_access_id` and will be removed in a future version of the provider",
			},

			"read_access_id": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceID,
				ExactlyOneOf: []string{"read_access_id", "write_access_id", "resource_id"},
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

		// TODO: Remove in 3.0
		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *schema.ResourceDiff, v interface{}) error {
			if d.HasChange("linked_service_name") {
				oldServiceName, newServiceName := d.GetChange("linked_service_name")

				// This is an unneeded field, if it is removed you can safely ignore it
				// as it's value can be(and is) derived via the 'read_access_id' field. It
				// is only here for backwards compatibility to avoid a breaking change
				if newServiceName.(string) != "" {
					// Ignore change if it's in case only
					if !strings.EqualFold(oldServiceName.(string), newServiceName.(string)) {
						d.ForceNew("linked_service_name")
					}
				}
			}

			if d.HasChange("workspace_id") {
				forceNew := true
				_, newWorkspaceName := d.GetChange("workspace_name")
				oldWorkspaceID, newWorkspaceID := d.GetChange("workspace_id")

				// If the workspcae ID has been removed, only do a force new if the new workspace name
				// and the old workspace ID points to different workspaces
				if oldWorkspaceID.(string) != "" && newWorkspaceName.(string) != "" && newWorkspaceID.(string) == "" {
					workspace, err := parse.LogAnalyticsWorkspaceID(oldWorkspaceID.(string))
					if err == nil {
						if workspace.WorkspaceName == newWorkspaceName.(string) {
							forceNew = false
						}
					}
				}

				if forceNew {
					d.ForceNew("workspace_id")
				}
			}

			if d.HasChange("workspace_name") {
				forceNew := true
				oldWorkspaceName, newWorkspaceName := d.GetChange("workspace_name")
				_, newWorkspaceID := d.GetChange("workspace_id")

				// If the workspcae name has been removed, only do a force new if the new workspace ID
				// and the old workspace name points to different workspaces
				if oldWorkspaceName.(string) != "" && newWorkspaceID.(string) != "" && newWorkspaceName.(string) == "" {
					workspace, err := parse.LogAnalyticsWorkspaceID(newWorkspaceID.(string))
					if err == nil {
						if workspace.WorkspaceName == oldWorkspaceName.(string) {
							forceNew = false
						}
					}
				}

				if forceNew {
					d.ForceNew("workspace_name")
				}
			}

			// TODO: Remove in 3.0
			if d.HasChange("resource_id") {
				if resourceID := d.Get("resource_id").(string); resourceID != "" {
					if _, err := validateAuto.AutomationAccountID(resourceID, "resource_id"); err != nil {
						return fmt.Errorf("'resource_id' must be an Automation Account resource ID, got %q", resourceID)
					}
				}
			}

			if d.HasChange("read_access_id") {
				if readAccessID := d.Get("read_access_id").(string); readAccessID != "" {
					if _, err := validateAuto.AutomationAccountID(readAccessID, "read_acces_id"); err != nil {
						return fmt.Errorf("'read_access_id' must be an Automation Account resource ID, got %q", readAccessID)
					}
				}
			}

			if d.HasChange("write_access_id") {
				if writeAccessID := d.Get("write_access_id").(string); writeAccessID != "" {
					if _, err := validate.LogAnalyticsClusterID(writeAccessID, "write_access_id"); err != nil {
						return fmt.Errorf("'write_access_id' must be a Log Analytics Cluster resource ID, got %q", writeAccessID)
					}
				}
			}

			return nil
		}),
	}
}

func resourceLogAnalyticsLinkedServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).LogAnalytics.LinkedServicesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics Linked Services creation.")

	// TODO: Remove in 3.0
	var tmpSpace parse.LogAnalyticsWorkspaceId
	var workspaceId string

	resourceGroup := d.Get("resource_group_name").(string)
	readAccess := d.Get("read_access_id").(string)
	writeAccess := d.Get("write_access_id").(string)
	linkedServiceName := d.Get("linked_service_name").(string)
	t := d.Get("tags").(map[string]interface{})

	if resourceId := d.Get("resource_id").(string); resourceId != "" {
		readAccess = resourceId
	}

	if workspaceName := d.Get("workspace_name").(string); workspaceName != "" {
		tmpSpace = parse.NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, workspaceName)
		workspaceId = tmpSpace.ID()
	} else {
		workspaceId = d.Get("workspace_id").(string)
	}

	workspace, err := parse.LogAnalyticsWorkspaceID(workspaceId)
	if err != nil {
		return fmt.Errorf("Linked Service (Resource Group %q) unable to parse workspace id: %+v", resourceGroup, err)
	}

	id := parse.NewLogAnalyticsLinkedServiceID(subscriptionId, resourceGroup, workspace.WorkspaceName, LogAnalyticsLinkedServiceType(readAccess))

	if linkedServiceName != "" {
		if !strings.EqualFold(linkedServiceName, LogAnalyticsLinkedServiceType(readAccess)) {
			return fmt.Errorf("Linked Service '%s/%s' (Resource Group %q): 'linked_service_name' %q does not match expected value of %q", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, linkedServiceName, LogAnalyticsLinkedServiceType(readAccess))
		}
	}

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

	if _, err = client.Get(ctx, resourceGroup, workspace.WorkspaceName, id.LinkedServiceName); err != nil {
		return fmt.Errorf("retrieving Linked Service '%s/%s' (Resource Group %q): %+v", workspace.WorkspaceName, id.LinkedServiceName, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceLogAnalyticsLinkedServiceRead(d, meta)
}

func resourceLogAnalyticsLinkedServiceRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("workspace_name", workspaceName)
	d.Set("linked_service_name", serviceType)

	if props := resp.LinkedServiceProperties; props != nil {
		d.Set("resource_id", props.ResourceID)
		d.Set("read_access_id", props.ResourceID)
		d.Set("write_access_id", props.WriteAccessResourceID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceLogAnalyticsLinkedServiceDelete(d *schema.ResourceData, meta interface{}) error {
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

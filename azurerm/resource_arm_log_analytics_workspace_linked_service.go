package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmLogAnalyticsWorkspaceLinkedService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmLogAnalyticsWorkspaceLinkedServiceCreateUpdate,
		Read:   resourceArmLogAnalyticsWorkspaceLinkedServiceRead,
		Update: resourceArmLogAnalyticsWorkspaceLinkedServiceCreateUpdate,
		Delete: resourceArmLogAnalyticsWorkspaceLinkedServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"workspace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAzureRmLogAnalyticsWorkspaceName,
			},

			"linked_service_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "automation",
				ValidateFunc: validation.NoZeroValues,
			},

			"linked_service_properties": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"tags": tagsSchema(),

			// Exported properties
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmLogAnalyticsWorkspaceLinkedServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).linkedServicesClient
	ctx := meta.(*ArmClient).StopContext
	log.Printf("[INFO] preparing arguments for AzureRM Log Analytics linked services creation.")

	resGroup := d.Get("resource_group_name").(string)

	workspaceName := d.Get("workspace_name").(string)
	lsName := d.Get("linked_service_name").(string)

	props := d.Get("linked_service_properties").(map[string]interface{})
	resourceID := props["resource_id"].(string)

	tags := d.Get("tags").(map[string]interface{})

	parameters := operationalinsights.LinkedService{
		Tags: expandTags(tags),
		LinkedServiceProperties: &operationalinsights.LinkedServiceProperties{
			ResourceID: &resourceID,
		},
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, workspaceName, lsName, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create request for Log Analytics Workspace Linked Service %q/%q (Resource Group %q): %+v", workspaceName, lsName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, workspaceName, lsName)
	if err != nil {
		return fmt.Errorf("Error retrieving Analytics Workspace Linked Service %q/%q (Resource Group %q): %+v", workspaceName, lsName, resGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Log Analytics Linked Service '%s' (resource group %s) ID", lsName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLogAnalyticsWorkspaceLinkedServiceRead(d, meta)

}

func resourceArmLogAnalyticsWorkspaceLinkedServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).linkedServicesClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	workspaceName := id.Path["workspaces"]
	lsName := id.Path["linkedservices"]

	resp, err := client.Get(ctx, resGroup, workspaceName, lsName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Log Analytics Linked Service '%s': %+v", lsName, err)
	}
	if resp.ID == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", *resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("workspace_name", workspaceName)
	d.Set("linked_service_name", lsName)

	linkedServiceProperties := flattenLogAnalyticsWorkspaceLinkedServiceProperties(resp.LinkedServiceProperties)
	if err := d.Set("linked_service_properties", linkedServiceProperties); err != nil {
		return fmt.Errorf("Error setting Log Analytics Linked Service Properties: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}

func flattenLogAnalyticsWorkspaceLinkedServiceProperties(input *operationalinsights.LinkedServiceProperties) interface{} {
	properties := make(map[string]interface{})

	// resource id linked service
	if resourceID := input.ResourceID; resourceID != nil {
		properties["resource_id"] = interface{}(*resourceID)
	}

	return interface{}(properties)
}

func resourceArmLogAnalyticsWorkspaceLinkedServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).linkedServicesClient
	ctx := meta.(*ArmClient).StopContext
	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	workspaceName := id.Path["workspaces"]
	lsName := id.Path["linkedservices"]

	resp, err := client.Delete(ctx, resGroup, workspaceName, lsName)

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error issuing AzureRM delete request for Log Analytics Linked Service '%s': %+v", lsName, err)
	}

	return nil
}

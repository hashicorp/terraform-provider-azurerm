package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2015-11-01-preview/operationalinsights"
	"github.com/hashicorp/terraform/helper/schema"
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

		Schema: map[string]*schema.Schema{
			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"linked_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "automation",
			},

			"linked_service_properties": {
				Type:     schema.TypeMap,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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

func resourceArmLogAnalyticsLinkedServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
		// Name: &name,
		Tags: expandTags(tags),
		LinkedServiceProperties: &operationalinsights.LinkedServiceProperties{
			ResourceID: &resourceID,
		},
	}

	_, err := client.CreateOrUpdate(ctx, resGroup, workspaceName, lsName, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, workspaceName, lsName)
	if err != nil {
		return err
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read Log Analytics Linked Service '%s' (resource group %s) ID", lsName, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmLogAnalyticsLinkedServiceRead(d, meta)

}

func resourceArmLogAnalyticsLinkedServiceRead(d *schema.ResourceData, meta interface{}) error {
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

	linkedServiceProperties := flattenLogAnalyticsLinkedServiceProperties(resp.LinkedServiceProperties)
	if err := d.Set("linked_service_properties", linkedServiceProperties); err != nil {
		return fmt.Errorf("Error flattening Log Analytics Linked Service Properties: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)
	return nil
}

func flattenLogAnalyticsLinkedServiceProperties(input *operationalinsights.LinkedServiceProperties) interface{} {
	properties := make(map[string]interface{}, 0)

	// resource id linked service
	if resourceID := input.ResourceID; resourceID != nil {
		properties["resource_id"] = interface{}(*resourceID)
	}

	return interface{}(properties)
}

func resourceArmLogAnalyticsLinkedServiceDelete(d *schema.ResourceData, meta interface{}) error {
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

package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApplicationSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApplicationSecurityGroupCreateUpdate,
		Read:   resourceArmApplicationSecurityGroupRead,
		Update: resourceArmApplicationSecurityGroupCreateUpdate,
		Delete: resourceArmApplicationSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmApplicationSecurityGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationSecurityGroupsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	tags := d.Get("tags").(map[string]interface{})

	securityGroup := network.ApplicationSecurityGroup{
		Location: utils.String(location),
		Tags:     expandTags(tags),
	}
	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, securityGroup)
	if err != nil {
		return fmt.Errorf("Error creating Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the Application Security Group %q (Resource Group %q) to finish creating: %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Application Security Group %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApplicationSecurityGroupRead(d, meta)
}

func resourceArmApplicationSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationSecurityGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["applicationSecurityGroups"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmApplicationSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationSecurityGroupsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["applicationSecurityGroups"]

	log.Printf("[DEBUG] Deleting Application Security Group %q (resource group %q)", name, resourceGroup)

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error issuing delete request for Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deletion of Application Security Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

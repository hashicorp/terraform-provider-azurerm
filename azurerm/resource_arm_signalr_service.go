package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/signalr/mgmt/2018-03-01-preview/signalr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSignalRService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSignalRServiceCreateOrUpdate,
		Read:   resourceArmSignalRServiceRead,
		Update: resourceArmSignalRServiceCreateOrUpdate,
		Delete: resourceArmSignalRServiceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Free_F1",
								"Standard_S1",
							}, false),
						},

						"capacity": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntInSlice([]int{1, 2, 5, 10, 20, 50, 100}),
						},
					},
				},
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"server_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSignalRServiceCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).signalRClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	sku := d.Get("sku").([]interface{})

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := &signalr.CreateParameters{
		Location: utils.String(location),
		Sku:      expandSignalRServiceSku(sku),
		Tags:     expandedTags,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating or updating SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the result of creating or updating SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("SignalR %q (Resource Group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmSignalRServiceRead(d, meta)
}

func resourceArmSignalRServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).signalRClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["SignalR"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] SignalR %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	if err = d.Set("sku", flattenSignalRServiceSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}
	if properties := resp.Properties; properties != nil {
		d.Set("hostname", properties.HostName)
		d.Set("ip_address", properties.ExternalIP)
		d.Set("public_port", properties.PublicPort)
		d.Set("server_port", properties.ServerPort)
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSignalRServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).signalRClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["SignalR"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
		return nil
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of SignalR %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandSignalRServiceSku(input []interface{}) *signalr.ResourceSku {
	v := input[0].(map[string]interface{})
	return &signalr.ResourceSku{
		Name:     utils.String(v["name"].(string)),
		Capacity: utils.Int32(int32(v["capacity"].(int))),
	}
}

func flattenSignalRServiceSku(input *signalr.ResourceSku) []interface{} {
	result := make(map[string]interface{})
	if input != nil {
		if input.Name != nil {
			result["name"] = *input.Name
		}
		if input.Capacity != nil {
			result["capacity"] = *input.Capacity
		}
	}
	return []interface{}{result}
}

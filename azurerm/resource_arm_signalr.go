package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/signalr/mgmt/2018-03-01-preview/signalr"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmSignalR() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSignalRCreateOrUpdate,
		Read:   resourceArmSignalRRead,
		Update: resourceArmSignalRCreateOrUpdate,
		Delete: resourceArmSignalRDelete,

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

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Free_F1",
					"Standard_S1",
				}, false),
			},

			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
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

func resourceArmSignalRCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).signalRClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	sku := d.Get("sku_name").(string)
	capacity := d.Get("capacity").(int)

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	parameters := &signalr.CreateParameters{
		Location: utils.String(location),
		Sku: &signalr.ResourceSku{
			Name:     utils.String(sku),
			Capacity: utils.Int32(int32(capacity)),
		},
		Tags: expandedTags,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating or updating SignalR %q (resource group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the result of creating or updating SignalR %q (resource group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("SignalR %q (resource group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmSignalRRead(d, meta)
}

func resourceArmSignalRRead(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("Error getting SignalR %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}
	if sku := resp.Sku; sku != nil {
		if sku.Name != nil {
			d.Set("sku_name", *sku.Name)
		}
		if sku.Capacity != nil {
			d.Set("capacity", *sku.Capacity)
		}
	}
	if properties := resp.Properties; properties != nil {
		if properties.HostName != nil {
			d.Set("hostname", *properties.HostName)
		}
		if properties.ExternalIP != nil {
			d.Set("ip_address", *properties.ExternalIP)
		}
		if properties.PublicPort != nil {
			d.Set("port", *properties.PublicPort)
		}
		if properties.ServerPort != nil {
			d.Set("server_port", *properties.ServerPort)
		}
	}
	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmSignalRDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting SignalR %q (resource group %q): %+v", name, resourceGroup, err)
		}
		return nil
	}
	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting SignalR %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

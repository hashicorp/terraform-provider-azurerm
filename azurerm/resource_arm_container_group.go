package azurerm

import (
	"github.com/Azure/azure-sdk-for-go/arm/containerinstance"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmContainerGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmContainerGroupCreate,
		Read:   resourceArmContainerGroupRead,
		Update: resourceArmContainerGroupCreate,
		Delete: resourceArmContainerGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cpu_cores": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"command_line": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"env_vars": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"ip_address_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"memory": {
				Type:     schema.TypeFloat,
				Optional: true,
				ForceNew: true,
			},

			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"tags": tagsSchema(),
		},
	}
}

// func containerSchema() *schema.Resource {
// 	return &schema.Resource{
// 		Type:schema
// 	}
// }

func resourceArmContainerGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	containterGroupsClient := client.containerGroupsClient

	// container group properties
	resGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)
	location := d.Get("location").(string)
	OSType := d.Get("os_type").(string)
	tags := d.Get("tags").(map[string]interface{})

	// per container properties
	// image := d.Get("image").(string)
	// cores := d.Get("cpu_cores").(int)
	// cmdLine := d.Get("command_line").(string)
	// envVars := d.Get("env_vars").(string)
	// IPAddressType := d.Get("ip_address_type").(string)
	// memory := d.Get("memory").(float32)
	// port := d.Get("port").(int)

	containerGroup := containerinstance.ContainerGroup{
		Name:     &name,
		Type:     &OSType,
		Location: &location,
		Tags:     expandTags(tags),
	}

	_, error := containterGroupsClient.CreateOrUpdate(resGroup, name, containerGroup)
	if error != nil {
		return error
	}

	return nil
}
func resourceArmContainerGroupRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}
func resourceArmContainerGroupDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

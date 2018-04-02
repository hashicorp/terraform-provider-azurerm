package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNotificationHubNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNotificationHubNamespaceCreateUpdate,
		Read:   resourceArmNotificationHubNamespaceRead,
		Update: resourceArmNotificationHubNamespaceCreateUpdate,
		Delete: resourceArmNotificationHubNamespaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"critical": {
				Type: schema.TypeBool,
				// TODO: is this a sensible default?
				Optional: true,
				Default:  false,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"namespace_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(notificationhubs.Messaging),
					string(notificationhubs.NotificationHub),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"tags": tagsSchema(),

			"servicebus_endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmNotificationHubNamespaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	tags := d.Get("tags").(map[string]interface{})

	namespaceType := d.Get("namespace_type").(string)
	critical := d.Get("critical").(bool)
	enabled := d.Get("enabled").(bool)

	parameters := notificationhubs.NamespaceCreateOrUpdateParameters{
		Location: utils.String(location),
		Tags:     expandTags(tags),
		NamespaceProperties: &notificationhubs.NamespaceProperties{
			Region:        utils.String(location),
			NamespaceType: notificationhubs.NamespaceType(namespaceType),
			Critical:      utils.Bool(critical),
			Enabled:       utils.Bool(enabled),
		},
	}
	_, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Notification Hub Namespace %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNotificationHubNamespaceRead(d, meta)
}

func resourceArmNotificationHubNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Notification Hub Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	//  TODO: SKU?

	if props := resp.NamespaceProperties; props != nil {
		d.Set("critical", props.Critical)
		d.Set("enabled", props.Enabled)
		d.Set("namespace_type", props.NamespaceType)
		d.Set("servicebus_endpoint", props.ServiceBusEndpoint)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmNotificationHubNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationNamespacesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["namespaces"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error deleting Notification Hub Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of Notification Hub Namespace %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

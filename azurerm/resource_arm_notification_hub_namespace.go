package azurerm

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform/helper/resource"
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

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(notificationhubs.Basic),
								string(notificationhubs.Free),
								string(notificationhubs.Standard),
							}, false),
						},
					},
				},
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

			// NOTE: skipping tags as there's a bug in the API where the Keys for Tags are returned in lower-case
			// Azure Rest API Specs issue: https://github.com/Azure/azure-sdk-for-go/issues/2239
			//"tags": tagsSchema(),

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
	location := azureRMNormalizeLocation(d.Get("location").(string))

	sku := expandNotificationHubNamespacesSku(d.Get("sku").([]interface{}))

	namespaceType := d.Get("namespace_type").(string)
	enabled := d.Get("enabled").(bool)

	parameters := notificationhubs.NamespaceCreateOrUpdateParameters{
		Location: utils.String(location),
		Sku:      sku,
		NamespaceProperties: &notificationhubs.NamespaceProperties{
			Region:        utils.String(location),
			NamespaceType: notificationhubs.NamespaceType(namespaceType),
			Enabled:       utils.Bool(enabled),
		},
	}
	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating Notification Hub Namesapce %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Notification Hub Namesapce %q (Resource Group %q): %+v", name, resourceGroup, err)
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
			log.Printf("[DEBUG] Notification Hub Namespace %q (Resource Group %q) was not found - removing from state!", name, resourceGroup)
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

	sku := flattenNotificationHubNamespacesSku(resp.Sku)
	if err := d.Set("sku", sku); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if props := resp.NamespaceProperties; props != nil {
		d.Set("enabled", props.Enabled)
		d.Set("namespace_type", props.NamespaceType)
		d.Set("servicebus_endpoint", props.ServiceBusEndpoint)
	}

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

	// the future returned from the Delete method is broken 50% of the time - let's poll ourselves for now
	// Related Bug: https://github.com/Azure/azure-sdk-for-go/issues/2254
	log.Printf("[DEBUG] Waiting for Notification Hub Namespace %q (Resource Group %q) to be deleted", name, resourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"200", "202"},
		Target:  []string{"404"},
		Refresh: notificationHubNamespaceStateRefreshFunc(ctx, client, resourceGroup, name),
		Timeout: 10 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Notification Hub %q (Resource Group %q) to be deleted: %s", name, resourceGroup, err)
	}

	return nil
}

func expandNotificationHubNamespacesSku(input []interface{}) *notificationhubs.Sku {
	v := input[0].(map[string]interface{})

	skuName := v["name"].(string)

	return &notificationhubs.Sku{
		Name: notificationhubs.SkuName(skuName),
	}
}

func flattenNotificationHubNamespacesSku(input *notificationhubs.Sku) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	output := map[string]interface{}{
		"name": string(input.Name),
	}
	outputs = append(outputs, output)
	return outputs
}

func notificationHubNamespaceStateRefreshFunc(ctx context.Context, client notificationhubs.NamespacesClient, resourceGroupName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(res.Response) {
				return nil, "", fmt.Errorf("Error retrieving Notification Hub Namespace %q (Resource Group %q): %s", name, resourceGroupName, err)
			}
		}

		// Note: this exists as the Delete API only seems to work some of the time
		// in this case we're going to try triggering the Deletion again, in-case it didn't work prior to this attepmpt
		// Upstream Bug: https://github.com/Azure/azure-sdk-for-go/issues/2254

		if _, err := client.Delete(ctx, resourceGroupName, name); err != nil {
			log.Printf("Error reissuing Notification Hub Namespace %q delete request (Resource Group %q): %+v", name, resourceGroupName, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

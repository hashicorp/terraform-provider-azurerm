package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmNotificationHubAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNotificationHubAuthorizationRuleCreateUpdate,
		Read:   resourceArmNotificationHubAuthorizationRuleRead,
		Update: resourceArmNotificationHubAuthorizationRuleCreateUpdate,
		Delete: resourceArmNotificationHubAuthorizationRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		// TODO: customizeDiff for send+listen when manage selected

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"notification_hub_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"namespace_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"manage": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"send": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"listen": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"primary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secondary_access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmNotificationHubAuthorizationRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubs.HubsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	notificationHubName := d.Get("notification_hub_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	manage := d.Get("manage").(bool)
	send := d.Get("send").(bool)
	listen := d.Get("listen").(bool)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", name, notificationHubName, namespaceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_notification_hub_authorization_rule", *existing.ID)
		}
	}

	locks.ByName(notificationHubName, notificationHubResourceName)
	defer locks.UnlockByName(notificationHubName, notificationHubResourceName)

	locks.ByName(namespaceName, notificationHubNamespaceResourceName)
	defer locks.UnlockByName(namespaceName, notificationHubNamespaceResourceName)

	parameters := notificationhubs.SharedAccessAuthorizationRuleCreateOrUpdateParameters{
		Properties: &notificationhubs.SharedAccessAuthorizationRuleProperties{
			Rights: expandNotificationHubAuthorizationRuleRights(manage, send, listen),
		},
	}

	if _, err := client.CreateOrUpdateAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", name, notificationHubName, namespaceName, resourceGroup, err)
	}

	read, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, name)
	if err != nil {
		return fmt.Errorf("Error retrieving Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", name, notificationHubName, namespaceName, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q) ID", name, notificationHubName, namespaceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmNotificationHubAuthorizationRuleRead(d, meta)

}

func resourceArmNotificationHubAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubs.HubsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	notificationHubName := id.Path["notificationHubs"]
	name := id.Path["AuthorizationRules"]

	resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Authorization Rule %q was not found in Notification Hub %q / Namespace %q / Resource Group %q", name, notificationHubName, namespaceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", name, notificationHubName, namespaceName, resourceGroup, err)
	}

	keysResp, err := client.ListKeys(ctx, resourceGroup, namespaceName, notificationHubName, name)
	if err != nil {
		return fmt.Errorf("Error Listing Access Keys for Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", name, notificationHubName, namespaceName, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("notification_hub_name", notificationHubName)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.SharedAccessAuthorizationRuleProperties; props != nil {
		manage, send, listen := flattenNotificationHubAuthorizationRuleRights(props.Rights)
		d.Set("manage", manage)
		d.Set("send", send)
		d.Set("listen", listen)
	}

	d.Set("primary_access_key", keysResp.PrimaryKey)
	d.Set("secondary_access_key", keysResp.SecondaryKey)

	return nil
}

func resourceArmNotificationHubAuthorizationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).notificationHubs.HubsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	notificationHubName := id.Path["notificationHubs"]
	name := id.Path["AuthorizationRules"]

	locks.ByName(notificationHubName, notificationHubResourceName)
	defer locks.UnlockByName(notificationHubName, notificationHubResourceName)

	locks.ByName(namespaceName, notificationHubNamespaceResourceName)
	defer locks.UnlockByName(namespaceName, notificationHubNamespaceResourceName)

	resp, err := client.DeleteAuthorizationRule(ctx, resourceGroup, namespaceName, notificationHubName, name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", name, notificationHubName, namespaceName, resourceGroup, err)
		}
	}

	return nil
}

func expandNotificationHubAuthorizationRuleRights(manage bool, send bool, listen bool) *[]notificationhubs.AccessRights {
	rights := make([]notificationhubs.AccessRights, 0)

	if manage {
		rights = append(rights, notificationhubs.Manage)
	}

	if send {
		rights = append(rights, notificationhubs.Send)
	}

	if listen {
		rights = append(rights, notificationhubs.Listen)
	}

	return &rights
}

func flattenNotificationHubAuthorizationRuleRights(input *[]notificationhubs.AccessRights) (manage bool, send bool, listen bool) {
	if input == nil {
		return
	}

	for _, right := range *input {
		switch right {
		case notificationhubs.Manage:
			manage = true
			continue
		case notificationhubs.Send:
			send = true
			continue
		case notificationhubs.Listen:
			listen = true
			continue
		}
	}

	return manage, send, listen
}

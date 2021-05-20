package notificationhub

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/notificationhubs/mgmt/2017-04-01/notificationhubs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/notificationhub/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceNotificationHubAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNotificationHubAuthorizationRuleCreateUpdate,
		Read:   resourceNotificationHubAuthorizationRuleRead,
		Update: resourceNotificationHubAuthorizationRuleCreateUpdate,
		Delete: resourceNotificationHubAuthorizationRuleDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NotificationHubAuthorizationRuleID(id)
			return err
		}),
		// TODO: customizeDiff for send+listen when manage selected

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

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

func resourceNotificationHubAuthorizationRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	notificationHubName := d.Get("notification_hub_name").(string)
	namespaceName := d.Get("namespace_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	manage := d.Get("manage").(bool)
	send := d.Get("send").(bool)
	listen := d.Get("listen").(bool)

	if d.IsNewResource() {
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

	return resourceNotificationHubAuthorizationRuleRead(d, meta)
}

func resourceNotificationHubAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NotificationHubAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetAuthorizationRule(ctx, id.ResourceGroup, id.NamespaceName, id.NotificationHubName, id.AuthorizationRuleName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Authorization Rule %q was not found in Notification Hub %q / Namespace %q / Resource Group %q", id.AuthorizationRuleName, id.NotificationHubName, id.NamespaceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", id.AuthorizationRuleName, id.NotificationHubName, id.NamespaceName, id.ResourceGroup, err)
	}

	keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.NamespaceName, id.NotificationHubName, id.AuthorizationRuleName)
	if err != nil {
		return fmt.Errorf("Error Listing Access Keys for Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", id.AuthorizationRuleName, id.NotificationHubName, id.NamespaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.AuthorizationRuleName)
	d.Set("notification_hub_name", id.NotificationHubName)
	d.Set("namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

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

func resourceNotificationHubAuthorizationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).NotificationHubs.HubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NotificationHubAuthorizationRuleID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NotificationHubName, notificationHubResourceName)
	defer locks.UnlockByName(id.NotificationHubName, notificationHubResourceName)

	locks.ByName(id.NamespaceName, notificationHubNamespaceResourceName)
	defer locks.UnlockByName(id.NamespaceName, notificationHubNamespaceResourceName)

	resp, err := client.DeleteAuthorizationRule(ctx, id.ResourceGroup, id.NamespaceName, id.NotificationHubName, id.AuthorizationRuleName)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Authorization Rule %q (Notification Hub %q / Namespace %q / Resource Group %q): %+v", id.AuthorizationRuleName, id.NotificationHubName, id.NamespaceName, id.ResourceGroup, err)
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
		rights = append(rights, notificationhubs.SendEnumValue)
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
		case notificationhubs.SendEnumValue:
			send = true
			continue
		case notificationhubs.Listen:
			listen = true
			continue
		}
	}

	return manage, send, listen
}

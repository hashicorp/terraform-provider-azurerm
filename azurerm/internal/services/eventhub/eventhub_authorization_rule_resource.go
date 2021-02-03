package eventhub

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/eventhub/mgmt/2018-01-01-preview/eventhub"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceEventHubAuthorizationRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceEventHubAuthorizationRuleCreateUpdate,
		Read:   resourceEventHubAuthorizationRuleRead,
		Update: resourceEventHubAuthorizationRuleCreateUpdate,
		Delete: resourceEventHubAuthorizationRuleDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: EventHubAuthorizationRuleSchemaFrom(map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateEventHubAuthorizationRuleName(),
			},

			"namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateEventHubNamespaceName(),
			},

			"eventhub_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateEventHubName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),
		}),

		CustomizeDiff: EventHubAuthorizationRuleCustomizeDiff,
	}
}

func resourceEventHubAuthorizationRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM EventHub Authorization Rule creation.")

	name := d.Get("name").(string)
	namespaceName := d.Get("namespace_name").(string)
	eventHubName := d.Get("eventhub_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing EventHub Authorization Rule %q (EventHub %q / Namespace %q / Resource Group %q): %s", name, eventHubName, namespaceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_eventhub_authorization_rule", *existing.ID)
		}
	}

	locks.ByName(eventHubName, eventHubResourceName)
	defer locks.UnlockByName(eventHubName, eventHubResourceName)

	locks.ByName(namespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(namespaceName, eventHubNamespaceResourceName)

	parameters := eventhub.AuthorizationRule{
		Name: &name,
		AuthorizationRuleProperties: &eventhub.AuthorizationRuleProperties{
			Rights: ExpandEventHubAuthorizationRuleRights(d),
		},
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		if _, err := client.CreateOrUpdateAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name, parameters); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error creating Authorization Rule %q (event Hub %q / Namespace %q / Resource Group %q): %+v", name, eventHubName, namespaceName, resourceGroup, err))
		}

		read, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name)
		if err != nil {
			if utils.ResponseWasNotFound(read.Response) {
				return resource.RetryableError(fmt.Errorf("Expected instance of the Authorization Rule %q (event Hub %q / Namespace %q / Resource Group %q) to be created but was in non existent state, retrying", name, eventHubName, namespaceName, resourceGroup))
			}
			return resource.NonRetryableError(fmt.Errorf("Expected instance of Authorization Rule %q (event Hub %q / Namespace %q / Resource Group %q) could not be found", name, eventHubName, namespaceName, resourceGroup))
		}

		if read.ID == nil {
			return resource.NonRetryableError(fmt.Errorf("Cannot read Authorization Rule %q (event Hub %q / Namespace %q / Resource Group %q) ID", name, eventHubName, namespaceName, resourceGroup))
		}

		d.SetId(*read.ID)

		return resource.NonRetryableError(resourceEventHubAuthorizationRuleRead(d, meta))
	})
}

func resourceEventHubAuthorizationRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["authorizationRules"]
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	eventHubName := id.Path["eventhubs"]

	resp, err := client.GetAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure EventHub Authorization Rule %s: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("eventhub_name", eventHubName)
	d.Set("namespace_name", namespaceName)
	d.Set("resource_group_name", resourceGroup)

	if properties := resp.AuthorizationRuleProperties; properties != nil {
		listen, send, manage := FlattenEventHubAuthorizationRuleRights(properties.Rights)
		d.Set("manage", manage)
		d.Set("listen", listen)
		d.Set("send", send)
	}

	keysResp, err := client.ListKeys(ctx, resourceGroup, namespaceName, eventHubName, name)
	if err != nil {
		return fmt.Errorf("Error making Read request on Azure EventHub Authorization Rule List Keys %s: %+v", name, err)
	}

	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("secondary_key", keysResp.SecondaryKey)
	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)
	d.Set("primary_connection_string_alias", keysResp.AliasPrimaryConnectionString)
	d.Set("secondary_connection_string_alias", keysResp.AliasSecondaryConnectionString)

	return nil
}

func resourceEventHubAuthorizationRuleDelete(d *schema.ResourceData, meta interface{}) error {
	eventhubClient := meta.(*clients.Client).Eventhub.EventHubsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["authorizationRules"]
	resourceGroup := id.ResourceGroup
	namespaceName := id.Path["namespaces"]
	eventHubName := id.Path["eventhubs"]

	locks.ByName(eventHubName, eventHubResourceName)
	defer locks.UnlockByName(eventHubName, eventHubResourceName)

	locks.ByName(namespaceName, eventHubNamespaceResourceName)
	defer locks.UnlockByName(namespaceName, eventHubNamespaceResourceName)

	if resp, err := eventhubClient.DeleteAuthorizationRule(ctx, resourceGroup, namespaceName, eventHubName, name); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error issuing Azure ARM delete request of EventHub Authorization Rule '%s': %+v", name, err)
		}
	}

	return nil
}

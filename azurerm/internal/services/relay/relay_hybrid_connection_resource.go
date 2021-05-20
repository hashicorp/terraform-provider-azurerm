package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmRelayHybridConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmRelayHybridConnectionCreateUpdate,
		Read:   resourceArmRelayHybridConnectionRead,
		Update: resourceArmRelayHybridConnectionCreateUpdate,
		Delete: resourceArmRelayHybridConnectionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.HybridConnectionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"relay_namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"requires_client_authorization": {
				Type:     schema.TypeBool,
				Default:  true,
				ForceNew: true,
				Optional: true,
			},
			"user_metadata": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmRelayHybridConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Relay Hybrid Connection creation.")

	resourceId := parse.NewHybridConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("relay_namespace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.NamespaceName, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Hybrid Connection %q (Namespace %q / Resource Group %q): %+v", resourceId.Name, resourceId.NamespaceName, resourceId.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_relay_hybrid_connection", resourceId.ID())
		}
	}

	requireClientAuthorization := d.Get("requires_client_authorization").(bool)
	userMetadata := d.Get("user_metadata").(string)

	parameters := relay.HybridConnection{
		HybridConnectionProperties: &relay.HybridConnectionProperties{
			RequiresClientAuthorization: &requireClientAuthorization,
			UserMetadata:                &userMetadata,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.NamespaceName, resourceId.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating Hybrid Connection %q (Namespace %q Resource Group %q): %+v", resourceId.Name, resourceId.NamespaceName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceArmRelayHybridConnectionRead(d, meta)
}

func resourceArmRelayHybridConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Hybrid Connection %q (Namespace %q / Resource Group %q): %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("relay_namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.HybridConnectionProperties; props != nil {
		d.Set("requires_client_authorization", props.RequiresClientAuthorization)
		d.Set("user_metadata", props.UserMetadata)
	}

	return nil
}

func resourceArmRelayHybridConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.HybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	rc, err := client.Delete(ctx, id.ResourceGroup, id.NamespaceName, id.Name)
	if err != nil {
		if response.WasNotFound(rc.Response) {
			return nil
		}

		return fmt.Errorf("deleting Hybrid Connection %q (Namespace %q / Resource Group %q): %+v", id.NamespaceName, id.NamespaceName, id.ResourceGroup, err)
	}

	log.Printf("[INFO] Waiting for Hybrid Connection %q (Namespace %q / Resource Group %q) to be deleted", id.Name, id.NamespaceName, id.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Deleted"},
		Refresh:    hybridConnectionDeleteRefreshFunc(ctx, client, id.ResourceGroup, id.NamespaceName, id.Name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Relay Hybrid Connection %q (Namespace %q Resource Group %q) to be deleted: %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	return nil
}

func hybridConnectionDeleteRefreshFunc(ctx context.Context, client *relay.HybridConnectionsClient, resourceGroupName string, relayNamespace string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, relayNamespace, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}

			return nil, "Error", fmt.Errorf("Error issuing read request in relayNamespaceDeleteRefreshFunc to Relay Hybrid Connection %q (Namespace %q Resource Group %q): %s", name, relayNamespace, resourceGroupName, err)
		}

		return res, "Pending", nil
	}
}

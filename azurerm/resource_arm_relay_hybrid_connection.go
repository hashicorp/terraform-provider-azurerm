package azurerm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmHybridConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmHybridConnectionCreateUpdate,
		Read:   resourceArmHybridConnectionRead,
		Update: resourceArmHybridConnectionCreateUpdate,
		Delete: resourceArmHybridConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"relay_namespace_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
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
				ValidateFunc: validate.NoEmptyStrings,
			},
		},
	}
}

func resourceArmHybridConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Relay Hybrid Connection creation.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	relayNamespace := d.Get("relay_namespace_name").(string)
	requireClientAuthroization := d.Get("requires_client_authorization").(bool)
	userMetadata := d.Get("user_metadata").(string)

	parameters := relay.HybridConnection{
		HybridConnectionProperties: &relay.HybridConnectionProperties{
			RequiresClientAuthorization: &requireClientAuthroization,
			UserMetadata:                &userMetadata,
		},
	}

	_, err := client.CreateOrUpdate(ctx, resourceGroup, relayNamespace, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Relay Hybrid Connection %q (Namespace %q Resource Group %q): %+v", name, relayNamespace, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, relayNamespace, name)
	if err != nil {
		return fmt.Errorf("Error issuing get request for Relay Hybrid Connection %q (Namespace %q Resource Group %q): %+v", name, relayNamespace, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read Relay Hybrid Connection %q (Namespace %q Resource group %s) ID", name, relayNamespace, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmHybridConnectionRead(d, meta)
}

func resourceArmHybridConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	relayNamespace := id.Path["namespaces"]
	name := id.Path["hybridConnections"]

	resp, err := client.Get(ctx, resourceGroup, relayNamespace, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Relay Hybrid Connection %q (Namespace %q Resource Group %q): %s", name, relayNamespace, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("relay_namespace_name", relayNamespace)

	if props := resp.HybridConnectionProperties; props != nil {
		d.Set("requires_client_authorization", props.RequiresClientAuthorization)
		d.Set("user_metadata", props.UserMetadata)
	}

	return nil
}

func resourceArmHybridConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	relayNamespace := id.Path["namespaces"]
	name := id.Path["hybridConnections"]

	log.Printf("[INFO] Waiting for Relay Hybrid Connection %q (Namespace %q Resource Group %q) to be deleted", name, relayNamespace, resourceGroup)
	rc, err := client.Delete(ctx, resourceGroup, relayNamespace, name)

	if err != nil {
		if response.WasNotFound(rc.Response) {
			return nil
		}

		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Deleted"},
		Refresh:    hybridConnectionDeleteRefreshFunc(ctx, client, resourceGroup, relayNamespace, name),
		MinTimeout: 15 * time.Second,
	}

	if features.SupportsCustomTimeouts() {
		stateConf.Timeout = d.Timeout(schema.TimeoutDelete)
	} else {
		stateConf.Timeout = 30 * time.Minute
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Relay Hybrid Connection %q (Namespace %q Resource Group %q) to be deleted: %s", name, relayNamespace, resourceGroup, err)
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

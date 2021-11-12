package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/sdk/2017-04-01/hybridconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmRelayHybridConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmRelayHybridConnectionCreateUpdate,
		Read:   resourceArmRelayHybridConnectionRead,
		Update: resourceArmRelayHybridConnectionCreateUpdate,
		Delete: resourceArmRelayHybridConnectionDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := hybridconnections.ParseHybridConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"relay_namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"requires_client_authorization": {
				Type:     pluginsdk.TypeBool,
				Default:  true,
				ForceNew: true,
				Optional: true,
			},
			"user_metadata": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceArmRelayHybridConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Relay Hybrid Connection creation.")

	id := hybridconnections.NewHybridConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("relay_namespace_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_relay_hybrid_connection", id.ID())
		}
	}

	requireClientAuthorization := d.Get("requires_client_authorization").(bool)
	userMetadata := d.Get("user_metadata").(string)

	parameters := hybridconnections.HybridConnection{
		Properties: &hybridconnections.HybridConnectionProperties{
			RequiresClientAuthorization: &requireClientAuthorization,
			UserMetadata:                &userMetadata,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating Hybrid Connection %q (Namespace %q Resource Group %q): %+v", id.Name, id.NamespaceName, id.ResourceGroup, err)
	}

	d.SetId(id.ID())
	return resourceArmRelayHybridConnectionRead(d, meta)
}

func resourceArmRelayHybridConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hybridconnections.ParseHybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("relay_namespace_name", id.NamespaceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("requires_client_authorization", props.RequiresClientAuthorization)
			d.Set("user_metadata", props.UserMetadata)
		}
	}

	return nil
}

func resourceArmRelayHybridConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.HybridConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hybridconnections.ParseHybridConnectionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Hybrid Connection %q (Namespace %q / Resource Group %q): %+v", id.NamespaceName, id.NamespaceName, id.ResourceGroup, err)
	}

	// we can't make use of the Future here due to a bug where 404 isn't tracked as Successful
	log.Printf("[INFO] Waiting for %s to be deleted", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Deleted"},
		Refresh:    hybridConnectionDeleteRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", *id, err)
	}

	return nil
}

func hybridConnectionDeleteRefreshFunc(ctx context.Context, client *hybridconnections.HybridConnectionsClient, id hybridconnections.HybridConnectionId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, "Deleted", nil
			}

			return nil, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return res, "Pending", nil
	}
}

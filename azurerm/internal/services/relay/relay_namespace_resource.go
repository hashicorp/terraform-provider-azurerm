package relay

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/relay/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"

	"github.com/Azure/azure-sdk-for-go/services/relay/mgmt/2017-04-01/relay"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceRelayNamespace() *schema.Resource {
	return &schema.Resource{
		Create: resourceRelayNamespaceCreateUpdate,
		Read:   resourceRelayNamespaceRead,
		Update: resourceRelayNamespaceCreateUpdate,
		Delete: resourceRelayNamespaceDelete,
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.NamespaceID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(6, 50),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(relay.Standard),
				}, false),
			},

			"metric_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"primary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_connection_string": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"primary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceRelayNamespaceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.NamespacesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Relay Namespace create/update.")

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	resourceId := parse.NewNamespaceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Relay Namespace %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_relay_namespace", resourceId.ID())
		}
	}

	parameters := relay.Namespace{
		Location: utils.String(location),
		Sku: &relay.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
			Tier: relay.SkuTier(d.Get("sku_name").(string)),
		},
		NamespaceProperties: &relay.NamespaceProperties{},
		Tags:                expandedTags,
	}

	future, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating Relay Namespace %q (Resource Group %q): %+v", resourceId.Name, resourceId.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for create/update of Relay Namespace %q (Resource Group %q) creation: %+v", resourceId.Name, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())
	return resourceRelayNamespaceRead(d, meta)
}

func resourceRelayNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.NamespacesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Relay Namespace %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
	}

	keysResp, err := client.ListKeys(ctx, id.ResourceGroup, id.Name, "RootManageSharedAccessKey")
	if err != nil {
		return fmt.Errorf("listing keys for Relay Namespace %q (Resource Group %q): %s", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if props := resp.NamespaceProperties; props != nil {
		d.Set("metric_id", props.MetricID)
	}

	d.Set("primary_connection_string", keysResp.PrimaryConnectionString)
	d.Set("primary_key", keysResp.PrimaryKey)
	d.Set("secondary_connection_string", keysResp.SecondaryConnectionString)
	d.Set("secondary_key", keysResp.SecondaryKey)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceRelayNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Relay.NamespacesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NamespaceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("")
	}

	// we can't make use of the Future here due to a bug where 404 isn't tracked as Successful
	log.Printf("[DEBUG] Waiting for Relay Namespace %q (Resource Group %q) to be deleted", id.Name, id.ResourceGroup)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Pending"},
		Target:     []string{"Deleted"},
		Refresh:    relayNamespaceDeleteRefreshFunc(ctx, client, id.ResourceGroup, id.Name),
		MinTimeout: 15 * time.Second,
		Timeout:    d.Timeout(schema.TimeoutDelete),
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Relay Namespace %q (Resource Group %q) to be deleted: %s", id.Name, id.ResourceGroup, err)
	}

	return nil
}

func relayNamespaceDeleteRefreshFunc(ctx context.Context, client *relay.NamespacesClient, resourceGroupName string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name)
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}

			return nil, "Error", fmt.Errorf("Error issuing read request in relayNamespaceDeleteRefreshFunc to Relay Namespace %q (Resource Group %q): %s", name, resourceGroupName, err)
		}

		return res, "Pending", nil
	}
}

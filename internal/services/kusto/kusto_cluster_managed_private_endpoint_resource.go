package kusto

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/kusto/mgmt/2022-02-01/kusto" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoClusterManagedPrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoClusterManagedPrivateEndpointCreateUpdate,
		Read:   resourceKustoClusterManagedPrivateEndpointRead,
		Update: resourceKustoClusterManagedPrivateEndpointCreateUpdate,
		Delete: resourceKustoClusterManagedPrivateEndpointDelete,

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoManagedPrivateEndpointV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedPrivateEndpointsID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"resource_group_name": commonschema.ResourceGroupName(),

			"cluster_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ClusterName,
			},

			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"private_link_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"private_link_resource_region": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"request_message": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceKustoClusterManagedPrivateEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewManagedPrivateEndpointsID(subscriptionId, d.Get("resource_group_name").(string), d.Get("cluster_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		managedPrivateEndpoint, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.ManagedPrivateEndpointName)
		if err != nil {
			if !utils.ResponseWasNotFound(managedPrivateEndpoint.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(managedPrivateEndpoint.Response) {
			return tf.ImportAsExistsError("azurerm_kusto_cluster_managed_private_endpoint", id.ID())
		}
	}

	managedPrivateEndpoint := kusto.ManagedPrivateEndpoint{
		ManagedPrivateEndpointProperties: &kusto.ManagedPrivateEndpointProperties{
			PrivateLinkResourceID: utils.String(d.Get("private_link_resource_id").(string)),
			GroupID:               utils.String(d.Get("group_id").(string)),
		},
	}

	if v, ok := d.GetOk("private_link_resource_region"); ok {
		managedPrivateEndpoint.PrivateLinkResourceRegion = utils.String(v.(string))
	}

	if v, ok := d.GetOk("request_message"); ok {
		managedPrivateEndpoint.RequestMessage = utils.String(v.(string))
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ClusterName, id.ManagedPrivateEndpointName, managedPrivateEndpoint)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoClusterManagedPrivateEndpointRead(d, meta)
}

func resourceKustoClusterManagedPrivateEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedPrivateEndpointsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ClusterName, id.ManagedPrivateEndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ManagedPrivateEndpointName)
	d.Set("cluster_name", id.ClusterName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("private_link_resource_id", resp.PrivateLinkResourceID)
	d.Set("group_id", resp.GroupID)

	if resp.PrivateLinkResourceRegion != nil {
		d.Set("private_link_resource_region", resp.PrivateLinkResourceRegion)
	}

	if resp.RequestMessage != nil {
		d.Set("request_message", resp.RequestMessage)
	}

	return nil
}

func resourceKustoClusterManagedPrivateEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ClusterManagedPrivateEndpointClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedPrivateEndpointsID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ClusterName, id.ManagedPrivateEndpointName)

	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

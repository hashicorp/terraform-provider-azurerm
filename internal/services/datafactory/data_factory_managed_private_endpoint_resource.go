package datafactory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/validate"
	networkParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDataFactoryManagedPrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDataFactoryManagedPrivateEndpointCreate,
		Read:   resourceDataFactoryManagedPrivateEndpointRead,
		Delete: resourceDataFactoryManagedPrivateEndpointDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ManagedPrivateEndpointID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryManagedPrivateEndpointName(),
			},

			"data_factory_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DataFactoryID,
			},

			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"subresource_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.PrivateLinkSubResourceName,
			},

			"fqdns": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func resourceDataFactoryManagedPrivateEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpointsClient
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := parse.DataFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	managedVirtualNetworkName, err := getManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, dataFactoryId.ResourceGroup, dataFactoryId.FactoryName)
	if err != nil {
		return err
	}
	if managedVirtualNetworkName == nil {
		return fmt.Errorf("managed Private endpoints are only available after managed virtual network for %s is enabled", dataFactoryId)
	}

	id := parse.NewManagedPrivateEndpointID(subscriptionId, dataFactoryId.ResourceGroup, dataFactoryId.FactoryName, *managedVirtualNetworkName, d.Get("name").(string))
	existing, err := getManagedPrivateEndpoint(ctx, client, id.ResourceGroup, id.FactoryName, *managedVirtualNetworkName, id.Name)
	if err != nil {
		return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
	}
	if existing != nil {
		return tf.ImportAsExistsError("azurerm_data_factory_managed_private_endpoint", id.ID())
	}

	targetResourceId := d.Get("target_resource_id").(string)
	subResourceName := d.Get("subresource_name").(string)
	fqdns := d.Get("fqdns").([]interface{})

	if _, err := networkParse.PrivateLinkServiceID(targetResourceId); err == nil {
		if len(subResourceName) > 0 {
			return fmt.Errorf("`subresource_name` should not be specified when target resource is `Private Link Service`")
		}

		if len(fqdns) == 0 {
			return fmt.Errorf("`fqdns` should be specified when target resource is `Private Link Service`")
		}
	} else {
		if len(strings.TrimSpace(subResourceName)) < 3 {
			return fmt.Errorf("`subresource_name` must be at least 3 character in length")
		}

		if len(fqdns) > 0 {
			return fmt.Errorf("`fqdns` should not be specified for the target resource: %q", targetResourceId)
		}
	}

	managedPrivateEndpoint := datafactory.ManagedPrivateEndpointResource{
		Properties: &datafactory.ManagedPrivateEndpoint{
			PrivateLinkResourceID: utils.String(targetResourceId),
		},
	}

	if len(subResourceName) > 0 {
		managedPrivateEndpoint.Properties.GroupID = utils.String(subResourceName)
	}

	if len(fqdns) > 0 {
		managedPrivateEndpoint.Properties.Fqdns = utils.ExpandStringSlice(fqdns)
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.FactoryName, id.ManagedVirtualNetworkName, id.Name, managedPrivateEndpoint, ""); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Provisioning"},
		Target:     []string{"Succeeded"},
		Refresh:    getManagedPrivateEndpointProvisionStatus(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    d.Timeout(pluginsdk.TimeoutDelete),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be created: %+v", id.ID(), err)
	}

	d.SetId(id.ID())

	return resourceDataFactoryManagedPrivateEndpointRead(d, meta)
}

func resourceDataFactoryManagedPrivateEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpointsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.ManagedVirtualNetworkName, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("data_factory_id", parse.NewDataFactoryID(subscriptionId, id.ResourceGroup, id.FactoryName).ID())

	if props := resp.Properties; props != nil {
		d.Set("target_resource_id", props.PrivateLinkResourceID)
		d.Set("subresource_name", props.GroupID)
		if props.Fqdns != nil {
			d.Set("fqdns", props.Fqdns)
		}
	}

	return nil
}

func resourceDataFactoryManagedPrivateEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpointsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.FactoryName, id.ManagedVirtualNetworkName, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

// if ManagedPrivateEndpoint not exist, get rest api will return 400 bad request
// invoke list rets api and then filter by name
func getManagedPrivateEndpoint(ctx context.Context, client *datafactory.ManagedPrivateEndpointsClient, resourceGroupName, factoryName, managedVirtualNetworkName, name string) (*datafactory.ManagedPrivateEndpointResource, error) {
	iter, err := client.ListByFactoryComplete(ctx, resourceGroupName, factoryName, managedVirtualNetworkName)
	if err != nil {
		return nil, err
	}
	for iter.NotDone() {
		managedPrivateEndpoint := iter.Value()
		if managedPrivateEndpoint.Name != nil && *managedPrivateEndpoint.Name == name {
			return &managedPrivateEndpoint, nil
		}

		if err := iter.NextWithContext(ctx); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func getManagedPrivateEndpointProvisionStatus(ctx context.Context, client *datafactory.ManagedPrivateEndpointsClient, id parse.ManagedPrivateEndpointId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id.ResourceGroup, id.FactoryName, id.ManagedVirtualNetworkName, id.Name, "")
		if err != nil || resp.Properties == nil || resp.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return resp, *resp.Properties.ProvisioningState, nil
	}
}

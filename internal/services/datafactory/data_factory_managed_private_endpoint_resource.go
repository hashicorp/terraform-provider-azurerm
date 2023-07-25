// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/managedprivateendpoints"
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
				ValidateFunc: factories.ValidateFactoryID,
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
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpoints
	managedVirtualNetworksClient := meta.(*clients.Client).DataFactory.ManagedVirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	dataFactoryId, err := factories.ParseFactoryID(d.Get("data_factory_id").(string))
	if err != nil {
		return err
	}

	managedVirtualNetworkName, err := getManagedVirtualNetworkName(ctx, managedVirtualNetworksClient, dataFactoryId.SubscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName)
	if err != nil {
		return err
	}
	if managedVirtualNetworkName == nil {
		return fmt.Errorf("managed Private endpoints are only available after managed virtual network for %s is enabled", dataFactoryId)
	}

	id := managedprivateendpoints.NewManagedPrivateEndpointID(subscriptionId, dataFactoryId.ResourceGroupName, dataFactoryId.FactoryName, *managedVirtualNetworkName, d.Get("name").(string))
	existing, err := getManagedPrivateEndpoint(ctx, client, id.SubscriptionId, id.ResourceGroupName, id.FactoryName, id.ManagedVirtualNetworkName, id.ManagedPrivateEndpointName)
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

	payload := managedprivateendpoints.ManagedPrivateEndpointResource{
		Properties: managedprivateendpoints.ManagedPrivateEndpoint{
			PrivateLinkResourceId: utils.String(targetResourceId),
		},
	}

	if len(subResourceName) > 0 {
		payload.Properties.GroupId = utils.String(subResourceName)
	}

	if len(fqdns) > 0 {
		payload.Properties.Fqdns = utils.ExpandStringSlice(fqdns)
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload, managedprivateendpoints.DefaultCreateOrUpdateOperationOptions()); err != nil {
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
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpoints
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, managedprivateendpoints.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ManagedPrivateEndpointName)
	d.Set("data_factory_id", factories.NewFactoryID(subscriptionId, id.ResourceGroupName, id.FactoryName).ID())

	if model := resp.Model; model != nil {
		props := model.Properties
		d.Set("target_resource_id", props.PrivateLinkResourceId)
		d.Set("subresource_name", props.GroupId)
		d.Set("fqdns", utils.FlattenStringSlice(props.Fqdns))
	}

	return nil
}

func resourceDataFactoryManagedPrivateEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataFactory.ManagedPrivateEndpoints
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := managedprivateendpoints.ParseManagedPrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

// if ManagedPrivateEndpoint not exist, get rest api will return 400 bad request
// invoke list rets api and then filter by name
func getManagedPrivateEndpoint(ctx context.Context, client *managedprivateendpoints.ManagedPrivateEndpointsClient, subscriptionId, resourceGroupName, factoryName, managedVirtualNetworkName, privateEndpointName string) (*managedprivateendpoints.ManagedPrivateEndpointResource, error) {
	managedVirtualNetworkId := managedprivateendpoints.NewManagedVirtualNetworkID(subscriptionId, resourceGroupName, factoryName, managedVirtualNetworkName)
	iter, err := client.ListByFactoryComplete(ctx, managedVirtualNetworkId)
	if err != nil {
		return nil, err
	}

	for _, item := range iter.Items {
		if item.Name != nil && *item.Name == privateEndpointName {
			return &item, nil
		}

	}
	return nil, nil
}

func getManagedPrivateEndpointProvisionStatus(ctx context.Context, client *managedprivateendpoints.ManagedPrivateEndpointsClient, id managedprivateendpoints.ManagedPrivateEndpointId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// TODO: it should be possible to remove this function https://github.com/hashicorp/go-azure-sdk/issues/307 has been fixed
		resp, err := client.Get(ctx, id, managedprivateendpoints.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		model := resp.Model
		if model == nil {
			return nil, "", fmt.Errorf("retrieving %s: `model` was nil", id)
		}

		if model.Properties.ProvisioningState == nil {
			return nil, "", fmt.Errorf("retrieving %s: `provisioningState` is nil", id)
		}

		return resp, *resp.Model.Properties.ProvisioningState, nil
	}
}

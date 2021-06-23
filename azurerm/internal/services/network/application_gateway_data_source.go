package network

import (
	"fmt"
	"time"

	msiparse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/identity"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type applicationGatewayDataSourceIdentity = identity.UserAssigned

func dataSourceApplicationGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApplicationGatewayRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"identity": applicationGatewayDataSourceIdentity{}.SchemaDataSource(),

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceApplicationGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ApplicationGatewaysClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewApplicationGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("location", location.NormalizeNilable(resp.Location))

	identity, err := flattenApplicationGatewayDataSourceIdentity(resp.Identity)
	if err != nil {
		return err
	}
	flattenedIdentity := applicationGatewayDataSourceIdentity{}.Flatten(identity)
	if err = d.Set("identity", flattenedIdentity); err != nil {
		return err
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func flattenApplicationGatewayDataSourceIdentity(input *network.ManagedServiceIdentity) (*identity.ExpandedConfig, error) {
	var config *identity.ExpandedConfig
	if input != nil {
		identityIds := make([]string, 0, len(input.UserAssignedIdentities))
		for id := range input.UserAssignedIdentities {
			parsedId, err := msiparse.UserAssignedIdentityIDInsensitively(id)
			if err != nil {
				return nil, err
			}
			identityIds = append(identityIds, parsedId.ID())
		}
		config = &identity.ExpandedConfig{
			Type:                    string(input.Type),
			UserAssignedIdentityIds: &identityIds,
		}
	}
	return config, nil
}

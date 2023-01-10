package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverDnsResolverDataSourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	VirtualNetworkId  string            `tfschema:"virtual_network_id"`
}

type PrivateDNSResolverDnsResolverDataSource struct{}

var _ sdk.DataSource = PrivateDNSResolverDnsResolverDataSource{}

func (r PrivateDNSResolverDnsResolverDataSource) ResourceType() string {
	return "azurerm_private_dns_resolver"
}

func (r PrivateDNSResolverDnsResolverDataSource) ModelObject() interface{} {
	return &PrivateDNSResolverDnsResolverDataSourceModel{}
}

func (r PrivateDNSResolverDnsResolverDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnsresolvers.ValidateDnsResolverID
}

func (r PrivateDNSResolverDnsResolverDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),
	}
}

func (r PrivateDNSResolverDnsResolverDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"virtual_network_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
}

func (r PrivateDNSResolverDnsResolverDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolversClient

			var state PrivateDNSResolverDnsResolverDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			var top int64 = 1
			resp, err := client.ListByResourceGroupCompleteMatchingPredicate(ctx,
				commonids.NewResourceGroupID(metadata.Client.Account.SubscriptionId, state.ResourceGroupName),
				dnsresolvers.ListByResourceGroupOperationOptions{Top: &top},
				dnsresolvers.DnsResolverOperationPredicate{Name: &state.Name})
			if err != nil || len(resp.Items) != int(top) {
				return fmt.Errorf("retrieving %s: %+v", state.Name, err)
			}

			model := resp.Items[0]
			id, err := dnsresolvers.ParseDnsResolverID(*model.Id)
			if err != nil {
				return err
			}

			state.Location = location.Normalize(model.Location)
			state.VirtualNetworkId = model.Properties.VirtualNetwork.Id

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}

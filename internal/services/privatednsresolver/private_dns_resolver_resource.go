// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatednsresolver

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PrivateDNSResolverDnsResolverModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	Tags              map[string]string `tfschema:"tags"`
	VirtualNetworkId  string            `tfschema:"virtual_network_id"`
}

type PrivateDNSResolverDnsResolverResource struct{}

var _ sdk.ResourceWithUpdate = PrivateDNSResolverDnsResolverResource{}

func (r PrivateDNSResolverDnsResolverResource) ResourceType() string {
	return "azurerm_private_dns_resolver"
}

func (r PrivateDNSResolverDnsResolverResource) ModelObject() interface{} {
	return &PrivateDNSResolverDnsResolverModel{}
}

func (r PrivateDNSResolverDnsResolverResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return dnsresolvers.ValidateDnsResolverID
}

func (r PrivateDNSResolverDnsResolverResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"virtual_network_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r PrivateDNSResolverDnsResolverResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PrivateDNSResolverDnsResolverResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model PrivateDNSResolverDnsResolverModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.PrivateDnsResolver.DnsResolversClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := dnsresolvers.NewDnsResolverID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &dnsresolvers.DnsResolver{
				Location: location.Normalize(model.Location),
				Properties: dnsresolvers.DnsResolverProperties{
					VirtualNetwork: dnsresolvers.SubResource{
						Id: model.VirtualNetworkId,
					},
				},
				Tags: &model.Tags,
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *properties, dnsresolvers.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PrivateDNSResolverDnsResolverResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolversClient

			id, err := dnsresolvers.ParseDnsResolverID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateDNSResolverDnsResolverModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *properties, dnsresolvers.CreateOrUpdateOperationOptions{}); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PrivateDNSResolverDnsResolverResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolversClient

			id, err := dnsresolvers.ParseDnsResolverID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := PrivateDNSResolverDnsResolverModel{
				Name:              id.DnsResolverName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			properties := &model.Properties
			if properties != nil {
				state.VirtualNetworkId = properties.VirtualNetwork.Id
			}
			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PrivateDNSResolverDnsResolverResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.PrivateDnsResolver.DnsResolversClient

			id, err := dnsresolvers.ParseDnsResolverID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, dnsresolvers.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
